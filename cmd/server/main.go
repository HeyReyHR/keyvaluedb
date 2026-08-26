package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/HeyReyHR/keyvaluedb/internal/config"
	"github.com/HeyReyHR/keyvaluedb/internal/initialization"
)

var ConfigFileName = os.Getenv("CONFIG_FILE_NAME")

// CONFIG_FILE_NAME=./config.yaml go run ./cmd/server
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := &config.Config{}
	if ConfigFileName != "" {
		var err error
		cfg, err = config.LoadConfig(ConfigFileName)
		if err != nil {
			log.Fatal(err)
		}
	}

	initializer, err := initialization.NewInitializer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err = initializer.StartDatabase(ctx); err != nil {
		log.Fatal(err)
	}
}
