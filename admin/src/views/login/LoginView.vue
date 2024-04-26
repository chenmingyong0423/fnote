<template>
  <div class="flex items-center justify-center h-screen bg-light-blue">
    <div class="w-96 p-8 bg-white rounded-lg shadow-md modal">
      <div class="flex justify-center mb-6">
        <img src="@/assets/logo.png" alt="Logo" class="w-60 h-30" />
      </div>
      <a-form
        :model="loginReq"
        name="normal_login"
        class="login-form"
        @finish="login"
        @finishFailed="onFinishFailed"
      >
        <a-form-item
          label=""
          name="username"
          :rules="[{ required: true, message: '请输入用户名！' }]"
        >
          <a-input v-model:value="loginReq.username">
            <template #prefix>
              <UserOutlined class="site-form-item-icon" />
            </template>
          </a-input>
        </a-form-item>

        <a-form-item
          label=""
          name="password"
          :rules="[{ required: true, message: '请输入密码！' }]"
        >
          <a-input-password v-model:value="loginReq.password">
            <template #prefix>
              <LockOutlined class="site-form-item-icon" />
            </template>
          </a-input-password>
        </a-form-item>

        <a-form-item>
          <a-button
            :disabled="disabled"
            type="primary"
            html-type="submit"
            class="login-form-button w-full h-10"
          >
            登录
          </a-button>
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { reactive, computed } from 'vue'
import { UserOutlined, LockOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import type { LoginRequest } from '@/interfaces/User'
import router from '@/router'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

const loginReq = reactive<LoginRequest>({
  username: '',
  password: ''
})

const onFinishFailed = (errorInfo: any) => {
  console.log('Failed:', errorInfo)
  message.error('登录失败，请检查用户名和密码的格式是否正确')
}
const disabled = computed(() => {
  return !(loginReq.username && loginReq.password)
})

const login = async () => {
  const success = await userStore.loginIn(loginReq)
  if (success) {
    await router.push('/')
  }
}

document.title = '登录 - 后台管理'
</script>

<style scoped>
#components-form-demo-normal-login .login-form {
  max-width: 300px;
}

#components-form-demo-normal-login .login-form-forgot {
  float: right;
}

#components-form-demo-normal-login .login-form-button {
  width: 100%;
}
</style>
