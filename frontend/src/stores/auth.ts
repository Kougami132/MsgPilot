import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api/auth'
import type { LoginRequest, TokenResponse, User } from '@/types/auth'

// 初始化用户信息的辅助函数
const initUser = (): User | null => {
  const savedUsername = localStorage.getItem('username')
  if (savedUsername && localStorage.getItem('token')) {
    return {
      id: 1,
      username: savedUsername,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    }
  }
  return null
}

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const token = ref<string>(localStorage.getItem('token') || '')
  const user = ref<User | null>(initUser()) // 立即从localStorage恢复用户信息
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
      localStorage.setItem('username', credentials.username)

      // 设置用户信息
      user.value = {
        id: 1,
        username: credentials.username,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      }

      // 尝试获取完整用户信息
      try {
        await fetchUserInfo()
      } catch (error) {
        console.log('获取用户详细信息失败，使用登录用户名')
      }
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
      localStorage.setItem('username', credentials.username)

      // 设置用户信息
      user.value = {
        id: 1,
        username: credentials.username,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      }

      // 尝试获取完整用户信息
      try {
        await fetchUserInfo()
      } catch (error) {
        console.log('获取用户详细信息失败，使用注册用户名')
      }
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
    localStorage.removeItem('username')
  }

  // 获取用户信息
  const fetchUserInfo = async (): Promise<void> => {
    try {
      const userInfo = await authApi.getCurrentUser()
      user.value = userInfo
      // 保存用户名到localStorage
      localStorage.setItem('username', userInfo.username)
    } catch (error) {
      console.error('获取用户信息失败:', error)
      // 如果获取失败，使用保存的用户名或默认值
      const savedUsername = localStorage.getItem('username')
      user.value = {
        id: 1,
        username: savedUsername || 'admin',
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
      // 从localStorage恢复用户名
      const savedUsername = localStorage.getItem('username')
      if (savedUsername) {
        user.value = {
          id: 1,
          username: savedUsername,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString()
        }
      }

      // 尝试获取完整用户信息
      try {
        await fetchUserInfo()
      } catch (error) {
        console.log('获取用户详细信息失败，使用保存的用户名')
      }
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
