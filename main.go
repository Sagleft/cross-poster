package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	swissknife "github.com/Sagleft/swiss-knife"
	"github.com/Sagleft/uchatbot-engine"
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
		app.setLastError(fmt.Errorf("connect messengers: %w", err))
	}

	swissknife.WaitForAppFinish()
}

func newSolution() *solution {
	return &solution{}
}

func (sol *solution) setLastError(err error) {
	sol.LastError = err
	color.Red(err.Error())
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

	// create chatbot to handle auto-reconnect
	chatBot, err := uchatbot.NewChatBot(uchatbot.ChatBotData{
		Config: utopiago.Config{
			Host:     sol.Config.Utopia.Host,
			Token:    sol.Config.Utopia.Token,
			Port:     sol.Config.Utopia.Port,
			WsPort:   defaultWsPort,
			Protocol: sol.Config.Utopia.Protocol,
		},
		Chats: []uchatbot.Chat{
			{ID: sol.Config.Utopia.ChannelID},
		},
		DisableEvents:    true,
		UseErrorCallback: true,
		ErrorCallback:    sol.setLastError,
	})
	if err != nil {
		return fmt.Errorf("setup Utopia: %w", err)
	}

	sol.setLastError(nil)
	sol.Messengers.Utopia = chatBot.GetClient()

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
