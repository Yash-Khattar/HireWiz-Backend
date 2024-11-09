package handlers

import (
	"context"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cld *cloudinary.Cloudinary

func InitCloudinary() error {
	var err error
	cld, err = cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		return err
	}
	return nil
}

func UploadResume(file *multipart.FileHeader) (string, error) {
	// Open the file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Upload to cloudinary
	uploadResult, err := cld.Upload.Upload(
		context.Background(),
		src,
		uploader.UploadParams{
			Folder:         "resumes",
			AllowedFormats: []string{"pdf"},
		})

	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
} 