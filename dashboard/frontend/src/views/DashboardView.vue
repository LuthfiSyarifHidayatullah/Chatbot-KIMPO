<script setup>
import { onMounted, onBeforeUnmount } from 'vue'
import { storeToRefs } from 'pinia'
import { useDashboardStore } from '../stores/dashboard'
import { useAuthStore } from '../stores/auth'
import StatsCards from '../components/StatsCards.vue'
import DailyChart from '../components/DailyChart.vue'
import HourlyChart from '../components/HourlyChart.vue'
import QRPanel from '../components/QRPanel.vue'
import MessageLog from '../components/MessageLog.vue'
import SendForm from '../components/SendForm.vue'
import TopUsers from '../components/TopUsers.vue'

const dash = useDashboardStore()
const auth = useAuthStore()
const { stats, messages, qr, loading } = storeToRefs(dash)

let refreshTimer = null

onMounted(async () => {
  await dash.loadAll()
  dash.subscribe()
  // Safety net: kalau websocket drop, polling pelan 30 detik sekali
  refreshTimer = setInterval(() => dash.loadAll(), 30000)
})

onBeforeUnmount(() => {
  if (refreshTimer) clearInterval(refreshTimer)
})
</script>

<template>
  <div class="page">
    <header>
      <h1>🤖 KIMPO Dashboard</h1>
      <div class="right">
        <span class="status">
          <span class="dot" :class="{ on: stats.connected }"></span>
          {{ stats.connected ? 'Terhubung' : 'Tidak terhubung' }}
        </span>
        <span class="user">{{ auth.user?.name }}</span>
        <button class="logout" @click="auth.logout()">Keluar</button>
      </div>
    </header>

    <main v-if="!loading">
      <StatsCards :stats="stats" />

      <div class="grid two">
        <DailyChart :data="stats.daily_stats" />
        <QRPanel :qr="qr" />
      </div>

      <div class="grid two">
        <HourlyChart :data="stats.hourly_stats" />
        <TopUsers :users="stats.top_users" />
      </div>

      <div class="grid two">
        <MessageLog :messages="messages" />
        <SendForm @sent="dash.loadAll" />
      </div>
    </main>

    <p v-else class="loading">Memuat data...</p>
  </div>
</template>

<style scoped>
.page { max-width: 1400px; margin: 0 auto; padding: 24px 32px; }
header {
  display: flex; justify-content: space-between; align-items: center;
  padding-bottom: 20px; border-bottom: 1px solid #334155; margin-bottom: 24px;
}
h1 {
  font-size: 22px; margin: 0;
  background: linear-gradient(90deg, #22d3ee, #a78bfa);
  -webkit-background-clip: text; -webkit-text-fill-color: transparent;
}
.right { display: flex; align-items: center; gap: 16px; }
.status {
  display: flex; align-items: center; gap: 8px;
  padding: 6px 12px; border-radius: 999px;
  background: #1e293b; border: 1px solid #334155; font-size: 13px;
}
.dot { width: 10px; height: 10px; border-radius: 50%; background: #ef4444; }
.dot.on { background: #10b981; box-shadow: 0 0 10px #10b981; }
.user { color: #94a3b8; font-size: 13px; }
.logout {
  background: transparent; color: #94a3b8; border: 1px solid #334155;
  padding: 6px 12px; border-radius: 6px; cursor: pointer; font-size: 13px;
}
.logout:hover { color: #e2e8f0; border-color: #64748b; }
.grid { display: grid; gap: 20px; margin-bottom: 20px; }
.grid.two { grid-template-columns: 2fr 1fr; }
@media (max-width: 900px) { .grid.two { grid-template-columns: 1fr; } }
.loading { color: #94a3b8; text-align: center; padding: 80px 0; }
</style>
