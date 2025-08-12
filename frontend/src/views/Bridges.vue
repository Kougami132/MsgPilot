<template>
  <div class="page-container">
    <div class="bridges-page">
      <!-- 页面头部 -->
      <div class="page-header">
        <h2>中转配置</h2>
        <el-button type="primary" @click="showCreateDialog = true">
          <el-icon><Plus /></el-icon>
          创建中转配置
        </el-button>
      </div>
      
      <!-- 中转配置列表 -->
      <el-card class="table-card">
        <el-table :data="bridges" :loading="loading" style="width: 100%">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column label="启用" width="80" align="center">
            <template #default="{ row }">
              <el-switch
                v-model="row.is_active"
                @change="handleToggleActive(row)"
              />
            </template>
          </el-table-column>
          <el-table-column prop="name" label="名称" min-width="150" />
          <el-table-column label="接收格式" width="120">
            <template #default="{ row }">
              {{ row.source_channel_type || '未知' }}
            </template>
          </el-table-column>
          <el-table-column label="发送渠道" width="120">
            <template #default="{ row }">
              {{ row.target_channel?.name || '未知' }}
            </template>
          </el-table-column>
          <el-table-column label="Ticket" width="120">
            <template #default="{ row }">
              <span class="ticket-display">{{ row.ticket || '未设置' }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="创建时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="220" fixed="right">
            <template #default="{ row }">
              <el-button-group>
                <el-button type="primary" size="small" @click="handleView(row)">
                  查看
                </el-button>
                <el-button type="warning" size="small" @click="handleEdit(row)">
                  编辑
                </el-button>
                <el-button type="danger" size="small" @click="handleDelete(row)">
                  删除
                </el-button>
              </el-button-group>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>
    
    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      :title="editingBridge ? '编辑中转配置' : '创建中转配置'"
      width="600px"
      @closed="handleDialogClose"
    >
      <el-form
        ref="bridgeFormRef"
        :model="bridgeForm"
        :rules="bridgeRules"
        label-width="100px"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="bridgeForm.name" placeholder="请输入中转配置名称" />
        </el-form-item>
        <el-form-item label="接收格式" prop="source_channel_type">
          <el-select
            v-model="bridgeForm.source_channel_type"
            placeholder="请选择接收格式类型"
            style="width: 100%"
          >
            <el-option
              v-for="adapter in adapters"
              :key="adapter"
              :label="adapter"
              :value="adapter"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="发送渠道" prop="target_channel_id">
          <el-select
            v-model="bridgeForm.target_channel_id"
            placeholder="请选择发送渠道"
            style="width: 100%"
          >
            <el-option
              v-for="channel in channels"
              :key="channel.id"
              :label="channel.name"
              :value="channel.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="Ticket" prop="ticket">
          <div class="ticket-input-group">
            <el-input
              v-model="bridgeForm.ticket"
              placeholder="请输入或生成ticket"
              style="flex: 1;"
            />
            <el-button
              type="primary"
              @click="generateTicket"
              style="margin-left: 8px;"
            >
              随机生成
            </el-button>
          </div>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          {{ editingBridge ? '更新' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
    
    <!-- 查看详情对话框 -->
    <el-dialog v-model="showViewDialog" title="中转配置详情" width="600px">
      <div v-if="viewingBridge" class="bridge-detail">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="ID">{{ viewingBridge.id }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="viewingBridge.is_active ? 'success' : 'danger'">
              {{ viewingBridge.is_active ? '启用' : '禁用' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="名称">{{ viewingBridge.name }}</el-descriptions-item>
          <el-descriptions-item label="接收格式">{{ viewingBridge.source_channel_type || '未知' }}</el-descriptions-item>
          <el-descriptions-item label="发送渠道">{{ viewingBridge.target_channel?.name || '未知' }}</el-descriptions-item>
          <el-descriptions-item label="Ticket">
            <div class="ticket-container">
              <span class="token-text">{{ viewingBridge.ticket || '未设置' }}</span>
              <el-button
                v-if="viewingBridge.ticket"
                type="primary"
                size="small"
                @click="copyToClipboard(viewingBridge.ticket)"
                style="margin-left: 8px;"
              >
                复制
              </el-button>
            </div>
          </el-descriptions-item>
          <el-descriptions-item label="接口地址">
            <div class="api-url-container">
              <span class="api-url-text">{{ getApiUrl(viewingBridge) }}</span>
              <el-button
                type="primary"
                size="small"
                @click="copyToClipboard(getApiUrl(viewingBridge))"
                style="margin-left: 8px;"
              >
                复制
              </el-button>
            </div>
          </el-descriptions-item>

          <el-descriptions-item label="创建时间">
            {{ formatDate(viewingBridge.created_at) }}
          </el-descriptions-item>
          <el-descriptions-item label="更新时间">
            {{ formatDate(viewingBridge.updated_at) }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { bridgeApi } from '@/api/bridge'
import { channelApi } from '@/api/channel'
import { adapterApi } from '@/api/adapter'
import type { Bridge, CreateBridgeRequest, Channel } from '@/types/api'

// 响应式数据
const loading = ref(false)
const submitting = ref(false)
const bridges = ref<Bridge[]>([])
const channels = ref<Channel[]>([])
const adapters = ref<string[]>([])

// 对话框状态
const showCreateDialog = ref(false)
const showViewDialog = ref(false)
const editingBridge = ref<Bridge | null>(null)
const viewingBridge = ref<Bridge | null>(null)

// 表单引用
const bridgeFormRef = ref<FormInstance>()

// 中转配置表单
const bridgeForm = reactive<CreateBridgeRequest>({
  name: '',
  source_channel_type: '',
  target_channel_id: null as any,
  ticket: '',
  is_active: true
})

// 表单验证规则
const bridgeRules: FormRules = {
  name: [
    { required: true, message: '请输入中转配置名称', trigger: 'blur' }
  ],
  source_channel_type: [
    { required: true, message: '请选择接收格式类型', trigger: 'change' }
  ],
  target_channel_id: [
    { required: true, message: '请选择发送渠道', trigger: 'change' }
  ]
}

// 格式化日期
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

// 获取接口地址
const getApiUrl = (bridge: Bridge) => {
  const baseUrl = window.location.origin // 获取当前域名
  const ticket = bridge.ticket || 'your-ticket'
  const sourceType = bridge.source_channel_type?.toLowerCase()

  switch (sourceType) {
    case 'onebot':
      return `${baseUrl}/api/onebot/${ticket}/send_msg`
    case 'bark':
      return `${baseUrl}/api/bark/${ticket}`
    case 'gotify':
      return `${baseUrl}/api/gotify/message?token=${ticket}`
    case 'pushdeer':
      return `${baseUrl}/api/pushdeer/message/push`
    case 'ntfy':
      return `${baseUrl}/api/ntfy/${ticket}`
    case 'webhook':
      return `${baseUrl}/api/webhook/${ticket}`
    default:
      return `${baseUrl}/api/${sourceType}/${ticket}`
  }
}



// 复制到剪贴板
const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch (error) {
    console.error('复制失败:', error)
    ElMessage.error('复制失败')
  }
}

// 生成随机ticket
const generateTicket = () => {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  let result = ''
  for (let i = 0; i < 8; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  bridgeForm.ticket = result
}

// 加载中转配置列表
const loadBridges = async () => {
  loading.value = true
  try {
    bridges.value = await bridgeApi.getList()
  } catch (error) {
    console.error('加载中转配置列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 加载渠道列表
const loadChannels = async () => {
  try {
    channels.value = await channelApi.getList()
  } catch (error) {
    console.error('加载渠道列表失败:', error)
  }
}

// 加载适配器adapters列表
const loadAdapters = async () => {
  try {
    const response = await adapterApi.getList()
    adapters.value = response.adapters
  } catch (error) {
    console.error('加载adapters列表失败:', error)
  }
}

// 切换启用状态
const handleToggleActive = async (bridge: Bridge) => {
  try {
    await bridgeApi.update(bridge.id, {
      name: bridge.name,
      source_channel_type: bridge.source_channel_type || '',
      target_channel_id: bridge.target_channel_id,
      is_active: bridge.is_active
    })
    ElMessage.success(bridge.is_active ? '已启用' : '已禁用')
  } catch (error) {
    // 如果更新失败，恢复原状态
    bridge.is_active = !bridge.is_active
    console.error('更新状态失败:', error)
  }
}

// 查看中转配置
const handleView = (bridge: Bridge) => {
  viewingBridge.value = bridge
  showViewDialog.value = true
}

// 编辑中转配置
const handleEdit = (bridge: Bridge) => {
  editingBridge.value = bridge
  bridgeForm.name = bridge.name
  bridgeForm.source_channel_type = bridge.source_channel_type || ''
  bridgeForm.target_channel_id = bridge.target_channel_id
  bridgeForm.ticket = bridge.ticket || ''
  bridgeForm.is_active = bridge.is_active
  showCreateDialog.value = true
}

// 删除中转配置
const handleDelete = async (bridge: Bridge) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除中转配置 "${bridge.name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await bridgeApi.delete(bridge.id)
    ElMessage.success('删除成功')
    loadBridges()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除中转配置失败:', error)
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!bridgeFormRef.value) return
  
  try {
    await bridgeFormRef.value.validate()
    submitting.value = true
    
    if (editingBridge.value) {
      await bridgeApi.update(editingBridge.value.id, bridgeForm)
      ElMessage.success('更新成功')
    } else {
      await bridgeApi.create(bridgeForm)
      ElMessage.success('创建成功')
    }
    
    showCreateDialog.value = false
    loadBridges()
  } catch (error) {
    console.error('提交失败:', error)
  } finally {
    submitting.value = false
  }
}

// 对话框关闭处理
const handleDialogClose = () => {
  editingBridge.value = null
  bridgeFormRef.value?.resetFields()
  bridgeForm.name = ''
  bridgeForm.source_channel_type = ''
  bridgeForm.target_channel_id = null as any
  bridgeForm.ticket = ''
  bridgeForm.is_active = true
}

onMounted(() => {
  loadBridges()
  loadChannels()
  loadAdapters()
})
</script>

<style scoped>
.bridges-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-header h2 {
  margin: 0;
  color: #303133;
}

.table-card {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.bridge-detail {
  padding: 16px 0;
}

.api-url-container,
.token-container {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.api-url-text,
.token-text {
  font-family: 'Courier New', monospace;
  background-color: #f5f7fa;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  color: #606266;
  word-break: break-all;
  flex: 1;
  min-width: 200px;
}

.ticket-input-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.ticket-display {
  font-family: 'Courier New', monospace;
  font-size: 12px;
  color: #606266;
  background-color: #f5f7fa;
  padding: 2px 6px;
  border-radius: 3px;
}
</style>
