package domain

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestScheduleUnmarshalling(t *testing.T) {
	dat, _ := os.ReadFile("example_schedule.json")
	schedule := &Schedule{}
	jsonErr := json.Unmarshal(dat, &schedule)
	if jsonErr != nil {
		log.Fatalln(jsonErr)
	}
	for _, date := range schedule.Dates {
		for _, game := range date.Games {
			fmt.Printf("Game %v on %v between %v and %v in status %v\n", game.GamePk, date.Date, game.Teams.Home.Team.Name, game.Teams.Away.Team.Name, game.Status.AbstractGameState)
			fmt.Println(game.Teams)
		}
	}
}
