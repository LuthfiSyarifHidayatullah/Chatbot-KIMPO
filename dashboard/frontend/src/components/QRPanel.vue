<script setup>
defineProps({ qr: { type: Object, required: true } })
</script>

<template>
  <div class="panel">
    <h2>Koneksi WhatsApp</h2>
    <div v-if="qr.connected" class="state">
      <div class="emoji">✅</div>
      <p class="ok">Terhubung</p>
      <p class="muted">Bot siap menerima pesan</p>
    </div>
    <div v-else-if="qr.qr_b64" class="state">
      <p class="muted">Pindai QR berikut di WhatsApp:</p>
      <img :src="`data:image/png;base64,${qr.qr_b64}`" alt="QR Code" />
    </div>
    <div v-else class="state">
      <div class="emoji">⏳</div>
      <p class="muted">Menunggu QR Code...</p>
    </div>
  </div>
</template>

<style scoped>
.panel { background: #1e293b; border: 1px solid #334155; border-radius: 14px; padding: 20px; text-align: center; }
h2 { margin: 0 0 16px; font-size: 16px; color: #e2e8f0; text-align: left; }
.state { padding: 20px 0; }
.emoji { font-size: 48px; }
.ok { color: #10b981; font-weight: 600; margin: 12px 0 4px; }
.muted { color: #94a3b8; font-size: 13px; }
img { width: 220px; height: 220px; margin-top: 12px; background: white; padding: 8px; border-radius: 12px; }
</style>
