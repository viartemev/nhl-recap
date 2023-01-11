package nhl

import (
	"fmt"
	"testing"
)

func TestNHL(t *testing.T) {
	nhlClient := NewNHLClient()
	games := nhlClient.Subscribe()
	for game := range games {
		fmt.Println(game)
	}
}
