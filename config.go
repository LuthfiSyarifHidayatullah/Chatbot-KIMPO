package main

import (
	"flag"
	"os"
	"strings"
)

// Config menyimpan seluruh konfigurasi runtime bot.
// Semua nilai bisa di-override via environment variable atau flag CLI.
type Config struct {
	// DashboardAddr alamat HTTP untuk API dashboard Go (mis. ":8080").
	DashboardAddr string

	// ReplyAPIURL endpoint Laravel lama yang memberi balasan chatbot.
	ReplyAPIURL string

	// LaravelWebhookURL endpoint Laravel untuk menerima event dari Go.
	// Kosongkan jika belum mau dipakai (mis. saat pakai dashboard HTML lama).
	LaravelWebhookURL string

	// WebhookAPIKey dikirim di header X-Api-Key saat push webhook ke Laravel.
	WebhookAPIKey string

	// InboundAPIKey diminta dari Laravel/Vue saat memanggil endpoint Go
	// (/api/send, /api/qr, dll). Kosongkan untuk mode dev tanpa auth.
	InboundAPIKey string

	// AllowedOrigins daftar origin yang diizinkan CORS (comma separated).
	// Contoh: "http://localhost:5173,https://dashboard.example.com"
	AllowedOrigins []string
}

var config Config

func loadConfig() {
	flag.StringVar(&config.DashboardAddr, "dashboard", getEnv("DASHBOARD_ADDR", ":8080"),
		"alamat HTTP dashboard (mis. :8080)")
	flag.StringVar(&config.ReplyAPIURL, "reply-url", getEnv("REPLY_API_URL", "http://36.67.17.105:8000/api/chatbot"),
		"URL API Laravel untuk generate balasan chatbot")
	flag.StringVar(&config.LaravelWebhookURL, "webhook-url", getEnv("LARAVEL_WEBHOOK_URL", ""),
		"URL webhook Laravel untuk menerima event bot")
	flag.StringVar(&config.WebhookAPIKey, "webhook-key", getEnv("WEBHOOK_API_KEY", ""),
		"API key untuk dikirim ke Laravel (header X-Api-Key)")
	flag.StringVar(&config.InboundAPIKey, "inbound-key", getEnv("INBOUND_API_KEY", ""),
		"API key yang diminta dari caller untuk API Go")

	origins := getEnv("ALLOWED_ORIGINS", "http://localhost:5173")
	flag.Func("cors", "daftar origin yang diizinkan (comma separated)", func(v string) error {
		origins = v
		return nil
	})

	flag.Parse()

	config.AllowedOrigins = splitCSV(origins)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
