package usecase

import (
	"slide-share/infrastructure/repository/firebase"
	"slide-share/model"
)

type ISpeakerUsecase interface {
	GetSpeakerList() ([]model.SpeakerResponse, error)
}

type speakerUsecase struct {
	sr firebase.ISpeakerRepository
}

func NewSpeakerUsecase(sr firebase.ISpeakerRepository) ISpeakerUsecase {
	return &speakerUsecase{sr: sr}
}

func (su *speakerUsecase) GetSpeakerList() ([]model.SpeakerResponse, error) {
	speakers, err := su.sr.GetSpeakerList()
	if err != nil {
		return nil, err
	}

	return speakers, nil
}
