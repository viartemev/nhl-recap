package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	tele "gopkg.in/telebot.v3"
	"log"
	"nhl-recap/nhl"
	"time"
)

type BotSettings struct {
	Token    string
	Settings *tele.Bot
}

func (b *BotSettings) Initialize() {
	var err error
	pref := tele.Settings{
		Token:  b.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b.Settings, err = tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
}

var (
	bot BotSettings
)

func init() {
	flag.StringVarP(&bot.Token, "token", "t", "", "Token for Telegram Bot API")
	flag.Parse()
	bot.Initialize()

}

func main() {
	bot.Settings.Handle("/games", func(c tele.Context) error {
		return c.Send(nhl.FetchGames(), &tele.SendOptions{ParseMode: tele.ModeMarkdown})
	})

	fmt.Println("Telegram bot NHL Recap starting...")
	bot.Settings.Start()
}
