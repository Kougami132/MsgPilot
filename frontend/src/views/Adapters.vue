<template>
  <div class="page-container">
    <div class="adapters-page">
      <!-- 页面头部 -->
      <div class="page-header">
        <h2>适配器</h2>
      </div>
      
      <!-- 适配器列表 -->
      <div class="adapters-grid">
        <!-- 适配器卡片 -->
        <el-card class="adapter-card" header="可用适配器">
          <div class="adapter-list">
            <div
              v-for="adapter in adapters"
              :key="adapter"
              class="adapter-item"
            >
              <div class="adapter-info">
                <el-icon size="24" color="#409EFF">
                  <Grid />
                </el-icon>
                <span class="adapter-name">{{ adapter }}</span>
              </div>
              <el-tag type="success">可用</el-tag>
            </div>
          </div>
        </el-card>
        
        <!-- 处理器卡片 -->
        <el-card class="adapter-card" header="消息处理器">
          <div class="adapter-list">
            <div
              v-for="handler in handlers"
              :key="handler"
              class="adapter-item"
            >
              <div class="adapter-info">
                <el-icon size="24" color="#67C23A">
                  <Setting />
                </el-icon>
                <span class="adapter-name">{{ handler }}</span>
              </div>
              <el-tag type="success">活跃</el-tag>
            </div>
          </div>
        </el-card>
      </div>
      
      <!-- 适配器说明 -->
      <el-card class="info-card" header="适配器说明">
        <div class="adapter-descriptions">
          <div class="description-item">
            <h4>OneBot</h4>
            <p>支持QQ机器人消息推送，兼容OneBot v11协议</p>
          </div>
          <div class="description-item">
            <h4>Bark</h4>
            <p>iOS设备推送服务，支持自定义通知内容和样式</p>
          </div>
          <div class="description-item">
            <h4>Gotify</h4>
            <p>开源的服务器推送服务，支持多平台客户端</p>
          </div>
          <div class="description-item">
            <h4>PushDeer</h4>
            <p>轻量级的推送服务，支持微信小程序和App</p>
          </div>
          <div class="description-item">
            <h4>Ntfy</h4>
            <p>简单的HTTP推送服务，无需注册即可使用</p>
          </div>
          <div class="description-item">
            <h4>Webhook</h4>
            <p>通用的HTTP回调接口，支持自定义请求格式</p>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Grid, Setting } from '@element-plus/icons-vue'
import { adapterApi } from '@/api/adapter'

// 响应式数据
const adapters = ref<string[]>([])
const handlers = ref<string[]>([])
const loading = ref(false)

// 加载适配器列表
const loadAdapters = async () => {
  loading.value = true
  try {
    const data = await adapterApi.getList()
    adapters.value = data.adapters || []
    handlers.value = data.handlers || []
  } catch (error) {
    console.error('加载适配器列表失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadAdapters()
})
</script>

<style scoped>
.adapters-page {
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

.adapters-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 20px;
}

.adapter-card,
.info-card {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.adapter-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.adapter-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background-color: #f5f7fa;
  border-radius: 6px;
}

.adapter-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.adapter-name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.adapter-descriptions {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.description-item h4 {
  margin: 0 0 8px 0;
  color: #303133;
  font-size: 16px;
}

.description-item p {
  margin: 0;
  color: #606266;
  font-size: 14px;
  line-height: 1.5;
}
</style>
