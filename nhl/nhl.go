package nhl

import (
	"context"
	log "github.com/sirupsen/logrus"
	"nhl-recap/nhl/domain"
	"nhl-recap/util"
	"sync"
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

type NHL struct {
	Fetcher util.Fetcher[*GameInfo]
}

func NewNHL() *NHL {
	return &NHL{Fetcher: &NHLFetcher{client: NewNHLClient(), uniqueGames: *util.NewSet[int]()}}
}

type NHLFetcher struct {
	client      NHLClient
	uniqueGames util.Set[int]
}

func (f *NHLFetcher) Fetch(ctx context.Context) chan *GameInfo {
	games := f.fetchScheduledGames(ctx)
	return util.Filter(ctx, games, func(info *GameInfo) bool { return f.uniqueGames.Add(info.GamePk) })
}

func (f *NHLFetcher) fetchScheduledGames(ctx context.Context) chan *GameInfo {
	out := make(chan *GameInfo)
	var wg sync.WaitGroup

	schedule, err := f.client.FetchSchedule()
	if err != nil {
		log.WithError(err).Error("Can't get schedule")
	}
	finishedGames := schedule.ExtractFinishedGames()
	log.Debugf("Got schedule, cotains %d finished games", len(finishedGames))

	for _, game := range finishedGames {
		wg.Add(1)
		go func(game domain.Games) {
			defer wg.Done()

			gameInfo := f.fetchGameInfo(game)
			log.Debugf("Got gameInfo %d", gameInfo.GamePk)
			if gameInfo != nil {
				select {
				case out <- gameInfo:
				case <-ctx.Done():
					return
				}
			}
		}(game)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func (f *NHLFetcher) fetchGameInfo(game domain.Games) *GameInfo {
	fetchedGame, err := f.client.FetchGame(game.GamePk)
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
