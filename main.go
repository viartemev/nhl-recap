package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"nhl-recap/nhl"
	"nhl-recap/telegram"
	"nhl-recap/util"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	ctx, cancel := context.WithCancel(context.Background())
	gracefulShutdown(cancel)

	nhl := nhl.NewNHL()
	subscription := util.NewSubscription(ctx, nhl.Fetcher, time.NewTicker(10*time.Minute))

	nhlRecapBot := telegram.InitializeBot()
	nhlRecapBot.HandleSubscription()
	nhlRecapBot.HandleUnsubscription()
	nhlRecapBot.SendSubscriptions(subscription)

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
