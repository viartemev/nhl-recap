package nhl

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"nhl-recap/nhl/domain"
	"sync"
	"time"
)

type GameInfo struct {
	GamePk   int
	Video    string
	HomeTeam *TeamInfo
	AwayTeam *TeamInfo
}

type TeamInfo struct {
	Name  string
	Score int
}

func RecapFetcher() <-chan *GameInfo {
	out := make(chan *GameInfo)
	gg := make(map[int]struct{})
	go func() {
		for range time.Tick(time.Minute) {
			log.Info("Fetching games info")
			g := fetchGames()
			for element := range g {
				if _, ok := gg[element.GamePk]; !ok {
					log.Debug(fmt.Sprintf("Sending game: %v", element))
					out <- element
					gg[element.GamePk] = struct{}{}
					//TODO clean temporal table
				}
			}
		}
	}()
	return out
}

func fetchGames() chan *GameInfo {
	nhlClient := NHLHTTPClient{client: &http.Client{Timeout: 3 * time.Second}}
	var wg sync.WaitGroup
	gamesInfo := make(chan *GameInfo)
	schedule, err := nhlClient.FetchSchedule()
	if err != nil {
		log.WithError(err).Error("Can't get schedule")
	}
	finishedGames := schedule.ExtractFinishedGames()
	for _, game := range finishedGames {
		wg.Add(1)
		go func(game domain.Games) {
			gi := fetchGameInfo(game)
			if gi != nil {
				gamesInfo <- gi
			}
			defer wg.Done()
		}(game)
	}
	go func() {
		wg.Wait()
		close(gamesInfo)
	}()
	return gamesInfo
}

func fetchGameInfo(game domain.Games) *GameInfo {
	nhlClient := NHLHTTPClient{client: &http.Client{Timeout: 3 * time.Second}}
	fetchedGame, err := nhlClient.FetchGame(game.GamePk)
	if err != nil {
		log.WithError(err).Error("Can't get game info")
	}
	video := fetchedGame.ExtractGameVideo()
	if len(video) == 0 {
		return nil
	}
	return &GameInfo{
		game.GamePk,
		video,
		&TeamInfo{game.Teams.Home.Team.Name, game.Teams.Home.Score},
		&TeamInfo{game.Teams.Away.Team.Name, game.Teams.Away.Score},
	}
}
