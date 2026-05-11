package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

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
	dashboardAddr := flag.String("dashboard", ":8080", "alamat HTTP dashboard (mis. :8080)")
	flag.Parse()

	ctx := context.Background()

	// Jalankan dashboard web terlebih dulu supaya QR bisa langsung muncul
	startDashboard(*dashboardAddr)

	dbLog := waLog.Stdout("Database", "INFO", true)
	container, err := sqlstore.New(ctx, "sqlite3", "file:store.db?_foreign_keys=on", dbLog)
	if err != nil {
		log.Fatal(err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		log.Fatal(err)
	}

	clientLog := waLog.Stdout("Client", "INFO", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)

	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(ctx)
		err = client.Connect()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			for evt := range qrChan {
				switch evt.Event {
				case "code":
					fmt.Println("Scan QR ini di WhatsApp (atau buka dashboard):")
					qr, _ := qrcode.New(evt.Code, qrcode.Medium)
					fmt.Println(qr.ToString(false))

					if b64, err := generateQRBase64(evt.Code); err == nil {
						metrics.setQR(evt.Code, b64)
					}
				case "success":
					fmt.Println("QR berhasil dipindai, login sukses.")
					metrics.setConnected(true)
				case "timeout":
					fmt.Println("QR timeout.")
				}
			}
		}()
	} else {
		err = client.Connect()
		if err != nil {
			log.Fatal(err)
		}
		metrics.setConnected(true)
	}

	select {}
}

// =========================
// EVENT HANDLER
// =========================
func eventHandler(evt interface{}) {
	switch v := evt.(type) {

	case *events.Connected:
		fmt.Println("WhatsApp terhubung.")
		metrics.setConnected(true)

	case *events.Disconnected:
		fmt.Println("WhatsApp terputus.")
		metrics.setConnected(false)

	case *events.LoggedOut:
		fmt.Println("WhatsApp logged out.")
		metrics.setConnected(false)

	case *events.Message:
		msg := ""
		if v.Message.Conversation != nil {
			msg = *v.Message.Conversation
		}

		fmt.Println("Pesan masuk:", msg)

		// Ambil jawaban dari API
		reply := getReplyFromAPI(msg)

		_, err := client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
			Conversation: &reply,
		})
		if err != nil {
			metrics.incError()
			fmt.Println("Gagal kirim balasan:", err)
		}

		// Catat ke dashboard
		metrics.logIncoming(v.Info.Chat.String(), msg, reply)
	}
}

// =========================
// FUNCTION API LARAVEL
// =========================
func getReplyFromAPI(message string) string {
	url := "http://36.67.17.105:8000/api/chatbot" // GANTI jika beda

	payload := map[string]string{
		"message": message,
	}

	jsonData, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		metrics.incError()
		return "Server sedang bermasalah"
	}
	defer resp.Body.Close()

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		metrics.incError()
		return "Server sedang bermasalah"
	}

	return result["reply"]
}
