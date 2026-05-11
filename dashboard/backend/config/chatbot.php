<?php

return [
    /*
     * Alamat dasar bot Go KIMPO. Dipakai oleh SendController & QR proxy.
     */
    'go_url' => env('CHATBOT_GO_URL', 'http://localhost:8080'),

    /*
     * API key yang harus dimasukkan bot Go di header X-Api-Key saat
     * mengirim webhook ke /api/webhook/message.
     */
    'webhook_key' => env('CHATBOT_WEBHOOK_KEY'),

    /*
     * API key yang dikirim Laravel saat memanggil endpoint Go.
     */
    'go_key' => env('CHATBOT_GO_KEY'),
];
