package nhl

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"nhl-recap/nhl/domain"
	"strings"
	"sync"
	"time"
)

var gamesGG = make(map[int]*GameInfo)

func extractGameVideo(game *domain.Game) (video string) {
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
	HomeTeam *TeamInfo
	AwayTeam *TeamInfo
}

type TeamInfo struct {
	Name  string
	Score int
}

func RecapFetcher(games chan *GameInfo) {
	for {
		//TODO fix schedule
		time.Sleep(10 * time.Second)
		log.Info("Fetching games")
		g := fetchGames()
		for key, element := range g {
			if _, ok := gamesGG[key]; !ok {
				gamesGG[key] = element
				log.Debug(fmt.Sprintf("Sending game: %v", element))
				games <- element
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
	schedule := &domain.Schedule{}
	_, err := resty.New().R().SetResult(schedule).Get("https://statsapi.web.nhl.com/api/v1/schedule")
	if err != nil {
		log.WithError(err).Error("Can't get schedule")
	}
	finishedGames := filterFinishedGames(schedule)
	for _, game := range finishedGames {
		wg.Add(1)
		go func(game domain.Games) {
			gamesInfo[game.GamePk] = fetchGameInfo(game)
			defer wg.Done()
		}(game)
	}
	wg.Wait()
	return gamesInfo
}

func fetchGameInfo(game domain.Games) *GameInfo {
	result := &domain.Game{}
	_, err := resty.New().R().SetResult(result).SetPathParams(map[string]string{
		"gamePk": string(rune(game.GamePk)),
	}).Get("https://statsapi.web.nhl.com/api/v1/game/{gamePk}/content")
	if err != nil {
		log.WithError(err).Error("Can't get game info")
	}
	video := extractGameVideo(result)
	return &GameInfo{
		video,
		&TeamInfo{game.Teams.Home.Team.Name, game.Teams.Home.Score},
		&TeamInfo{game.Teams.Away.Team.Name, game.Teams.Away.Score},
	}
}

func filterFinishedGames(schedule *domain.Schedule) (finishedGames []domain.Games) {
	for _, date := range schedule.Dates {
		for _, game := range date.Games {
			if game.Status.AbstractGameState == "Final" {
				finishedGames = append(finishedGames, game)
			}
		}
	}
	return
}
