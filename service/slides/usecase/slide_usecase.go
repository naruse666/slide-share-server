package usecase

import (
	"fmt"
	"slide-share/infrastructure/adapter"
	"slide-share/infrastructure/repository/firebase"
	"slide-share/infrastructure/repository/storage"
	"slide-share/model"
	"slide-share/utils"
)

type ISlideUsecase interface {
	GetNewestSlideGroup() (*model.SlideGroupResponse, error)
	GetSlideGroupByPage(page int) ([]*model.SlideGroupResponse, error)
	GetSlideGroup(slideGroupID string) (*model.SlideGroupResponse, error)
	GetSlideGroups() ([]string, error)
	CreateSlideGroup(slideGroup *model.SlideGroup) (string, error)
	GetSlide(slideGroupID string, slideID string) (*model.SlideResponse, error)
	UpdateSlide(slideGroupID string, slideID string, slide *model.Slide) error
	UploadSlideBySlidesURL(SlideUploadBySlidesURL *model.SlideUploadBySlidesURL) error
}

type slideUsecase struct {
	sr firebase.ISlideRepository
	tr storage.IThumbnailRepository
	sa adapter.ISlidesAdapter
	da adapter.IDriveAdapter
}

func NewSlideUsecase(sr firebase.ISlideRepository, da adapter.IDriveAdapter, sa adapter.ISlidesAdapter, tr storage.IThumbnailRepository) ISlideUsecase {
	return &slideUsecase{sr: sr, da: da, sa: sa, tr: tr}
}

func (su *slideUsecase) GetNewestSlideGroup() (*model.SlideGroupResponse, error) {
	slideGroup, err := su.sr.GetNewestSlideGroup()
	if err != nil {
		return nil, err
	}

	return slideGroup, nil
}

func (su *slideUsecase) GetSlideGroupByPage(page int) ([]*model.SlideGroupResponse, error) {
	slideGroups, err := su.sr.GetSlideGroupByPage(page)
	if err != nil {
		return nil, err
	}

	return slideGroups, nil
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

func (su *slideUsecase) UpdateSlide(slideGroupID string, slideID string, slide *model.Slide) error {
	err := su.sr.UpdateSlide(slideGroupID, slideID, slide)
	if err != nil {
		return err
	}

	return nil
}

func (su *slideUsecase) UploadSlideBySlidesURL(SlideUploadBySlidesURL *model.SlideUploadBySlidesURL) error {
	// Slides URLからIDを取得
	slidesId, err := utils.ExtractSlideIDFromURL(SlideUploadBySlidesURL.SlidesURL)
	if err != nil {
		fmt.Println("error extracting slide id from url: ", err)
		return err
	}
	// Drive APIを使用し、IDからPDFを取得
	pdfData, err := su.da.DownloadSlideAsPDF(slidesId)
	if err != nil {
		return err
	}

	// Slides APIを使用し、IDからサムネイル画像を取得
	thumbnailData, err := su.sa.FetchSlideThumbnail(slidesId)
	if err != nil {
		return err
	}

	// Drive APIを使用し、PDFをDriveにアップロード
	pdfURL, err := su.da.UploadPDFToDrive(pdfData, SlideUploadBySlidesURL.DriveID, SlideUploadBySlidesURL.Title)
	if err != nil {
		return err
	}

	// Firebase Cloud Storageを使用し、サムネイル画像をアップロード
	thumbnailURL, err := su.tr.UploadThumbnail(thumbnailData, SlideUploadBySlidesURL.ID, SlideUploadBySlidesURL.GroupID)
	if err != nil {
		return err
	}

	// PDFのURL, サムネイル画像のURL, スライド情報を Firestoreに保存
	slide := &model.Slide{
		ID:                  SlideUploadBySlidesURL.ID,
		IsPublish:           SlideUploadBySlidesURL.IsPublish,
		Title:               SlideUploadBySlidesURL.Title,
		DrivePDFURL:         pdfURL,
		StorageThumbnailURL: thumbnailURL,
		GoogleSlideShareURL: "",
		GroupID:             SlideUploadBySlidesURL.GroupID,
		SpeakerID:           SlideUploadBySlidesURL.SpeakerID,
	}

	_, err = su.sr.CreateSlide(SlideUploadBySlidesURL.GroupID, slide)
	if err != nil {
		return err
	}

	return nil
}
