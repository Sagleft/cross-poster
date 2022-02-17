package main

import (
	utopiago "github.com/Sagleft/utopialib-go"
	"github.com/gin-gonic/gin"
	"gopkg.in/tucnak/telebot.v2"
)

type solution struct {
	Config     appConfig
	Gin        *gin.Engine
	Messengers messengerClients
}

type messengerClients struct {
	Utopia   *utopiago.UtopiaClient
	Telegram *telebot.Bot
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
