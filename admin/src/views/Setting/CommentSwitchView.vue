<template>
  <div>
    <a-descriptions title="评论配置" :column="1" bordered>
      <a-descriptions-item label="是否开启评论">
        <div>
          <a-switch v-model:checked="data.enable_comment" @change="save" />
        </div>
      </a-descriptions-item>
    </a-descriptions>
  </div>
</template>

<script lang="ts" setup>
import axios from '@/http/axios'
import type { IResponse } from '@/interfaces/Common'
import type { CommentConfig, OwnerConfig, SeoConfig } from '@/interfaces/Config'
import { message } from 'ant-design-vue'
import { ref } from 'vue'

const data = ref<CommentConfig>({
  enable_comment: false
})

const getCommentSwitch = async () => {
  try {
    const response = await axios.get<IResponse<OwnerConfig>>('/admin/configs/comment')
    console.log(response)
    data.value.enable_comment = response.data.data.enable_comment || false
  } catch (error) {
    console.log(error)
    message.error('获取信息失败')
  }
}
getCommentSwitch()

const save = async () => {
  try {
    const response = await axios.put<IResponse<SeoConfig>>('/admin/configs/comment', {
      enable_comment: data.value.enable_comment
    })
    if (response.data.code === 200) {
      message.success('保存成功')
    } else {
      message.error(response.data.message)
    }
  } catch (error) {
    console.log(error)
    message.error('保存失败')
  }
  await getCommentSwitch()
}
</script>
