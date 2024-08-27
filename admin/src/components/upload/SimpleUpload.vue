<template>
  <div class="clearfix">
    <a-upload
      v-model:file-list="fileList"
      :action="action"
      @change="handleChange"
      :before-upload="beforeUpload"
      :headers="{ Authorization: props.authorization }"
      name="file"
    >
      <a-button>
        <upload-outlined></upload-outlined>
        {{ props.label }}
      </a-button>
    </a-upload>
  </div>
</template>
<script lang="ts" setup>
import { ref } from 'vue'
import { UploadOutlined } from '@ant-design/icons-vue'
import { message, type UploadChangeParam, type UploadProps } from 'ant-design-vue'

const props = defineProps({
  action: {
    type: String,
    required: true
  },
  label: {
    type: String,
    default: '上传文件'
  },
  authorization: String,
  fileTypes: Array,
  maxSize: Number
})

const emit = defineEmits(['success:imageUrl'])

const fileList = ref<UploadProps['fileList']>([])

const handleChange = async (info: UploadChangeParam) => {
  if (info.file.status === undefined) {
    fileList.value = []
  }

  if (info.file.status === 'done') {
    console.log(123)
    emit('success:imageUrl', info.file.response.data.url)
    fileList.value = []
  }

  if (info.file.status === 'error') {
    message.error('上传失败')
  }
}

const beforeUpload = (file: any) => {
  if (props.fileTypes) {
    if (!props.fileTypes.includes(file.type)) {
      const types = props.fileTypes.join(',')
      message.error(`你只能上传 ${types} 文件。`)
      return false
    }
  }
  if (props.maxSize) {
    if (file.size > props.maxSize) {
      message.error(`你只能上传小于 ${props.maxSize}MB 的文件。`)
      return false
    }
  }
  return true
}
</script>
