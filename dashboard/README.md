# KIMPO Dashboard (Laravel 12 + Vue 3)

Dashboard interaktif untuk chatbot WhatsApp KIMPO. Arsitektur **hybrid**:
data historis disimpan di Laravel, event realtime dikirim via **Reverb + Laravel Echo**,
dan frontend Vue 3 bisa kirim pesan manual lewat Laravel yang mem-proxy ke bot Go.

## Arsitektur

```
  [Pesan WA masuk]
         │
         ▼
  [Go eventHandler] ──reply ke WA──► user
         │
         └──POST webhook──►  [Laravel WebhookController]
                                     │
                                     ├──simpan MySQL (chatbot_messages)
                                     └──broadcast via Reverb
                                              │
                                              ▼
                                    [Vue via Laravel Echo]
                                    (update realtime tanpa polling)

  [Vue: kirim pesan] ──► [Laravel SendController] ──► [Go /api/send] ──► WhatsApp
```

Tiga komponen berjalan independen, bisa di-deploy di server berbeda:

| Komponen | Folder | Port default | Catatan |
|---|---|---|---|
| Go bot   | `../` (root repo) | `:8080` | Menangani WA + API Go internal |
| Laravel  | `backend/`        | `:8000` | REST API + Reverb WS |
| Vue      | `frontend/`       | `:5173` | SPA admin |
| Reverb   | dalam proses Laravel | `:8080` (Reverb)* | WebSocket broadcasting |

\* Reverb default juga `:8080`. **Ganti salah satu** agar tidak bentrok dengan Go bot —
misalnya pakai `REVERB_PORT=8081` di `.env`.

---

## 1. Setup Go Bot

Di root repo `Chatbot-KIMPO/`:

```bash
# Salin env atau pakai flag CLI. Lihat config.go untuk daftar lengkap.
export DASHBOARD_ADDR=":8080"
export LARAVEL_WEBHOOK_URL="http://localhost:8000/api/webhook/message"
export WEBHOOK_API_KEY="ganti_dengan_key_rahasia_1"
export INBOUND_API_KEY="ganti_dengan_key_rahasia_2"
export ALLOWED_ORIGINS="http://localhost:5173"

go run .
```

`WEBHOOK_API_KEY` harus sama dengan `CHATBOT_WEBHOOK_KEY` di Laravel.
`INBOUND_API_KEY` harus sama dengan `CHATBOT_GO_KEY` di Laravel.

## 2. Setup Laravel Backend

```bash
# Buat proyek Laravel 12 baru, lalu copy isi backend/ ke dalamnya:
composer create-project laravel/laravel dashboard-kimpo "^12.0"
cp -r backend/app dashboard-kimpo/
cp -r backend/config dashboard-kimpo/
cp -r backend/database dashboard-kimpo/
cp -r backend/routes dashboard-kimpo/
cp backend/.env.example dashboard-kimpo/.env

cd dashboard-kimpo

# Install Sanctum + Reverb
php artisan install:api
composer require laravel/reverb
php artisan reverb:install

# Merge manual: lihat backend/bootstrap/app.php.snippet dan salin alias
# middleware 'webhook.key' + ->statefulApi() ke bootstrap/app.php.

php artisan key:generate
php artisan migrate
php artisan db:seed --class=AdminUserSeeder   # buat user admin pertama (opsional)

# Jalankan 3 proses paralel (pakai tmux atau 3 terminal):
php artisan serve                 # :8000  (HTTP API)
php artisan reverb:start          # :8080  (WebSocket) -- GANTI PORT jika bentrok Go
php artisan queue:work            # untuk proses broadcast yang implements ShouldBroadcast
```

## 3. Setup Vue Frontend

```bash
cd frontend
cp .env.example .env
# Edit .env: pastikan VITE_API_BASE_URL dan VITE_REVERB_* cocok dengan Laravel

npm install
npm run dev      # buka http://localhost:5173
```

Login dengan akun admin yang dibuat seeder. Setelah login, dashboard akan otomatis
terhubung ke private channel `dashboard` untuk update realtime.

---

## Alur Realtime (Reverb + Echo)

1. Go bot menerima pesan WA → panggil `sendWebhook("message", ...)`.
2. Laravel `WebhookController` menyimpan ke `chatbot_messages`, lalu
   `broadcast(new MessageReceived($msg))->toOthers()`.
3. Job broadcast dieksekusi `queue:work` → dikirim ke Reverb → ke semua client Echo.
4. Pinia store `dashboard.js` mendengarkan `.message.received` dan
   langsung menambahkan pesan baru ke UI tanpa refresh.

Event lain: `.connection.changed` (status online/offline), `.qr.updated` (QR baru).

## Keamanan

