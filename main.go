package main

import (
	"fmt"
	"nhl-recap/nhl"
	"nhl-recap/telegram"
)

func main() {
	bot := telegram.InitializeBot()

	games := make(chan string)

	telegram.HandleSubscription(bot)
	telegram.HandleGames(bot)

	go nhl.RecapFetcher(games)
	go telegram.SendSubscriptions(bot, games)

	fmt.Println("Telegram bot NHL Recap starting...")
	bot.Start()
}
