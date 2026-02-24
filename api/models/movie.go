package models

type Movie struct {
	Id           int64    `json:"id"`
	Rating       int      `json:"rating"`
	Movie        string   `json:"movie"`
	Year         int      `json:"year"`
	Country      string   `json:"country"`
	RatingBall   float64  `json:"rating_ball"`
	Overview     string   `json:"overview"`
	Director     string   `json:"director"`
	Screenwriter []string `json:"screenwriter"`
	Actors       []string `json:"actors"`
	UrlLogo      string   `json:"url_logo"`
}

type GetMoviesListResp struct {
	Movies []Movie `json:"movies"`
	Count  int64   `json:"count"`
}
