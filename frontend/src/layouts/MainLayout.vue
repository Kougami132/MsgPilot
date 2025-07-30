<template>
  <el-container class="main-layout">
    <!-- 侧边栏 -->
    <el-aside :width="sidebarWidth" class="sidebar">
      <Sidebar />
    </el-aside>
    
    <!-- 主内容区域 -->
    <el-container>
      <!-- 头部 -->
      <el-header class="header">
        <Header />
      </el-header>
      
      <!-- 内容区域 -->
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAppStore } from '@/stores/app'
import Sidebar from '@/components/layout/Sidebar.vue'
import Header from '@/components/layout/Header.vue'

const appStore = useAppStore()

const sidebarWidth = computed(() => {
  return appStore.sidebarCollapsed ? '64px' : '200px'
})
</script>

<style scoped>
.main-layout {
  height: 100vh;
}

.sidebar {
  background-color: #304156;
  transition: width 0.3s;
}

.header {
  background-color: #fff;
  border-bottom: 1px solid #e4e7ed;
  padding: 0;
  height: 60px !important;
  line-height: 60px;
}

.main-content {
  background-color: #f0f2f5;
  padding: 0;
}
</style>
