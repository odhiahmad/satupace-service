package helper

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"

	_ "github.com/lib/pq"
)

var client *whatsmeow.Client

// InitWhatsApp menginisialisasi client WhatsApp
func InitWhatsApp() {
	dbURI := fmt.Sprintf(
		"user=%s password=%s host=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	container, err := sqlstore.New(
		context.Background(),
		"postgres",
		dbURI,
		waLog.Stdout("WhatsApp", "INFO", true),
	)
	if err != nil {
		log.Fatalf("❌ Gagal buat store WhatsApp: %v", err)
	}

	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		log.Fatalf("❌ Gagal ambil device WhatsApp: %v", err)
	}

	client = whatsmeow.NewClient(deviceStore, waLog.Stdout("Client", "INFO", true))

	// Tambahkan event handler
	client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Connected:
			log.Println("✅ Terhubung ke WhatsApp sebagai:", client.Store.ID.User)
		case *events.Disconnected:
			log.Println("❌ Terputus dari WhatsApp")
		case *events.LoggedOut:
			log.Println("🔓 Session logout, kamu harus scan ulang QR")
		case *events.PairSuccess:
			log.Println("✅ Pairing sukses:", v.ID.User)
		default:
			// Event lainnya bisa diabaikan atau dilogging
			// fmt.Printf("📥 Event lain: %T\n", v)
		}
	})

	// Jika belum login (session kosong)
	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			log.Fatalf("❌ Gagal connect WhatsApp: %v", err)
		}

		for evt := range qrChan {
			switch evt.Event {
			case "code":
				fmt.Println("📲 Scan QR Code berikut:")
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			case "success":
				fmt.Println("✅ QR Code discan, login sukses")
			case "timeout", "error":
				fmt.Println("❌ Gagal login:", evt.Event)
				return
			}
		}
	} else {
		// Sudah login sebelumnya
		err = client.Connect()
		if err != nil {
			log.Fatalf("❌ Gagal reconnect WhatsApp: %v", err)
		}
		log.Println("✅ WhatsApp client reconnected:", client.Store.ID.User)
	}
}

// SendOTPViaWhatsApp mengirim pesan OTP ke nomor WhatsApp
func SendOTPViaWhatsApp(phone string, message string) error {
	if client == nil || !client.IsConnected() {
		return fmt.Errorf("❌ WhatsApp client belum aktif")
	}

	jid := types.NewJID(formatPhoneNumber(phone), types.DefaultUserServer)

	msg := &waProto.Message{
		Conversation: proto.String(message),
	}

	_, err := client.SendMessage(context.Background(), jid, msg)
	if err != nil {
		log.Printf("❌ Gagal kirim pesan ke %s: %v", phone, err)
	}
	return err
}

// formatPhoneNumber membersihkan dan mengubah ke format JID
func formatPhoneNumber(phone string) string {
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.TrimPrefix(phone, "+")
	if strings.HasPrefix(phone, "08") {
		phone = "62" + phone[1:]
	}
	if !strings.HasPrefix(phone, "62") {
		log.Printf("⚠️ Nomor bukan format Indonesia: %s", phone)
	}
	return phone
}
