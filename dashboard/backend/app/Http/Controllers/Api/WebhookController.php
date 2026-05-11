<?php

namespace App\Http\Controllers\Api;

use App\Events\ConnectionChanged;
use App\Events\MessageReceived;
use App\Events\QrUpdated;
use App\Http\Controllers\Controller;
use App\Models\ChatbotMessage;
use App\Models\ChatbotState;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

/**
 * Menerima push event dari bot Go KIMPO.
 *
 * Format body:
 *  {
 *    "event": "message|connection|qr",
 *    "timestamp": "2026-05-11T10:00:00Z",
 *    "payload": { ... }
 *  }
 */
class WebhookController extends Controller
{
    public function __invoke(Request $request): JsonResponse
    {
        $data = $request->validate([
            'event'     => 'required|string|in:message,connection,qr',
            'timestamp' => 'required|date',
            'payload'   => 'required|array',
        ]);

        return match ($data['event']) {
            'message'    => $this->handleMessage($data),
            'connection' => $this->handleConnection($data),
            'qr'         => $this->handleQr($data),
        };
    }

    private function handleMessage(array $data): JsonResponse
    {
        $payload = validator($data['payload'], [
            'from'      => 'required|string|max:120',
            'direction' => 'required|in:in,out',
            'message'   => 'required|string',
            'reply'     => 'nullable|string',
        ])->validate();

        $msg = ChatbotMessage::create([
            'jid'         => $payload['from'],
            'direction'   => $payload['direction'],
            'message'     => $payload['message'],
            'reply'       => $payload['reply'] ?? null,
            'occurred_at' => $data['timestamp'],
        ]);

        broadcast(new MessageReceived($msg))->toOthers();

        return response()->json(['ok' => true, 'id' => $msg->id]);
    }

    private function handleConnection(array $data): JsonResponse
    {
        $payload = validator($data['payload'], [
            'connected' => 'required|boolean',
        ])->validate();

        $state = ChatbotState::current();
        $state->update([
            'connected'     => $payload['connected'],
            'last_event_at' => now(),
            // Saat connected, QR tidak relevan lagi
            'qr_b64'        => $payload['connected'] ? null : $state->qr_b64,
        ]);

        broadcast(new ConnectionChanged($payload['connected']))->toOthers();

        return response()->json(['ok' => true]);
    }

    private function handleQr(array $data): JsonResponse
    {
        $payload = validator($data['payload'], [
            'qr_b64' => 'required|string',
        ])->validate();

        $state = ChatbotState::current();
        $state->update([
            'qr_b64'        => $payload['qr_b64'],
            'connected'     => false,
            'last_event_at' => now(),
        ]);

        broadcast(new QrUpdated($payload['qr_b64']))->toOthers();

        return response()->json(['ok' => true]);
    }
}
