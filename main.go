package main

import (
	log "github.com/sirupsen/logrus"
	"nhl-recap/nhl"
	"nhl-recap/telegram"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	nhlRecapBot := telegram.InitializeBot()
	games := nhl.RecapFetcher()

	nhlRecapBot.HandleSubscription()
	nhlRecapBot.HandleUnsubscription()

	nhlRecapBot.ShowImage()

	nhlRecapBot.SendSubscriptions(games)

	log.Info("NHL Recap telegram nhlRecapBot is starting...")
	nhlRecapBot.Start()
}
