package domain

import "strings"

type Game struct {
	Link  string `json:"link"`
	Media Media
}

func (g *Game) ExtractGameVideo() (video string) {
	for _, media := range g.Media.Epg {
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

type Media struct {
	Epg []Epg
}

type Epg struct {
	Title string `json:"title"`
	Items []EpgItem
}

type EpgItem struct {
	Playbacks []Playback
}
type Playback struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
