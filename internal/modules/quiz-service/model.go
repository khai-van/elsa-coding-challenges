package quizservice

import "time"

type Quiz struct {
	ID        string     `json:"id"`
	Questions []Question `json:"-"`
	Start     time.Time  `json:"start"`
}

type Question struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Answer      string    `json:"-"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
}
