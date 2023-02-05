package main

import (
	"context"
	"github.com/ninja-way/pc-store/internal/repository/postgres"
	"github.com/ninja-way/pc-store/internal/service"
	"github.com/ninja-way/pc-store/internal/transport"
	"log"
)

func main() {
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
