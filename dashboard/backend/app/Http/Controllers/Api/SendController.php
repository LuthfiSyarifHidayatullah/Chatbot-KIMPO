<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Http;

/**
 * Proxy untuk kirim pesan manual dari Vue dashboard ke bot Go.
 * Alur: Vue -> Laravel (dengan Sanctum auth) -> Go /api/send -> WhatsApp.
 */
class SendController extends Controller
{
    public function __invoke(Request $request): JsonResponse
    {
        $data = $request->validate([
            'to'      => 'required|string|max:120',
            'message' => 'required|string|max:4096',
        ]);

        $response = Http::withHeaders([
                'X-Api-Key' => config('chatbot.go_key'),
            ])
            ->timeout(10)
            ->acceptJson()
            ->post(config('chatbot.go_url').'/api/send', $data);

        if ($response->failed()) {
            return response()->json([
                'ok'    => false,
                'error' => $response->body() ?: 'Gagal menghubungi bot Go',
            ], $response->status() ?: 502);
        }

        return response()->json(['ok' => true, 'result' => $response->json()]);
    }
}
