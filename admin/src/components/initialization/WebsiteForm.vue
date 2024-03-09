<template>
  <div>
    <a-form :model="form" name="normal_login" class="login-form">
      <a-form-item label="" name="username">
        <a-input v-model:value="form.username">
          <template #prefix>
            <UserOutlined class="site-form-item-icon" />
          </template>
        </a-input>
      </a-form-item>

      <a-form-item label="" name="password" :rules="[{ required: true, message: '请输入密码！' }]">
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
</template>

<script setup lang="ts">
import { LockOutlined, UserOutlined } from '@ant-design/icons-vue'

const props = defineProps<{
  form: {
    type: object as PropType<>,
  }
}>()
</script>
