package main

import (
	"context"
	"fmt"
	"log"

	waLog "go.mau.fi/whatsmeow/util/log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/skip2/go-qrcode"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"

	waProto "go.mau.fi/whatsmeow/binary/proto"
)

var client *whatsmeow.Client

func main() {
	ctx := context.Background()

	// Database
	dbLog := waLog.Stdout("Database", "INFO", true)
	container, err := sqlstore.New(ctx, "sqlite3", "file:store.db?_foreign_keys=on", dbLog)
	if err != nil {
		log.Fatal(err)
	}

	// Device
	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Client
	clientLog := waLog.Stdout("Client", "INFO", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)

	client.AddEventHandler(eventHandler)

	// Login / QR
	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(ctx)
		err = client.Connect()
		if err != nil {
			log.Fatal(err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				fmt.Println("Scan QR ini di WhatsApp:\n")

				qr, _ := qrcode.New(evt.Code, qrcode.Medium)
				fmt.Println(qr.ToString(false))
			}
		}
	} else {
		err = client.Connect()
		if err != nil {
			log.Fatal(err)
		}
	}

	// Keep running
	select {}
}

// =========================
// EVENT HANDLER (MENU INTERAKTIF)
// =========================
func eventHandler(evt interface{}) {
	switch v := evt.(type) {

	case *events.Message:
		msg := ""

		// Ambil pesan
		if v.Message.Conversation != nil {
			msg = *v.Message.Conversation
		}

		fmt.Println("Pesan masuk:", msg)

		var reply string

		// =========================
		// LOGIC MENU
		// =========================
		switch msg {

		case "1":
			reply = "📄 *Informasi Layanan*\n\n" +
				"1. Fasilitasi Zoom meeting\n" +
				"2. Fasilitasi Dokumentasi\n" +
				"3. Pembuatan Email Pemkab \n" +
				"4. Pembuatan TTE \n\n" +
				"Ketik 0 untuk kembali ke menu"

		case "2":
			reply = "📑 *Syarat Administrasi Fasilitasi*\n\n" +
				"- Surat Permintaan Fasilitasi \n" +
				"- \n\n" +
				"Ketik 0 untuk kembali ke menu"

		case "3":
			reply = "📞 *Kontak Layanan *\n\n" +
				"Telp: 0812-xxxx-xxxx\n" +
				"Alamat: Kantor Diskominfo Kab.Bengkayang \n\n" +
				"Ketik 0 untuk kembali ke menu"
			
		case "4":
			reply = "🌐 *Website Bengkayang *\n\n" +
				"https://bengkayangkab.go.id/"

		case "0":
			reply = mainMenu()

		default:
			reply = mainMenu()
		}

		client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
			Conversation: &reply,
		})
	}
}

// =========================
// MENU UTAMA
// =========================
func mainMenu() string {
	return "Halo 👋 Selamat datang \n\n" +
		"Silakan pilih layanan:\n" +
		"1. Informasi Layanan\n" +
		"2. Syarat Administrasi\n" +
		"3. Kontak\n\n" +
		"4. Website Bengkayang\n\n" +
		"Ketik angka untuk memilih."
}
