<template>
  <a-card title="公告配置">
    <template #extra>
      <div class="flex gap-x-3">
        <a-tooltip title="刷新数据">
          <a-button shape="circle" :icon="h(ReloadOutlined)" :loading="loading" @click="getData" />
        </a-tooltip>
      </div>
    </template>
    <a-spin :spinning="loading">
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
    </a-spin>
    <div style="margin-top: 10px">
      <a-button v-if="!editable" @click="editable = true">编辑</a-button>
      <div v-else>
        <a-button @click="cancel" style="margin-right: 5px">取消</a-button>
        <a-button type="primary" @click="save">保存</a-button>
      </div>
    </div>
  </a-card>
</template>
<script lang="ts" setup>
import {
  type FrontPostCountConfig,
  GetFrontPostCount,
  UpdateFrontPostCount
} from '@/interfaces/Config'
import { h, ref } from 'vue'
import { message } from 'ant-design-vue'
import { ReloadOutlined } from '@ant-design/icons-vue'

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

const loading = ref(false)

const save = async () => {
  try {
    loading.value = true
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
  } finally {
    loading.value = false
  }
}

const changed = (e: any) => {
  data.value.count = Number(e.target.value)
}
</script>

<style scoped></style>
