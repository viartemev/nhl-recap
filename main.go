package main

import (
	"example/client"
	"example/domain"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

type GameInfo struct {
	Title string
	Video string
}

func main() {
	//scheduler := gocron.NewScheduler(time.UTC)
	//_, _ = scheduler.Every(5).Seconds().Do(fetchGames)
	//scheduler.StartBlocking()
	fetchGames()
}

func fetchGames() {
	start := time.Now()
	var wg sync.WaitGroup
	gamesVideos := make(map[int]*GameInfo)
	schedule := client.HttpGet[domain.Schedule]("https://statsapi.web.nhl.com/api/v1/schedule")
	finishedGames := filterFinishedGames(schedule)

	for _, games := range finishedGames {
		wg.Add(1)
		go func(games domain.Games) {
			gameInfo := client.HttpGet[domain.Game]("https://statsapi.web.nhl.com/api/v1/game/" + fmt.Sprintf("%v", games.GamePk) + "/content")
			video := extractGameVideo(gameInfo)
			title := fmt.Sprintf("%v vs %v: %v - %v", games.Teams.Home.Team.Name, games.Teams.Away.Team.Name, games.Teams.Home.Score, games.Teams.Away.Score)
			gamesVideos[games.GamePk] = &GameInfo{title, video}
			defer wg.Done()
		}(games)
	}
	wg.Wait()
	for _, info := range gamesVideos {
		fmt.Printf("Game %v %v\n", info.Title, info.Video)
	}
	log.Printf("Operations took %s", time.Since(start))
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
