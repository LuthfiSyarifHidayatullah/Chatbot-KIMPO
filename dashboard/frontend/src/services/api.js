import axios from 'axios'

/**
 * Axios instance terpusat.
 * - baseURL dari VITE_API_BASE_URL
 * - withCredentials true supaya cookie Sanctum ikut terkirim
 * - interceptor otomatis redirect ke /login saat 401
 */
const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  withCredentials: true,
  withXSRFToken: true,
  headers: {
    'Accept': 'application/json',
    'X-Requested-With': 'XMLHttpRequest',
  },
})

api.interceptors.response.use(
  (r) => r,
  (err) => {
    if (err.response?.status === 401 && window.location.pathname !== '/login') {
      window.location.href = '/login'
    }
    return Promise.reject(err)
  },
)

/** Panggil sebelum login agar CSRF cookie ter-set. */
export async function csrf() {
  await api.get('/sanctum/csrf-cookie')
}

export default api
