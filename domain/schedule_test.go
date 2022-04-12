package domain

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScheduleUnmarshalling(t *testing.T) {
	dat, _ := os.ReadFile("example_schedule.json")
	schedule := &Schedule{}
	jsonErr := json.Unmarshal(dat, &schedule)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, schedule)
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
