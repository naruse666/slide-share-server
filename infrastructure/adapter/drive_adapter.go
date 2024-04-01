package adapter

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"google.golang.org/api/drive/v3"
)

type IDriveAdapter interface {
	CreateFolder(folderName string) (string, error)
	DownloadSlideAsPDF(slideID string) ([]byte, error)
	UploadPDFToDrive(pdfData []byte, folderID string, fileName string) (string, error)
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

func (da *DriveAdapter) DownloadSlideAsPDF(slideID string) ([]byte, error) {
	resp, err := da.service.Files.Export(slideID, "application/pdf").Download()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	pdfData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return pdfData, nil
}

func (da *DriveAdapter) UploadPDFToDrive(pdfData []byte, folderID string, fileName string) (string, error) {
	pdfMetadata := &drive.File{
		Name:    fileName,
		Parents: []string{folderID},
	}

	file, err := da.service.Files.Create(pdfMetadata).Media(bytes.NewReader(pdfData)).Do()
	if err != nil {
		return "", err
	}

	permission := &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}
	_, err = da.service.Permissions.Create(file.Id, permission).Do()
	if err != nil {
		return "", err
	}

	publicURL := fmt.Sprintf("https://drive.google.com/file/d/%s/preview", file.Id)

	return publicURL, nil
}
