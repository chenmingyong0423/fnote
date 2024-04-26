<template>
  <div class="flex items-center justify-center h-screen bg-light-blue">
    <div class="w-96 p-8 bg-white rounded-lg shadow-md modal">
      <div class="flex justify-center mb-6">
        <img src="@/assets/logo.png" alt="Logo" class="w-80 h-40" />
      </div>
      <div class="flex justify-center mb-6">
        {{ stepInfo[step - 1] }}
      </div>
      <a-form
        ref="formRef"
        :model="formState"
        name="normal_login"
        class="login-form"
        @finish="submit"
        @finishFailed="onFinishFailed"
        :labelCol="{ span: labelCols[step - 1] }"
      >
        <div v-show="step == 1">
          <a-form-item
            label=""
            :name="['admin', 'username']"
            :rules="[{ required: true, message: '请输入用户名！' }]"
          >
            <a-input v-model:value="formState.admin.username">
              <template #prefix>
                <UserOutlined class="site-form-item-icon" />
              </template>
            </a-input>
          </a-form-item>

          <a-form-item
            label=""
            :name="['admin', 'password']"
            :rules="[{ required: true, message: '请输入密码！' }]"
          >
            <a-input-password v-model:value="formState.admin.password">
              <template #prefix>
                <LockOutlined class="site-form-item-icon" />
              </template>
            </a-input-password>
          </a-form-item>

          <a-form-item>
            <a-button
              type="primary"
              class="login-form-button w-full h-10"
              @click="
                validate(
                  [
                    ['admin', 'username'],
                    ['admin', 'password']
                  ],
                  () => {
                    step++
                  }
                )
              "
            >
              下一步
            </a-button>
          </a-form-item>
        </div>
        <div v-show="step == 2">
          <a-form-item
            label="站点名称"
            name="website_name"
            :rules="[{ required: true, message: '请输入站点昵称' }]"
          >
            <a-input v-model:value="formState.website_name"></a-input>
          </a-form-item>

          <a-form-item
            label="站长昵称"
            name="website_owner"
            :rules="[{ required: true, message: '请输入站长昵称' }]"
          >
            <a-input v-model:value="formState.website_owner"></a-input>
          </a-form-item>
          <a-form-item
            label="站长简介"
            name="website_owner_profile"
            :rules="[{ required: true, message: '请输入站长简介' }]"
          >
            <a-textarea
              v-model:value="formState.website_owner_profile"
              :auto-size="{ minRows: 2, maxRows: 3 }"
            >
            </a-textarea>
          </a-form-item>

          <a-form-item>
            <a-button type="primary" class="login-form-button w-40% h-10" @click="step--">
              上一步
            </a-button>
            <a-button
              type="primary"
              class="login-form-button w-40% h-10 float-right"
              @click="
                validate(
                  ['website_name', 'website_owner', 'website_owner_email', 'website_owner_profile'],
                  () => {
                    step++
                  }
                )
              "
            >
              下一步
            </a-button>
          </a-form-item>
        </div>
        <div v-show="step == 3">
          <a-form-item
            label="站点 Logo"
            name="website_icon"
            :rules="[{ required: true, message: '请输入上传站点 logo' }]"
          >
            <StaticUpload
              :image-url="formState.website_icon"
              @update:imageUrl="(value) => handleWebsiteIcon(value)"
            />
          </a-form-item>

          <a-form-item
            label="站长头像"
            name="website_owner_avatar"
            :rules="[{ required: true, message: '请输入上传站长头像' }]"
          >
            <StaticUpload
              :image-url="formState.website_owner_avatar"
              @update:imageUrl="(value) => handleOwnerAvatar(value)"
            />
          </a-form-item>

          <a-form-item>
            <a-button type="primary" class="login-form-button w-40% h-10" @click="step--">
              上一步
            </a-button>
            <a-button
              type="primary"
              class="login-form-button w-40% h-10 float-right"
              @click="
                validate(
                  [
                    'website_icon',
                    'website_owner_avatar',
                    'website_owner_email',
                    'website_owner_profile'
                  ],
                  () => {
                    step++
                  }
                )
              "
            >
              下一步
            </a-button>
          </a-form-item>
        </div>

        <div v-show="step == 4">
          <a-form-item
            label="host"
            :name="['email_server', 'host']"
            :rules="[{ required: true, message: '请输入邮件服务的 host' }]"
          >
            <a-input v-model:value="formState.email_server.host" />
          </a-form-item>

          <a-form-item
            label="port"
            :name="['email_server', 'port']"
            :rules="[{ required: true, message: '请输入邮件服务的 port' }]"
          >
            <a-input v-model:value="formState.email_server.port" />
          </a-form-item>

          <a-form-item
            label="username"
            :name="['email_server', 'username']"
            :rules="[{ required: true, message: '请输入邮件服务的 username' }]"
          >
            <a-input v-model:value="formState.email_server.username" />
          </a-form-item>

          <a-form-item
            label="password"
            :name="['email_server', 'password']"
            :rules="[{ required: true, message: '请输入邮件服务的 password' }]"
          >
            <a-input v-model:value="formState.email_server.password" />
          </a-form-item>

          <a-form-item
            label="email"
            :name="['email_server', 'email']"
            :rules="[{ required: true, message: '请输入接收邮件的 email' }]"
            tooltip="接收通知时使用"
          >
            <a-input v-model:value="formState.email_server.email" />
          </a-form-item>

          <a-form-item>
            <a-button type="primary" class="login-form-button w-40% h-10" @click="step--">
              上一步
            </a-button>
            <a-button
              type="primary"
              class="login-form-button w-40% h-10 float-right"
              @click="initWebsite"
            >
              初始化
            </a-button>
          </a-form-item>
        </div>
      </a-form>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { reactive, ref, type UnwrapRef, toRaw } from 'vue'
