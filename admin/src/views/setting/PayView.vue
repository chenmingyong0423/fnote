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
          ><a-input v-model:value="formState.image" placeholder="请输入二维码路径" />
          <a-upload
            v-model:file-list="fileList"
            name="file"
            list-type="picture-card"
            class="avatar-uploader m-y-5"
            :show-upload-list="false"
            action="http://localhost:8080/admin/files/upload"
            :before-upload="beforeUpload"
            @change="handleChange"
          >
            <img v-if="imageUrl" :src="imageUrl" alt="avatar" width="250" height="150" />
            <div v-else>
              <loading-outlined v-if="loading"></loading-outlined>
              <plus-outlined v-else></plus-outlined>
              <div class="ant-upload-text">Upload</div>
            </div>
          </a-upload>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
  <a-table :columns="columns" :data-source="list">
    <template #bodyCell="{ column, record }">
      <template v-if="column.key === 'image'">
        <a-image :width="200" :src="record.image" />
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
import type { PayConfig, PayConfigRequest } from '@/interfaces/Config'
import type { IBaseResponse, IListData, IResponse } from '@/interfaces/Common'
import axios from '@/http/axios'
import {
  type FormInstance,
  message,
  type UploadChangeParam,
  type UploadProps
} from 'ant-design-vue'

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
const imageUrl = ref<string>('')

const getPayConfig = async () => {
  const res: any = await axios.get<IResponse<IListData<PayConfig>>>('/admin/configs/pay')
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
      .then(async (values) => {
        try {
          // 提交 body 参数 values
          const response = await axios.post<IBaseResponse>('/admin/configs/pay', formState)
          if (response.data.code !== 200) {
            message.error(response.data.message)
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
          message.error('添加失败')
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
    const response = await axios.delete<IBaseResponse>(
      `/admin/configs/pay/${record.name}?image=${record.image}`
    )
    if (response.data.code === 200) {
      message.success('删除成功')
      await getPayConfig()
    } else {
      message.error(response.data.message)
    }
  } catch (error) {
    console.log(error)
    message.error('删除失败')
  }
}

// 二维码上传
function getBase64(img: Blob, callback: (base64Url: string) => void) {
  const reader = new FileReader()
  reader.addEventListener('load', () => callback(reader.result as string))
  reader.readAsDataURL(img)
}

const fileList = ref([])
const loading = ref<boolean>(false)

const handleChange = (info: UploadChangeParam) => {
  if (info.file.status === 'uploading') {
    loading.value = true
    return
  }
  if (info.file.status === 'done') {
    // Get this url from response in real world.
    getBase64(info.file.originFileObj, (base64Url: string) => {
      imageUrl.value = base64Url
      loading.value = false
      // Get this url from response in real world.
      postReq.cover_img = info.file.response.data.url
      message.success('上传成功')
    })
  }
  if (info.file.status === 'error') {
    loading.value = false
    message.error('upload error')
  }
}

const beforeUpload = (file: UploadProps['fileList'][number]) => {
  const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png'
  if (!isJpgOrPng) {
    message.error('You can only upload JPG file!')
  }
  const isLt2M = file.size / 1024 / 1024 < 2
  if (!isLt2M) {
    message.error('Image must smaller than 2MB!')
  }
  return isJpgOrPng && isLt2M
}
</script>
