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
import { type CommentConfig, GetComment, UpdateComment } from '@/interfaces/Config'
import { message } from 'ant-design-vue'
import { ref } from 'vue'

const data = ref<CommentConfig>({
  enable_comment: false
})

const getCommentConfig = async () => {
  try {
    const response: any = await GetComment()
    data.value.enable_comment = response.data.data.enable_comment || false
  } catch (error) {
    console.log(error)
  }
}
getCommentConfig()

const save = async () => {
  try {
    const response: any = await UpdateComment({
      enable_comment: data.value.enable_comment
    })
    if (response.data.code === 0) {
      message.success('保存成功')
    } else {
      message.error(response.data.message)
    }
  } catch (error) {
    console.log(error)
  }
  await getCommentConfig()
}
</script>
