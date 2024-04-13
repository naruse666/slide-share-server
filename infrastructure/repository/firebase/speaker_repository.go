package firebase

import (
	"context"
	"fmt"
	"slide-share/model"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type ISpeakerRepository interface {
	GetSpeakerByID(speakerID string) (*model.SpeakerWithSlideResponse, error)
	GetSpeakerList() ([]model.SpeakerResponse, error)
}

type SpeakerRepository struct {
	client *firestore.Client
}

func NewSpeakerRepository(client *firestore.Client) ISpeakerRepository {
	return &SpeakerRepository{client: client}
}

func (sr *SpeakerRepository) GetSpeakerByID(speakerID string) (*model.SpeakerWithSlideResponse, error) {
	iter := sr.client.Collection("users").Where("SpeakerID", "==", speakerID).Documents(context.Background())
	speakerDoc, err := iter.Next()
	if err == iterator.Done {
		fmt.Println("No more items in iterator - No document found")
		return nil, fmt.Errorf("no speaker document found for ID: %s", speakerID)
	}
	if err != nil {
		fmt.Printf("error getting speaker document: %v\n", err)
		return nil, err
	}

	var speaker model.SpeakerResponse
	speakerDoc.DataTo(&speaker)

	slides, err := sr.client.CollectionGroup("slides").Where("SpeakerID", "==", speakerID).Documents(context.Background()).GetAll()
	if err != nil {
		fmt.Printf("error getting slide collection: %v", err)
		return nil, err
	}

	var slideCollection []model.Slide
	for _, slide := range slides {
		var s model.Slide
		if err := slide.DataTo(&s); err != nil {
			fmt.Printf("error decoding slide data: %v\n", err)
			continue
		}
		slideCollection = append(slideCollection, s)
	}

	speakerWithSlide := model.SpeakerWithSlideResponse{
		SpeakerID:   speaker.SpeakerID,
		DisplayName: speaker.DisplayName,
		Image:       speaker.Image,
		School:      speaker.School,
		Course:      speaker.Course,
		SlideList:   slideCollection,
	}

	return &speakerWithSlide, nil
}

func (sr *SpeakerRepository) GetSpeakerList() ([]model.SpeakerResponse, error) {
	speakers, err := sr.client.Collection("users").Where("Role", "!=", "user").OrderBy("CreatedAt", firestore.Asc).Documents(context.Background()).GetAll()
	if err != nil {
		fmt.Printf("error getting speaker collection: %v", err)
	}

	var speakerCollection []model.SpeakerResponse
	for _, speaker := range speakers {
		var s model.SpeakerResponse
		speaker.DataTo(&s)
		speakerCollection = append(speakerCollection, s)
	}

	return speakerCollection, nil
}
