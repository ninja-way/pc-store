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
	"github.com/ninja-way/pc-store/pkg/hash"
	log "github.com/sirupsen/logrus"
)

//	@title		Computer store API
//	@version	v1.6.0
//
//	@host		localhost:8080
//	@BasePath	/

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

	c := cache.New()
	h := hash.NewSHA1Hasher(cfg.DB.HashSalt)

	// init services and http handler
	compStore := service.NewComputersStore(c, cfg, db)
	usersService := service.NewUsers(db, h)
	handler := transport.NewHandler(compStore, usersService)

	// setup and run pcstore
	transport.NewServer(
		fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		handler.InitRouter(cfg),
	).Run()
}
