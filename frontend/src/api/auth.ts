import request from '@/utils/request'
import type { LoginRequest, RegisterRequest, ChangePasswordRequest, TokenResponse, User } from '@/types/auth'

export const authApi = {
  // 登录
  login: (data: LoginRequest): Promise<TokenResponse> => {
    return request.post('/auth/login', data)
  },

  // 注册
  register: (data: RegisterRequest): Promise<TokenResponse> => {
    return request.post('/auth/register', data)
  },

  // 刷新token
  refreshToken: (token: string): Promise<TokenResponse> => {
    return request.post('/auth/refresh', { token })
  },

  // 修改密码
  changePassword: (data: ChangePasswordRequest): Promise<void> => {
    return request.post('/auth/changePassword', data)
  },

  // 获取当前用户信息
  getCurrentUser: (): Promise<User> => {
    return request.get('/auth/me')
  }
}
