package main

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
)

func storedFilePath(filename string) string {
	dir_path := os.Getenv("DOCS_PATH")
	return path.Join(dir_path, filename)
}

func createStoredFile(filename string, data string) error {
	filepath := storedFilePath(filename)
	err := os.WriteFile(filepath, []byte(data), 0777)
	return err
}

func readStoredFile(filename string) (string, error) {
	filepath := storedFilePath(filename)
	data, err := os.ReadFile(filepath)
	return string(data), err
}

func getTagIds(tags []string, tx *sql.Tx) ([]int, error) {
	tagIds := []int{}

	for _, tag := range tags {
		checkResult := tx.QueryRow("SELECT id FROM tags WHERE name = $1;", tag)

		var tagId int
		err := checkResult.Scan(&tagId)

		if err != nil {
			if err == sql.ErrNoRows {
				_, err := tx.Exec("INSERT INTO tags (name) VALUES ($1);", tag)
				if err != nil {
					return []int{}, err
				}

				checkResult = tx.QueryRow("SELECT id FROM tags WHERE name = $1;", tag)
				err = checkResult.Scan(&tagId)
				if err != nil {
					return []int{}, err
				}
			} else {
				return []int{}, err
			}
		}

		tagIds = append(tagIds, tagId)
	}

	return tagIds, nil
}

func UploadPost(authorId int, title string, tags []string, draft bool,
	archived bool, body string, db *sql.DB) (int, error) {

	// --- create the file ---

	fileId, err := uuid.NewRandom()
	if err != nil {
		return -1, err
	}
	filename := fileId.String() + ".md"

	err = createStoredFile(filename, body)
	if err != nil {
		return -1, err
	}

	// --- create database entry ---

	// begin transaction
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	// upload timestamp
	timestamp := time.Now()

	// tags relations
	tagIds, err := getTagIds(tags, tx)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	// insert into posts
	_, err = tx.
		Exec("INSERT INTO posts (author_id, upload_timestamp, "+
			"modified_timestamp, draft, archived, title, source_ref) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7)",
			authorId, timestamp, timestamp, draft, archived, title, filename)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	// get post id
	postIdResult := tx.
		QueryRow("SELECT id FROM posts WHERE source_ref = $1;", filename)
	var postId int
	err = postIdResult.Scan(&postId)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	// insert into tags
	for _, tagId := range tagIds {
		_, err = tx.Exec("INSERT INTO tags_to_posts (tag_id, post_id) "+
			"VALUES ($1, $2);", tagId, postId)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	// end transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	return postId, nil
}

type PostData struct {
	Id           int    `json:"id"`
	AuthorHandle string `json:"author_handle"`

	Title string   `json:"title"`
	Tags  []string `json:"tags"`

	UploadTimestamp   string `json:"upload_timestamp"`
	ModifiedTimestamp string `json:"modified_timestamp"`

	Draft    bool `json:"draft"`
	Archived bool `json:"archived"`
}

func GetAllPosts(db *sql.DB) ([]PostData, error) {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		tx.Rollback()
		return []PostData{}, err
	}

	// --- get all the data ---

	allPosts := make(map[int]PostData)

	allPostsRows, err := tx.
		Query("select posts.id, users.handle, posts.title, " +
			"posts.upload_timestamp, posts.modified_timestamp, " +
			"posts.draft, posts.archived from posts join users " +
			"on posts.author_id = users.id;")
	defer allPostsRows.Close()
	if err != nil {
		tx.Rollback()
		return []PostData{}, err
	}

	for allPostsRows.Next() {
		var postData PostData

		var uploadTimestamp time.Time
		var modifiedTimestamp time.Time

		err = allPostsRows.
			Scan(&postData.Id, &postData.AuthorHandle, &postData.Title,
				&uploadTimestamp, &modifiedTimestamp,
				&postData.Draft, &postData.Archived)

		postData.UploadTimestamp = uploadTimestamp.Format(time.RFC3339)
		postData.ModifiedTimestamp = modifiedTimestamp.Format(time.RFC3339)

		if err != nil {
			tx.Rollback()
			return []PostData{}, err
		}

		allPosts[postData.Id] = postData
	}

	allTagsRows, err := tx.
		Query("select tags_to_posts.post_id, tags.name from tags_to_posts " +
			"join tags on tags_to_posts.tag_id=tags.id;")
	defer allTagsRows.Close()
	if err != nil {
		tx.Rollback()
		return []PostData{}, err
	}

	for allTagsRows.Next() {
		var postId int
		var tag string

		err = allTagsRows.Scan(&postId, &tag)
		if err != nil {
			tx.Rollback()
			return []PostData{}, err
		}

		postData, ok := allPosts[postId]
		if !ok {
			tx.Rollback()
			return []PostData{}, errors.New("cms: id mismatch in posts")
		}
		postData.Tags = append(postData.Tags, tag)
		allPosts[postId] = postData
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return []PostData{}, err
	}

	allPostsList := make([]PostData, 0, len(allPosts))
	for _, value := range allPosts {
		allPostsList = append(allPostsList, value)
	}
	return allPostsList, nil
}
