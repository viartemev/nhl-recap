package nhl

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"nhl-recap/nhl/domain"
)

var clientIsNil = errors.New("client can't be nil")
var clientError = errors.New("client error")

type NHLClient interface {
	FetchGame(gameId int) *domain.Game
	FetchSchedule() *domain.Schedule
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
	if err == nil || resp.StatusCode != 200 {
		defer resp.Body.Close()
	} else {
		log.WithError(err).Error("Can't get schedule")
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
