package nhl

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"nhl-recap/nhl/domain"
	"strconv"
	"sync"
	"time"
)

var gamesGG = make(map[int]*GameInfo)

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

func GetSchedule() string {
	return "schedule"
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
	var mutex sync.Mutex
	var gamesInfo = make(map[int]*GameInfo)
	schedule := &domain.Schedule{}
	_, err := resty.New().R().SetResult(schedule).Get("https://statsapi.web.nhl.com/api/v1/schedule")
	if err != nil {
		log.WithError(err).Error("Can't get schedule")
	}
	finishedGames := schedule.ExtractFinishedGames()
	for _, game := range finishedGames {
		wg.Add(1)
		go func(game domain.Games) {
			mutex.Lock()
			gamesInfo[game.GamePk] = fetchGameInfo(game)
			defer wg.Done()
			defer mutex.Unlock()
		}(game)
	}
	wg.Wait()
	return gamesInfo
}

func fetchGameInfo(game domain.Games) *GameInfo {
	result := &domain.Game{}
	_, err := resty.New().R().SetResult(result).SetPathParams(map[string]string{
		"gamePk": strconv.Itoa(game.GamePk),
	}).Get("https://statsapi.web.nhl.com/api/v1/game/{gamePk}/content")
	if err != nil {
		log.WithError(err).Error("Can't get game info")
	}
	video := result.ExtractGameVideo()
	return &GameInfo{
		video,
		&TeamInfo{game.Teams.Home.Team.Name, game.Teams.Home.Score},
		&TeamInfo{game.Teams.Away.Team.Name, game.Teams.Away.Score},
	}
}