import { UserOutlined, LockOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { Init, type InitReq } from '@/interfaces/Config'
import router from '@/router'
import { useUserStore } from '@/stores/user'
import StaticUpload from '@/components/upload/StaticUpload.vue'
import type { NamePath } from 'ant-design-vue/es/form/interface'

document.title = '内容发布统计 - 后台管理'

const userStore = useUserStore()

const formState: UnwrapRef<InitReq> = reactive({
  website_name: '',
  website_icon: '',
  website_owner: '',
  website_owner_profile: '',
  website_owner_avatar: '',
  email_server: {
    host: '',
    port: 0,
    username: '',
    password: '',
    email: ''
  },
  admin: {
    username: '',
    password: ''
  }
})

const step = ref(1)

const stepInfo = ['管理员信息', '站点信息', '站点信息', '邮件服务器信息']
const labelCols = [0, 7, 7, 7]

const onFinishFailed = (errorInfo: any) => {
  console.log('Failed:', errorInfo)
  message.error('初始化失败')
}

const submit = async () => {
  message.success('初始化成功')
}

const formRef = ref()
const handleWebsiteIcon = (value: string) => {
  formState.website_icon = value
  // 现在手动触发验证
  formRef.value?.validateFields(['website_icon'])
}
const handleOwnerAvatar = (value: string) => {
  formState.website_owner_avatar = value
  // 现在手动触发验证
  formRef.value?.validateFields(['website_owner_avatar'])
}

const validate = (fields: NamePath[] | string, callback: () => void) => {
  formRef.value?.validateFields(fields).then(() => callback())
}

const initWebsite = () => {
  formRef.value
    .validate()
    .then(async () => {
      formState.email_server.port = Number(formState.email_server.port)
      const response: any = await Init(toRaw(formState))
      if (response.data.code !== 0) {
        message.error(response.data.message)
        return
      }
      userStore.initialization = true
      await router.push('/login')
    })
    .catch((error: any) => {
      console.log('error', error)
    })
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
