package domain

import "fmt"

type Schedule struct {
	Dates []DateGames
}

type DateGames struct {
	Date  string `json:"date"`
	Games []Games
}

type Games struct {
	GamePk int `json:"gamePk"`
	Teams  Teams
	Status struct {
		AbstractGameState string `json:"abstractGameState"`
	}
}

type Teams struct {
	Away TeamScore `json:"away"`
	Home TeamScore `json:"home"`
}

func (t *Teams) TeamsAndWinner() string {
	if t.Home.Score > t.Away.Score {
		return fmt.Sprintf("👑%v vs %v", t.Home.Team.Name, t.Away.Team.Name)
	} else {
		return fmt.Sprintf("%v vs 👑%v", t.Home.Team.Name, t.Away.Team.Name)
	}
}

type TeamScore struct {
	Team  Team
	Score int `json:"score"`
}

type Team struct {
	Name string `json:"name"`
}
