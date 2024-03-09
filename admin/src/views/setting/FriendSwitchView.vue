<template>
  <div>
    <a-descriptions title="友链配置" :column="1" bordered>
      <a-descriptions-item label="是否开启友链申请">
        <div>
          <a-switch v-model:checked="data.enable_friend_commit" @change="save" />
        </div>
      </a-descriptions-item>
    </a-descriptions>
  </div>
</template>

<script lang="ts" setup>
import { type FriendConfig, GetFriend, UpdateFriend } from '@/interfaces/Config'
import { message } from 'ant-design-vue'
import { ref } from 'vue'

const data = ref<FriendConfig>({
  enable_friend_commit: false
})

const getCommentConfig = async () => {
  try {
    const response: any = await GetFriend()
    data.value.enable_friend_commit = response.data.data.enable_friend_commit || false
  } catch (error) {
    console.log(error)
  }
}
getCommentConfig()

const save = async () => {
  try {
    const response: any = await UpdateFriend({
      enable_friend_commit: data.value.enable_friend_commit
    })
    if (response.data.code === 0) {
      message.success('保存成功')
    } else {
      message.error(response.data.data.message)
    }
  } catch (error) {
    console.log(error)
  }
  await getCommentConfig()
}
</script>
