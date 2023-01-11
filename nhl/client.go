package nhl

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"nhl-recap/nhl/domain"
	"time"
)

var clientIsNil = errors.New("client can't be nil")
var clientError = errors.New("client error")

type NHLClient interface {
	FetchGame(gameId int) (*domain.Game, error)
	FetchSchedule() (*domain.Schedule, error)
}

func NewNHLClient() NHLClient {
	return &NHLHTTPClient{client: &http.Client{Timeout: 3 * time.Second}}
}

type NHLHTTPClient struct {
	client HTTPClient
}

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

func (nhl *NHLHTTPClient) FetchSchedule() (*domain.Schedule, error) {
	if nhl.client == nil {
		return nil, clientIsNil
	}
	uri := "https://statsapi.web.nhl.com/api/v1/schedule"
	resp, err := nhl.client.Get(uri)
	if err == nil {
		defer resp.Body.Close()
	} else {
		log.WithError(err).Error("Can't get the schedule")
		return nil, clientError
	}
	schedule := &domain.Schedule{}
	errJson := json.NewDecoder(resp.Body).Decode(schedule)
	if errJson != nil {
		log.WithError(err).Error("Can't parse schedule response to json")
		return nil, errJson
	}
	return schedule, nil
}

func (nhl *NHLHTTPClient) FetchGame(gameId int) (*domain.Game, error) {
	if nhl.client == nil {
		return nil, clientIsNil
	}
	uri := fmt.Sprintf("https://statsapi.web.nhl.com/api/v1/game/%d/content", gameId)
	resp, err := nhl.client.Get(uri)
	if err == nil {
		defer resp.Body.Close()
	} else {
		log.WithError(err).Error("Can't get the game")
		return nil, clientError
	}
	game := &domain.Game{}
	errJson := json.NewDecoder(resp.Body).Decode(game)
	if errJson != nil {
		log.WithError(err).Error("Can't parse game response to json")
		return nil, errJson
	}
	return game, nil
}
