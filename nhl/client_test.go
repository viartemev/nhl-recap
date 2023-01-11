package nhl

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

type MockHttpClient struct {
}

func (c *MockHttpClient) Get(url string) (*http.Response, error) {
	if url == "https://statsapi.web.nhl.com/api/v1/schedule" {
		data, _ := os.ReadFile("schedule_response.json")
		resp := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBuffer(data)),
		}
		return resp, nil
	} else if strings.HasPrefix(url, "https://statsapi.web.nhl.com/api/v1/game") {
		data, _ := os.ReadFile("game_response.json")
		resp := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBuffer(data)),
		}
		return resp, nil
	}
	return nil, errors.New("no mock response for that request")
}

func TestNHLHTTPClient_FetchSchedule_Client_Is_Nil(t *testing.T) {
	client := NHLHTTPClient{client: nil}
	schedule, err := client.FetchSchedule()
	if schedule != nil {
		t.Errorf("Schedule should be nil")
	}
	if !errors.Is(err, clientIsNil) {
		t.Errorf("Should be error is nil error")
	}
}

func TestNHLHTTPClient_FetchSchedule_Timeout_Error(t *testing.T) {
	httpClient := &http.Client{Timeout: time.Millisecond}
	client := NHLHTTPClient{client: httpClient}
	schedule, err := client.FetchSchedule()
	if !errors.Is(err, clientError) {
		t.Errorf("Should be client error")
	}
	if schedule != nil {
		t.Errorf("Schedule should be nil")
	}
}

func TestNHLHTTPClient_FetchSchedule_Real_Data(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	httpClient := &http.Client{}
	client := NHLHTTPClient{client: httpClient}
	schedule, err := client.FetchSchedule()
	if err != nil {
		t.Errorf("Error should be nil")
	}
	if schedule == nil {
		t.Errorf("Schedule shouldn't be nil")
	}
	if len(schedule.Dates) == 0 {
		t.Errorf("Array shouldn't be empty")
	}
	for _, date := range schedule.Dates {
		assert.NotEmpty(t, date.Date, "Date should be empty")
		for _, game := range date.Games {
			assert.NotEmpty(t, game.GamePk, "GamePk should not be empty")
			assert.NotEmpty(t, game.Status.AbstractGameState, "Game state should not be empty")
			assert.NotEmpty(t, game.Teams.Away.Team.Name, "Team name should not be empty")
			assert.NotEmpty(t, game.Teams.Away.Score, "Score should not be empty")
			assert.NotEmpty(t, game.Teams.Home.Team.Name, "Team name should not be empty")
			assert.NotEmpty(t, game.Teams.Home.Score, "Score should not be empty")
		}
	}
}

func TestNHLHTTPClient_FetchSchedule_Mock_Client_Marshalling(t *testing.T) {
	client := NHLHTTPClient{client: &MockHttpClient{}}
	schedule, _ := client.FetchSchedule()
	schedule, err := client.FetchSchedule()
	if err != nil {
		t.Errorf("Error should be nil")
	}
	if schedule == nil {
		t.Errorf("Schedule shouldn't be nil")
	}
	if len(schedule.Dates) == 0 {
		t.Errorf("Array shouldn't be empty")
	}
	for _, date := range schedule.Dates {
		assert.NotEmpty(t, date.Date, "Date should be empty")
		for _, game := range date.Games {
			assert.NotEmpty(t, game.GamePk, "GamePk should not be empty")
			assert.NotEmpty(t, game.Status.AbstractGameState, "Game state should not be empty")
			assert.NotEmpty(t, game.Teams.Away.Team.Name, "Team name should not be empty")
			assert.NotEmpty(t, game.Teams.Away.Score, "Score should not be empty")
			assert.NotEmpty(t, game.Teams.Home.Team.Name, "Team name should not be empty")
			assert.NotEmpty(t, game.Teams.Home.Score, "Score should not be empty")
		}
	}
}

func TestNHLHTTPClient_FetchGames_Client_Is_Nil(t *testing.T) {
	client := NHLHTTPClient{client: nil}
	key := 42
	game, err := client.FetchGame(key)
	if game != nil {
		t.Errorf("Game should be nil")
	}
	if !errors.Is(err, clientIsNil) {
		t.Errorf("Should be error is nil error")
	}
}

func TestNHLHTTPClient_FetchGame_Timeout_Error(t *testing.T) {
	httpClient := &http.Client{Timeout: time.Millisecond}
	client := NHLHTTPClient{client: httpClient}
	key := 2022020648
	schedule, err := client.FetchGame(key)
	if !errors.Is(err, clientError) {
		t.Errorf("Should be client error")
	}
	if schedule != nil {
		t.Errorf("Schedule should be nil")
	}
}

func TestNHLHTTPClient_FetchGame_Real_Data(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	httpClient := &http.Client{}
	client := NHLHTTPClient{client: httpClient}
	key := 2021020988
	game, err := client.FetchGame(key)
	if err != nil {
		t.Errorf("Error should be nil")
	}
	if game == nil {
		t.Errorf("Game shouldn't be nil")
	}
	assert.Equal(t, game.Link, "/api/v1/game/2021020988/content", "Game link is not equal")
	assert.NotNil(t, game.Media)
	assert.NotEmpty(t, game.Media.Epg)
	for _, epg := range game.Media.Epg {
		assert.NotNil(t, epg.Title, "Title should be not null")
		for _, item := range epg.Items {
			assert.NotNil(t, item)
			for _, playback := range item.Playbacks {
				assert.NotEmpty(t, playback.Name, "Game name should not be empty")
				assert.NotEmpty(t, playback.Url, "Url name should not be empty")
			}
		}
	}
}

func TestNHLHTTPClient_FetchGame_Mock_Client_Marshalling(t *testing.T) {
	client := NHLHTTPClient{client: &MockHttpClient{}}
	key := 2021020988
	game, err := client.FetchGame(key)
	if err != nil {
		t.Errorf("Error should be nil")
	}
	if game == nil {
		t.Errorf("Game shouldn't be nil")
	}
	assert.Equal(t, game.Link, "/api/v1/game/2021020988/content", "Game link is not equal")
	assert.NotNil(t, game.Media)
	assert.NotEmpty(t, game.Media.Epg)
	for _, epg := range game.Media.Epg {
		assert.NotNil(t, epg.Title, "Title should be not null")
		for _, item := range epg.Items {
			assert.NotNil(t, item)
			for _, playback := range item.Playbacks {
				assert.NotEmpty(t, playback.Name, "Game name should not be empty")
				assert.NotEmpty(t, playback.Url, "Url name should not be empty")
			}
		}
	}
}
