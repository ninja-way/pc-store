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
)

const (
	ConfigDir  = "configs"
	ConfigFile = "main"
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

	config.SetupLogger(cfg)
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

	// setup and run pcstore
	transport.NewServer(
		fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		handler.InitRouter(),
	).Run()
}
