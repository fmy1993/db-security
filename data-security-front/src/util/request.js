import axios from 'axios'

const request = axios.create({
  baseURL: '/api'
})

request.interceptors.request.use(
  config => {
    if (localStorage.getItem('X_CSRF_Token')) {
      config.headers = {
        'X-CSRFToken': localStorage.getItem('X_CSRF_Token')
      }
      if (localStorage.getItem('Authorization')) {
        config.headers.Authorization = localStorage.getItem('Authorization')
      }
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

export default request
