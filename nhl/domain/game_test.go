package domain

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGamesUnmarshalling(t *testing.T) {
	dat, _ := os.ReadFile("example_game.json")
	game := &Game{}
	jsonErr := json.Unmarshal(dat, &game)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, game)
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
