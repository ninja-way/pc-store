package main

import (
	"context"
	"github.com/ninja-way/pc-store/internal/middleware"
	"github.com/ninja-way/pc-store/internal/repository/postgres"
	"github.com/ninja-way/pc-store/internal/server"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	// connect to database
	db, err := postgres.Init(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// setup server with logging
	s := server.New(":8080", middleware.Logging(http.DefaultServeMux), db)

	// start server
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
