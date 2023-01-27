package nhl

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	d "nhl-recap/domain"
	"nhl-recap/nhl/domain"
	"nhl-recap/nhl/logos"
	"nhl-recap/util"
)

var scheduleError = errors.New("can't get schedule")

type NHL struct {
	Fetcher util.Fetcher[*d.GameInfo]
}

func NewNHL() *NHL {
	return &NHL{Fetcher: &NHLFetcher{client: NewNHLClient(), uniqueGames: util.NewSet[int](), scoreCardGenerator: NewScoreCardGenerator(logos.LoadLogos())}}
}

type NHLFetcher struct {
	client             NHLClient
	uniqueGames        *util.Set[int]
	scoreCardGenerator ScoreCardGenerator
}

func (f *NHLFetcher) Fetch(ctx context.Context) (chan *d.GameInfo, error) {
	log.Info("Requesting nhl info")
	schedule, err := f.client.FetchSchedule()
	if err != nil {
		return nil, scheduleError
	}
	finishedGames := schedule.ExtractFinishedGames()
	log.Infof("Got %d finished games", len(finishedGames))
	//TODO fix errors in channel
	fetchGame := func(games domain.ScheduleGame) *d.GameInfo { return f.fetchGameInfo(games) }
	games := util.FanIn(ctx, finishedGames, fetchGame)
	uniqueGame := func(info *d.GameInfo) bool { return f.uniqueGames.Add(info.GamePk) }
	notNil := func(info *d.GameInfo) bool { return info != nil }
	return util.Filter(ctx, games, util.And(notNil, uniqueGame)), nil
}

func (f *NHLFetcher) fetchGameInfo(game domain.ScheduleGame) *d.GameInfo {
	fetchedGame, err := f.client.FetchGame(game.GamePk)
	if err != nil {
		log.WithError(err).Error("Can't get game info")
		return nil
	}
	video := fetchedGame.ExtractGameVideo()
	if len(video) == 0 {
		return nil
	}
	return &d.GameInfo{GamePk: game.GamePk, Video: video, ScoreCard: f.scoreCardGenerator.GenerateScoreCard(game)}
}
