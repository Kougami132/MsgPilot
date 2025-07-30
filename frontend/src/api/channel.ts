import request from '@/utils/request'
import type { Channel, CreateChannelRequest } from '@/types/api'

export const channelApi = {
  // 创建通道
  create: (data: CreateChannelRequest): Promise<Channel> => {
    return request.post('/channel/create', data)
  },

  // 获取通道列表
  getList: (): Promise<Channel[]> => {
    return request.get('/channel/list')
  },

  // 根据ID获取通道
  getById: (id: number): Promise<Channel> => {
    return request.get(`/channel/get/${id}`)
  },

  // 更新通道
  update: (id: number, data: Partial<CreateChannelRequest>): Promise<Channel> => {
    return request.put(`/channel/update/${id}`, data)
  },

  // 删除通道
  delete: (id: number): Promise<void> => {
    return request.delete(`/channel/delete/${id}`)
  },

  // 测试推送
  testPush: (data: Channel): Promise<any> => {
    return request.post('/channel/test', data)
  }
}
