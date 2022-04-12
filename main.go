package main

import (
	"bytes"
	"fmt"
	"log"
	"nhl-recap/client"
	"nhl-recap/domain"
	"os"
	"strings"
	"sync"
	"time"

	tele "gopkg.in/telebot.v3"
)

func main() {
	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	b.Handle("/games", func(c tele.Context) error {
		return c.Send(fetchGames(), &tele.SendOptions{ParseMode: tele.ModeMarkdown})
	})

	b.Start()
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
