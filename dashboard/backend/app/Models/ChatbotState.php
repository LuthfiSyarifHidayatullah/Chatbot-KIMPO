<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class ChatbotState extends Model
{
    protected $fillable = [
        'connected',
        'qr_b64',
        'last_event_at',
    ];

    protected $casts = [
        'connected'     => 'boolean',
        'last_event_at' => 'datetime',
    ];

    /**
     * Ambil baris state global (selalu id=1). Buat baru jika belum ada.
     */
    public static function current(): self
    {
        return static::firstOrCreate(['id' => 1], [
            'connected' => false,
        ]);
    }
}
