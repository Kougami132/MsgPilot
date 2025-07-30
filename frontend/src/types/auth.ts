// 认证相关类型定义

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  password: string
}

export interface ChangePasswordRequest {
  username: string
  old_password: string
  new_password: string
}

export interface TokenResponse {
  access_token: string
  expiry: number
}

export interface User {
  id: number
  username: string
  created_at: string
  updated_at: string
}
