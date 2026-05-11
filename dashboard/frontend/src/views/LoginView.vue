<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

const auth = useAuthStore()
const router = useRouter()

async function submit() {
  error.value = ''
  loading.value = true
  try {
    await auth.login(email.value, password.value)
    router.push({ name: 'dashboard' })
  } catch (e) {
    error.value = e.response?.data?.message || 'Gagal login'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-wrap">
    <form class="login-card" @submit.prevent="submit">
      <h1>🤖 KIMPO Dashboard</h1>
      <p class="muted">Masuk untuk mengakses dashboard admin</p>

      <label>Email</label>
      <input v-model="email" type="email" required autocomplete="email" />

      <label>Password</label>
      <input v-model="password" type="password" required autocomplete="current-password" />

      <p v-if="error" class="err">{{ error }}</p>

      <button :disabled="loading">{{ loading ? 'Memproses...' : 'Masuk' }}</button>
    </form>
  </div>
</template>

<style scoped>
.login-wrap {
  min-height: 100vh;
  display: grid;
  place-items: center;
  background: linear-gradient(135deg, #0f172a, #1e293b);
}
.login-card {
  background: #1e293b;
  padding: 32px;
  border-radius: 14px;
  width: 360px;
  border: 1px solid #334155;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
h1 { margin: 0 0 4px; font-size: 20px; color: #e2e8f0; }
.muted { color: #94a3b8; margin: 0 0 12px; font-size: 13px; }
label { font-size: 12px; color: #94a3b8; margin-top: 6px; }
input {
  background: #273449; color: #e2e8f0; border: 1px solid #334155;
  padding: 10px 12px; border-radius: 8px; font-size: 14px;
}
input:focus { outline: none; border-color: #22d3ee; }
.err { color: #ef4444; font-size: 13px; margin: 4px 0 0; }
button {
  margin-top: 12px;
  background: linear-gradient(90deg, #22d3ee, #a78bfa);
  color: #0f172a; border: 0; padding: 12px; border-radius: 8px;
  font-weight: 600; cursor: pointer;
}
button:disabled { opacity: 0.6; cursor: not-allowed; }
</style>
