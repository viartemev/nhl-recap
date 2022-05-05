package telegram

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	tele "gopkg.in/telebot.v3"
	"nhl-recap/nhl"
	"nhl-recap/util"
	"time"
)

type NHLRecapBot struct {
	*tele.Bot
}

var users = util.NewSet[int64]()

func InitializeBot() *NHLRecapBot {
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
		log.WithError(err).Error("Can't start bot")
		return nil
	}
	return &NHLRecapBot{bot}
}

func (bot *NHLRecapBot) SendSubscriptions(messages <-chan *nhl.GameInfo) {
	go func() {
		for {
			message := <-messages
			users.Range(func(user int64) {
				_, err := bot.Send(&tele.User{ID: user}, gameInfoToTelegramMessage(message), &tele.SendOptions{ParseMode: tele.ModeMarkdown})
				if err != nil {
					log.WithError(err).Error("Can't send a message")
				}
				log.Debug("A message has sent")
			})
		}
	}()
}

func (bot *NHLRecapBot) HandleSubscription() {
	bot.Handle("/subscribe", func(c tele.Context) error {
		recipient := c.Sender().ID
		users.Add(recipient)
		log.WithFields(log.Fields{"user": recipient}).Info("User subscribed")
		return c.Send("Successfully subscribed")
	})
}

func (bot *NHLRecapBot) HandleUnsubscription() {
	bot.Handle("/unsubscribe", func(c tele.Context) error {
		recipient := c.Sender().ID
		users.Delete(recipient)
		log.WithFields(log.Fields{"user": recipient}).Info("User unsubscribed")
		return c.Send("Successfully unsubscribed")
	})
}

func (bot *NHLRecapBot) HandleSchedule() {
	bot.Handle("/schedule", func(c tele.Context) error {
		log.Debug("Schedule was requested")
		return c.Send(nhl.GetSchedule(), &tele.SendOptions{ParseMode: tele.ModeMarkdown})
	})
}

func (bot *NHLRecapBot) HandleGames() {
	bot.Handle("/games", func(c tele.Context) error {
		log.Debug("Games were requested")
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
