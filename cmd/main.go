package main

import (
	"bank/internal/config"
	"bank/internal/database"
	"bank/internal/server"
	"os"
	"os/signal"
	"syscall"

	"context"
	"log"
	"sync"

	"github.com/joho/godotenv"
)

// @title Bank API
// @version 1.0
// @description This is a bank API on REST API.

func init() {
	if err := godotenv.Load("config.env"); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	cfg := config.NewConfig()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	go database.Broadcast(ctx, &wg, &cfg.DbConfig)
	go server.Handle(&cfg.ServerConfig)

	wg.Wait()
}
