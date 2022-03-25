package domain

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

func TestGamesUnmarshalling(t *testing.T) {
	dat, _ := os.ReadFile("example_game.json")
	game := &Game{}
	jsonErr := json.Unmarshal(dat, &game)
	if jsonErr != nil {
		log.Fatalln(jsonErr)
	}
	if game == nil {
		log.Fatalln("Game is nil")
	}
	if game.Link != "/api/v1/game/2021020988/content" {
		log.Fatalln("Link is wrong")
	}
}
