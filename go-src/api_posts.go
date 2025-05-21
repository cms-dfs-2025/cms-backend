package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (handler ServerContext) HandlePostsGetAll(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

func (handler ServerContext) HandlePostsGetBody(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}
