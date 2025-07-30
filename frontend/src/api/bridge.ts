import request from '@/utils/request'
import type { Bridge, CreateBridgeRequest } from '@/types/api'

export const bridgeApi = {
  // 创建中转配置
  create: (data: CreateBridgeRequest): Promise<Bridge> => {
    return request.post('/bridge/create', data)
  },

  // 获取中转配置列表
  getList: (): Promise<Bridge[]> => {
    return request.get('/bridge/list')
  },

  // 根据ID获取中转配置
  getById: (id: number): Promise<Bridge> => {
    return request.get(`/bridge/get/${id}`)
  },

  // 更新中转配置
  update: (id: number, data: Partial<CreateBridgeRequest>): Promise<Bridge> => {
    return request.put(`/bridge/update/${id}`, data)
  },

  // 删除中转配置
  delete: (id: number): Promise<void> => {
    return request.delete(`/bridge/delete/${id}`)
  }
}
