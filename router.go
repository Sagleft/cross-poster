package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/tucnak/telebot.v2"
)

func (sol *solution) setupRoutes() error {
	sol.loadTemplates("./templates/*")
	sol.Gin.Static("/assets", "./assets")

	sol.Gin.GET("/", func(c *gin.Context) {
		sol.renderTemplate(
			c,
			http.StatusOK,
			"home.html",
			sol.getDefaultRequestHeaders(),
		)
	})
	sol.Gin.NoRoute(func(c *gin.Context) {
		sol.renderTemplate(
			c,
			http.StatusNotFound,
			"404.html",
			sol.getDefaultRequestHeaders(),
		)
	})
	sol.Gin.POST("/send", func(c *gin.Context) {
		sol.handleMessageRequest(c)
	})

	return nil
}

func (sol *solution) handleMessageRequest(c *gin.Context) {
	/*var msg message
	err := c.BindJSON(&msg)
	if err != nil {
		handleRequestError(c, err)
		return
	}*/

	msg := c.PostForm("post_text")
	if msg == "" {
		handleRequestError(c, errors.New("Post message is empty"))
		return
	}

	sendToTelegram := c.PostForm("post_telegram") == "1"
	sendToUtopia := c.PostForm("post_utopia") == "1"
	if !sendToTelegram && !sendToUtopia {
		handleRequestError(c, errors.New("No messenger is selected"))
		return
	}

	if sendToTelegram {
		sol.sendTelegramPost(msg, c)
	}
	if sendToUtopia {
		sol.sendUtopiaPost(msg, c)
	}
	handleRequestSuccess(c)
}

func (sol *solution) sendTelegramPost(postText string, c *gin.Context) {
	_, err := sol.Messengers.Telegram.Send(telebot.ChatID(sol.Config.Telegram.ChatID), postText, telebot.ModeMarkdown)
	if err != nil {
		handleRequestError(c, err)
	}
}

func (sol *solution) sendUtopiaPost(postText string, c *gin.Context) {
	_, err := sol.Messengers.Utopia.SendChannelMessage(sol.Config.Utopia.ChannelID, postText)
	if err != nil {
		handleRequestError(c, err)
	}
}

func handleRequestError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, response{
		Status:    "error",
		ErrorInfo: err.Error(),
	})
}

func handleRequestSuccess(c *gin.Context) {
	c.JSON(http.StatusOK, response{
		Status: "success",
	})
}

func (sol *solution) getDefaultRequestHeaders() gin.H {
	return gin.H{
		"version": sol.Config.FrontEnd.Version,
	}
}

func (sol *solution) loadTemplates(pattern string) {
	sol.Gin.LoadHTMLGlob(pattern)
}

func (sol *solution) renderTemplate(c *gin.Context, code int, name string, obj interface{}) {
	c.HTML(code, name, obj)
}
