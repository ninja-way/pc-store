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
	"github.com/ninja-way/pc-store/internal/transport/rabbitmq"
	"github.com/ninja-way/pc-store/pkg/hash"
	log "github.com/sirupsen/logrus"
)

//	@title						Computer store API
//	@version					v1.8.0
//
//	@host						localhost:8080
//	@BasePath					/
//
//	@securityDefinitions.apiKey	BearerAuth
//	@in							header
//	@name						Authorization

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
	h := hash.NewSHA1Hasher(cfg.Service.HashSalt)

	// init services and http handler
	auditClient, err := rabbitmq.NewClient(cfg.Audit.URI)
	if err != nil {
		log.Fatal(err)
	}
	defer auditClient.CloseConnection()

	usersService := service.NewUsers(db, db, auditClient, h, []byte(cfg.Service.TokenSecret), cfg.Auth.TokenTTL)
	compStore := service.NewComputersStore(c, cfg, db, auditClient)
	handler := transport.NewHandler(compStore, usersService)

	// setup and run pcstore
	transport.NewServer(
		fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		handler.InitRouter(cfg),
	).Run()
}
