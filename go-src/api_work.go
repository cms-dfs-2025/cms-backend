package main

import (
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

func (handler ServerContext) HandleWorkGetAll(c *gin.Context) {
	allPosts, err := GetAllPosts(handler.db)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": allPosts})
}
