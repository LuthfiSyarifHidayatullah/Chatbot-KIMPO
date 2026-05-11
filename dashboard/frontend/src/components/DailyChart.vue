<script setup>
import { computed } from 'vue'
import { Line } from 'vue-chartjs'
import {
  Chart, LineElement, PointElement, LinearScale, CategoryScale, Tooltip, Legend, Filler,
} from 'chart.js'

Chart.register(LineElement, PointElement, LinearScale, CategoryScale, Tooltip, Legend, Filler)

const props = defineProps({ data: { type: Array, required: true } })

const chartData = computed(() => ({
  labels: props.data.map((d) => d.date.slice(5)),
  datasets: [{
    label: 'Pesan Masuk',
    data: props.data.map((d) => d.count),
    borderColor: '#22d3ee',
    backgroundColor: 'rgba(34,211,238,0.15)',
    tension: 0.35,
    fill: true,
  }],
}))

const options = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: { legend: { labels: { color: '#e2e8f0' } } },
  scales: {
    x: { ticks: { color: '#94a3b8' }, grid: { color: 'rgba(148,163,184,0.1)' } },
    y: { ticks: { color: '#94a3b8' }, grid: { color: 'rgba(148,163,184,0.1)' }, beginAtZero: true },
  },
}
</script>

<template>
  <div class="panel">
    <h2>Statistik 7 Hari Terakhir</h2>
    <div class="wrap"><Line :data="chartData" :options="options" /></div>
  </div>
</template>

<style scoped>
.panel { background: #1e293b; border: 1px solid #334155; border-radius: 14px; padding: 20px; }
h2 { margin: 0 0 16px; font-size: 16px; color: #e2e8f0; }
.wrap { height: 260px; position: relative; }
</style>
