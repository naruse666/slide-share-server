package model

import "time"

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Image        string    `json:"image"`
	Role         string    `json:"role"`
	SpeakerID    string    `json:"speaker_id"`
	DisplayName  string    `json:"display_name"`
	School       string    `json:"school"`
	Course       string    `json:"course"`
	IsTopDisplay bool      `json:"is_top_display"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SignInUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Image string `json:"image"`
}

type JWTPayload struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}

type SpeakerResponse struct {
	SpeakerID    string    `json:"speaker_id"`
	DisplayName  string    `json:"display_name"`
	Image        string    `json:"image"`
	School       string    `json:"school"`
	Course       string    `json:"course"`
	IsTopDisplay bool      `json:"is_top_display"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SpeakerWithSlideResponse struct {
	SpeakerID   string  `json:"speaker_id"`
	DisplayName string  `json:"display_name"`
	Image       string  `json:"image"`
	School      string  `json:"school"`
	Course      string  `json:"course"`
	SlideList   []Slide `json:"slide_list"`
}
