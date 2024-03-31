package usecase

import (
	"slide-share/infrastructure/adapter"
	"slide-share/infrastructure/repository/firebase"
	"slide-share/model"
)

type ISlideUsecase interface {
	GetNewestSlideGroup() (*model.SlideGroupResponse, error)
	GetSlideGroup(slideGroupID string) (*model.SlideGroupResponse, error)
	GetSlideGroups() ([]string, error)
	CreateSlideGroup(slideGroup *model.SlideGroup) (string, error)
	GetSlide(slideGroupID string, slideID string) (*model.SlideResponse, error)
}

type slideUsecase struct {
	sr firebase.ISlideRepository
	da adapter.IDriveAdapter
}

func NewSlideUsecase(sr firebase.ISlideRepository, da adapter.IDriveAdapter) ISlideUsecase {
	return &slideUsecase{sr: sr, da: da}
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

func (su *slideUsecase) GetSlideGroups() ([]string, error) {
	slideGroups, err := su.sr.GetSlideGroups()
	if err != nil {
		return nil, err
	}

	return slideGroups, nil
}

func (su *slideUsecase) CreateSlideGroup(slideGroup *model.SlideGroup) (string, error) {
	DriveID, err := su.da.CreateFolder(slideGroup.Title)
	if err != nil {
		return "", err
	}

	slideGroup.DriveID = DriveID
	_, err = su.sr.CreateSlideGroup(slideGroup)
	if err != nil {
		return "", err
	}

	return DriveID, nil
}

func (su *slideUsecase) GetSlide(slideGroupID string, slideID string) (*model.SlideResponse, error) {
	slide, err := su.sr.GetSlide(slideGroupID, slideID)
	if err != nil {
		return nil, err
	}

	return slide, nil
}
