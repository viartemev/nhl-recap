package nhl

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"nhl-recap/nhl/domain"
	"nhl-recap/util"
)

var scheduleError = errors.New("can't get schedule")

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

type NHL struct {
	Fetcher util.Fetcher[*GameInfo]
}

func NewNHL() *NHL {
	return &NHL{Fetcher: &NHLFetcher{client: NewNHLClient(), uniqueGames: util.NewSet[int]()}}
}

type NHLFetcher struct {
	client      NHLClient
	uniqueGames *util.Set[int]
}

func (f *NHLFetcher) Fetch(ctx context.Context) (chan *GameInfo, error) {
	log.Info("Requesting nhl info")
	schedule, err := f.client.FetchSchedule()
	if err != nil {
		return nil, scheduleError
	}
	finishedGames := schedule.ExtractFinishedGames()
	log.Infof("Got %d finished games", len(finishedGames))
	//TODO fix errors in channel
	fetchGame := func(games domain.Games) *GameInfo { return f.fetchGameInfo(games) }
	games := util.FanIn(ctx, finishedGames, fetchGame)
	uniqueGame := func(info *GameInfo) bool { return f.uniqueGames.Add(info.GamePk) }
	notNil := func(info *GameInfo) bool { return info != nil }
	return util.Filter(ctx, games, util.And(notNil, uniqueGame)), nil
}

func (f *NHLFetcher) fetchGameInfo(game domain.Games) *GameInfo {
	fetchedGame, err := f.client.FetchGame(game.GamePk)
	if err != nil {
		log.WithError(err).Error("Can't get game info")
		return nil
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
