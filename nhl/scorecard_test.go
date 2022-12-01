package nhl

import (
    "bytes"
    log "github.com/sirupsen/logrus"
    "image"
    "image/png"
    "os"
    "testing"
    "time"
)

func TestGamesUnmarshalling(t *testing.T) {
	game := &GameInfo{
		GamePk:   0,
		Video:    "Hello",
		HomeTeam: &TeamInfo{
            Name:  "Home Team",
            Score: 1,
        },
		AwayTeam: &TeamInfo{
            Name:  "Away Team",
            Score: 3,
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
    out, err2 := os.Create(time.Now().Format("2006-01-02 15:04:05.000000") + ".png")
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
