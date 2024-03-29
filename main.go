package main

import (
	"errors"
	"log"
	"time"

	utopiago "github.com/Sagleft/utopialib-go/v2"
	"github.com/gin-gonic/gin"
	"gopkg.in/tucnak/telebot.v2"
)

func main() {

	app := newSolution()
	err := checkErrors(
		app.parseConfig,
		app.initGin,
		app.connectMessengers,
		app.setupRoutes,
		app.runGin,
	)
	if err != nil {
		log.Fatalln(err)
	}

}

func newSolution() *solution {
	return &solution{}
}

func (sol *solution) initGin() error {
	sol.Gin = gin.Default()
	return nil
}

func (sol *solution) runGin() error {
	go func() {
		time.Sleep(time.Millisecond * 400)
		err := openBrowserURL("http://127.0.0.1:" + sol.Config.BindPort)
		if err != nil {
			log.Fatalln(err)
		}
	}()
	return sol.Gin.Run(":" + sol.Config.BindPort)
}

func (sol *solution) connectMessengers() error {
	return checkErrors(
		sol.connectUtopia,
		sol.connectTelegram,
	)
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
