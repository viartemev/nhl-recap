package nhl

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"image/color"
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
	settings := GeneratorSettings{
		Width:      300,
		Height:     100,
		Background: color.White,
	}
	return &NHL{Fetcher: &NHLFetcher{
		client:             NewNHLClient(),
		uniqueGames:        util.NewSet[int](),
		scoreCardGenerator: NewScoreCardGenerator(logos.LoadLogos(), settings)},
	}
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
	gamesChannel := make([]<-chan *d.GameInfo, 0)
	for _, game := range finishedGames {
		gamesChannel = append(gamesChannel, f.fetchGameInfo(game))
	}
	games := util.FanIn[*d.GameInfo](ctx, gamesChannel...)
	uniqueGame := func(info *d.GameInfo) bool { return f.uniqueGames.Add(info.GamePk) }
	notNil := func(info *d.GameInfo) bool { return info != nil }
	return util.Filter(ctx, games, util.And(notNil, uniqueGame)), nil
}

func (f *NHLFetcher) fetchGameInfo(game domain.ScheduleGame) <-chan *d.GameInfo {
	out := make(chan *d.GameInfo)
	go func() {
		defer close(out)
		fetchedGame, err := f.client.FetchGame(game.GamePk)
		if err != nil {
			log.WithError(err).Error("Can't get game info")
			//TODO what should I do with channel?
			return
		}
		video := fetchedGame.ExtractGameVideo()
		if len(video) == 0 {
			//TODO what should I do with channel?
			return
		}
		out <- &d.GameInfo{GamePk: game.GamePk, Video: video, ScoreCard: f.scoreCardGenerator.GenerateScoreCard(game)}
	}()
	return out
}
