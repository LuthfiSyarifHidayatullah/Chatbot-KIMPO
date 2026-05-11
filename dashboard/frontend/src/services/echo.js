import Echo from 'laravel-echo'
import Pusher from 'pusher-js'
import api from './api'

window.Pusher = Pusher

/**
 * Inisialisasi Laravel Echo dengan Reverb.
 * Dipanggil setelah user login agar auth endpoint dapat mengakses
 * session user untuk otorisasi private channel.
 */
let echoInstance = null

export function getEcho() {
  if (echoInstance) return echoInstance

  echoInstance = new Echo({
    broadcaster: 'reverb',
    key: import.meta.env.VITE_REVERB_APP_KEY,
    wsHost: import.meta.env.VITE_REVERB_HOST,
    wsPort: import.meta.env.VITE_REVERB_PORT,
    wssPort: import.meta.env.VITE_REVERB_PORT,
    forceTLS: (import.meta.env.VITE_REVERB_SCHEME || 'http') === 'https',
    enabledTransports: ['ws', 'wss'],
    authorizer: (channel) => ({
      authorize: (socketId, callback) => {
        api.post('/broadcasting/auth', {
          socket_id: socketId,
          channel_name: channel.name,
        })
          .then((r) => callback(null, r.data))
          .catch((e) => callback(e, null))
      },
    }),
  })

  return echoInstance
}

export function resetEcho() {
  if (echoInstance) {
    echoInstance.disconnect()
    echoInstance = null
  }
}
