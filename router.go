package main

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

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
	sol.Gin.POST("/upload", func(c *gin.Context) {
		fileHandler, err := c.FormFile("files[]")
		if err != nil {
			handleRequestError(c, errors.New("failed to get file from context: "+err.Error()))
			return
		}
		fileExtension := filepath.Ext(fileHandler.Filename)
		switch fileExtension {
		default:
			handleRequestError(c, errors.New("unknown file extension: "+fileExtension))
			return
		case ".png":
			break
		case ".jpg":
			break
		case ".jpeg":
			break
		}

		imageFilename := "image" + fileExtension
		err = c.SaveUploadedFile(fileHandler, imageFilename)
		if err != nil {
			handleRequestError(c, errors.New("failed to save uploaded file: "+err.Error()))
			return
		}
		handleRequestSuccess(c, imageFilename)
	})

	return nil
}

func (sol *solution) handleMessageRequest(c *gin.Context) {
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
	imageFilename := c.PostForm("imagename")
	if strings.Contains(imageFilename, "/") {
		handleRequestError(c, errors.New("invail image name given"))
		return
	}

	if sendToTelegram {
		sol.sendTelegramPost(msg, imageFilename, c)
	}
	if sendToUtopia {
		sol.sendUtopiaPost(msg, imageFilename, c)
	}
	handleRequestSuccess(c, nil)
}

func (sol *solution) sendTelegramPost(postText string, imageFilename string, c *gin.Context) {
	var msg interface{}
	if imageFilename != "" {
		msg = &telebot.Photo{
			File:    telebot.FromDisk(imageFilename),
			Caption: postText,
		}
	} else {
		msg = postText
	}

	_, err := sol.Messengers.Telegram.Send(
		telebot.ChatID(sol.Config.Telegram.ChatID),
		msg,
		telebot.ModeMarkdown,
	)
	if err != nil {
		handleRequestError(c, errors.New("failed to send post to Telegram: "+err.Error()))
	}
}

func (sol *solution) sendUtopiaPost(postText string, imageFilename string, c *gin.Context) {
	if imageFilename != "" {
		imageBytes, err := ioutil.ReadFile(imageFilename)
		if err != nil {
			handleRequestError(c, errors.New("failed to read uploaded image: "+err.Error()))
			return
		}

		_, err = sol.Messengers.Utopia.SendChannelPicture(
			sol.Config.Utopia.ChannelID,
			base64.StdEncoding.EncodeToString(imageBytes),
			postText,
			"image.jpg",
		)
		if err != nil {
			handleRequestError(c, errors.New("failed to send post with image to Utopia: "+err.Error()))
			return
		}
	}

	// send plain text
	_, err := sol.Messengers.Utopia.SendChannelMessage(sol.Config.Utopia.ChannelID, postText)
	if err != nil {
		handleRequestError(c, errors.New("failed to send post to Utopia: "+err.Error()))
	}
}

func handleRequestError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, response{
		Status:    "error",
		ErrorInfo: err.Error(),
	})
}

func handleRequestSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, response{
		Status: "success",
		Data:   data,
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
