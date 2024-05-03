<template>
  <div>
    <a-descriptions title="友链配置" :column="1" bordered>
      <a-descriptions-item label="是否开启友链申请">
        <div>
          <a-switch v-model:checked="enableFriendCommit" @change="save" />
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="友链页面介绍">
        <div>
          <div class="mb-5">
            <a-button
              type="primary"
              v-if="editorMode === 'preview'"
              @click="editorMode = 'editable'"
              >编辑</a-button
            >
            <div v-else class="flex gap-x-2">
              <a-button type="primary" @click="saveIntroduction">保存</a-button>
              <a-button @click="refreshData()">取消</a-button>
            </div>
          </div>
          <v-md-editor
            v-model="introduction"
            height="400px"
            :disabled-menus="[]"
            @upload-image="handleUploadImage"
            :mode="editorMode"
          />
        </div>
      </a-descriptions-item>
    </a-descriptions>
  </div>
</template>

<script lang="ts" setup>
import {
  type FriendConfigVO,
  GetFriend,
  UpdateFriendIntroduction,
  UpdateFriendSwitch
} from '@/interfaces/Config'
import { message } from 'ant-design-vue'
import { ref } from 'vue'
import { FileUpload } from '@/interfaces/File'
import type { IBaseResponse, IResponse } from '@/interfaces/Common'

const enableFriendCommit = ref(false)
const introduction = ref('')

const getCommentConfig = async () => {
  try {
    const response: any = await GetFriend()
    const result: IResponse<FriendConfigVO> = response.data
    if (result.code === 0) {
      enableFriendCommit.value = result.data?.enable_friend_commit || enableFriendCommit.value
      introduction.value = result.data?.introduction || introduction.value
    } else {
      message.error('获取友链配置失败：' + result.message)
    }
  } catch (error) {
    console.log(error)
  }
}
getCommentConfig()

const save = async () => {
  try {
    const response: any = await UpdateFriendSwitch({
      enable_friend_commit: enableFriendCommit.value
    })
    const result: IBaseResponse = response.data
    if (result.code === 0) {
      message.success('保存成功')
    } else {
      message.error(response.data.data.message)
    }
  } catch (error) {
    console.log(error)
  }
  await getCommentConfig()
}

// md 图片上传
const handleUploadImage = async (_event: any, insertImage: any, files: any) => {
  try {
    const formData = new FormData()
    formData.append('file', files[0])
    try {
      const res: any = await FileUpload(formData)
      if (res.data.code !== 0) {
        message.error(res.data.message)
        return
      }
      insertImage({
        url: res.data.data.url,
        desc: '请在此添加图片描述'
      })
    } catch (error) {
      message.error(error)
    }
  } catch (error) {
    console.log(error)
  }
}

const editorMode = ref('preview')

const refreshData = async () => {
  editorMode.value = 'preview'
  await getCommentConfig()
}

const saveIntroduction = async () => {
  try {
    const response: any = await UpdateFriendIntroduction({
      introduction: introduction.value
    })
    const result: IBaseResponse = response.data
    if (result.code === 0) {
      message.success('保存成功')
      editorMode.value = 'preview'
    } else {
      message.error(response.data.data.message)
    }
  } catch (error) {
    console.log(error)
  }
  await getCommentConfig()
}
</script>
