package nhl

import (
	"bytes"
	"fmt"
	"nhl-recap/client"
	domain2 "nhl-recap/nhl/domain"
	"strings"
	"sync"
	"time"
)

func extractGameVideo(game domain2.Game) (video string) {
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

type GameInfo struct {
	Title string
	Video string
}

func RecapFetcher(games chan string) {
	for {
		games <- "game info"
		time.Sleep(5 * time.Second)
	}
}

func FetchGames() string {
	var wg sync.WaitGroup
	gamesVideos := make(map[int]*GameInfo)
	schedule := client.HttpGet[domain2.Schedule]("https://statsapi.web.nhl.com/api/v1/schedule")
	finishedGames := filterFinishedGames(schedule)
	if len(finishedGames) == 0 {
		return "There are no finished games"
	}
	for _, games := range finishedGames {
		wg.Add(1)
		go func(games domain2.Games) {
			gameInfo := client.HttpGet[domain2.Game]("https://statsapi.web.nhl.com/api/v1/game/" + fmt.Sprintf("%v", games.GamePk) + "/content")
			video := extractGameVideo(gameInfo)
			title := fmt.Sprintf("*%s*\n🥅🏒 %v - %v ", games.Teams.TeamsAndWinner(), games.Teams.Home.Score, games.Teams.Away.Score)
			gamesVideos[games.GamePk] = &GameInfo{title, video}
			defer wg.Done()
		}(games)
	}
	wg.Wait()

	var buffer bytes.Buffer
	for _, info := range gamesVideos {
		buffer.WriteString(fmt.Sprintf("%v[Recap](%v)\n", info.Title, info.Video))
	}
	return buffer.String()
}

func filterFinishedGames(schedule domain2.Schedule) (finishedGames []domain2.Games) {
	for _, date := range schedule.Dates {
		for _, game := range date.Games {
			if game.Status.AbstractGameState == "Final" {
				finishedGames = append(finishedGames, game)
			}
		}
	}
	return
}
