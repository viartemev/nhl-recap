package domain

type Schedule struct {
	Dates []struct {
		Date  string `json:"date"`
		Games []Games
	}
}

type Games struct {
	GamePk int `json:"gamePk"`
	Teams  Teams
	Status struct {
		AbstractGameState string `json:"abstractGameState"`
	}
}

type Teams struct {
	Away struct {
		Team  Team
		Score int `json:"score"`
	} `json:"away"`
	Home struct {
		Team  Team
		Score int `json:"score"`
	} `json:"home"`
}

type Team struct {
	Name string `json:"name"`
}
