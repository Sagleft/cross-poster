package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	swissknife "github.com/Sagleft/swiss-knife"
	utopiago "github.com/Sagleft/utopialib-go/v2"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"gopkg.in/tucnak/telebot.v2"
)

func main() {

	app := newSolution()
	err := checkErrors(
		app.parseConfig,
		app.initGin,
		app.setupRoutes,
	)
	if err != nil {
		log.Fatalln(err)
	}

	go app.runGin()

	if err := app.connectMessengers(); err != nil {
		app.LastError = err
		color.Red(err.Error())
	}

	swissknife.WaitForAppFinish()
}

func newSolution() *solution {
	return &solution{}
}

func (sol *solution) initGin() error {
	sol.Gin = gin.Default()
	return nil
}

func (sol *solution) runGin() {
	go sol.Gin.Run(":" + sol.Config.BindPort)

	time.Sleep(time.Millisecond * 400)

	url := "http://127.0.0.1:" + sol.Config.BindPort
	if err := openBrowserURL(url); err != nil {
		log.Fatalln(err)
	}
}

func (sol *solution) connectMessengers() error {
	if err := sol.connectUtopia(); err != nil {
		return fmt.Errorf("connect Utopia: %w", err)
	}

	if err := sol.connectTelegram(); err != nil {
		return fmt.Errorf("connect Telegram: %w", err)
	}
	return nil
}

func (sol *solution) connectUtopia() error {
	if sol.Config.Utopia.Token == "" {
		return errors.New("utopia token is not set in " + configJSONPath)
	}

	sol.Messengers.Utopia = utopiago.NewUtopiaClient(utopiago.Config{
		Host:     sol.Config.Utopia.Host,
		Token:    sol.Config.Utopia.Token,
		Port:     sol.Config.Utopia.Port,
		WsPort:   defaultWsPort,
		Protocol: sol.Config.Utopia.Protocol,
	})

	if !sol.Messengers.Utopia.CheckClientConnection() {
		return errors.New("failed to connect to Utopia messenger")
	}
	return nil
}

func (sol *solution) connectTelegram() error {
	if sol.Config.Telegram.Token == "" {
		return errors.New("telegram token is not set in " + configJSONPath)
	}

	var err error
	sol.Messengers.Telegram, err = telebot.NewBot(telebot.Settings{
		Token:  sol.Config.Telegram.Token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return errors.New("failed to connect to telegram bot: " + err.Error())
	}
	return nil
}
