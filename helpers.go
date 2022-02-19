package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
)

func (sol *solution) parseConfig() error {
	// parse config file
	if _, err := os.Stat(configJSONPath); os.IsNotExist(err) {
		return errors.New("failed to find config file")
	}
	jsonBytes, err := ioutil.ReadFile(configJSONPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, &sol.Config)
}

type errorFunc func() error

func checkErrors(errChecks ...errorFunc) error {
	for _, errFunc := range errChecks {
		err := errFunc()
		if err != nil {
			return err
		}
	}
	return nil
}

func openBrowserURL(urlAddress string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", urlAddress).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", urlAddress).Start()
	case "darwin":
		return exec.Command("open", urlAddress).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}
