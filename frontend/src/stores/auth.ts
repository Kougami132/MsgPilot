import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api/auth'
import type { LoginRequest, TokenResponse, User } from '@/types/auth'

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const token = ref<string>(localStorage.getItem('token') || '')
  const user = ref<User | null>(null)
  const loading = ref(false)

  // 计算属性
  const isAuthenticated = computed(() => !!token.value)

  // 登录
  const login = async (credentials: LoginRequest): Promise<void> => {
    loading.value = true
    try {
      const response = await authApi.login(credentials)
      token.value = response.access_token
      localStorage.setItem('token', response.access_token)
      localStorage.setItem('token_expiry', response.expiry.toString())
      
      // 获取用户信息
      await fetchUserInfo()
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  // 注册
  const register = async (credentials: LoginRequest): Promise<void> => {
    loading.value = true
    try {
      const response = await authApi.register(credentials)
      token.value = response.access_token
      localStorage.setItem('token', response.access_token)
      localStorage.setItem('token_expiry', response.expiry.toString())
      
      // 获取用户信息
      await fetchUserInfo()
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  // 登出
  const logout = (): void => {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('token_expiry')
  }

  // 获取用户信息
  const fetchUserInfo = async (): Promise<void> => {
    try {
      const userInfo = await authApi.getCurrentUser()
      user.value = userInfo
    } catch (error) {
      console.error('获取用户信息失败:', error)
      // 如果获取失败，设置默认用户信息
      user.value = {
        id: 1,
        username: 'admin',
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      }
    }
  }

  // 检查token是否过期
  const checkTokenExpiry = (): boolean => {
    const expiry = localStorage.getItem('token_expiry')
    if (!expiry) return false
    
    const expiryTime = parseInt(expiry) * 1000 // 转换为毫秒
    const now = Date.now()
    
    if (now >= expiryTime) {
      logout()
      return false
    }
    return true
  }

  // 初始化时检查token
  const initialize = async (): Promise<void> => {
    if (token.value && checkTokenExpiry()) {
      await fetchUserInfo()
    } else {
      logout()
    }
  }

  return {
    // 状态
    token: readonly(token),
    user: readonly(user),
    loading: readonly(loading),
    
    // 计算属性
    isAuthenticated,
    
    // 方法
    login,
    register,
    logout,
    fetchUserInfo,
    checkTokenExpiry,
    initialize
  }
})
