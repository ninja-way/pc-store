package main

import (
	"context"
	"github.com/ninja-way/pc-store/internal/repository/postgres"
	"github.com/ninja-way/pc-store/internal/service"
	"github.com/ninja-way/pc-store/internal/transport"
	"github.com/ninja-way/pc-store/internal/transport/middleware"
	"log"
	"net/http"
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

	// setup server with logging
	srv := transport.NewServer(
		":8080",
		middleware.Logging(http.DefaultServeMux),
		handler,
	)

	// start server
	if err = srv.Run(); err != nil {
		log.Fatal(err)
	}
}
