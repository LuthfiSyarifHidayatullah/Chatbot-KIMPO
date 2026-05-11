package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// =========================
// OUTBOUND WEBHOOK KE LARAVEL
// =========================
//
// Setiap event penting (pesan masuk, status koneksi, QR baru) di-push ke
// endpoint Laravel supaya dashboard Vue bisa simpan ke DB dan broadcast
// via Reverb.
//
// Di sisi Laravel, endpoint harus memverifikasi header X-Api-Key.

var webhookHTTP = &http.Client{Timeout: 5 * time.Second}

type webhookEnvelope struct {
	Event     string      `json:"event"`      // "message", "connection", "qr"
	Timestamp time.Time   `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}

type messagePayload struct {
	From      string `json:"from"`
	Direction string `json:"direction"` // "in" | "out"
	Message   string `json:"message"`
	Reply     string `json:"reply,omitempty"`
}

type connectionPayload struct {
	Connected bool `json:"connected"`
}

type qrPayload struct {
	QRCodeB64 string `json:"qr_b64"`
}

// sendWebhook posts an event to the Laravel dashboard. Fails silently
// (only logs) so that bot operation is never blocked by dashboard issues.
func sendWebhook(event string, payload interface{}) {
	if config.LaravelWebhookURL == "" {
		return
	}

	go func() {
		body, err := json.Marshal(webhookEnvelope{
			Event:     event,
			Timestamp: time.Now(),
			Payload:   payload,
		})
		if err != nil {
			fmt.Println("webhook marshal error:", err)
			return
		}

		req, err := http.NewRequest(http.MethodPost, config.LaravelWebhookURL, bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("webhook request error:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Api-Key", config.WebhookAPIKey)

		resp, err := webhookHTTP.Do(req)
		if err != nil {
			fmt.Println("webhook send error:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			fmt.Printf("webhook non-2xx response: %d\n", resp.StatusCode)
		}
	}()
}
