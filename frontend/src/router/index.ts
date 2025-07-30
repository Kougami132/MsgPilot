import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/Login.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/',
      name: 'Layout',
      component: () => import('@/layouts/MainLayout.vue'),
      meta: { requiresAuth: true },
      redirect: '/dashboard',
      children: [
        {
          path: '/dashboard',
          name: 'Dashboard',
          component: () => import('@/views/Dashboard.vue'),
          meta: { title: '仪表板', icon: 'DataBoard' }
        },
        {
          path: '/messages',
          name: 'Messages',
          component: () => import('@/views/Messages.vue'),
          meta: { title: '消息管理', icon: 'ChatDotRound' }
        },
        {
          path: '/channels',
          name: 'Channels',
          component: () => import('@/views/Channels.vue'),
          meta: { title: '渠道管理', icon: 'Connection' }
        },
        {
          path: '/bridges',
          name: 'Bridges',
          component: () => import('@/views/Bridges.vue'),
          meta: { title: '中转配置', icon: 'Share' }
        },
        {
          path: '/adapters',
          name: 'Adapters',
          component: () => import('@/views/Adapters.vue'),
          meta: { title: '适配器', icon: 'Grid' }
        },
        {
          path: '/profile',
          name: 'Profile',
          component: () => import('@/views/Profile.vue'),
          meta: { title: '个人中心', icon: 'User' }
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFound.vue')
    }
  ]
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth !== false && !authStore.isAuthenticated) {
    next('/login')
  } else if (to.path === '/login' && authStore.isAuthenticated) {
    next('/')
  } else {
    next()
  }
})

export default router
