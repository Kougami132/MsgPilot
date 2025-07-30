import request from '@/utils/request'
import type { Message, CreateMessageRequest, PaginationParams } from '@/types/api'

export const messageApi = {
  // 创建消息
  create: (data: CreateMessageRequest): Promise<Message> => {
    return request.post('/message/create', data)
  },

  // 获取消息列表
  getList: (params?: PaginationParams): Promise<Message[]> => {
    return request.get('/message/list', { params })
  },

  // 根据ID获取消息
  getById: (id: number): Promise<Message> => {
    return request.get(`/message/get/${id}`)
  },

  // 更新消息
  update: (id: number, data: Partial<CreateMessageRequest>): Promise<Message> => {
    return request.put(`/message/update/${id}`, data)
  },

  // 删除消息
  delete: (id: number): Promise<void> => {
    return request.delete(`/message/delete/${id}`)
  }
}
