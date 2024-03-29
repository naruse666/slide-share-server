package firebase

import (
	"context"
	"fmt"
	"slide-share/model"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type ISlideRepository interface {
	GetNewestSlideGroup() (*model.SlideGroupResponse, error)
	GetSlideGroup(slideGroupID string) (*model.SlideGroupResponse, error)
	GetSlide(slideGroupID string, slideID string) (*model.SlideResponse, error)
}

type SlideRepository struct {
	client *firestore.Client
}

func NewSlideRepository(client *firestore.Client) ISlideRepository {
	return &SlideRepository{client: client}
}

func (sr *SlideRepository) GetNewestSlideGroup() (*model.SlideGroupResponse, error) {
	ctx := context.Background()
	newestGroup := sr.client.Collection("slide_group").OrderBy("PresentationAt", firestore.Desc).Limit(1).Documents(ctx)
	newestGroupSnapshot, err := newestGroup.Next()

	if err != nil {
		fmt.Printf("error getting newest slide group: %v", err)
		return nil, err
	}

	if !newestGroupSnapshot.Exists() {
		fmt.Printf("newest slide group does not exist")
		return nil, nil
	}

	slides := newestGroupSnapshot.Ref.Collection("slides").Documents(ctx)
	var slideList []model.SlideResponse
	for {
		slideSnapshot, err := slides.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			fmt.Printf("error getting slide: %v", err)
			return nil, err
		}

		var slide model.Slide
		if err := slideSnapshot.DataTo(&slide); err != nil {
			fmt.Printf("error converting slide data: %v", err)
			return nil, err
		}

		speaker, err := NewSpeakerRepository(sr.client).GetSpeakerByID(slide.SpeakerID)
		if err != nil {
			fmt.Printf("error getting speaker by ID: %v", err)
			return nil, err
		}

		slideList = append(slideList, model.SlideResponse{
			ID:                  slide.ID,
			Title:               slide.Title,
			DrivePDFURL:         slide.DrivePDFURL,
			StorageThumbnailURL: slide.StorageThumbnailURL,
			GoogleSlideShareURL: slide.GoogleSlideShareURL,
			SpeakerID:           speaker.SpeakerID,
			SpeakerName:         speaker.DisplayName,
			SpeakerImage:        speaker.Image,
		})
	}

	var slideGroup model.SlideGroup
	if err := newestGroupSnapshot.DataTo(&slideGroup); err != nil {
		fmt.Printf("error converting slide group data: %v", err)
		return nil, err
	}

	return &model.SlideGroupResponse{
		ID:             slideGroup.ID,
		Title:          slideGroup.Title,
		DriveID:        slideGroup.DriveID,
		PresentationAt: slideGroup.PresentationAt,
		SlideList:      slideList,
	}, nil
}

func (sr *SlideRepository) GetSlideGroup(slideGroupID string) (*model.SlideGroupResponse, error) {
	ctx := context.Background()
	slideGroupDoc, err := sr.client.Collection("slide_group").Doc(slideGroupID).Get(ctx)
	if err != nil {
		fmt.Printf("error getting slide group document: %v", err)
		return nil, err
	}

	if !slideGroupDoc.Exists() {
		fmt.Printf("slide group document does not exist")
		return nil, nil
	}

	slides := slideGroupDoc.Ref.Collection("slides").Documents(ctx)
	var slideList []model.SlideResponse
	for {
		slideSnapshot, err := slides.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			fmt.Printf("error getting slide: %v", err)
			return nil, err
		}

		var slide model.Slide
		if err := slideSnapshot.DataTo(&slide); err != nil {
			fmt.Printf("error converting slide data: %v", err)
			return nil, err
		}

		speaker, err := NewSpeakerRepository(sr.client).GetSpeakerByID(slide.SpeakerID)
		if err != nil {
			fmt.Printf("error getting speaker by ID: %v", err)
			return nil, err
		}

		slideList = append(slideList, model.SlideResponse{
			ID:                  slide.ID,
			Title:               slide.Title,
			DrivePDFURL:         slide.DrivePDFURL,
			StorageThumbnailURL: slide.StorageThumbnailURL,
			GoogleSlideShareURL: slide.GoogleSlideShareURL,
			SpeakerID:           speaker.SpeakerID,
			SpeakerName:         speaker.DisplayName,
			SpeakerImage:        speaker.Image,
		})
	}

	var slideGroup model.SlideGroup
	if err := slideGroupDoc.DataTo(&slideGroup); err != nil {
		fmt.Printf("error converting slide group data: %v", err)
		return nil, err
	}

	return &model.SlideGroupResponse{
		ID:             slideGroup.ID,
		Title:          slideGroup.Title,
		DriveID:        slideGroup.DriveID,
		PresentationAt: slideGroup.PresentationAt,
		SlideList:      slideList,
	}, nil
}

func (sr *SlideRepository) GetSlide(slideGroupID string, slideID string) (*model.SlideResponse, error) {
	ctx := context.Background()
	slideDoc, err := sr.client.Collection("slide_group").Doc(slideGroupID).Collection("slides").Doc(slideID).Get(ctx)
	if err != nil {
		fmt.Printf("error getting slide document: %v", err)
		return nil, err
	}

	if !slideDoc.Exists() {
		fmt.Printf("slide document does not exist")
		return nil, nil
	}

	var slide model.Slide
	if err := slideDoc.DataTo(&slide); err != nil {
		fmt.Printf("error converting slide data: %v", err)
		return nil, err
	}

	speaker, err := NewSpeakerRepository(sr.client).GetSpeakerByID(slide.SpeakerID)
	if err != nil {
		fmt.Printf("error getting speaker by ID: %v", err)
		return nil, err
	}

	return &model.SlideResponse{
		ID:                  slide.ID,
		Title:               slide.Title,
		DrivePDFURL:         slide.DrivePDFURL,
		StorageThumbnailURL: slide.StorageThumbnailURL,
		GoogleSlideShareURL: slide.GoogleSlideShareURL,
		SpeakerID:           speaker.SpeakerID,
		SpeakerName:         speaker.DisplayName,
		SpeakerImage:        speaker.Image,
	}, nil
}