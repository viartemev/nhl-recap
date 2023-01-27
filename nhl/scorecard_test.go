package nhl

import (
	"bytes"
	"image"
	"image/png"
	"nhl-recap/nhl/domain"
	"nhl-recap/nhl/logos"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestGamesUnmarshalling(t *testing.T) {
	game := domain.ScheduleGame{GamePk: 123, Status: struct {
		AbstractGameState string `json:"abstractGameState"`
	}(struct{ AbstractGameState string }{AbstractGameState: "FINAL"}),
		Teams: domain.Teams(struct {
			Away domain.TeamScore
			Home domain.TeamScore
		}{Away: domain.TeamScore(struct {
			Team  domain.Team
			Score int
		}{Team: domain.Team{Name: "PIT"}, Score: 7}), Home: domain.TeamScore(struct {
			Team  domain.Team
			Score int
		}{Team: domain.Team{Name: "ANA"}, Score: 2})}),
	}
	scoreCard := NewScoreCardGenerator(logos.LoadLogos())
	save(scoreCard.GenerateScoreCard(game))
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
