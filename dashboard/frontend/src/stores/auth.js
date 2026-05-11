import { defineStore } from 'pinia'
import api, { csrf } from '../services/api'
import { resetEcho } from '../services/echo'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    ready: false,
  }),
  actions: {
    async fetchUser() {
      try {
        const { data } = await api.get('/api/user')
        this.user = data
      } catch {
        this.user = null
      } finally {
        this.ready = true
      }
    },
    async login(email, password) {
      await csrf()
      await api.post('/login', { email, password })
      await this.fetchUser()
    },
    async logout() {
      await api.post('/logout')
      this.user = null
      resetEcho()
    },
  },
})
