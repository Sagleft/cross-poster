package main

import (
	"log"
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
	sol.Gin.GET("/pairs", func(c *gin.Context) {
		pairs, err := sol.getExchangePairs()
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to get pairs: "+err.Error())
			return
		}
		log.Println(pairs)

		sol.renderTemplate(
			c,
			http.StatusOK,
			"pairs.html",
			gin.H{
				"pairs": pairs,
			},
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
