<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration {
    public function up(): void
    {
        Schema::create('chatbot_messages', function (Blueprint $table) {
            $table->id();
            $table->string('jid', 120)->index();          // JID WhatsApp pengirim/penerima
            $table->enum('direction', ['in', 'out'])->index();
            $table->text('message');
            $table->text('reply')->nullable();
            $table->timestamp('occurred_at')->index();     // waktu dari bot (bukan created_at)
            $table->timestamps();

            $table->index(['direction', 'occurred_at']);
        });

        Schema::create('chatbot_states', function (Blueprint $table) {
            // single-row table untuk status global bot
            $table->id();
            $table->boolean('connected')->default(false);
            $table->longText('qr_b64')->nullable();
            $table->timestamp('last_event_at')->nullable();
            $table->timestamps();
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('chatbot_messages');
        Schema::dropIfExists('chatbot_states');
    }
};
