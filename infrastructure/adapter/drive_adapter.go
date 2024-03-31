package adapter

import (
	"os"

	"google.golang.org/api/drive/v3"
)

type IDriveAdapter interface {
	CreateFolder(folderName string) (string, error)
}

type DriveAdapter struct {
	service *drive.Service
}

func NewDriveAdapter(service *drive.Service) IDriveAdapter {
	return &DriveAdapter{service: service}
}

func (da *DriveAdapter) CreateFolder(folderName string) (string, error) {
	folderMetadata := &drive.File{
		Name:     folderName,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{os.Getenv("DRIVE_ROOT_FOLDER_ID")},
	}

	folder, err := da.service.Files.Create(folderMetadata).Do()
	if err != nil {
		return "", err
	}

	return folder.Id, nil
}
