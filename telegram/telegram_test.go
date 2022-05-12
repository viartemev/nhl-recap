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
			Name:  "Calgary Flames",
			Score: 0,
		},
		AwayTeam: &nhl.TeamInfo{
			Name:  "Dallas Stars",
			Score: 0,
		},
	}
	telegramMessage := GameInfoToTelegramMessage(game)
	fmt.Println(telegramMessage)

	//TODO write test
}
