<script setup>
defineProps({ messages: { type: Array, required: true } })

function fmt(iso) {
  if (!iso) return '-'
  return new Date(iso).toLocaleTimeString('id-ID', {
    hour: '2-digit', minute: '2-digit', second: '2-digit',
  })
}
</script>

<template>
  <div class="panel">
    <h2>Log Pesan Terakhir</h2>
    <div class="wrap">
      <table v-if="messages.length">
        <thead>
          <tr>
            <th>Waktu</th>
            <th>Arah</th>
            <th>Kontak</th>
            <th>Pesan</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="m in messages" :key="m.id">
            <td>{{ fmt(m.occurred_at) }}</td>
            <td>
              <span class="badge" :class="m.direction">
                {{ m.direction === 'in' ? 'MASUK' : 'KELUAR' }}
              </span>
            </td>
            <td class="jid">{{ m.jid }}</td>
            <td>
              {{ m.message }}
              <div v-if="m.reply" class="reply">↳ {{ m.reply }}</div>
            </td>
          </tr>
        </tbody>
      </table>
      <p v-else class="muted">Belum ada pesan</p>
    </div>
  </div>
</template>

<style scoped>
.panel { background: #1e293b; border: 1px solid #334155; border-radius: 14px; padding: 20px; }
h2 { margin: 0 0 16px; font-size: 16px; color: #e2e8f0; }
.wrap { max-height: 420px; overflow-y: auto; }
table { width: 100%; border-collapse: collapse; font-size: 13px; }
th, td { text-align: left; padding: 8px; border-bottom: 1px solid #334155; color: #e2e8f0; }
th { color: #94a3b8; font-size: 11px; text-transform: uppercase; letter-spacing: 0.5px; }
.badge { padding: 2px 8px; border-radius: 6px; font-size: 11px; font-weight: 600; }
.badge.in { background: rgba(34,211,238,0.15); color: #22d3ee; }
.badge.out { background: rgba(167,139,250,0.15); color: #a78bfa; }
.jid { font-family: monospace; font-size: 11px; color: #94a3b8; }
.reply { color: #94a3b8; font-style: italic; margin-top: 2px; font-size: 12px; }
.muted { color: #94a3b8; font-size: 13px; text-align: center; padding: 40px 0; }
</style>
