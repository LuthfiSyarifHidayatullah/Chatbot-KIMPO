<?php

use Illuminate\Support\Facades\Broadcast;

/*
|--------------------------------------------------------------------------
| Broadcast Channels - KIMPO Dashboard
|--------------------------------------------------------------------------
| Channel private 'dashboard' hanya boleh didengarkan oleh user yang login
| (Sanctum). Di masa depan bisa ditambah pengecekan role admin.
*/

Broadcast::channel('dashboard', function ($user) {
    return $user !== null;
});
