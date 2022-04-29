package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	tele "gopkg.in/telebot.v3"
	"log"
	"nhl-recap/nhl"
	"time"
)

func InitializeBot(token string) *tele.Bot {
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

var (
	newBot *tele.Bot
)

func init() {
	var token string
	flag.StringVarP(&token, "token", "t", "", "Token for Telegram Bot API")
	flag.Parse()
	newBot = InitializeBot(token)
}

type Item struct {
	Message   string
	Recipient tele.Recipient
}

func sendMessages(messages chan Item) {
	for {
		message := <-messages
		_, err := newBot.Send(message.Recipient, message.Message)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Message sent")
	}
}

func main() {
	var users = make([]int64, 1)
	messages := make(chan Item)

	newBot.Handle("/games", func(c tele.Context) error {
		return c.Send(nhl.FetchGames(), &tele.SendOptions{ParseMode: tele.ModeMarkdown})
	})

	newBot.Handle("/subscribe", func(c tele.Context) error {
		recipient := c.Sender().ID
		users = append(users, recipient)
		return c.Send("Successfully subscribed")
	})

	go sendMessages(messages)

	go func(messages chan Item) {
		for {
			messages <- Item{Recipient: &tele.User{ID: 111067917}, Message: "Hello!"}
			time.Sleep(5 * time.Second)
		}
	}(messages)

	fmt.Println("Telegram bot NHL Recap starting...")
	newBot.Start()
}
