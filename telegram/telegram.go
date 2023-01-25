package telegram

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	tele "gopkg.in/telebot.v3"
	"nhl-recap/nhl"
	"nhl-recap/util"
	"strconv"
	"time"
)

type TelegramUsers struct {
	Users *util.Set[int64]
}

type NHLRecapBot struct {
	*tele.Bot
	*TelegramUsers
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
	users := &TelegramUsers{util.NewSet[int64]()}
	return &NHLRecapBot{bot, users}
}

func (bot *NHLRecapBot) SendSubscriptions(subscription util.Subscription[*nhl.GameInfo]) {
	go func() {
		for message := range subscription.Updates() {
			scoreCard := tele.FromReader(bytes.NewReader(nhl.GenerateScoreCard(message)))
			bot.Users.Range(sendScoreCard(scoreCard, message, bot))
		}
	}()
}

func sendScoreCard(scoreCard tele.File, message *nhl.GameInfo, bot *NHLRecapBot) func(user int64) {
	return func(user int64) {
		photo := &tele.Photo{File: scoreCard}
		senderOptions := &tele.SendOptions{ReplyMarkup: &tele.ReplyMarkup{InlineKeyboard: [][]tele.InlineButton{
			{
				tele.InlineButton{
					Unique: strconv.Itoa(message.GamePk),
					Text:   fmt.Sprintf("%s v.s. %s", message.HomeTeam.Name, message.AwayTeam.Name),
					URL:    message.Video,
				},
			},
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

func (bot *NHLRecapBot) ShowImage() {
	bot.Handle("/show", func(context tele.Context) error {
		photo := &tele.Photo{
			File: tele.FromReader(bytes.NewReader(nhl.GenerateScoreCard(nil))),
		}
		return context.Send(photo, &tele.SendOptions{ReplyMarkup: &tele.ReplyMarkup{InlineKeyboard: [][]tele.InlineButton{
			{
				tele.InlineButton{
					Unique: "foo_btn",
					Text:   "Watch",
					URL:    "https://wsczoominwestus.prod-cdn.clipro.tv/publish/7884995/12395948/58e70ac9-77f9-4111-9c90-53216088cce0.mp4",
				},
			},
		}}})
	})
}
