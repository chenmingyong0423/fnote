<template>
  <div>
    <a-button type="primary" @click="visible = true">新增支付二维码</a-button>
    <a-modal
      v-model:open="visible"
      title="新增支付二维码"
      ok-text="提交"
      cancel-text="取消"
      @ok="addPay"
    >
      <a-form ref="formRef" :model="formState" layout="vertical" name="form_in_modal">
        <a-form-item name="name" label="名称" :rules="[{ required: true, message: '请输入名称' }]">
          <a-input v-model:value="formState.name" />
        </a-form-item>
        <a-form-item
          name="image"
          label="二维码"
          :rules="[{ required: true, message: '请上传二维码' }]"
        >
          <StaticUpload
            :image-url="formState.image"
            @update:imageUrl="(value) => (formState.image = value)"
            :authorization="userStore.token"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
  <a-table :columns="columns" :data-source="list">
    <template #bodyCell="{ column, record }">
      <template v-if="column.key === 'image'">
        <a-image :width="200" :src="serverHost + record.image" />
      </template>
      <template v-else-if="column.dataIndex === 'operation'">
        <a-popconfirm v-if="list.length" title="确认删除？" @confirm="deletePay(record)">
          <a>删除</a>
        </a-popconfirm>
      </template>
    </template>
  </a-table>
</template>

<script lang="ts" setup>
import { reactive, ref } from 'vue'
import {
  AddPay,
  DeletePay,
  GetPay,
  type PayConfig,
  type PayConfigRequest
} from '@/interfaces/Config'
import { type FormInstance, message } from 'ant-design-vue'
import StaticUpload from '@/components/upload/StaticUpload.vue'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const serverHost = import.meta.env.VITE_API_HOST

const columns = [
  {
    title: '二维码',
    dataIndex: 'image',
    key: 'image'
  },
  {
    title: '名称',
    dataIndex: 'name',
    key: 'name'
  },
  {
    title: 'operation',
    dataIndex: 'operation',
    key: 'operation'
  }
]
const list = ref<PayConfig[]>([])

const getPayConfig = async () => {
  const res: any = await GetPay()
  list.value = res.data.data.list || []
}
getPayConfig()

const formRef = ref<FormInstance>()
const visible = ref(false)
const formState = reactive<PayConfigRequest>({
  name: '',
  image: ''
})
const addPay = () => {
  if (formRef.value) {
    formRef.value
      .validateFields()
      .then(async () => {
        try {
          const response: any = await AddPay(formState)
          console.log(response)
          if (response.data.code !== 0) {
            message.error(response.data.data.message)
            return
          }
          message.success('添加成功')
          visible.value = false
          if (formRef.value) {
            formRef.value.resetFields()
          }
          await getPayConfig()
        } catch (error) {
          console.log(error)
        }
      })
      .catch((info) => {
        console.log('Validate Failed:', info)
        message.warning('请检查表单是否填写正确')
      })
  }
}

const deletePay = async (record: PayConfig) => {
  try {
    const response: any = await DeletePay(record.name, record.image)
    if (response.data.code === 0) {
      message.success('删除成功')
      await getPayConfig()
    } else {
      message.error(response.data.data.message)
    }
  } catch (error) {
    console.log(error)
  }
}
</script>
