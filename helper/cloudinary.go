package helper

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadBase64ToCloudinary(base64Str string, folder string) (string, error) {
	if base64Str == "" {
		return "", nil
	}

	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLD_NAME"),
		os.Getenv("CLD_API_KEY"),
		os.Getenv("CLD_API_SECRET"),
	)
	if err != nil {
		return "", fmt.Errorf("init cloudinary gagal: %w", err)
	}

	ctx := context.Background()
	uploadResp, err := cld.Upload.Upload(ctx, base64Str, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		return "", fmt.Errorf("upload cloudinary gagal: %w", err)
	}

	log.Printf("✅ Upload base64 sukses: %s", uploadResp.SecureURL)
	return uploadResp.SecureURL, nil
}
