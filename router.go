package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (sol *solution) setupRoutes() error {
	sol.loadTemplates("./templates/*")
	sol.Gin.Static("/assets", "./assets")

	sol.Gin.GET("/", func(c *gin.Context) {
		sol.renderTemplate(
			c,
			http.StatusOK,
			"home.html",
			gin.H{},
		)
	})
	sol.Gin.NoRoute(func(c *gin.Context) {
		sol.renderTemplate(
			c,
			http.StatusNotFound,
			"404.html",
			gin.H{},
		)
	})

	return nil
}

func (sol *solution) loadTemplates(pattern string) {
	sol.Gin.LoadHTMLGlob(pattern)
}

func (sol *solution) renderTemplate(c *gin.Context, code int, name string, obj interface{}) {
	c.HTML(code, name, obj)
}
