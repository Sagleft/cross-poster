package main

import (
	utopiago "github.com/Sagleft/utopialib-go/v2"
	"github.com/gin-gonic/gin"
	"gopkg.in/tucnak/telebot.v2"
)

type solution struct {
	Config     appConfig
	Gin        *gin.Engine
	Messengers messengerClients

	LastError error
}

type messengerClients struct {
	Utopia   utopiago.Client
	Telegram *telebot.Bot
}

type appConfig struct {
	Utopia   utopiaConfig   `json:"utopia"`
	Telegram telegramConfig `json:"telegram"`
	BindPort string         `json:"bind_port"`
	FrontEnd frontendConfig `json:"frontend"`
}

type frontendConfig struct {
	Version string `json:"version"`
}

type utopiaConfig struct {
	Host      string `json:"host"`
	Protocol  string `json:"protocol"`
	Port      int    `json:"port"`
	Token     string `json:"token"`
	ChannelID string `json:"channel_id"`
}

type telegramConfig struct {
	Token      string `json:"token"`
	ChatID     int64  `json:"chat_id"`
	SilentMode bool   `json:"silent_mode"`
}

type response struct {
	Status    string `json:"status"`
	ErrorInfo string `json:"error"`
}
