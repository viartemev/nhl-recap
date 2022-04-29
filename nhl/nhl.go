package nhl

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"nhl-recap/client"
	"nhl-recap/nhl/domain"
	"strings"
	"sync"
	"time"
)

var gamesGG = make(map[int]*GameInfo)

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

type GameInfo struct {
	Video    string
	HomeTeam struct {
		Name  string
		Score int
	}
	AwayTeam struct {
		Name  string
		Score int
	}
}

func RecapFetcher(games chan *GameInfo) {
	for {
		//TODO fix schedule
		time.Sleep(30 * time.Second)
		log.Info("Fetching games")
		g := fetchGames()
		for key, element := range g {
			if _, ok := gamesGG[key]; !ok {
				gamesGG[key] = element
				log.Debug(fmt.Sprintf("Sending game: %v", element))
				games <- element //fmt.Sprintf("%v[Recap](%v)\n", element.Title, element.Video)
			}
			//TODO remove old events
		}
	}
}

func GetGames() []*GameInfo {
	gms := make([]*GameInfo, 0, len(gamesGG))
	for _, gm := range gamesGG {
		gms = append(gms, gm)
	}
	return gms
}

func fetchGames() map[int]*GameInfo {
	var wg sync.WaitGroup
	var gamesInfo = make(map[int]*GameInfo)
	schedule := client.HttpGet[domain.Schedule]("https://statsapi.web.nhl.com/api/v1/schedule")
	finishedGames := filterFinishedGames(schedule)
	for _, games := range finishedGames {
		wg.Add(1)
		go func(games domain.Games) {
			gameInfo := client.HttpGet[domain.Game]("https://statsapi.web.nhl.com/api/v1/game/" + fmt.Sprintf("%v", games.GamePk) + "/content")
			video := extractGameVideo(gameInfo)
			//title := fmt.Sprintf("*%s*\n🥅🏒 %v - %v ", games.Teams.TeamsAndWinner(), games.Teams.Home.Score, games.Teams.Away.Score)
			gamesInfo[games.GamePk] = &GameInfo{
				Video: video,
				HomeTeam: struct {
					Name  string
					Score int
				}{Name: games.Teams.Home.Team.Name, Score: games.Teams.Home.Score},
				AwayTeam: struct {
					Name  string
					Score int
				}{Name: games.Teams.Away.Team.Name, Score: games.Teams.Away.Score},
			}
			defer wg.Done()
		}(games)
	}
	wg.Wait()
	return gamesInfo
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
