<template>
  <div class="sidebar-container">
    <!-- Logo区域 -->
    <div class="logo-container">
      <div class="logo">
        <el-icon v-if="!appStore.sidebarCollapsed" size="24" color="#409EFF">
          <Message />
        </el-icon>
        <span v-if="!appStore.sidebarCollapsed" class="logo-text">MsgPilot</span>
        <el-icon v-else size="24" color="#409EFF">
          <Message />
        </el-icon>
      </div>
    </div>
    
    <!-- 菜单 -->
    <el-menu
      :default-active="activeMenu"
      :collapse="appStore.sidebarCollapsed"
      :unique-opened="true"
      background-color="#304156"
      text-color="#bfcbd9"
      active-text-color="#409EFF"
      router
      class="sidebar-menu"
    >
      <template v-for="route in menuRoutes" :key="route.path">
        <el-menu-item :index="route.path">
          <el-icon>
            <component :is="route.meta?.icon" />
          </el-icon>
          <template #title>{{ route.meta?.title }}</template>
        </el-menu-item>
      </template>
    </el-menu>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { Message } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const appStore = useAppStore()

// 当前激活的菜单
const activeMenu = computed(() => route.path)

// 菜单路由
const menuRoutes = computed(() => {
  const routes = router.getRoutes()
  const layoutRoute = routes.find(r => r.name === 'Layout')
  
  if (layoutRoute?.children) {
    return layoutRoute.children.filter(child => 
      child.meta?.title && child.path !== '/profile'
    )
  }
  
  return []
})
</script>

<style scoped>
.sidebar-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.logo-container {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-bottom: 1px solid #434a50;
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
}

.logo-text {
  font-size: 18px;
  font-weight: bold;
  color: #409EFF;
}

.sidebar-menu {
  flex: 1;
  border: none;
}

.sidebar-menu:not(.el-menu--collapse) {
  width: 200px;
}
</style>
