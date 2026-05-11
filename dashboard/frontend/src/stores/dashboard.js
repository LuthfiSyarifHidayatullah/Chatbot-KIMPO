import { defineStore } from 'pinia'
import api from '../services/api'
import { getEcho } from '../services/echo'

export const useDashboardStore = defineStore('dashboard', {
  state: () => ({
    stats: {
      total_messages: 0,
      incoming_messages: 0,
      outgoing_messages: 0,
      unique_users: 0,
      connected: false,
      daily_stats: [],
      hourly_stats: Array(24).fill(0),
      top_users: [],
    },
    messages: [],
    qr: { connected: false, qr_b64: null },
    loading: false,
    subscribed: false,
  }),

  actions: {
    async loadAll() {
      this.loading = true
      try {
        const [s, m, q] = await Promise.all([
          api.get('/api/dashboard/stats'),
          api.get('/api/dashboard/messages?limit=50'),
          api.get('/api/dashboard/qr'),
        ])
        this.stats = s.data
        this.messages = m.data
        this.qr = q.data
      } finally {
        this.loading = false
      }
    },

    /**
     * Subscribe ke private-channel 'dashboard' untuk realtime updates.
     * Setiap event langsung memperbarui state lokal + optimistis tambah statistik.
     */
    subscribe() {
      if (this.subscribed) return
      const echo = getEcho()
      const channel = echo.private('dashboard')

      channel.listen('.message.received', (e) => {
        this.messages.unshift(e)
        if (this.messages.length > 100) this.messages.pop()

        this.stats.total_messages++
        if (e.direction === 'in') this.stats.incoming_messages++
        else this.stats.outgoing_messages++

        // update hourly bucket
        const hour = new Date(e.occurred_at).getHours()
        this.stats.hourly_stats[hour] = (this.stats.hourly_stats[hour] || 0) + 1
      })

      channel.listen('.connection.changed', (e) => {
        this.stats.connected = e.connected
        this.qr.connected = e.connected
        if (e.connected) this.qr.qr_b64 = null
      })

      channel.listen('.qr.updated', (e) => {
        this.qr.qr_b64 = e.qr_b64
        this.qr.connected = false
      })

      this.subscribed = true
    },

    async sendMessage(to, message) {
      const { data } = await api.post('/api/dashboard/send', { to, message })
      return data
    },
  },
})
