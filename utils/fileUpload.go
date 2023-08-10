package utils

import (
	"context"
	"time"

	"github.com/chiboycalix/hotel-booking-system-backend/common"
	"github.com/chiboycalix/hotel-booking-system-backend/models"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type mediaUpload interface {
	FileUpload(file models.File) (string, error)
	RemoteUpload(url models.Url) (string, error)
}

type media struct{}

func NewMediaUpload() mediaUpload {
	return &media{}
}

func (*media) FileUpload(file models.File) (string, error) {
	//validate
	err := validate.Struct(file)
	if err != nil {
		return "", err
	}

	//upload
	uploadUrl, err := ImageUploadHelper(file.File)
	if err != nil {
		return "", err
	}
	return uploadUrl, nil
}
func (*media) RemoteUpload(url models.Url) (string, error) {
	err := validate.Struct(url)
	if err != nil {
		return "", err
	}

	uploadUrl, errUrl := ImageUploadHelper(url.Url)
	if errUrl != nil {
		return "", err
	}
	return uploadUrl, nil
}

func ImageUploadHelper(input interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//create cloudinary instance
	cld, err := cloudinary.NewFromParams(common.CloudinaryCloudName(), common.CloudinaryAPIKey(), common.CloudinaryAPISecret())
	if err != nil {
		return "", err
	}

	//upload file
	uploadParam, err := cld.Upload.Upload(ctx, input, uploader.UploadParams{Folder: common.CloudinaryUploadFolder()})
	if err != nil {
		return "", err
	}
	return uploadParam.SecureURL, nil
}
