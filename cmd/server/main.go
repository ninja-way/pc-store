package main

import (
	"github.com/ninja-way/pc-store/internal/model"
	"github.com/ninja-way/pc-store/internal/repository/cache"
	"github.com/ninja-way/pc-store/internal/server"
	"log"
)

func main() {
	db := cache.Init()
	if err := db.AddComputer(model.PC{Name: "Test PC", Price: 19999}); err != nil {
		log.Fatal(err)
	}

	s := server.New(":8080", nil, db)

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
