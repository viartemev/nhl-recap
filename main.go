package main

import (
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
	start := time.Now()
	var wg sync.WaitGroup
	gamesVideos := make(map[int]*GameInfo)
	schedule := httpGet[domain.Schedule]("https://statsapi.web.nhl.com/api/v1/schedule")
	finishedGames := filterFinishedGames(schedule)

	for _, game := range finishedGames {
		wg.Add(1)
		go func(game domain.Games) {
			gameInfo := httpGet[domain.Game]("https://statsapi.web.nhl.com/api/v1/game/" + fmt.Sprintf("%v", game.GamePk) + "/content")
			video := extractGameVideo(gameInfo)
			title := fmt.Sprintf("%v vs %v: %v - %v", game.Teams.Home.Team.Name, game.Teams.Away.Team.Name, game.Teams.Home.Score, game.Teams.Away.Score)
			gamesVideos[game.GamePk] = &GameInfo{title, video}
			defer wg.Done()
		}(game)
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
