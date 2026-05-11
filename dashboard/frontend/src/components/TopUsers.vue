<script setup>
defineProps({ users: { type: Array, required: true } })
</script>

<template>
  <div class="panel">
    <h2>Pengguna Teraktif</h2>
    <ul v-if="users.length" class="list">
      <li v-for="u in users" :key="u.jid">
        <span class="jid">{{ u.jid }}</span>
        <strong>{{ u.count }}</strong>
      </li>
    </ul>
    <p v-else class="muted">Belum ada data</p>
  </div>
</template>

<style scoped>
.panel { background: #1e293b; border: 1px solid #334155; border-radius: 14px; padding: 20px; }
h2 { margin: 0 0 16px; font-size: 16px; color: #e2e8f0; }
.list { list-style: none; padding: 0; margin: 0; }
.list li {
  display: flex; justify-content: space-between; align-items: center;
  padding: 10px 0; border-bottom: 1px solid #334155; font-size: 13px;
}
.list li:last-child { border-bottom: none; }
.jid { color: #94a3b8; font-family: monospace; font-size: 12px; }
strong { color: #22d3ee; }
.muted { color: #94a3b8; font-size: 13px; }
</style>
