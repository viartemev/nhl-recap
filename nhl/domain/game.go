package domain

type Game struct {
	Link  string `json:"link"`
	Media Media
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
