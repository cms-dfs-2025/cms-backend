package main

import (
	"context"
	"database/sql"
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
