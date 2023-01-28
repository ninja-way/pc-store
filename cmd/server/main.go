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
	db, err := postgres.Init(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//db := cache.Init()
	//if err := db.AddComputer(model.PC{Name: "Test PC", Price: 19999}); err != nil {
	//	log.Fatal(err)
	//}

	s := server.New(":8080", middleware.Logging(http.DefaultServeMux), db)

	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
