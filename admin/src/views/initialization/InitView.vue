<template>
  <div class="theme-page flex items-center justify-center h-screen">
    <div class="theme-panel w-96 p-8 rounded-lg modal">
      <div class="flex justify-center mb-6">
        <img src="@/assets/logo.png" alt="Logo" class="w-80 h-40" />
      </div>
      <div class="flex justify-center mb-6">
        {{ stepInfo[step - 1] }}
      </div>
      <div class="mb-6">
        <a-button
          class="w-full h-10"
          :loading="importing"
          :disabled="initializing"
          @click="selectBackupFile"
        >
          <template #icon><UploadOutlined /></template>
          {{ importing ? '导入中...' : '从备份文件导入' }}
        </a-button>
      </div>
      <a-form
        ref="formRef"
        :model="formState"
        name="normal_login"
        class="login-form"
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
              :disabled="importing"
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
            <a-button
              type="primary"
              class="login-form-button w-40% h-10"
              :disabled="importing"
              @click="step--"
            >
              上一步
            </a-button>
            <a-button
              type="primary"
              class="login-form-button w-40% h-10 float-right"
              :disabled="importing"
              @click="
                validate(['website_name', 'website_owner', 'website_owner_profile'], () => {
                  step++
                })
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
            <a-button
              type="primary"
              class="login-form-button w-40% h-10"
              :disabled="importing || initializing"
              @click="step--"
            >
              上一步
            </a-button>
            <a-button
              type="primary"
              class="login-form-button w-40% h-10 float-right"
              :loading="initializing"
              :disabled="importing"
              @click="
                validate(['website_icon', 'website_owner_avatar'], () => {
                  initWebsite()
                })
              "
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
import { UserOutlined, LockOutlined, UploadOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { Init, isInit, type InitReq } from '@/interfaces/Config'
import router from '@/router'
import { useUserStore } from '@/stores/user'
import StaticUpload from '@/components/upload/StaticUpload.vue'
import type { NamePath } from 'ant-design-vue/es/form/interface'
import { Recovery } from '@/interfaces/Backup'
import { toErrorMessage } from '@/utils/error'

document.title = '内容发布统计 - 后台管理'

const userStore = useUserStore()

const formState: UnwrapRef<InitReq> = reactive({
  website_name: '',
  website_icon: '',
  website_owner: '',
  website_owner_profile: '',
  website_owner_avatar: '',
  admin: {
    username: '',
    password: ''
  }
})

const step = ref(1)
const importing = ref(false)
const initializing = ref(false)

const stepInfo = ['管理员信息', '站点信息', '站点信息']
const labelCols = [0, 7, 7]

const onFinishFailed = (errorInfo: any) => {
  console.log('Failed:', errorInfo)
  message.error('初始化失败')
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

const refreshInitStatus = async () => {
  const response: any = await isInit()
  if (response.data.code !== 0) {
    return false
  }
  const initStatus = Boolean(response.data.data?.initStatus)
  userStore.initialization = initStatus
  return initStatus
}

const selectBackupFile = () => {
  if (importing.value || initializing.value) return

  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.zip'
  input.onchange = async (event: Event) => {
    const target = event.target as HTMLInputElement
    const file = target.files?.[0]
    if (!file) return

    const formData = new FormData()
    formData.append('file', file)

    importing.value = true
    const hide = message.loading('正在导入备份，请稍候...', 0)
    try {
      await Recovery(formData)
      const initialized = await refreshInitStatus()
      if (!initialized) {
        message.warning('导入成功，但备份中未包含已初始化配置，请继续手动初始化')
        return
      }
      message.success('导入成功，初始化已完成')
      await router.replace('/login')
    } catch (error) {
      message.error(toErrorMessage(error, '导入失败，请稍后再试'))
    } finally {
      hide()
      importing.value = false
    }
  }
  input.click()
}

const initWebsite = () => {
  initializing.value = true
  formRef.value
    .validate()
    .then(async () => {
      const response: any = await Init(toRaw(formState))
      if (response.data.code !== 0) {
        message.error(response.data.message)
        return
      }
      userStore.initialization = true
      message.success('初始化成功')
      await router.replace('/login')
    })
    .catch((error: any) => {
      console.log('error', error)
    })
    .finally(() => {
      initializing.value = false
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
