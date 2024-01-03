<template>
  <div>
    <a-descriptions title="首页展示文章数量配置" :column="1" bordered>
      <a-descriptions-item label="数量">
        <div>
          <a-input
            v-if="editable"
            v-model:value="data.count"
            style="margin: -5px 0"
            @change="changed"
          />
          <template v-else>
            {{ data.count }}
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
import axios from '@/http/axios'
import type { IBaseResponse, IResponse } from '@/interfaces/Common'
import type { FrontPostCountConfig } from '@/interfaces/Config'
import { ref } from 'vue'
import { message } from 'ant-design-vue'

const editable = ref<boolean>(false)

const data = ref<FrontPostCountConfig>({
  count: 0
})

const getEmail = async () => {
  try {
    const response = await axios.get<IResponse<FrontPostCountConfig>>(
      '/admin/configs/front-post-count'
    )
    data.value = response.data.data || {}
  } catch (error) {
    console.log(error)
    message.error('获取信息失败')
  }
}
getEmail()

const cancel = () => {
  editable.value = false
  getEmail()
}

const save = async () => {
  try {
    const response = await axios.put<IBaseResponse>('/admin/configs/front-post-count', {
      count: data.value.count
    })
    if (response.data.code === 200) {
      message.success('保存成功')
      await getEmail()
      editable.value = false
    } else {
      message.error(response.data.message)
    }
  } catch (error) {
    console.log(error)
    message.error('保存失败')
  }
}

const changed = (e: any) => {
  data.value.count = Number(e.target.value)
}
</script>

<style scoped></style>
