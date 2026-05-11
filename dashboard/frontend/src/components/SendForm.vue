<script setup>
import { ref } from 'vue'
import { useDashboardStore } from '../stores/dashboard'

const emit = defineEmits(['sent'])
const dash = useDashboardStore()

const to = ref('')
const message = ref('')
const loading = ref(false)
const alert = ref(null)

async function submit() {
  alert.value = null
  loading.value = true
  try {
    await dash.sendMessage(to.value.trim(), message.value.trim())
    alert.value = { type: 'success', text: 'Pesan berhasil dikirim!' }
    message.value = ''
    emit('sent')
  } catch (e) {
    alert.value = {
      type: 'error',
      text: e.response?.data?.error || e.response?.data?.message || 'Gagal kirim',
    }
  } finally {
    loading.value = false
    setTimeout(() => (alert.value = null), 4000)
  }
}
</script>

<template>
  <div class="panel">
    <h2>Kirim Pesan Manual</h2>
    <form @submit.prevent="submit">
      <input v-model="to" placeholder="Nomor (628xxx) atau JID lengkap" required />
      <textarea v-model="message" placeholder="Tulis pesan..." required></textarea>
      <button :disabled="loading">
        {{ loading ? 'Mengirim...' : 'Kirim Pesan' }}
      </button>
      <p v-if="alert" class="alert" :class="alert.type">{{ alert.text }}</p>
    </form>
  </div>
</template>

<style scoped>
.panel { background: #1e293b; border: 1px solid #334155; border-radius: 14px; padding: 20px; }
h2 { margin: 0 0 16px; font-size: 16px; color: #e2e8f0; }
form { display: flex; flex-direction: column; gap: 12px; }
input, textarea {
  background: #273449; color: #e2e8f0; border: 1px solid #334155;
  padding: 10px 12px; border-radius: 8px; font-size: 14px; font-family: inherit;
}
textarea { min-height: 90px; resize: vertical; }
input:focus, textarea:focus { outline: none; border-color: #22d3ee; }
button {
  background: linear-gradient(90deg, #22d3ee, #a78bfa);
  color: #0f172a; border: 0; padding: 10px 16px; border-radius: 8px;
  font-weight: 600; cursor: pointer;
}
button:disabled { opacity: 0.5; cursor: not-allowed; }
.alert { padding: 10px 14px; border-radius: 8px; font-size: 13px; margin: 0; }
.alert.success { background: rgba(16,185,129,0.15); color: #10b981; }
.alert.error { background: rgba(239,68,68,0.15); color: #ef4444; }
</style>
