<template>
  <div class="flex items-center justify-center h-screen bg-light-blue">
    <div class="w-96 p-8 bg-white rounded-lg shadow-md modal">
      <div class="flex justify-center mb-6">
        <img src="@/assets/logo.png" alt="Logo" class="w-60 h-30" />
      </div>
      <a-form
        ref="formRef"
        :model="initReq"
        name="normal_login"
        class="login-form"
        @finish="submit"
        @finishFailed="onFinishFailed"
        labelAlign="left"
      >
        <div v-show="step == 1">
          <a-form-item
            label="站点名称"
            name="website_name"
            :rules="[{ required: true, message: '请输入站点昵称！' }]"
          >
            <a-input v-model:value="initReq.website_name">
            </a-input>
          </a-form-item>

          <a-form-item
            label="站长昵称"
            name="website_owner"
            :rules="[{ required: true, message: '请输入站长昵称！' }]"
          >
            <a-input v-model:value="initReq.website_owner">
            </a-input>
          </a-form-item>
          <a-form-item
            label="域名"
            name="website_domain"
            :rules="[{ required: true, message: '请输入域名！' }]"
          >
            <a-input v-model:value="initReq.website_domain">
            </a-input>
          </a-form-item>
          <a-form-item
            label="站长简介"
            name="website_owner_profile"
            :rules="[{ required: true, message: '请输入站长简介！' }]"
          >
            <a-textarea v-model:value="initReq.website_owner_profile" :auto-size="{ minRows: 2, maxRows: 3}">
            </a-textarea>
          </a-form-item>

          <a-form-item
            label="站点Logo"
            name="website_icon"
            :rules="[{ required: true, message: '请输入上传站点 logo！' }]"
          >
            <AvatarUpload :image-url="initReq.website_icon" @update:imageUrl="value => handleWebsiteIcon(value)"/>
          </a-form-item>

          <a-form-item
            label="站长头像"
            name="website_owner_avatar"
            :rules="[{ required: true, message: '请输入上传站长头像！' }]"
          >
            <AvatarUpload :imageUrl="initReq.website_owner_avatar" @update:imageUrl="value => handleOwnerAvatar(value)"/>
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
            label="seo - 站点标题"
            :name="['seo', 'title']"
            :rules="[{ required: true, message: '请输入 seo - 站点标题！' }]"
          >
            <a-textarea v-model:value="initReq.seo.title">
            </a-textarea>
          </a-form-item>

          <a-form-item
            label="seo- 站点描述"
            :name="['seo', 'description']"
            :rules="[{ required: true, message: '请输入 seo- 站点描述！' }]"
          >
            <a-input v-model:value="initReq.seo.description">
            </a-input>
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
  website_name: "",
  website_icon: "",
  website_owner: "",
  website_owner_profile: "",
  website_owner_avatar: "",
  website_domain: "",
  seo: {
    title: "",
    description: "",
    og_title: "",
    og_image: "",
    baidu_site_verification: "",
    keywords: "",
    author: "",
    robots: ""
  }
})

const step = ref(1)

const onFinishFailed = (errorInfo: any) => {
  console.log('Failed:', errorInfo)
  message.error('初始化失败')
}

const submit = async () => {
  message.success('初始化成功')
}

const formRef = ref();
const handleWebsiteIcon = (value:string) => {
  initReq.website_icon = value;
  // 现在手动触发验证
  formRef.value?.validateFields(['website_icon']);
};
const handleOwnerAvatar = (value:string) => {
  initReq.website_owner_avatar = value;
  // 现在手动触发验证
  formRef.value?.validateFields(['website_owner_avatar']);
};


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
