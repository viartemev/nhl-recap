package telegram

import (
	"fmt"
	flag "github.com/spf13/pflag"
	tele "gopkg.in/telebot.v3"
	"log"
	"nhl-recap/nhl"
	"time"
)

var users = make([]int64, 0)

func InitializeBot() *tele.Bot {
	var token string
	flag.StringVarP(&token, "token", "t", "", "Token for Telegram Bot API")
	flag.Parse()

	var err error
	var bot *tele.Bot
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err = tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return bot
}

func SendSubscriptions(bot *tele.Bot, messages chan string) {
	for {
		message := <-messages
		for _, user := range users {
			_, err := bot.Send(&tele.User{ID: user}, message)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Message sent")
		}
	}
}

func HandleSubscription(bot *tele.Bot) {
	bot.Handle("/subscribe", func(c tele.Context) error {
		recipient := c.Sender().ID
		users = append(users, recipient)
		return c.Send("Successfully subscribed")
	})
}

func HandleGames(bot *tele.Bot) {
	bot.Handle("/games", func(c tele.Context) error {
		return c.Send(nhl.FetchGames(), &tele.SendOptions{ParseMode: tele.ModeMarkdown})
	})
}

type Item struct {
	Message   string
	Recipient tele.Recipient
}