- **Go → Laravel** (webhook): header `X-Api-Key` harus cocok dengan `CHATBOT_WEBHOOK_KEY`.
  Gunakan `hash_equals()` (sudah di `VerifyWebhookKey` middleware).
- **Laravel → Go** (proxy send): header `X-Api-Key` dengan `CHATBOT_GO_KEY`.
  Go memeriksa di middleware `withAPIMiddleware` → `InboundAPIKey`.
- **Vue → Laravel**: Sanctum SPA (cookie + XSRF). Set `SANCTUM_STATEFUL_DOMAINS` ke
  origin Vue.
- **Private channel `dashboard`**: hanya user login yang bisa subscribe
  (lihat `routes/channels.php`).

## Deployment (ringkas)

- Go bot: build binary (`go build -o kimpo-bot`) dan jalankan via `systemd` atau pm2.
- Laravel: deploy seperti biasa, pastikan `queue:work` dan `reverb:start` di-supervisor-kan.
- Vue: `npm run build` → serve folder `dist/` dengan nginx atau serve via Laravel juga bisa.

Bila Laravel + Vue di-host di domain berbeda dari Go bot, pastikan:
- `ALLOWED_ORIGINS` di Go memuat origin Vue **dan** origin Laravel.
- HTTPS untuk Reverb (set `REVERB_SCHEME=https`, gunakan reverse proxy).


---

## Struktur File

```
dashboard/
├── backend/                                 # Laravel 12 (drop-in ke project baru)
│   ├── app/
│   │   ├── Events/
│   │   │   ├── MessageReceived.php          # broadcast saat pesan baru
│   │   │   ├── ConnectionChanged.php        # broadcast saat status koneksi berubah
│   │   │   └── QrUpdated.php                # broadcast saat QR code baru
│   │   ├── Http/
│   │   │   ├── Controllers/Api/
│   │   │   │   ├── WebhookController.php    # terima event dari Go (X-Api-Key)
│   │   │   │   ├── DashboardController.php  # stats, messages, qr (Sanctum)
│   │   │   │   └── SendController.php       # proxy kirim pesan → Go
│   │   │   └── Middleware/VerifyWebhookKey.php
│   │   └── Models/{ChatbotMessage,ChatbotState}.php
│   ├── bootstrap/app.php.snippet            # PANDUAN merge ke bootstrap/app.php
│   ├── config/chatbot.php
│   ├── database/migrations/..._create_chatbot_messages_table.php
│   ├── routes/{api,channels}.php
│   └── .env.example
│
└── frontend/                                # Vue 3 + Vite
    ├── package.json                         # deps: pinia, axios, chart.js, laravel-echo
    ├── vite.config.js
    ├── index.html
    ├── .env.example
    └── src/
        ├── main.js
        ├── App.vue
        ├── router/index.js                  # guard: redirect ke /login kalau belum auth
        ├── services/
        │   ├── api.js                       # axios + Sanctum SPA cookie
        │   └── echo.js                      # Laravel Echo + Reverb
        ├── stores/
        │   ├── dashboard.js                 # state + Echo subscriptions
        │   └── auth.js                      # login/logout/user
        ├── views/
        │   ├── DashboardView.vue
        │   └── LoginView.vue
        └── components/
            ├── StatsCards.vue
            ├── DailyChart.vue               # line chart 7 hari
            ├── HourlyChart.vue              # bar chart 24 jam
            ├── QRPanel.vue
            ├── MessageLog.vue
            ├── SendForm.vue
            └── TopUsers.vue
```

## Troubleshooting

**Port bentrok 8080.** Go bot default `:8080`, Reverb juga default `:8080`.
Set `REVERB_PORT=8081` di `.env` Laravel dan `VITE_REVERB_PORT=8081` di Vue.

**401 saat Vue panggil API.** Pastikan `SANCTUM_STATEFUL_DOMAINS` di Laravel
memuat `localhost:5173`, dan panggil `csrf()` sebelum login (sudah otomatis di `auth.js`).

**Event broadcast tidak sampai ke Vue.** Cek:
1. `BROADCAST_CONNECTION=reverb` di `.env` Laravel
2. `php artisan queue:work` sedang berjalan (event `ShouldBroadcast` masuk queue)
3. `php artisan reverb:start` sedang berjalan
4. Browser DevTools → Network → WS: koneksi `reverb` ter-establish & auth 200

**Webhook Go gagal (401).** Samakan `WEBHOOK_API_KEY` (Go) dengan
`CHATBOT_WEBHOOK_KEY` (Laravel). Cek di log Laravel (`storage/logs/laravel.log`).

**CORS error.** Tambahkan origin Vue ke `ALLOWED_ORIGINS` Go bot, dan origin Vue
ke `config/cors.php` Laravel (`supports_credentials: true`).
