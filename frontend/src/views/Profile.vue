<template>
  <div class="page-container">
    <div class="profile-page">
      <!-- 页面头部 -->
      <div class="page-header">
        <h2>个人中心</h2>
      </div>
      
      <!-- 用户信息 -->
      <el-card class="profile-card" header="用户信息">
        <div class="user-info">
          <div class="avatar-section">
            <el-avatar :size="80">
              <el-icon><User /></el-icon>
            </el-avatar>
            <el-button type="primary" size="small" style="margin-top: 12px;">
              更换头像
            </el-button>
          </div>
          <div class="info-section">
            <el-descriptions :column="1" border>
              <el-descriptions-item label="用户ID">
                {{ authStore.user?.id || 1 }}
              </el-descriptions-item>
              <el-descriptions-item label="用户名">
                {{ authStore.user?.username || 'admin' }}
              </el-descriptions-item>
              <el-descriptions-item label="创建时间">
                {{ formatDate(authStore.user?.created_at || new Date().toISOString()) }}
              </el-descriptions-item>
            </el-descriptions>
          </div>
        </div>
      </el-card>
      
      <!-- 修改密码 -->
      <el-card class="password-card" header="修改密码">
        <el-form
          ref="passwordFormRef"
          :model="passwordForm"
          :rules="passwordRules"
          label-width="100px"
          style="max-width: 400px;"
        >
          <el-form-item label="当前密码" prop="old_password">
            <el-input
              v-model="passwordForm.old_password"
              type="password"
              placeholder="请输入当前密码"
              show-password
            />
          </el-form-item>
          <el-form-item label="新密码" prop="new_password">
            <el-input
              v-model="passwordForm.new_password"
              type="password"
              placeholder="请输入新密码"
              show-password
            />
          </el-form-item>
          <el-form-item label="确认密码" prop="confirm_password">
            <el-input
              v-model="passwordForm.confirm_password"
              type="password"
              placeholder="请再次输入新密码"
              show-password
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="submitting" @click="handleChangePassword">
              修改密码
            </el-button>
            <el-button @click="resetPasswordForm">重置</el-button>
          </el-form-item>
        </el-form>
      </el-card>
      

    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted, type FormInstance, type FormRules } from 'vue'
import { ElMessage } from 'element-plus'
import { User } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'
import { authApi } from '@/api/auth'
import type { ChangePasswordRequest } from '@/types/auth'

const authStore = useAuthStore()
const appStore = useAppStore()

// 表单引用
const passwordFormRef = ref<FormInstance>()

// 提交状态
const submitting = ref(false)

// 密码表单
const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

// 密码验证规则
const passwordRules: FormRules = {
  old_password: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.new_password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 格式化日期
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

// 修改密码
const handleChangePassword = async () => {
  if (!passwordFormRef.value) return
  
  try {
    await passwordFormRef.value.validate()
    submitting.value = true
    
    const changePasswordData: ChangePasswordRequest = {
      username: authStore.user?.username || 'admin',
      old_password: passwordForm.old_password,
      new_password: passwordForm.new_password
    }
    
    await authApi.changePassword(changePasswordData)
    ElMessage.success('密码修改成功')
    resetPasswordForm()
  } catch (error) {
    console.error('修改密码失败:', error)
  } finally {
    submitting.value = false
  }
}

// 重置密码表单
const resetPasswordForm = () => {
  passwordFormRef.value?.resetFields()
  passwordForm.old_password = ''
  passwordForm.new_password = ''
  passwordForm.confirm_password = ''
}

// 页面加载时获取用户信息
onMounted(async () => {
  try {
    await authStore.fetchUserInfo()
  } catch (error) {
    console.error('获取用户信息失败:', error)
  }
})
</script>

<style scoped>
.profile-page {
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

.profile-card,
.password-card {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.user-info {
  display: flex;
  gap: 40px;
  align-items: flex-start;
}

.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.info-section {
  flex: 1;
}
</style>
