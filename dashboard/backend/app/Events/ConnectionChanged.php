<?php

namespace App\Events;

use Illuminate\Broadcasting\InteractsWithSockets;
use Illuminate\Broadcasting\PrivateChannel;
use Illuminate\Contracts\Broadcasting\ShouldBroadcast;
use Illuminate\Foundation\Events\Dispatchable;
use Illuminate\Queue\SerializesModels;

class ConnectionChanged implements ShouldBroadcast
{
    use Dispatchable, InteractsWithSockets, SerializesModels;

    public function __construct(public bool $connected)
    {
    }

    public function broadcastOn(): array
    {
        return [new PrivateChannel('dashboard')];
    }

    public function broadcastAs(): string
    {
        return 'connection.changed';
    }

    public function broadcastWith(): array
    {
        return ['connected' => $this->connected];
    }
}
