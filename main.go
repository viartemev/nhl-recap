package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"nhl-recap/nhl"
	"nhl-recap/telegram"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	ctx, cancel := context.WithCancel(context.Background())
	gracefulShutdown(cancel)

	nhlRecapBot := telegram.InitializeBot()
	nhl := nhl.NewNHL()
	games := nhl.Subscribe(ctx)

	nhlRecapBot.HandleSubscription()
	nhlRecapBot.HandleUnsubscription()

	nhlRecapBot.ShowImage()

	nhlRecapBot.SendSubscriptions(games)

	log.Info("NHL Recap telegram nhlRecapBot is starting...")
	nhlRecapBot.Start()
}

func gracefulShutdown(cancel context.CancelFunc) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	signal.Notify(s, syscall.SIGTERM)
	go func() {
		<-s
		log.Info("Shutting down gracefully...")
		cancel()
		os.Exit(0)
	}()
}
