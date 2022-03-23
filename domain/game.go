package domain

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
