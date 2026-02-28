package helper

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

// faceServiceDetectResponse adalah response dari Face Service Detection API
type faceServiceDetectResponse struct {
	Result []struct {
		Pose *struct {
			Pitch float32 `json:"pitch"`
			Roll  float32 `json:"roll"`
			Yaw   float32 `json:"yaw"`
		} `json:"pose"`
		Box struct {
			Probability float32 `json:"probability"`
		} `json:"box"`
	} `json:"result"`
}

// faceServiceVerifyResponse adalah response dari Face Service Verification API
type faceServiceVerifyResponse struct {
	Result []struct {
		FaceMatches []struct {
			Similarity float32 `json:"similarity"`
		} `json:"face_matches"`
	} `json:"result"`
}

func faceServiceBaseURL() string {
	return strings.TrimRight(os.Getenv("FACE_SERVICE_URL"), "/")
}

func decodeBase64ToBytes(b64 string) ([]byte, error) {
	// Hapus data URI prefix jika ada (misal "data:image/jpeg;base64,...")
	if idx := strings.Index(b64, ","); idx != -1 {
		b64 = b64[idx+1:]
	}
	return base64.StdEncoding.DecodeString(b64)
}

func downloadImageBytes(url string) ([]byte, error) {
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("gagal download gambar: %w", err)
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// DetectFrontFace memvalidasi bahwa gambar base64 mengandung wajah tampak depan.
// Memanggil Face Service Detection API.
func DetectFrontFace(base64Str string) error {
	apiKey := os.Getenv("FACE_SERVICE_API_KEY")
	if apiKey == "" {
		// Face Service not configured — skip face validation
		return nil
	}

	imageBytes, err := decodeBase64ToBytes(base64Str)
	if err != nil {
		return fmt.Errorf("gagal decode gambar: %w", err)
	}

	// Buat multipart body
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "photo.jpg")
	if err != nil {
		return fmt.Errorf("gagal membuat form: %w", err)
	}
	if _, err = part.Write(imageBytes); err != nil {
		return fmt.Errorf("gagal menulis gambar: %w", err)
	}
	writer.Close()

	url := faceServiceBaseURL() + "/api/v1/detection/detect"
	req, err := http.NewRequest(http.MethodPost, url, &body)
	if err != nil {
		return fmt.Errorf("gagal membuat request: %w", err)
	}
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("gagal menghubungi Face Service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("face service detection error (%d): %s", resp.StatusCode, string(raw))
	}

	var result faceServiceDetectResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("gagal parse response Face Service: %w", err)
	}

	if len(result.Result) == 0 {
		return fmt.Errorf("tidak ada wajah yang terdeteksi pada foto")
	}

	// Validasi pose: wajah harus tampak depan
	face := result.Result[0]
	if face.Pose != nil {
		yaw := face.Pose.Yaw
		pitch := face.Pose.Pitch
		log.Printf("[FaceService] pose — yaw=%.2f pitch=%.2f", yaw, pitch)
		if yaw < -25 || yaw > 25 || pitch < -25 || pitch > 25 {
			return fmt.Errorf("wajah harus menghadap lurus ke depan (tampak depan), hindari sudut miring")
		}
	}

	return nil
}

// VerifyFaces membandingkan gambar kamera (base64) dengan foto tersimpan (URL Cloudinary).
// Mengembalikan similarity (0.0–1.0) dan apakah wajah cocok (threshold 0.80).
func VerifyFaces(cameraBase64 string, storedPhotoURL string) (similarity float32, matched bool, err error) {
	apiKey := os.Getenv("FACE_SERVICE_API_KEY")
	if apiKey == "" {
		// Face Service not configured — skip face verification
		return 0, false, nil
	}

	sourceBytes, err := decodeBase64ToBytes(cameraBase64)
	if err != nil {
		return 0, false, fmt.Errorf("gagal decode gambar kamera: %w", err)
	}

	targetBytes, err := downloadImageBytes(storedPhotoURL)
	if err != nil {
		return 0, false, fmt.Errorf("gagal mengambil foto verifikasi: %w", err)
	}

	// Buat multipart body dengan source_image dan target_image
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	srcPart, err := writer.CreateFormFile("source_image", "source.jpg")
	if err != nil {
		return 0, false, fmt.Errorf("gagal membuat source form: %w", err)
	}
	if _, err = srcPart.Write(sourceBytes); err != nil {
		return 0, false, fmt.Errorf("gagal menulis source image: %w", err)
	}

	tgtPart, err := writer.CreateFormFile("target_image", "target.jpg")
	if err != nil {
		return 0, false, fmt.Errorf("gagal membuat target form: %w", err)
	}
	if _, err = tgtPart.Write(targetBytes); err != nil {
		return 0, false, fmt.Errorf("gagal menulis target image: %w", err)
	}
	writer.Close()

	url := faceServiceBaseURL() + "/api/v1/verification/verify"
	req, err := http.NewRequest(http.MethodPost, url, &body)
	if err != nil {
		return 0, false, fmt.Errorf("gagal membuat request: %w", err)
	}
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, false, fmt.Errorf("gagal menghubungi Face Service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return 0, false, fmt.Errorf("face service verification error (%d): %s", resp.StatusCode, string(raw))
	}

	var result faceServiceVerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, false, fmt.Errorf("gagal parse response Face Service: %w", err)
	}

	if len(result.Result) == 0 || len(result.Result[0].FaceMatches) == 0 {
		return 0, false, nil
	}

	sim := result.Result[0].FaceMatches[0].Similarity
	const threshold float32 = 0.80
	return sim, sim >= threshold, nil
}
