package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/ninja-way/cache-ninja/pkg/cache"
	"github.com/ninja-way/pc-store/internal/config"
	"github.com/ninja-way/pc-store/internal/repository/postgres"
	"github.com/ninja-way/pc-store/internal/service"
	"github.com/ninja-way/pc-store/internal/transport"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	ConfigDir  = "configs"
	ConfigFile = "main"
	LogFile    = "pcstore.log"
)

var cfg *config.Config

func init() {
	// load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.WithField("loading .env file", err).Fatal()
	}

	// init config
	cfg, err = config.New(ConfigDir, ConfigFile)
	if err != nil {
		log.WithField("init config", err).Fatal()
	}

	if cfg.Environment == "prod" {
		log.SetFormatter(&log.JSONFormatter{})

		logf, err := os.OpenFile(LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			log.Fatal("failed open log file")
		}
		log.SetOutput(logf)

		log.SetLevel(log.WarnLevel)
	} else {
		log.SetFormatter(&log.TextFormatter{})
		log.SetOutput(os.Stdout)
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	ctx := context.Background()

	// connect to database
	db, err := postgres.Connect(ctx, &cfg.DB)
	if err != nil {
		log.WithField("connect to postgres", err).Fatal()
	}
	defer db.Close(ctx)

	// init service and http handler
	c := cache.New()
	compStore := service.NewComputersStore(c, cfg, db)
	handler := transport.NewHandler(compStore)

	// setup and run server
	transport.NewServer(
		fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		handler.InitRouter(),
	).Run()
}
