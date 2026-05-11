<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

/**
 * Middleware untuk memverifikasi webhook dari bot Go.
 * Bot Go wajib mengirim header X-Api-Key yang cocok dengan
 * config('chatbot.webhook_key').
 */
class VerifyWebhookKey
{
    public function handle(Request $request, Closure $next): Response
    {
        $expected = config('chatbot.webhook_key');

        if (! $expected) {
            abort(500, 'CHATBOT_WEBHOOK_KEY belum dikonfigurasi di .env');
        }

        $provided = $request->header('X-Api-Key');

        if (! hash_equals((string) $expected, (string) $provided)) {
            abort(401, 'Invalid API key');
        }

        return $next($request);
    }
}
