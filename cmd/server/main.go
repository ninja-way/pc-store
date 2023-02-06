package main

import (
	"context"
	"fmt"
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
	cfg, err := config.New(CONFIG_DIR, CONFIG_FILE)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)

	ctx := context.Background()

	// connect to database
	db, err := postgres.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)

	// init service and handler
	compStore := service.NewComputersStore(db)
	handler := transport.NewHandler(compStore)

	// setup and run server
	transport.NewServer(
		":8080",
		handler.InitRouter(),
	).Run()
}
