package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aaegamysta/listen-2-max-payne/internal/bot"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panicf("error loading .env file: %v", err)
	}
	env := os.Getenv("DEVELOPMENT_ENV")
	ctx := context.Background()
	bot := bot.New(ctx, env)
	bot.Run(ctx)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
}
