package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestGamesUnmarshalling(t *testing.T) {
	dat, _ := os.ReadFile("game.json")
	game := &Game{}
	jsonErr := json.Unmarshal(dat, &game)
	if jsonErr != nil {
		log.Fatalln(jsonErr)
	}
	fmt.Println(game.Link)
}

func TestScheduleUnmarshalling(t *testing.T) {
	dat, _ := os.ReadFile("schedule.json")
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
