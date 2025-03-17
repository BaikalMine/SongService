package models

// Song представляет структуру песни.
type Song struct {
	ID          int    `json:"id" example:"1"`
	Group       string `json:"group" example:"Muse"`
	Song        string `json:"song" example:"Supermassive Black Hole"`
	ReleaseDate string `json:"releaseDate" example:"16.07.2006"`
	Lyrics      string `json:"text" example:"Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?"`
	Link        string `json:"link" example:"https://www.youtube.com/watch?v=Xsp3_a-PMTw"`
}
