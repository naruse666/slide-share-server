package usecase

import (
	"slide-share/infrastructure/repository/firebase"
	"slide-share/model"
)

type ISlideUsecase interface {
	GetNewestSlideGroup() (*model.SlideGroupResponse, error)
	GetSlideGroup(slideGroupID string) (*model.SlideGroupResponse, error)
	GetSlide(slideGroupID string, slideID string) (*model.SlideResponse, error)
}

type slideUsecase struct {
	sr firebase.ISlideRepository
}

func NewSlideUsecase(sr firebase.ISlideRepository) ISlideUsecase {
	return &slideUsecase{sr: sr}
}

func (su *slideUsecase) GetNewestSlideGroup() (*model.SlideGroupResponse, error) {
	slideGroup, err := su.sr.GetNewestSlideGroup()
	if err != nil {
		return nil, err
	}

	return slideGroup, nil
}

func (su *slideUsecase) GetSlideGroup(slideGroupID string) (*model.SlideGroupResponse, error) {
	slideGroup, err := su.sr.GetSlideGroup(slideGroupID)
	if err != nil {
		return nil, err
	}

	return slideGroup, nil
}

func (su *slideUsecase) GetSlide(slideGroupID string, slideID string) (*model.SlideResponse, error) {
	slide, err := su.sr.GetSlide(slideGroupID, slideID)
	if err != nil {
		return nil, err
	}

	return slide, nil
}
