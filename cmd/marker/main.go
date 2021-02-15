package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/ilnurkhal/marker"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/rs/zerolog"
)

func main() {
	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()

	if lvl, ok := os.LookupEnv("LOG_LEVEL"); ok {
		logLevel, err := strconv.Atoi(lvl)
		if err != nil {
			logger.Error().
				Err(err).
				Str("log_level", lvl).
				Msg("Error occurred while parsing LOG_LEVEL")

		} else {
			logger = logger.Level(zerolog.Level(logLevel))
		}
	}

	config, err := marker.GetNewConfig()
	if err != nil {
		logger.Fatal().
			Msg("Bad config")
	}

	clientSet, err := marker.GetNewK8sClient()
	if err != nil {
		fmt.Println("Error:", err)
	}

	transport := httptransport.New(
		os.Getenv("NETBOX_HOST"),
		client.DefaultBasePath,
		[]string{"https"})
	transport.DefaultAuthentication = httptransport.APIKeyAuth(
		"Authorization",
		"header",
		fmt.Sprintf("Token %s", os.Getenv("NETBOX_TOKEN")))
	netBoxClient := client.New(transport, nil)

	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)
	signal.Notify(signalChan,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGINT,
	)

	m := marker.GetNewMarker(
		clientSet,
		netBoxClient,
		&config,
		&logger,
	)

	go m.Run(context)

	for {
		select {
		case receivedSignal := <-signalChan:
			logger.Warn().
				Str("signal", receivedSignal.String()).
				Msg("Got signal")
			cancel()

		}
	}
}
