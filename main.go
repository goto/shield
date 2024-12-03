package main

import (
	"fmt"
	"os"

	shieldlogger "github.com/raystack/shield/logger"

	"github.com/raystack/shield/cmd"
	"github.com/raystack/shield/config"

	_ "github.com/authzed/authzed-go/proto/authzed/api/v0"
)

func main() {
	appConfig := config.Load()
	logger := shieldlogger.InitLogger(appConfig)

	if err := cmd.New(logger, appConfig).Execute(); err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
}
