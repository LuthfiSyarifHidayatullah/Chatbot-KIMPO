<script setup>
import { computed } from 'vue'
import { Bar } from 'vue-chartjs'
import {
  Chart, BarElement, LinearScale, CategoryScale, Tooltip, Legend,
} from 'chart.js'

Chart.register(BarElement, LinearScale, CategoryScale, Tooltip, Legend)

const props = defineProps({ data: { type: Array, required: true } })

const chartData = computed(() => ({
  labels: Array.from({ length: 24 }, (_, i) => `${i}:00`),
  datasets: [{
    label: 'Pesan',
    data: props.data,
    backgroundColor: 'rgba(167,139,250,0.6)',
    borderColor: '#a78bfa',
    borderWidth: 1,
    borderRadius: 4,
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
    <h2>Distribusi Pesan per Jam</h2>
    <div class="wrap"><Bar :data="chartData" :options="options" /></div>
  </div>
</template>

<style scoped>
.panel { background: #1e293b; border: 1px solid #334155; border-radius: 14px; padding: 20px; }
h2 { margin: 0 0 16px; font-size: 16px; color: #e2e8f0; }
.wrap { height: 260px; position: relative; }
</style>
