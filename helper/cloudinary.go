package helper

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadBase64ToCloudinary(base64Str string, folder string) (string, error) {
	if base64Str == "" {
		return "", nil
	}

	// Cloudinary SDK requires data URI format: "data:image/jpeg;base64,<data>"
	// Flutter sends raw base64 without this prefix, so we add it here.
	if !strings.HasPrefix(base64Str, "data:") {
		base64Str = "data:image/jpeg;base64," + base64Str
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

	log.Printf("✅ Upload Cloudinary berhasil: %s", uploadResp.SecureURL)
	return uploadResp.SecureURL, nil
}
