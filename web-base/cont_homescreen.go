package main

import (
	"github.com/codenoid/go-web-boilerplate/web-base/structs/user"
	"github.com/gin-gonic/gin"
)

func homescreenHTML(c *gin.Context) {
	userData, _ := c.Get("userData")
	c.HTML(200, "homescreen", gin.H{
		"user": userData.(user.User),
	})
}
