import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  // 侧边栏状态
  const sidebarCollapsed = ref(false)
  
  // 主题模式
  const isDark = ref(false)
  
  // 页面加载状态
  const pageLoading = ref(false)

  // 切换侧边栏
  const toggleSidebar = () => {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  // 切换主题
  const toggleTheme = () => {
    isDark.value = !isDark.value
    // 这里可以添加主题切换逻辑
    document.documentElement.classList.toggle('dark', isDark.value)
  }

  // 设置页面加载状态
  const setPageLoading = (loading: boolean) => {
    pageLoading.value = loading
  }

  return {
    // 状态
    sidebarCollapsed: readonly(sidebarCollapsed),
    isDark: readonly(isDark),
    pageLoading: readonly(pageLoading),
    
    // 方法
    toggleSidebar,
    toggleTheme,
    setPageLoading
  }
})
