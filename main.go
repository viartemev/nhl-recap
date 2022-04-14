package main

import (
	"bytes"
	"fmt"
	flag "github.com/spf13/pflag"
	tele "gopkg.in/telebot.v3"
	"log"
	"nhl-recap/client"
	"nhl-recap/domain"
	"strings"
	"sync"
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
		return c.Send(fetchGames(), &tele.SendOptions{ParseMode: tele.ModeMarkdown})
	})
	bot.Settings.Start()
}

type GameInfo struct {
	Title string
	Video string
}

func fetchGames() string {
	var wg sync.WaitGroup
	gamesVideos := make(map[int]*GameInfo)
	schedule := client.HttpGet[domain.Schedule]("https://statsapi.web.nhl.com/api/v1/schedule")
	finishedGames := filterFinishedGames(schedule)
	if len(finishedGames) == 0 {
		return "There are no finished games"
	}
	for _, games := range finishedGames {
		wg.Add(1)
		go func(games domain.Games) {
			gameInfo := client.HttpGet[domain.Game]("https://statsapi.web.nhl.com/api/v1/game/" + fmt.Sprintf("%v", games.GamePk) + "/content")
			video := extractGameVideo(gameInfo)
			title := fmt.Sprintf("*%v vs %v*\n🥅🏒 %v - %v ", games.Teams.Home.Team.Name, games.Teams.Away.Team.Name, games.Teams.Home.Score, games.Teams.Away.Score)
			gamesVideos[games.GamePk] = &GameInfo{title, video}
			defer wg.Done()
		}(games)
	}
	wg.Wait()

	var buffer bytes.Buffer
	for _, info := range gamesVideos {
		//fmt.Printf("[%v](%v)\n", info.Title, info.Video)
		buffer.WriteString(fmt.Sprintf("%v[Recap](%v)\n", info.Title, info.Video))
	}
	return buffer.String()
}

func filterFinishedGames(schedule domain.Schedule) (finishedGames []domain.Games) {
	for _, date := range schedule.Dates {
		for _, game := range date.Games {
			if game.Status.AbstractGameState == "Final" {
				finishedGames = append(finishedGames, game)
			}
		}
	}
	return
}

func extractGameVideo(game domain.Game) (video string) {
	for _, media := range game.Media.Epg {
		if media.Title == "Recap" {
			for _, item := range media.Items {
				for _, playback := range item.Playbacks {
					if strings.Contains(playback.Name, "FLASH_1800K") {
						video = playback.Url
					}
				}
			}
		}
	}
	return
}
