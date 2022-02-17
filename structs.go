package main

import (
	"github.com/gin-gonic/gin"
)

type solution struct {
	Config appConfig
	Gin    *gin.Engine
}

type appConfig struct {
	Utopia   utopiaConfig   `json:"utopia"`
	Telegram telegramConfig `json:"telegram"`
}

type utopiaConfig struct {
	Host     string `json:"host"`
	Protocol string `json:"protocol"`
	Port     int    `json:"port"`
	Token    string `json:"token"`
}

type telegramConfig struct {
	Token  string `json:"token"`
	ChatID int64  `json:"chat_id"`
}
