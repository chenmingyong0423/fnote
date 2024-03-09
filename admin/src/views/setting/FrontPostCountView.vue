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
import {
  type FrontPostCountConfig,
  GetFrontPostCount,
  UpdateFrontPostCount
} from '@/interfaces/Config'
import { ref } from 'vue'
import { message } from 'ant-design-vue'

const editable = ref<boolean>(false)

const data = ref<FrontPostCountConfig>({
  count: 0
})

const getData = async () => {
  try {
    const response: any = await GetFrontPostCount()
    data.value = response.data.data || {}
  } catch (error) {
    console.log(error)
  }
}
getData()

const cancel = () => {
  editable.value = false
  getData()
}

const save = async () => {
  try {
    const response: any = await UpdateFrontPostCount({
      count: data.value.count
    })
    if (response.data.code === 0) {
      message.success('保存成功')
      await getData()
      editable.value = false
    } else {
      message.error(response.data.message)
    }
  } catch (error) {
    console.log(error)
  }
}

const changed = (e: any) => {
  data.value.count = Number(e.target.value)
}
</script>

<style scoped></style>
