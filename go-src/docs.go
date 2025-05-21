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

func writeStoredFile(filename string, data string) error {
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

	err = writeStoredFile(filename, body)
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

// body, is_auth_needed
func GetPostBody(db *sql.DB, id int) (string, bool, error) {
	var filename string
	var draft bool
	var archived bool
	result := db.QueryRow(
		"SELECT posts.source_ref, posts.draft, posts.archived "+
			"FROM posts WHERE posts.id = $1;", id)
	err := result.Scan(&filename, &draft, &archived)
	if err != nil {
		return "", true, err
	}

	body, err := readStoredFile(filename)
	if err != nil {
		return "", true, err
	}

	needs_auth := draft || archived

	return body, needs_auth, nil
}

func ModifyPost(db *sql.DB, id int, title *string, draft *bool, archived *bool,
	tags *[]string, body *string) error {
	if title == nil && draft == nil && archived == nil &&
		tags == nil && body == nil {
		return nil
	}

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		tx.Rollback()
		return err
	}

	var title_value string
	var draft_value bool
	var archived_value bool
	var filename string
	result := tx.QueryRow(
		"SELECT posts.title, posts.draft, posts.archived, posts.source_ref "+
			"FROM posts WHERE posts.id = $1;", id)
	err = result.Scan(&title_value, &draft_value, &archived_value, &filename)
	if err != nil {
		tx.Rollback()
		return err
	}

	if title != nil {
		title_value = *title
	}

	if draft != nil {
		draft_value = *draft
	}

	if archived != nil {
		archived_value = *archived
	}

	timestamp := time.Now()

	_, err = tx.Exec("UPDATE posts SET modified_timestamp = $1, title = $2, "+
		"draft = $3, archived = $4 WHERE id = $5;", timestamp, title_value,
		draft_value, archived_value, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if body != nil {
		err = writeStoredFile(filename, *body)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if tags != nil {
		tagIds, err := getTagIds(*tags, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		// remove tags
		_, err = tx.Exec("DELETE FROM tags_to_posts WHERE post_id = $1;", id)
		if err != nil {
			tx.Rollback()
			return err
		}

		// add tags
		for _, tagId := range tagIds {
			_, err = tx.Exec("INSERT INTO tags_to_posts (tag_id, post_id) "+
				"VALUES ($1, $2);", tagId, id)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func DeletePost(db *sql.DB, id int) error {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM tags_to_posts WHERE post_id = $1;", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	result, err := tx.Exec("DELETE FROM posts WHERE id = $1;", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if n == 0 {
		tx.Rollback()
		return sql.ErrNoRows
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
