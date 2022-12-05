package nhl

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestGamesUnmarshalling(t *testing.T) {
	game := &GameInfo{
		GamePk: 0,
		Video:  "link",
		HomeTeam: &TeamInfo{
			Name:  "COL",
			Score: 7,
		},
		AwayTeam: &TeamInfo{
			Name:  "PHI",
			Score: 1,
		},
	}
	scoreCard := GenerateScoreCard(game)
	save(scoreCard)
}

func save(imgBytes []byte) {
	img, _, err1 := image.Decode(bytes.NewReader(imgBytes))
	if err1 != nil {
		log.Fatal(err1)
	}
	out, err2 := os.Create("scorecard.png")
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.Fatal(err2)
		}
	}(out)

	err3 := png.Encode(out, img)
	if err3 != nil {
		log.Fatal(err3)
	}
}
