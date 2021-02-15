package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog"
)

func main() {
	context := context.Background()
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()

	if lvl, ok := os.LookupEnv("LOG_LEVEL"); ok {
		logLevel, err := strconv.Atoi(lvl)
		if err != nil {
			logger.Error().Err(err).Msgf("Error occurred while parsing LOG_LEVEL. Value: %v", lvl)
		} else {
			logger = logger.Level(zerolog.Level(logLevel))
		}
	}

	config, err := GetNewConfig()
	if err != nil {
		logger.Fatal().Msg("Bad config")
	}

	clientSet, err := GetNewK8sClient()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
