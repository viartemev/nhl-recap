package telegram

import (
	"fmt"
	"nhl-recap/nhl"
	"testing"
)

func TestGameInfoToTelegramMessage(t *testing.T) {
	game := &nhl.GameInfo{
		Video: "video-link",
		HomeTeam: &nhl.TeamInfo{
			Name:  "Team 1",
			Score: 10,
		},
		AwayTeam: &nhl.TeamInfo{
			Name:  "Team 2",
			Score: 1,
		},
	}
	telegramMessage := GameInfoToTelegramMessage(game)
	fmt.Println(telegramMessage)

	//TODO write test
}
