package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/ninja-way/pc-store/internal/config"
	"github.com/ninja-way/pc-store/internal/repository/postgres"
	"github.com/ninja-way/pc-store/internal/service"
	"github.com/ninja-way/pc-store/internal/transport"
	"log"
)

const (
	CONFIG_DIR  = "configs"
	CONFIG_FILE = "main"
)

func main() {
	// load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// init config
	cfg, err := config.New(CONFIG_DIR, CONFIG_FILE)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// connect to database
	db, err := postgres.Connect(ctx, &cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)

	// init service and http handler
	compStore := service.NewComputersStore(db)
	handler := transport.NewHandler(compStore)

	// setup and run server
	transport.NewServer(
		fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		handler.InitRouter(),
	).Run()
}
