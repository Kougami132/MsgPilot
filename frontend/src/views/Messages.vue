<template>
  <div class="page-container">
    <div class="messages-page">
      <!-- 页面头部 -->
      <div class="page-header">
        <h2>消息管理</h2>
      </div>
      
      <!-- 搜索和筛选 -->
      <el-card class="search-card">
        <el-form :model="searchForm" inline>
          <el-form-item label="搜索" style="margin-bottom: 0;">
            <el-input
              v-model="searchForm.keyword"
              placeholder="请输入标题或内容关键词"
              clearable
              style="width: 250px"
            />
          </el-form-item>
          <el-form-item label="状态" style="margin-bottom: 0;">
            <el-select
              v-model="searchForm.status"
              placeholder="请选择状态"
              clearable
              style="width: 120px"
            >
              <el-option label="待发送" :value="0" />
              <el-option label="发送中" :value="1" />
              <el-option label="已发送" :value="2" />
              <el-option label="失败" :value="3" />
            </el-select>
          </el-form-item>
          <el-form-item style="margin-bottom: 0;">
            <el-button @click="handleReset">
              <el-icon><Refresh /></el-icon>
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>
      
      <!-- 消息列表 -->
      <el-card class="table-card">
        <el-table
          :data="filteredMessages"
          :loading="loading"
          style="width: 100%"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="title" label="标题" min-width="150" />
          <el-table-column prop="content" label="内容" min-width="200" show-overflow-tooltip />
          <el-table-column label="渠道" width="150">
            <template #default="{ row }">
              {{ getChannelName(row) }}
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="创建时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="160" fixed="right">
            <template #default="{ row }">
              <el-button-group>
                <el-button type="primary" size="small" @click="handleView(row)">
                  查看
                </el-button>
                <el-button type="danger" size="small" @click="handleDelete(row)">
                  删除
                </el-button>
              </el-button-group>
            </template>
          </el-table-column>
        </el-table>
        
        <!-- 批量操作 -->
        <div v-if="selectedMessages.length > 0" class="batch-actions">
          <span>已选择 {{ selectedMessages.length }} 项</span>
          <el-button type="danger" size="small" @click="handleBatchDelete">
            批量删除
          </el-button>
        </div>
      </el-card>
    </div>



    <!-- 查看详情对话框 -->
    <el-dialog v-model="showViewDialog" title="消息详情" width="600px">
      <div v-if="viewingMessage" class="message-detail">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="ID">{{ viewingMessage.id }}</el-descriptions-item>
          <el-descriptions-item label="标题">{{ viewingMessage.title }}</el-descriptions-item>
          <el-descriptions-item label="内容">{{ viewingMessage.content }}</el-descriptions-item>
          <el-descriptions-item label="渠道">{{ getChannelName(viewingMessage) }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(viewingMessage.status)">
              {{ getStatusText(viewingMessage.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item v-if="viewingMessage.error_message" label="错误信息">
            {{ viewingMessage.error_message }}
          </el-descriptions-item>
          <el-descriptions-item label="发送时间">
            {{ formatDate(viewingMessage.created_at) }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { messageApi } from '@/api/message'
import type { Message } from '@/types/api'

// 响应式数据
const loading = ref(false)
const messages = ref<Message[]>([])
const selectedMessages = ref<Message[]>([])

// 对话框状态
const showViewDialog = ref(false)
const viewingMessage = ref<Message | null>(null)

// 搜索表单
const searchForm = reactive({
  keyword: '',
  status: undefined as number | undefined
})

// 获取状态类型
const getStatusType = (status: number) => {
  switch (status) {
    case 0: // 待发送
      return 'info'
    case 1: // 发送中
      return 'warning'
    case 2: // 已发送
      return 'success'
    case 3: // 失败
      return 'danger'
    default:
      return 'info'
  }
}

// 获取状态文本
const getStatusText = (status: number) => {
  switch (status) {
    case 0:
      return '待发送'
    case 1:
      return '发送中'
    case 2:
      return '已发送'
    case 3:
      return '失败'
    default:
      return '未知'
  }
}

// 获取渠道名称
const getChannelName = (message: Message) => {
  if (message.bridge?.target_channel?.name) {
    return message.bridge.target_channel.name
  }
  return `Bridge-${message.bridge_id}`
}

// 格式化日期
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

// 加载消息列表
const loadMessages = async () => {
  loading.value = true
  try {
    messages.value = await messageApi.getList()
  } catch (error) {
    console.error('加载消息列表失败:', error)
  } finally {
    loading.value = false
  }
}



// 过滤后的消息列表
const filteredMessages = computed(() => {
  let result = messages.value

  // 按关键词筛选（匹配标题或内容）
  if (searchForm.keyword.trim()) {
    const keyword = searchForm.keyword.toLowerCase()
    result = result.filter(message =>
      message.title.toLowerCase().includes(keyword) ||
      message.content.toLowerCase().includes(keyword)
    )
  }

  // 按状态筛选
  if (searchForm.status !== undefined) {
    result = result.filter(message => message.status === searchForm.status)
  }

  // 按创建时间降序排列（最新的在前面）
  result = result.sort((a, b) =>
    new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  )

  return result
})

// 重置搜索
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = undefined
}

// 选择变化处理
const handleSelectionChange = (selection: Message[]) => {
  selectedMessages.value = selection
}

// 查看消息
const handleView = (message: Message) => {
  viewingMessage.value = message
  showViewDialog.value = true
}



// 删除消息
const handleDelete = async (message: Message) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除消息 "${message.title}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    await messageApi.delete(message.id)
    ElMessage.success('删除成功')
    loadMessages()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除消息失败:', error)
    }
  }
}

// 批量删除
const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedMessages.value.length} 条消息吗？`,
      '确认批量删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    await Promise.all(
      selectedMessages.value.map(message => messageApi.delete(message.id))
    )

    ElMessage.success('批量删除成功')
    selectedMessages.value = []
    loadMessages()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量删除失败:', error)
    }
  }
}





onMounted(() => {
  loadMessages()
})
</script>

<style scoped>
.messages-page {
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

.search-card,
.table-card {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.batch-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 16px;
  padding: 12px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.message-detail {
  padding: 16px 0;
}
</style>
