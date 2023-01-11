package nhl

import (
	"context"
	log "github.com/sirupsen/logrus"
	"nhl-recap/nhl/domain"
	"nhl-recap/util"
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

type NHL struct {
	client    NHLClient
	frequency time.Duration
}

func NewNHL() *NHL {
	return &NHL{client: NewNHLClient(), frequency: 30 * time.Second}
}

func (nhl *NHL) Subscribe(ctx context.Context) <-chan *GameInfo {
	out := make(chan *GameInfo)
	ticker := time.NewTicker(nhl.frequency)
	go func() {
		defer close(out)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Debug("Tick, starting to request games")
				uniqueGames := unique(ctx, nhl.fetchScheduledGames(ctx))
				for uniqueGame := range uniqueGames {
					select {
					case out <- uniqueGame:
					case <-ctx.Done():
						return
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func unique(ctx context.Context, in <-chan *GameInfo) <-chan *GameInfo {
	out := make(chan *GameInfo)
	set := util.NewSet[int]()

	go func() {
		defer close(out)
		for element := range in {
			if !set.Add(element.GamePk) {
				continue
			}
			select {
			case out <- element:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func (nhl *NHL) fetchScheduledGames(ctx context.Context) chan *GameInfo {
	out := make(chan *GameInfo)
	var wg sync.WaitGroup

	schedule, err := nhl.client.FetchSchedule()
	if err != nil {
		log.WithError(err).Error("Can't get schedule")
	}
	finishedGames := schedule.ExtractFinishedGames()
	log.Debugf("Got schedule, cotains %d finished games", len(finishedGames))

	for _, game := range finishedGames {
		wg.Add(1)
		go func(game domain.Games) {
			defer wg.Done()

			gameInfo := nhl.fetchGameInfo(game)
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

func (nhl *NHL) fetchGameInfo(game domain.Games) *GameInfo {
	fetchedGame, err := nhl.client.FetchGame(game.GamePk)
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
