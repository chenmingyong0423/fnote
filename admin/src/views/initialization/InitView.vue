<template>
  <div class="flex items-center justify-center h-screen bg-light-blue">
    <div class="w-96 p-8 bg-white rounded-lg shadow-md modal">
      <div class="flex justify-center mb-6">
        <img src="@/assets/logo.png" alt="Logo" class="w-60 h-30" />
      </div>
      <a-form
        :model="initReq"
        name="normal_login"
        class="login-form"
        @finish="submit"
        @finishFailed="onFinishFailed"
      >
        <div v-show="step == 1">
          <a-form-item
            label="站点名称"
            name="blog_name"
            :rules="[{ required: true, message: '请输入站点昵称！' }]"
          >
            <a-input v-model:value="initReq.blog_name">
            </a-input>
          </a-form-item>

          <a-form-item
            label="站长昵称"
            name="nickname"
            :rules="[{ required: true, message: '请输入站长昵称！' }]"
          >
            <a-input v-model:value="initReq.nickname">
            </a-input>
          </a-form-item>
          <a-form-item
            label="站长简介"
            name="nickname"
            :rules="[{ required: true, message: '请输入站长简介！' }]"
          >
            <a-input v-model:value="initReq.profile">
            </a-input>
          </a-form-item>

          <a-form-item
            label="站点 Logo"
            name="icon"
            :rules="[{ required: true, message: '请输入上传站点 logo！' }]"
          >
            <AvatarUpload :image-url="initReq.icon"/>
          </a-form-item>

          <a-form-item>
            <a-button
              type="primary"
              class="login-form-button w-full h-10"
              @click="step++"
            >
              下一步
            </a-button>
          </a-form-item>
        </div>
        <div v-show="step == 2">
          <a-form-item
            label="seo - 站点描述"
            name="blog_name"
            :rules="[{ required: true, message: '请输入站点昵称！' }]"
          >
            <a-textarea v-model:value="initReq.blog_name">
            </a-textarea>
          </a-form-item>

          <a-form-item
            label="站长昵称"
            name="nickname"
            :rules="[{ required: true, message: '请输入站长昵称！' }]"
          >
            <a-input v-model:value="initReq.nickname">
            </a-input>
          </a-form-item>
          <a-form-item
            label="站长简介"
            name="nickname"
            :rules="[{ required: true, message: '请输入站长简介！' }]"
          >
            <a-input v-model:value="initReq.profile">
            </a-input>
          </a-form-item>

          <a-form-item
            label="站点 Logo"
            name="icon"
            :rules="[{ required: true, message: '请输入上传站点 logo！' }]"
          >
            <AvatarUpload :image-url="initReq.icon"/>
          </a-form-item>

          <a-form-item>
            <a-button
              type="primary"
              class="login-form-button w-40% h-10"
              @click="step--"
            >
              上一步
            </a-button>
            <a-button
              type="primary"
              class="login-form-button w-40% h-10 float-right"
              @click="step++"
            >
              下一步
            </a-button>
          </a-form-item>
        </div>
<!--        <a-form-item-->
<!--          label=""-->
<!--          name="username"-->
<!--          :rules="[{ required: true, message: '请输入用户名！' }]"-->
<!--        >-->
<!--          <a-input v-model:value="loginReq.username">-->
<!--            <template #prefix>-->
<!--              <UserOutlined class="site-form-item-icon" />-->
<!--            </template>-->
<!--          </a-input>-->
<!--        </a-form-item>-->

<!--        <a-form-item-->
<!--          label=""-->
<!--          name="password"-->
<!--          :rules="[{ required: true, message: '请输入密码！' }]"-->
<!--        >-->
<!--          <a-input-password v-model:value="loginReq.password">-->
<!--            <template #prefix>-->
<!--              <LockOutlined class="site-form-item-icon" />-->
<!--            </template>-->
<!--          </a-input-password>-->
<!--        </a-form-item>-->

<!--        <a-form-item>-->
<!--          <a-button-->
<!--            :disabled="disabled"-->
<!--            type="primary"-->
<!--            html-type="submit"-->
<!--            class="login-form-button w-full h-10"-->
<!--          >-->
<!--            登录-->
<!--          </a-button>-->
<!--        </a-form-item>-->
      </a-form>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { reactive, computed, ref } from 'vue'
import { UserOutlined, LockOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import type { InitReq } from '@/interfaces/Config'
import router from '@/router'
import { useUserStore } from '@/stores/user'
import AvatarUpload from '@/components/upload/AvatarUpload.vue'

const userStore = useUserStore()

const initReq = reactive<InitReq>({
  blog_name: "",
  nickname: "",
  profile: "",
  icon: ""
})

const step = ref(1)

const onFinishFailed = (errorInfo: any) => {
  console.log('Failed:', errorInfo)
  message.error('初始化失败')
}

const submit = async () => {
  message.success('初始化成功')
}

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
