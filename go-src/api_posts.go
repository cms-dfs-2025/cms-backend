package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type publicPostData struct {
	Id    int      `json:"id"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}

func (handler ServerContext) HandlePostsGetAll(c *gin.Context) {
	allPosts, err := GetAllPosts(handler.db)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	filteredAllPosts := make([]publicPostData, 0, len(allPosts))
	for _, value := range allPosts {
		needs_auth := value.Archived || value.Draft

		if !needs_auth {
			filteredAllPosts = append(filteredAllPosts, publicPostData{
				Id:    value.Id,
				Title: value.Title,
				Tags:  value.Tags,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"posts": filteredAllPosts})
}

func (handler ServerContext) HandlePostsGetBody(c *gin.Context) {
	var body getBodyBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	data, needs_auth, err := GetPostBody(handler.db, *body.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusNotFound)
			return
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	if needs_auth {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{"body": data})
}
