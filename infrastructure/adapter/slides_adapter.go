package adapter

import (
	"fmt"
	"io"
	"net/http"

	slides "google.golang.org/api/slides/v1"
)

type ISlidesAdapter interface {
	FetchSlideThumbnail(slideID string) ([]byte, error)
}

type SlidesAdapter struct {
	service *slides.Service
}

func NewSlidesAdapter(service *slides.Service) ISlidesAdapter {
	return &SlidesAdapter{service: service}
}

func (sa *SlidesAdapter) FetchSlideThumbnail(slideID string) ([]byte, error) {
	presentation, err := sa.service.Presentations.Get(slideID).Do()
	if err != nil {
		return nil, err
	}

	if len(presentation.Slides) == 0 {
		return nil, fmt.Errorf("no slides found in presentation")
	}

	firstSlideID := presentation.Slides[0].ObjectId

	thumbnail, err := sa.service.Presentations.Pages.GetThumbnail(slideID, firstSlideID).Do()
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(thumbnail.ContentUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	thumbnailData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return thumbnailData, nil
}
