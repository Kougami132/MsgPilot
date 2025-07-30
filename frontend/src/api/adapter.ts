import request from '@/utils/request'

export const adapterApi = {
  // 获取适配器列表
  getList: (): Promise<{
    adapters: string[]
    handlers: string[]
  }> => {
    return request.get('/adapter/list')
  }
}
