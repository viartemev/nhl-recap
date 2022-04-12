package domain

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

type TeamScore struct {
	Team  Team
	Score int `json:"score"`
}

type Team struct {
	Name string `json:"name"`
}
