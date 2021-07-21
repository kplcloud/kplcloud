import axios from 'axios'
import store from '@/store'
import storage from 'store'
import { VueAxios } from './axios'
import { ACCESS_TOKEN } from '@/store/mutation-types'

// 创建 axios 实例
const request = axios.create({
  // API 请求的默认前缀
  baseURL: process.env.VUE_APP_API_BASE_URL,
  timeout: 6000 // 请求超时时间
})

// 异常拦截处理器
const errorHandler = (error) => {
  if (error.response) {
    const data = error.response.data
    // 从 localstorage 获取 token
    const token = storage.get(ACCESS_TOKEN)
    if (error.response.status === 403) {
      store.state.app.toast.message = data.message
      store.state.app.toast.show = true
      store.state.app.toast.color = 'error'
    }
    if (error.response.status === 401 && !(data.result && data.result.isLogin)) {
      store.state.app.toast.message = 'Authorization verification failed'
      store.state.app.toast.show = true
      store.state.app.toast.color = 'error'
      if (token) {
        store.dispatch('Logout').then(() => {
          setTimeout(() => {
            window.location.reload()
          }, 1500)
        })
      }
    }
  }
  return Promise.reject(error)
}

// request interceptor
request.interceptors.request.use(config => {
  const token = storage.get(ACCESS_TOKEN)
  // 如果 token 存在
  // 让每个请求携带自定义 token 请根据实际情况自行修改
  if (token) {
    config.headers['Authorization'] = token
    config.headers['content-type'] = 'application/json'
  }
  config.mode = 'cors'
  return config
}, errorHandler)

// response interceptor
request.interceptors.response.use((response) => {
  const res = response.data
  if (!res) {
    store.state.app.toast.message = '网络错误,请重试'
    store.state.app.toast.show = true
    store.state.app.toast.color = 'error'
    return Promise.reject('网络错误,请重试')
  }
  if (!res.success) {
    if (res.code === 501) {
      // 跳去授权登录
      storage.remove(ACCESS_TOKEN)
      window.location.href = '/'
    }
    if (res.code === 1005) {
      // 跳去授权登录
      storage.remove(ACCESS_TOKEN)
      window.location.href = '/'
    }
    store.state.app.toast.message = res.message
    store.state.app.toast.show = true
    store.state.app.toast.color = 'error'
    return Promise.reject(res.message)
  }
  return res
}, errorHandler)

const installer = {
  vm: {},
  install(Vue) {
    Vue.use(VueAxios, request)
  }
}

export default request

export {
  installer as VueAxios,
  request as axios
}
