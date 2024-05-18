<template>
  <a-card title="评论配置">
    <template #extra>
      <div class="flex gap-x-3">
        <a-tooltip title="刷新数据">
          <a-button
            shape="circle"
            :icon="h(ReloadOutlined)"
            :loading="loading"
            @click="getCommentConfig"
          />
        </a-tooltip>
      </div>
    </template>
    <a-spin :spinning="loading">
      <a-descriptions :column="1" bordered>
        <a-descriptions-item label="是否开启评论">
          <div>
            <a-switch v-model:checked="data.enable_comment" @change="save" />
          </div>
        </a-descriptions-item>
      </a-descriptions>
    </a-spin>
  </a-card>
</template>

<script lang="ts" setup>
import { type CommentConfig, GetComment, UpdateComment } from '@/interfaces/Config'
import { message } from 'ant-design-vue'
import { h, ref } from 'vue'
import { ReloadOutlined } from '@ant-design/icons-vue'

const data = ref<CommentConfig>({
  enable_comment: false
})

const loading = ref(false)

const getCommentConfig = async () => {
  try {
    loading.value = true
    const response: any = await GetComment()
    data.value.enable_comment = response.data.data.enable_comment || false
  } catch (error) {
    console.log(error)
  } finally {
    loading.value = false
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
