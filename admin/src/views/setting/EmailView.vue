<template>
  <div>
    <a-descriptions title="邮件配置" :column="1" bordered>
      <a-descriptions-item label="host - 邮件服务器主机名">
        <div>
          <a-input v-if="editable" v-model:value="data.host" style="margin: -5px 0" />
          <template v-else>
            {{ data.host }}
          </template>
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="port - 邮件服务器端口号">
        <div>
          <a-input
            v-if="editable"
            v-model:value="data.port"
            style="margin: -5px 0"
            @change="portChanged"
          />
          <template v-else>
            {{ data.port }}
          </template>
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="username - 授权邮箱账号的用户名">
        <div>
          <a-input v-if="editable" v-model:value="data.username" style="margin: -5px 0" />
          <template v-else>
            {{ data.username }}
          </template>
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="password - 授权密码">
        <div>
          <a-input v-if="editable" v-model:value="data.password" style="margin: -5px 0" />
          <template v-else>
            {{ data.password }}
          </template>
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="email - 邮箱地址（用于接收邮件）">
        <div>
          <a-input v-if="editable" v-model:value="data.email" style="margin: -5px 0" />
          <template v-else>
            {{ data.email }}
          </template>
        </div>
      </a-descriptions-item>
    </a-descriptions>
    <div style="margin-top: 10px">
      <a-button v-if="!editable" type="primary" @click="editable = true">编辑</a-button>
      <div v-else>
        <a-button type="primary" @click="cancel" style="margin-right: 5px">取消</a-button>
        <a-button type="primary" @click="save">保存</a-button>
      </div>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { type EmailConfig, GetEmail, UpdateEmail } from '@/interfaces/Config'
import { ref } from 'vue'
import { message } from 'ant-design-vue'

const editable = ref<boolean>(false)

const data = ref<EmailConfig>({
  host: '',
  port: 0,
  username: '',
  password: '',
  email: ''
})

const getEmail = async () => {
  try {
    const response: any = await GetEmail()
    data.value = response.data || {}
  } catch (error) {
    console.log(error)
  }
}
getEmail()

const cancel = () => {
  editable.value = false
  getEmail()
}

const save = async () => {
  try {
    const response: any = await UpdateEmail({
      host: data.value.host,
      port: data.value.port,
      username: data.value.username,
      password: data.value.password,
      email: data.value.email
    })
    if (response.code === 0) {
      message.success('保存成功')
      await getEmail()
      editable.value = false
    } else {
      message.error(response.message)
    }
  } catch (error) {
    console.log(error)
  }
}

const portChanged = (e: any) => {
  data.value.port = Number(e.target.value)
}
</script>

<style scoped></style>
