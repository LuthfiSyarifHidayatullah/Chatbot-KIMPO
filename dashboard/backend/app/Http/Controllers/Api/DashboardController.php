<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\ChatbotMessage;
use App\Models\ChatbotState;
use Illuminate\Http\JsonResponse;
use Illuminate\Support\Carbon;
use Illuminate\Support\Facades\DB;

class DashboardController extends Controller
{
    /**
     * GET /api/dashboard/stats
     * Agregat untuk kartu statistik & chart.
     */
    public function stats(): JsonResponse
    {
        $total     = ChatbotMessage::count();
        $incoming  = ChatbotMessage::where('direction', 'in')->count();
        $outgoing  = ChatbotMessage::where('direction', 'out')->count();
        $uniqueUsr = ChatbotMessage::where('direction', 'in')->distinct('jid')->count('jid');

        // 7 hari terakhir (termasuk hari ini)
        $start  = Carbon::today()->subDays(6);
        $daily  = ChatbotMessage::select(
                DB::raw('DATE(occurred_at) as date'),
                DB::raw('COUNT(*) as count')
            )
            ->where('occurred_at', '>=', $start)
            ->where('direction', 'in')
            ->groupBy('date')
            ->pluck('count', 'date');

        $dailyStats = [];
        for ($i = 6; $i >= 0; $i--) {
            $d = Carbon::today()->subDays($i)->toDateString();
            $dailyStats[] = ['date' => $d, 'count' => (int) ($daily[$d] ?? 0)];
        }

        // Distribusi per jam (24 jam terakhir).
        // Pakai ekspresi sesuai driver agar portable MySQL/Postgres/SQLite.
        $driver = DB::connection()->getDriverName();
        $hourExpr = match ($driver) {
            'sqlite' => "CAST(strftime('%H', occurred_at) AS INTEGER)",
            'pgsql'  => 'EXTRACT(HOUR FROM occurred_at)',
            default  => 'HOUR(occurred_at)', // mysql/mariadb
        };

        $hourlyRaw = ChatbotMessage::select(
                DB::raw("$hourExpr as hour"),
                DB::raw('COUNT(*) as count')
            )
            ->where('occurred_at', '>=', now()->subDay())
            ->groupBy(DB::raw($hourExpr))
            ->pluck('count', 'hour')
            ->toArray();

        $hourly = array_fill(0, 24, 0);
        foreach ($hourlyRaw as $h => $c) {
            $hourly[(int) $h] = (int) $c;
        }

        // Top 5 users
        $topUsers = ChatbotMessage::select('jid', DB::raw('COUNT(*) as count'))
            ->where('direction', 'in')
            ->groupBy('jid')
            ->orderByDesc('count')
            ->limit(5)
            ->get();

        $state = ChatbotState::current();

        return response()->json([
            'total_messages'    => $total,
            'incoming_messages' => $incoming,
            'outgoing_messages' => $outgoing,
            'unique_users'      => $uniqueUsr,
            'connected'         => (bool) $state->connected,
            'daily_stats'       => $dailyStats,
            'hourly_stats'      => $hourly,
            'top_users'         => $topUsers,
            'last_event_at'     => $state->last_event_at?->toIso8601String(),
        ]);
    }

    /**
     * GET /api/dashboard/messages?limit=50
     */
    public function messages(): JsonResponse
    {
        $limit = (int) request('limit', 50);
        $limit = max(1, min($limit, 200));

        $rows = ChatbotMessage::orderByDesc('occurred_at')
            ->limit($limit)
            ->get();

        return response()->json($rows);
    }

    /**
     * GET /api/dashboard/qr
     * Ambil QR saat ini. Status bisa: connected | waiting | ready.
     */
    public function qr(): JsonResponse
    {
        $state = ChatbotState::current();

        return response()->json([
            'connected' => (bool) $state->connected,
            'qr_b64'    => $state->qr_b64,
        ]);
    }
}
