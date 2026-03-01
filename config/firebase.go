package config

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

var FCMClient *messaging.Client

func SetupFirebase() {
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	conf := &firebase.Config{ProjectID: projectID}

	var opt option.ClientOption

	// Prioritas 1: JSON langsung dari env (untuk production/server)
	if credJSON := os.Getenv("FIREBASE_CREDENTIALS_JSON"); credJSON != "" {
		opt = option.WithCredentialsJSON([]byte(credJSON))
	} else {
		// Prioritas 2: path ke file (untuk local development)
		credPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
		if credPath == "" {
			credPath = "./run-sync-firebase.json"
		}
		opt = option.WithCredentialsFile(credPath)
	}

	app, err := firebase.NewApp(context.Background(), conf, opt)
	if err != nil {
		log.Fatalf("❌ Gagal inisialisasi Firebase: %v", err)
	}

	FCMClient, err = app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("❌ Gagal mendapatkan FCM client: %v", err)
	}

	log.Println("✅ Firebase FCM berhasil diinisialisasi")
}
