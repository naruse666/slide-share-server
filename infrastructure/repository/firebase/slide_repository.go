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
	GetSlideGroupByPage(page int) ([]*model.SlideGroupResponse, error)
	GetSlideGroup(slideGroupID string) (*model.SlideGroupResponse, error)
	GetSlideGroups() ([]string, error)
	CreateSlideGroup(slideGroup *model.SlideGroup) (string, error)
	GetSlide(slideGroupID string, slideID string) (*model.SlideResponse, error)
	UpdateSlide(slideGroupID string, slideID string, slide *model.Slide) error
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
			IsPublish:           slide.IsPublish,
			DrivePDFURL:         slide.DrivePDFURL,
			StorageThumbnailURL: slide.StorageThumbnailURL,
			GoogleSlideShareURL: slide.GoogleSlideShareURL,
			GroupID:             slide.GroupID,
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

func (sr *SlideRepository) GetSlideGroupByPage(page int) ([]*model.SlideGroupResponse, error) {
	ctx := context.Background()
	slideGroups := sr.client.Collection("slide_group").OrderBy("PresentationAt", firestore.Desc).Offset((page - 1) * 5).Limit(5).Documents(ctx)

	var slideGroupList []*model.SlideGroupResponse
	for {
		slideGroupDoc, err := slideGroups.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			fmt.Printf("error getting slide group document: %v", err)
			return nil, err
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
				IsPublish:           slide.IsPublish,
				Title:               slide.Title,
				DrivePDFURL:         slide.DrivePDFURL,
				StorageThumbnailURL: slide.StorageThumbnailURL,
				GoogleSlideShareURL: slide.GoogleSlideShareURL,
				GroupID:             slide.GroupID,
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

		slideGroupList = append(slideGroupList, &model.SlideGroupResponse{
			ID:             slideGroup.ID,
			Title:          slideGroup.Title,
			DriveID:        slideGroup.DriveID,
			PresentationAt: slideGroup.PresentationAt,
			SlideList:      slideList,
		})
	}

	return slideGroupList, nil
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
			IsPublish:           slide.IsPublish,
			Title:               slide.Title,
			DrivePDFURL:         slide.DrivePDFURL,
			StorageThumbnailURL: slide.StorageThumbnailURL,
			GoogleSlideShareURL: slide.GoogleSlideShareURL,
			GroupID:             slide.GroupID,
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

func (sr *SlideRepository) GetSlideGroups() ([]string, error) {
	ctx := context.Background()
	iter := sr.client.Collection("slide_group").Documents(ctx)
	var slideGroupIDs []string
	for {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			fmt.Printf("error getting slide group document: %v", err)
			return nil, err
		}

		slideGroupIDs = append(slideGroupIDs, doc.Ref.ID)
	}

	return slideGroupIDs, nil
}

func (sr *SlideRepository) CreateSlideGroup(slideGroup *model.SlideGroup) (string, error) {
	ctx := context.Background()
	slideGroupRef := sr.client.Collection("slide_group").Doc(slideGroup.ID)
	_, err := slideGroupRef.Set(ctx, slideGroup)
	if err != nil {
		fmt.Printf("error setting slide group: %v", err)
		return "", err
	}

	return slideGroup.ID, nil
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
		IsPublish:           slide.IsPublish,
		Title:               slide.Title,
		DrivePDFURL:         slide.DrivePDFURL,
		StorageThumbnailURL: slide.StorageThumbnailURL,
		GoogleSlideShareURL: slide.GoogleSlideShareURL,
		SpeakerID:           speaker.SpeakerID,
		SpeakerName:         speaker.DisplayName,
		SpeakerImage:        speaker.Image,
	}, nil
}

func (sr *SlideRepository) UpdateSlide(slideGroupID string, slideID string, slide *model.Slide) error {
	ctx := context.Background()
	slideRef := sr.client.Collection("slide_group").Doc(slideGroupID).Collection("slides").Doc(slideID)
	_, err := slideRef.Set(ctx, slide)
	if err != nil {
		fmt.Printf("error setting slide: %v", err)
		return err
	}

	return nil
}
