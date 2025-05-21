package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type uploadBody struct {
	Title    *string   `json:"title" binding:"required"`
	Tags     *[]string `json:"tags" binding:"required"`
	Draft    *bool     `json:"draft" binding:"required"`
	Archived *bool     `json:"archived" binding:"required"`
	Body     *string   `json:"body" binding:"required"`
}

func (handler ServerContext) HandleWorkUpload(c *gin.Context) {
	var body uploadBody
	err := c.ShouldBindJSON(&body)

	if err != nil {
		log.Printf("/api/work/upload error: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}

	user := c.MustGet("user").(UserRow)
	postId, err := UploadPost(user.Id, *body.Title, *body.Tags, *body.Draft,
		*body.Archived, *body.Body, handler.db)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": postId})
}

type modifyBody struct {
	Id       *int      `json:"id" binding:"required"`
	Title    *string   `json:"title"`
	Tags     *[]string `json:"tags"`
	Draft    *bool     `json:"draft"`
	Archived *bool     `json:"archived"`
	Body     *string   `json:"body"`
}

func (handler ServerContext) HandleWorkModify(c *gin.Context) {
	var body modifyBody
	err := c.ShouldBindJSON(&body)

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = ModifyPost(handler.db, *body.Id, body.Title, body.Draft, body.Archived,
		body.Tags, body.Body)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "Post with id not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "Internal error"})
		}
	}

	c.Status(http.StatusOK)
}

type idBody struct {
	Id *int `json:"id" binding:"required"`
}

func (handler ServerContext) HandleWorkDelete(c *gin.Context) {
	var body idBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = DeletePost(handler.db, *body.Id)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "Post with id not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "Internal error"})
		}
	}

	c.Status(http.StatusOK)
}

func (handler ServerContext) HandleWorkGetBody(c *gin.Context) {
	var body idBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	data, _, err := GetPostBody(handler.db, *body.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusNotFound)
			return
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"body": data})
}

func (handler ServerContext) HandleWorkGetAll(c *gin.Context) {
	allPosts, err := GetAllPosts(handler.db)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": allPosts})
}
