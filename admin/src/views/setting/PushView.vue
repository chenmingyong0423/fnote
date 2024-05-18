<template>
  <a-card title="推送配置">
    <template #extra>
      <div class="flex gap-x-3">
        <a-tooltip title="刷新数据">
          <a-button
            shape="circle"
            :icon="h(ReloadOutlined)"
            :loading="loading"
            @click="getConfig"
          />
        </a-tooltip>
      </div>
    </template>
    <a-spin :spinning="loading">
      <a-descriptions title="百度推送配置" :column="1" bordered>
        <a-descriptions-item label="站点">
          <div>
            <a-input
              v-if="baiduEditable"
              v-model:value="baiduPushCfg.site"
              style="margin: -5px 0"
            />
            <template v-else>
              {{ baiduPushCfg.site }}
            </template>
          </div>
        </a-descriptions-item>
        <a-descriptions-item label="token">
          <div>
            <a-input
              v-if="baiduEditable"
              v-model:value="baiduPushCfg.token"
              style="margin: -5px 0"
            />
            <template v-else>
              {{ baiduPushCfg.token }}
            </template>
          </div>
        </a-descriptions-item>
      </a-descriptions>
      <div style="margin-top: 10px">
        <a-button v-if="!baiduEditable" @click="baiduEditable = true">编辑</a-button>
        <div v-else>
          <a-button @click="cancel4Baidu" style="margin-right: 5px">取消</a-button>
          <a-button type="primary" @click="save4Baidu">保存</a-button>
        </div>
      </div>
    </a-spin>
  </a-card>
</template>

<script lang="ts" setup>
import { h, ref } from 'vue'
import {
  type BaiduPushConfig,
  GetBaiduPushConfig,
  UpdateBaiduPushConfig
} from '@/interfaces/Config'
import { message } from 'ant-design-vue'
import { ReloadOutlined } from '@ant-design/icons-vue'

const baiduPushCfg = ref<BaiduPushConfig>({
  site: '',
  token: ''
})
const baiduEditable = ref(false)

const loading = ref(false)
const getConfig = async () => {
  try {
    loading.value = true
    const res: any = await GetBaiduPushConfig()
    baiduPushCfg.value = res.data.data || {}
  } catch (e) {
    console.log(e)
  } finally {
    loading.value = false
  }
}
getConfig()

const cancel4Baidu = async () => {
  baiduEditable.value = false
  await getConfig()
}

const save4Baidu = async () => {
  try {
    const res: any = await UpdateBaiduPushConfig(baiduPushCfg.value.site, baiduPushCfg.value.token)
    console.log(res)
    if (!res || res.data.code != 0) {
      message.error('更新失败！')
    }
    baiduEditable.value = false
    message.success('更新成功！')
  } catch (err) {
    message.error('更新失败！')
    console.log(err)
  }
}
</script>
