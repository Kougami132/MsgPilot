<template>
  <div class="login-container">
    <div class="login-box">
      <!-- Logo和标题 -->
      <div class="login-header">
        <div class="logo">
          <el-icon size="48" color="#409EFF">
            <Message />
          </el-icon>
        </div>
        <h1 class="title">MsgPilot</h1>
        <p class="subtitle">消息管理系统</p>
      </div>
      
      <!-- 登录表单 -->
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        size="large"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
            prefix-icon="User"
            clearable
          />
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            prefix-icon="Lock"
            show-password
            clearable
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            class="login-button"
            :loading="authStore.loading"
            @click="handleLogin"
          >
            {{ isRegisterMode ? '注册' : '登录' }}
          </el-button>
        </el-form-item>
        
        <div class="form-footer">
          <el-button
            type="text"
            @click="toggleMode"
            class="toggle-mode"
          >
            {{ isRegisterMode ? '已有账号？去登录' : '没有账号？去注册' }}
          </el-button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { Message } from '@element-plus/icons-vue'
import type { LoginRequest } from '@/types/auth'

const router = useRouter()
const authStore = useAuthStore()

// 表单引用
const loginFormRef = ref<FormInstance>()

// 是否为注册模式
const isRegisterMode = ref(false)

// 表单数据
const loginForm = reactive<LoginRequest>({
  username: '',
  password: ''
})

// 表单验证规则
const loginRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ]
}

// 切换登录/注册模式
const toggleMode = () => {
  isRegisterMode.value = !isRegisterMode.value
  // 清空表单
  loginFormRef.value?.resetFields()
}

// 处理登录/注册
const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    await loginFormRef.value.validate()
    
    if (isRegisterMode.value) {
      await authStore.register(loginForm)
      ElMessage.success('注册成功')
    } else {
      await authStore.login(loginForm)
      ElMessage.success('登录成功')
    }
    
    router.push('/')
  } catch (error) {
    console.error('登录/注册失败:', error)
  }
}

// 组件挂载时初始化认证状态
onMounted(async () => {
  await authStore.initialize()
})
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-box {
  width: 400px;
  padding: 40px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.logo {
  margin-bottom: 16px;
}

.title {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
  margin: 0 0 8px 0;
}

.subtitle {
  font-size: 14px;
  color: #909399;
  margin: 0;
}

.login-form {
  margin-top: 20px;
}

.login-button {
  width: 100%;
  height: 44px;
  font-size: 16px;
}

.form-footer {
  text-align: center;
  margin-top: 20px;
}

.toggle-mode {
  color: #409EFF;
  font-size: 14px;
}

.toggle-mode:hover {
  color: #66b1ff;
}
</style>
