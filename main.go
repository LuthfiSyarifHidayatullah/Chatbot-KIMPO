package main

import (
	"bytes"
	"context"
	"encoding/json"
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
	ctx := context.Background()

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

	select {}
}

// =========================
// EVENT HANDLER (AMBIL DARI API)
// =========================
func eventHandler(evt interface{}) {
	switch v := evt.(type) {

	case *events.Message:
		msg := ""

		if v.Message.Conversation != nil {
			msg = *v.Message.Conversation
		}

		fmt.Println("Pesan masuk:", msg)

		// 🔥 Ambil jawaban dari Laravel
		reply := getReplyFromAPI(msg)

		client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
			Conversation: &reply,
		})
	}
}

// =========================
// FUNCTION API LARAVEL
// =========================
func getReplyFromAPI(message string) string {
	url := "http://36.67.17.105:8000/api/chatbot" // ⚠️ GANTI jika beda

	payload := map[string]string{
		"message": message,
	}

	jsonData, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return "Server sedang bermasalah"
	}

	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	return result["reply"]
}
