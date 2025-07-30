<template>
  <div class="page-container">
    <div class="dashboard">
      <!-- 统计卡片 -->
      <div class="stats-grid">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon size="32" color="#409EFF">
                <ChatDotRound />
              </el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.totalMessages }}</div>
              <div class="stat-label">总消息数</div>
            </div>
          </div>
        </el-card>
        
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon size="32" color="#67C23A">
                <Connection />
              </el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.totalChannels }}</div>
              <div class="stat-label">渠道数量</div>
            </div>
          </div>
        </el-card>
        
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon size="32" color="#E6A23C">
                <Share />
              </el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.totalBridges }}</div>
              <div class="stat-label">中转配置</div>
            </div>
          </div>
        </el-card>
        
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon size="32" color="#F56C6C">
                <Grid />
              </el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.totalAdapters }}</div>
              <div class="stat-label">适配器数量</div>
            </div>
          </div>
        </el-card>
      </div>
      
      <!-- 最近消息 -->
      <el-card class="recent-messages" header="最近消息">
        <el-table :data="recentMessages" style="width: 100%">
          <el-table-column prop="title" label="标题" />
          <el-table-column prop="content" label="内容" show-overflow-tooltip />
          <el-table-column prop="status" label="状态">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="创建时间">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
        </el-table>
      </el-card>
      
      <!-- 系统信息 -->
      <el-card class="system-info" header="系统信息">
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">系统版本:</span>
            <span class="info-value">v1.0.0</span>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ChatDotRound, Connection, Share, Grid } from '@element-plus/icons-vue'
import { messageApi } from '@/api/message'
import { channelApi } from '@/api/channel'
import { bridgeApi } from '@/api/bridge'
import { adapterApi } from '@/api/adapter'
import type { Message } from '@/types/api'

// 统计数据
const stats = ref({
  totalMessages: 0,
  totalChannels: 0,
  totalBridges: 0,
  totalAdapters: 0
})

// 最近消息
const recentMessages = ref<Message[]>([])

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

// 格式化日期
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

// 加载统计数据
const loadStats = async () => {
  try {
    const [messages, channels, bridges, adapters] = await Promise.all([
      messageApi.getList(),
      channelApi.getList(),
      bridgeApi.getList(),
      adapterApi.getList()
    ])
    
    stats.value = {
      totalMessages: messages.length,
      totalChannels: channels.length,
      totalBridges: bridges.length,
      totalAdapters: adapters.adapters.length + adapters.handlers.length
    }
    
    // 获取最近的5条消息
    recentMessages.value = messages.slice(0, 5)
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}



onMounted(() => {
  loadStats()
})
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
}

.stat-card {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  flex-shrink: 0;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
  line-height: 1;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 4px;
}

.recent-messages,
.system-info {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
}

.info-label {
  font-size: 14px;
  color: #606266;
}

.info-value {
  font-size: 14px;
  color: #303133;
  font-weight: 500;
}
</style>
