<?php

use App\Http\Controllers\Api\DashboardController;
use App\Http\Controllers\Api\SendController;
use App\Http\Controllers\Api\WebhookController;
use Illuminate\Support\Facades\Route;

/*
|--------------------------------------------------------------------------
| API Routes - KIMPO Dashboard
|--------------------------------------------------------------------------
*/

// Dipanggil oleh bot Go. Dilindungi header X-Api-Key (lihat VerifyWebhookKey).
Route::middleware('webhook.key')->group(function () {
    Route::post('/webhook/message', WebhookController::class);
});

// Dipanggil oleh Vue SPA. Dilindungi Sanctum (user admin login).
Route::middleware('auth:sanctum')->prefix('dashboard')->group(function () {
    Route::get('/stats',    [DashboardController::class, 'stats']);
    Route::get('/messages', [DashboardController::class, 'messages']);
    Route::get('/qr',       [DashboardController::class, 'qr']);
    Route::post('/send',    SendController::class);
});

// Endpoint user terdaftar (dipakai frontend untuk cek siapa yang login).
Route::middleware('auth:sanctum')->get('/user', fn (\Illuminate\Http\Request $r) => $r->user());
