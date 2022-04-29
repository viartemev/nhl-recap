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

	bot := telegram.InitializeBot()

	games := make(chan string)

	telegram.HandleSubscription(bot)
	telegram.HandleGames(bot)

	go nhl.RecapFetcher(games)
	go telegram.SendSubscriptions(bot, games)
	log.Info("NHL Recap telegram bot is starting...")
	bot.Start()
}
