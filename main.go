package main

import (
	"errors"
	"log"

	utopiago "github.com/Sagleft/utopialib-go"
	"github.com/gin-gonic/gin"
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
	return sol.Gin.Run()
}

func (sol *solution) connectMessengers() error {
	return checkErrors(
		sol.connectUtopia,
		sol.connectTelegram,
	)
}

func (sol *solution) connectUtopia() error {
	sol.Messengers.Utopia = &utopiago.UtopiaClient{
		Protocol: sol.Config.Utopia.Protocol,
		Token:    sol.Config.Utopia.Token,
		Host:     sol.Config.Utopia.Host,
		Port:     sol.Config.Utopia.Port,
	}
	if !sol.Messengers.Utopia.CheckClientConnection() {
		return errors.New("failed to connect to Utopia messenger")
	}
	return nil
}

func (sol *solution) connectTelegram() error {
	// TODO
	return nil
}
