package main

import (
	"log"

	"user-auth-app/internal/config"
	"user-auth-app/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}
	cfg := config.Load()
	srv := server.NewServer(cfg)
	srv.Start()
}
