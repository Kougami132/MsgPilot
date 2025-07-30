<template>
  <div class="page-container">
    <div class="channels-page">
      <!-- 页面头部 -->
      <div class="page-header">
        <h2>渠道管理</h2>
        <el-button type="primary" @click="showCreateDialog = true">
          <el-icon><Plus /></el-icon>
          创建渠道
        </el-button>
      </div>

      <!-- 渠道列表 -->
      <el-card class="table-card">
        <el-table :data="channels" :loading="loading" style="width: 100%">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="name" label="名称" min-width="150" />
          <el-table-column prop="type" label="类型" width="120" />
          <el-table-column prop="created_at" label="创建时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="280" fixed="right">
            <template #default="{ row }">
              <el-button-group>
                <el-button type="primary" size="small" @click="handleView(row)">
                  查看
                </el-button>
                <el-button type="success" size="small" @click="handleTest(row)">
                  测试
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
      :title="editingChannel ? '编辑渠道' : '创建渠道'"
      width="600px"
      @closed="handleDialogClose"
    >
      <el-form
        ref="channelFormRef"
        :model="channelForm"
        :rules="channelRules"
        label-width="80px"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="channelForm.name" placeholder="请输入渠道名称" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select
            v-model="channelForm.type"
            placeholder="请选择渠道类型"
            style="width: 100%"
            @change="handleTypeChange"
          >
            <el-option
              v-for="adapter in adapters"
              :key="adapter"
              :label="adapter"
              :value="adapter"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="配置" prop="config">
          <el-input
            v-model="configJson"
            type="textarea"
            :rows="6"
            placeholder="请输入JSON格式的配置"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          {{ editingChannel ? '更新' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
    
    <!-- 查看详情对话框 -->
    <el-dialog v-model="showViewDialog" title="渠道详情" width="600px">
      <div v-if="viewingChannel" class="channel-detail">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="ID">{{ viewingChannel.id }}</el-descriptions-item>
          <el-descriptions-item label="名称">{{ viewingChannel.name }}</el-descriptions-item>
          <el-descriptions-item label="类型">{{ viewingChannel.type }}</el-descriptions-item>
          <el-descriptions-item label="配置">
            <pre>{{ JSON.stringify(viewingChannel.config, null, 2) }}</pre>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">
            {{ formatDate(viewingChannel.created_at) }}
          </el-descriptions-item>
          <el-descriptions-item label="更新时间">
            {{ formatDate(viewingChannel.updated_at) }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, type FormInstance, type FormRules } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { channelApi } from '@/api/channel'
import { adapterApi } from '@/api/adapter'
import type { Channel, CreateChannelRequest } from '@/types/api'

// 响应式数据
const loading = ref(false)
const submitting = ref(false)
const channels = ref<Channel[]>([])
const adapters = ref<string[]>([])

// 对话框状态
const showCreateDialog = ref(false)
const showViewDialog = ref(false)
const editingChannel = ref<Channel | null>(null)
const viewingChannel = ref<Channel | null>(null)

// 表单引用
const channelFormRef = ref<FormInstance>()

// 渠道表单
const channelForm = reactive<CreateChannelRequest>({
  name: '',
  type: '',
  config: {}
})

// 配置JSON字符串
const configJson = ref('')

// 表单验证规则
const channelRules: FormRules = {
  name: [
    { required: true, message: '请输入渠道名称', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择渠道类型', trigger: 'change' }
  ]
}

// 格式化日期
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

// 加载渠道列表
const loadChannels = async () => {
  loading.value = true
  try {
    channels.value = await channelApi.getList()
  } catch (error) {
    console.error('加载渠道列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 加载适配器列表
const loadAdapters = async () => {
  try {
    const response = await adapterApi.getList()
    adapters.value = response.adapters
  } catch (error) {
    console.error('加载适配器列表失败:', error)
  }
}

// 查看渠道
const handleView = (channel: Channel) => {
  viewingChannel.value = channel
  showViewDialog.value = true
}

// 测试渠道
const handleTest = async (channel: Channel) => {
  try {
    await channelApi.testPush(channel)
    ElMessage.success('测试推送成功')
  } catch (error) {
    console.error('测试推送失败:', error)
  }
}

// 编辑渠道
const handleEdit = (channel: Channel) => {
  editingChannel.value = channel
  channelForm.name = channel.name
  channelForm.type = channel.type
  channelForm.config = channel.config
  configJson.value = JSON.stringify(channel.config, null, 2)
  showCreateDialog.value = true
}

// 删除渠道
const handleDelete = async (channel: Channel) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除渠道 "${channel.name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    await channelApi.delete(channel.id)
    ElMessage.success('删除成功')
    loadChannels()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除渠道失败:', error)
    }
  }
}

// 类型变化处理
const handleTypeChange = () => {
  // 根据类型设置默认配置
  const defaultConfigs: Record<string, any> = {
    email: {
      smtp_host: '',
      smtp_port: 587,
      username: '',
      password: '',
      from: ''
    },
    wechat: {
      app_id: '',
      app_secret: '',
      template_id: ''
    },
    webhook: {
      url: '',
      method: 'POST',
      headers: {}
    }
  }

  const defaultConfig = defaultConfigs[channelForm.type] || {}
  configJson.value = JSON.stringify(defaultConfig, null, 2)
}

// 提交表单
const handleSubmit = async () => {
  if (!channelFormRef.value) return

  try {
    await channelFormRef.value.validate()

    // 解析配置JSON
    try {
      channelForm.config = JSON.parse(configJson.value || '{}')
    } catch (error) {
      ElMessage.error('配置JSON格式错误')
      return
    }

    submitting.value = true

    if (editingChannel.value) {
      await channelApi.update(editingChannel.value.id, channelForm)
      ElMessage.success('更新成功')
    } else {
      await channelApi.create(channelForm)
      ElMessage.success('创建成功')
    }

    showCreateDialog.value = false
    loadChannels()
  } catch (error) {
    console.error('提交失败:', error)
  } finally {
    submitting.value = false
  }
}

// 对话框关闭处理
const handleDialogClose = () => {
  editingChannel.value = null
  channelFormRef.value?.resetFields()
  channelForm.name = ''
  channelForm.type = ''
  channelForm.config = {}
  configJson.value = ''
}

onMounted(() => {
  loadChannels()
  loadAdapters()
})
</script>

<style scoped>
.channels-page {
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

.channel-detail {
  padding: 16px 0;
}

.channel-detail pre {
  background-color: #f5f7fa;
  padding: 12px;
  border-radius: 4px;
  font-size: 12px;
  line-height: 1.4;
  max-height: 200px;
  overflow-y: auto;
}
</style>
