package main

import (
	"github.com/Sagleft/uexchange-go"
	"github.com/gin-gonic/gin"
)

type solution struct {
	Config         appConfig
	Gin            *gin.Engine
	ExchangeClient *uexchange.Client
}

type appConfig struct {
	Exchange exchangeConfig `json:"exchange"`
}

type exchangeConfig struct {
	PublicKey string `json:"pubkey"`
	Password  string `json:"password"`
}
