<?php

namespace App\Events;

use App\Models\ChatbotMessage;
use Illuminate\Broadcasting\Channel;
use Illuminate\Broadcasting\InteractsWithSockets;
use Illuminate\Broadcasting\PrivateChannel;
use Illuminate\Contracts\Broadcasting\ShouldBroadcast;
use Illuminate\Foundation\Events\Dispatchable;
use Illuminate\Queue\SerializesModels;

class MessageReceived implements ShouldBroadcast
{
    use Dispatchable, InteractsWithSockets, SerializesModels;

    public function __construct(public ChatbotMessage $message)
    {
    }

    public function broadcastOn(): array
    {
        return [new PrivateChannel('dashboard')];
    }

    public function broadcastAs(): string
    {
        return 'message.received';
    }

    public function broadcastWith(): array
    {
        return [
            'id'          => $this->message->id,
            'jid'         => $this->message->jid,
            'direction'   => $this->message->direction,
            'message'     => $this->message->message,
            'reply'       => $this->message->reply,
            'occurred_at' => $this->message->occurred_at?->toIso8601String(),
        ];
    }
}
