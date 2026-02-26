package helper

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

// compreFaceDetectResponse adalah response dari CompreFace Detection API
type compreFaceDetectResponse struct {
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

// compreFaceVerifyResponse adalah response dari CompreFace Verification API
type compreFaceVerifyResponse struct {
	Result []struct {
		FaceMatches []struct {
			Similarity float32 `json:"similarity"`
		} `json:"face_matches"`
	} `json:"result"`
}

func compreFaceBaseURL() string {
	return strings.TrimRight(os.Getenv("COMPREFACE_URL"), "/")
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
// Memanggil CompreFace Detection API.
func DetectFrontFace(base64Str string) error {
	apiKey := os.Getenv("COMPREFACE_DETECTION_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("COMPREFACE_DETECTION_API_KEY belum dikonfigurasi")
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

	url := compreFaceBaseURL() + "/api/v1/detection/detect"
	req, err := http.NewRequest(http.MethodPost, url, &body)
	if err != nil {
		return fmt.Errorf("gagal membuat request: %w", err)
	}
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("gagal menghubungi CompreFace: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("CompreFace detection error (%d): %s", resp.StatusCode, string(raw))
	}

	var result compreFaceDetectResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("gagal parse response CompreFace: %w", err)
	}

	if len(result.Result) == 0 {
		return fmt.Errorf("tidak ada wajah yang terdeteksi pada foto")
	}

	// Validasi pose: wajah harus tampak depan
	face := result.Result[0]
	if face.Pose != nil {
		yaw := face.Pose.Yaw
		pitch := face.Pose.Pitch
		if yaw < -25 || yaw > 25 || pitch < -25 || pitch > 25 {
			return fmt.Errorf("wajah harus menghadap lurus ke depan (tampak depan), hindari sudut miring")
		}
	}

	return nil
}

// VerifyFaces membandingkan gambar kamera (base64) dengan foto tersimpan (URL Cloudinary).
// Mengembalikan similarity (0.0–1.0) dan apakah wajah cocok (threshold 0.80).
func VerifyFaces(cameraBase64 string, storedPhotoURL string) (similarity float32, matched bool, err error) {
	apiKey := os.Getenv("COMPREFACE_VERIFICATION_API_KEY")
	if apiKey == "" {
		return 0, false, fmt.Errorf("COMPREFACE_VERIFICATION_API_KEY belum dikonfigurasi")
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

	url := compreFaceBaseURL() + "/api/v1/verification/verify"
	req, err := http.NewRequest(http.MethodPost, url, &body)
	if err != nil {
		return 0, false, fmt.Errorf("gagal membuat request: %w", err)
	}
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, false, fmt.Errorf("gagal menghubungi CompreFace: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return 0, false, fmt.Errorf("CompreFace verification error (%d): %s", resp.StatusCode, string(raw))
	}

	var result compreFaceVerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, false, fmt.Errorf("gagal parse response CompreFace: %w", err)
	}

	if len(result.Result) == 0 || len(result.Result[0].FaceMatches) == 0 {
		return 0, false, nil
	}

	sim := result.Result[0].FaceMatches[0].Similarity
	const threshold float32 = 0.80
	return sim, sim >= threshold, nil
}
