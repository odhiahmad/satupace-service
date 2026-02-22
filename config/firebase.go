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
	credPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	if credPath == "" {
		// default: cari di root project
		credPath = "./run-sync-a470f-firebase-adminsdk-fbsvc-6dedfbd10f.json"
	}

	opt := option.WithCredentialsFile(credPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("❌ Gagal inisialisasi Firebase: %v", err)
	}

	FCMClient, err = app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("❌ Gagal mendapatkan FCM client: %v", err)
	}

	log.Println("✅ Firebase FCM berhasil diinisialisasi")
}
