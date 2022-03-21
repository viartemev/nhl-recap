package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type GameInfo struct {
	Title string
	Video string
}

func main() {
	start := time.Now()
	var wg sync.WaitGroup
	gamesVideos := make(map[int]*GameInfo)
	schedule := fetchSchedule()
	finishedGames := filterFinishedGames(schedule)

	for _, game := range finishedGames {
		wg.Add(1)
		go func(game Games) {
			gameInfo := fetchGameInfo(game.GamePk)
			video := extractGameVideo(gameInfo)
			title := fmt.Sprintf("%v vs %v", game.Teams.Home.Team.Name, game.Teams.Away.Team.Name)
			gamesVideos[game.GamePk] = &GameInfo{title, video}
			defer wg.Done()
		}(game)
	}
	wg.Wait()
	for _, info := range gamesVideos {
		fmt.Printf("Game %v %v\n", info.Title, info.Video)
	}
	log.Printf("Operations took %s", time.Since(start))
}

func filterFinishedGames(schedule Schedule) (finishedGames []Games) {
	for _, date := range schedule.Dates {
		for _, game := range date.Games {
			if game.Status.AbstractGameState == "Final" {
				finishedGames = append(finishedGames, game)
			}
		}
	}
	return
}

func extractGameVideo(game Game) (video string) {
	for _, media := range game.Media.Epg {
		if media.Title == "Recap" {
			for _, item := range media.Items {
				for _, playback := range item.Playbacks {
					if strings.Contains(playback.Name, "FLASH_1800K") {
						video = playback.Url
					}
				}
			}
		}
	}
	return
}

func fetchGameInfo(gamePk int) Game {
	client := http.Client{
		Timeout: time.Second * 2,
	}
	url := "https://statsapi.web.nhl.com/api/v1/game/" + fmt.Sprintf("%v", gamePk) + "/content"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	if res.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(res.Body)
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	game := &Game{}
	jsonErr := json.Unmarshal(body, &game)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return *game
}

type Game struct {
	Link  string `json:"link"`
	Media struct {
		Epg []struct {
			Title string `json:"title"`
			Items []struct {
				Playbacks []struct {
					Name string `json:"name"`
					Url  string `json:"url"`
				}
			}
		}
	}
}

func fetchSchedule() Schedule {
	client := http.Client{
		Timeout: time.Second * 2,
	}
	url := "https://statsapi.web.nhl.com/api/v1/schedule"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	if res.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(res.Body)
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	schedule := &Schedule{}
	jsonErr := json.Unmarshal(body, &schedule)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return *schedule
}

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
		Team Team
	} `json:"away"`
	Home struct {
		Team Team
	} `json:"home"`
}

type Team struct {
	Name string `json:"name"`
}
