<script setup>
defineProps({ stats: { type: Object, required: true } })
</script>

<template>
  <div class="cards">
    <div class="card">
      <h3>Total Pesan</h3>
      <div class="value">{{ stats.total_messages }}</div>
    </div>
    <div class="card accent">
      <h3>Pesan Masuk</h3>
      <div class="value">{{ stats.incoming_messages }}</div>
    </div>
    <div class="card success">
      <h3>Pesan Terkirim</h3>
      <div class="value">{{ stats.outgoing_messages }}</div>
    </div>
    <div class="card warning">
      <h3>Pengguna Unik</h3>
      <div class="value">{{ stats.unique_users }}</div>
    </div>
  </div>
</template>

<style scoped>
.cards { display: grid; gap: 16px; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); margin-bottom: 24px; }
.card {
  background: #1e293b; border: 1px solid #334155; border-radius: 14px; padding: 20px;
}
h3 {
  margin: 0 0 8px; font-size: 12px; color: #94a3b8; text-transform: uppercase; letter-spacing: 0.5px;
}
.value { font-size: 30px; font-weight: 700; color: #22d3ee; }
.card.accent .value { color: #a78bfa; }
.card.success .value { color: #10b981; }
.card.warning .value { color: #f59e0b; }
</style>
