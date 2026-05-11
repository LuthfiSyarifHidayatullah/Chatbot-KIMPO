package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/skip2/go-qrcode"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

// =========================
// METRICS STORE
// =========================

type MessageLog struct {
	Time      time.Time `json:"time"`
	From      string    `json:"from"`
	Direction string    `json:"direction"` // "in" or "out"
	Message   string    `json:"message"`
	Reply     string    `json:"reply,omitempty"`
}

type DailyStat struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type Metrics struct {
	mu sync.RWMutex

	TotalMessages    int
	IncomingMessages int
	OutgoingMessages int
	ErrorCount       int

	UniqueUsers map[string]int // jid -> message count
	HourlyStats [24]int        // messages per hour (last 24h window, index = hour of day)
	DailyStats  map[string]int // "YYYY-MM-DD" -> count

	RecentMessages []MessageLog

	Connected  bool
	QRCodeB64  string // base64 PNG of current QR (empty when connected)
	LastQRText string

	StartedAt time.Time
}

var metrics = &Metrics{
	UniqueUsers:    make(map[string]int),
	DailyStats:     make(map[string]int),
	RecentMessages: []MessageLog{},
	StartedAt:      time.Now(),
}

func (m *Metrics) logIncoming(from, message, reply string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	m.TotalMessages++
	m.IncomingMessages++
	m.OutgoingMessages++ // reply is sent right after
	m.UniqueUsers[from]++
	m.HourlyStats[now.Hour()]++
	m.DailyStats[now.Format("2006-01-02")]++

	entry := MessageLog{
		Time:      now,
		From:      from,
		Direction: "in",
		Message:   message,
		Reply:     reply,
	}
	m.RecentMessages = append([]MessageLog{entry}, m.RecentMessages...)
	if len(m.RecentMessages) > 100 {
		m.RecentMessages = m.RecentMessages[:100]
	}
}

func (m *Metrics) logOutgoing(to, message string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.TotalMessages++
	m.OutgoingMessages++

	entry := MessageLog{
		Time:      time.Now(),
		From:      to,
		Direction: "out",
		Message:   message,
	}
	m.RecentMessages = append([]MessageLog{entry}, m.RecentMessages...)
	if len(m.RecentMessages) > 100 {
		m.RecentMessages = m.RecentMessages[:100]
	}
}

func (m *Metrics) incError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ErrorCount++
}

func (m *Metrics) setConnected(v bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Connected = v
	if v {
		m.QRCodeB64 = ""
		m.LastQRText = ""
	}
}

func (m *Metrics) setQR(code string, pngB64 string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.LastQRText = code
	m.QRCodeB64 = pngB64
}

// =========================
// HTTP SERVER
// =========================

func startDashboard(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/api/stats", handleStats)
	mux.HandleFunc("/api/messages", handleMessages)
	mux.HandleFunc("/api/qr", handleQR)
	mux.HandleFunc("/api/send", handleSend)

	fmt.Printf("\n🚀 Dashboard KIMPO aktif di: http://localhost%s\n\n", addr)

	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil {
			fmt.Println("Dashboard server error:", err)
		}
	}()
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, dashboardHTML)
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	metrics.mu.RLock()
	defer metrics.mu.RUnlock()

	// top users
	type userStat struct {
		JID   string `json:"jid"`
		Count int    `json:"count"`
	}
	topUsers := []userStat{}
	for k, v := range metrics.UniqueUsers {
		topUsers = append(topUsers, userStat{JID: k, Count: v})
	}
	sort.Slice(topUsers, func(i, j int) bool {
		return topUsers[i].Count > topUsers[j].Count
	})
	if len(topUsers) > 5 {
		topUsers = topUsers[:5]
	}

	// daily stats (last 7 days)
	daily := []DailyStat{}
	for i := 6; i >= 0; i-- {
		d := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		daily = append(daily, DailyStat{Date: d, Count: metrics.DailyStats[d]})
	}

	uptime := time.Since(metrics.StartedAt).Round(time.Second).String()

	resp := map[string]interface{}{
		"total_messages":    metrics.TotalMessages,
		"incoming_messages": metrics.IncomingMessages,
		"outgoing_messages": metrics.OutgoingMessages,
		"error_count":       metrics.ErrorCount,
		"unique_users":      len(metrics.UniqueUsers),
		"connected":         metrics.Connected,
		"hourly_stats":      metrics.HourlyStats,
		"daily_stats":       daily,
		"top_users":         topUsers,
		"uptime":            uptime,
		"started_at":        metrics.StartedAt.Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleMessages(w http.ResponseWriter, r *http.Request) {
	metrics.mu.RLock()
	defer metrics.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics.RecentMessages)
}

func handleQR(w http.ResponseWriter, r *http.Request) {
	metrics.mu.RLock()
	defer metrics.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"connected": metrics.Connected,
		"qr_b64":    metrics.QRCodeB64,
	})
}

type sendRequest struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

func handleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req sendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.To == "" || req.Message == "" {
		http.Error(w, "to and message are required", http.StatusBadRequest)
		return
	}
	if client == nil || !client.IsConnected() {
		http.Error(w, "whatsapp client not connected", http.StatusServiceUnavailable)
		return
	}

	jid, err := types.ParseJID(req.To)
	if err != nil {
		// try append server if user only supplied a phone number
		jid, err = types.ParseJID(req.To + "@s.whatsapp.net")
		if err != nil {
			http.Error(w, "invalid JID", http.StatusBadRequest)
			return
		}
	}

	_, err = client.SendMessage(context.Background(), jid, &waProto.Message{
		Conversation: &req.Message,
	})
	if err != nil {
		metrics.incError()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metrics.logOutgoing(jid.String(), req.Message)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "sent"})
}

// =========================
// QR HELPERS
// =========================

// generateQRBase64 returns a base64-encoded PNG string (without data URI prefix).
func generateQRBase64(code string) (string, error) {
	png, err := qrcode.Encode(code, qrcode.Medium, 320)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(png), nil
}
