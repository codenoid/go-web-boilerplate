package main

import (
	"html/template"
	"os"

	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.Static("/assets", "./public/assets")
	r.Use(mainMiddleware)

	r.HTMLRender = ginview.New(goview.Config{
		Root:         "views",
		Extension:    ".html",
		Master:       "layouts/master",
		Funcs:        template.FuncMap{},
		DisableCache: true,
	})

	r.GET("/auth/login", loginHTML)
	r.POST("/auth/login", loginVerify)
	r.GET("/auth/logout", logout)

	// root
	r.GET("/", homescreenHTML)

	r.Run(os.Getenv("BIND_ADDR")) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
