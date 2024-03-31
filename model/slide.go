package model

import "time"

type SlideGroup struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	DriveID        string    `json:"drive_id"`
	PresentationAt time.Time `json:"presentation_at"`
}

type Slide struct {
	ID                  string `json:"id"`
	IsPublish           bool   `json:"is_publish"`
	Title               string `json:"title"`
	DrivePDFURL         string `json:"drive_pdf_url"`
	StorageThumbnailURL string `json:"storage_thumbnail_url"`
	GoogleSlideShareURL string `json:"google_slide_share_url"`
	GroupID             string `json:"group_id"`
	SpeakerID           string `json:"speaker_id"`
}

type SlideGroupResponse struct {
	ID             string          `json:"id"`
	Title          string          `json:"title"`
	DriveID        string          `json:"drive_id"`
	PresentationAt time.Time       `json:"presentation_at"`
	SlideList      []SlideResponse `json:"slide_list"`
}

type SlideResponse struct {
	ID                  string `json:"id"`
	IsPublish           bool   `json:"is_publish"`
	Title               string `json:"title"`
	DrivePDFURL         string `json:"drive_pdf_url"`
	StorageThumbnailURL string `json:"storage_thumbnail_url"`
	GoogleSlideShareURL string `json:"google_slide_share_url"`
	GroupID             string `json:"group_id"`
	SpeakerID           string `json:"speaker_id"`
	SpeakerName         string `json:"speaker_name"`
	SpeakerImage        string `json:"speaker_image"`
}
