// API 通用类型定义

export interface ApiResponse<T = any> {
  data?: T
  message?: string
  error?: string
}

export interface PaginationParams {
  page?: number
  size?: number
}

export interface PaginationResponse<T> {
  data: T[]
  total: number
  page: number
  size: number
}

// 消息相关类型
export interface Message {
  id: number
  title: string
  content: string
  status: number // 0=待发送，1=发送中，2=已发送，3=失败
  error_message?: string
  bridge_id: number
  bridge?: Bridge
  created_at: string
  updated_at: string
}

export interface CreateMessageRequest {
  title: string
  content: string
  bridge_id: number
}

// 渠道相关类型
export interface Channel {
  id: number
  name: string
  type: string
  config: Record<string, any>
  created_at: string
  updated_at: string
}

export interface CreateChannelRequest {
  name: string
  type: string
  config: Record<string, any>
}

// 中转配置相关类型
export interface Bridge {
  id: number
  name: string
  source_channel_id: number
  source_channel_type?: string
  target_channel_id: number
  source_channel?: Channel
  target_channel?: Channel
  ticket?: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface CreateBridgeRequest {
  name: string
  source_channel_type: string
  target_channel_id: number | null
  ticket?: string
  is_active: boolean
}
