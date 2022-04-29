package telegram

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	tele "gopkg.in/telebot.v3"
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

func SendSubscriptions(bot *tele.Bot, messages chan *nhl.GameInfo) {
	for {
		message := <-messages
		for _, user := range users {
			_, err := bot.Send(&tele.User{ID: user}, gameInfoToTelegramMessage(message), &tele.SendOptions{ParseMode: tele.ModeMarkdown})
			if err != nil {
				log.Error("Can't send a message", err)
			}
			log.Debug("A message has sent")
		}
	}
}

func HandleSubscription(bot *tele.Bot) {
	bot.Handle("/subscribe", func(c tele.Context) error {
		recipient := c.Sender().ID
		users = append(users, recipient)
		log.WithFields(log.Fields{"user": recipient}).Info("User subscribed")
		return c.Send("Successfully subscribed")
	})
}

func HandleGames(bot *tele.Bot) {
	bot.Handle("/games", func(c tele.Context) error {
		var buffer bytes.Buffer
		games := nhl.GetGames()
		for _, game := range games {
			buffer.WriteString(gameInfoToTelegramMessage(game))
		}
		return c.Send(buffer.String(), &tele.SendOptions{ParseMode: tele.ModeMarkdown})
	})
}

func gameInfoToTelegramMessage(game *nhl.GameInfo) string {
	return fmt.Sprintf("(%v vs %v) %v - %v video: %v \n", game.HomeTeam.Name, game.AwayTeam.Name, game.HomeTeam.Score, game.AwayTeam.Score, game.Video)
}

type Item struct {
	Recipient tele.Recipient
	Message   string
}
