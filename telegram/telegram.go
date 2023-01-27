package telegram

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	tele "gopkg.in/telebot.v3"
	d "nhl-recap/domain"
	"nhl-recap/util"
	"strconv"
	"time"
)

type NHLRecapBot struct {
	*tele.Bot
	Users *util.Set[int64]
}

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
	return &NHLRecapBot{bot, util.NewSet[int64]()}
}

func (bot *NHLRecapBot) SendSubscriptions(subscription util.Subscription[*d.GameInfo]) {
	go func() {
		for message := range subscription.Updates() {
			bot.Users.Range(sendScoreCard(message, bot))
		}
	}()
}

func sendScoreCard(message *d.GameInfo, bot *NHLRecapBot) func(user int64) {
	return func(user int64) {
		photo := &tele.Photo{File: tele.FromReader(bytes.NewReader(message.ScoreCard))}
		senderOptions := &tele.SendOptions{ReplyMarkup: &tele.ReplyMarkup{InlineKeyboard: [][]tele.InlineButton{
			{tele.InlineButton{Unique: strconv.Itoa(message.GamePk), Text: "Watch", URL: message.Video}},
		}}}
		_, err := bot.Send(&tele.User{ID: user}, photo, senderOptions)
		if err != nil {
			log.WithError(err).Error("Can't send the message")
		} else {
			log.Debug("The message has sent")
		}
	}
}

func (bot *NHLRecapBot) HandleSubscription() {
	bot.Handle("/subscribe", func(c tele.Context) error {
		recipient := c.Sender().ID
		bot.Users.Add(recipient)
		log.WithFields(log.Fields{"user": recipient}).Info("User subscribed")
		return c.Send("Successfully subscribed")
	})
}

func (bot *NHLRecapBot) HandleUnsubscription() {
	bot.Handle("/unsubscribe", func(c tele.Context) error {
		recipient := c.Sender().ID
		bot.Users.Delete(recipient)
		log.WithFields(log.Fields{"user": recipient}).Info("User unsubscribed")
		return c.Send("Successfully unsubscribed")
	})
}
