<template>
  <a-card title="站点验证配置">
    <template #extra>
      <div class="flex gap-x-3">
        <a-tooltip title="刷新数据">
          <a-button shape="circle" :icon="h(ReloadOutlined)" :loading="loading" @click="get" />
        </a-tooltip>
      </div>
    </template>
    <div class="color-gray">示例：key = baidu_site_verification, value = xxxxxx</div>
    <a-form
      class="my-4"
      :model="formState"
      name="horizontal_login"
      layout="inline"
      autocomplete="off"
      @finish="onFinish"
      @finishFailed="onFinishFailed"
    >
      <a-form-item label="key" name="key" :rules="[{ required: true, message: '请输入 key!' }]">
        <a-input v-model:value="formState.key"></a-input>
      </a-form-item>

      <a-form-item
        label="value"
        name="value"
        :rules="[{ required: true, message: '请输入 value!' }]"
      >
        <a-input v-model:value="formState.value"></a-input>
      </a-form-item>

      <a-form-item
        label="描述"
        name="description"
        :rules="[{ required: true, message: '请输入描述!' }]"
      >
        <a-input v-model:value="formState.description"></a-input>
      </a-form-item>

      <a-form-item>
        <a-button html-type="submit">添加</a-button>
      </a-form-item>
    </a-form>
    <a-spin :spinning="loading">
      <a-table bordered :data-source="dataSource" :columns="columns">
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'operation'">
            <a-popconfirm
              v-if="dataSource.length"
              title="Sure to delete?"
              @confirm="onDelete(record.key)"
            >
              <a>删除</a>
            </a-popconfirm>
          </template>
        </template>
      </a-table>
    </a-spin>
  </a-card>
</template>
<script lang="ts" setup>
import { h, reactive, ref } from 'vue'
import type { Ref } from 'vue'
import {
  AddThirdPartySiteVerification,
  DeleteThirdPartySiteVerification,
  GetThirdPartySiteVerification,
  type ThirdPartySiteVerification,
  type ThirdPartySiteVerificationRequest
} from '@/interfaces/Config'
import { message } from 'ant-design-vue'
import { ReloadOutlined } from '@ant-design/icons-vue'

const columns = [
  {
    title: 'key',
    dataIndex: 'key'
  },
  {
    title: 'value',
    dataIndex: 'value'
  },
  {
    title: '描述',
    dataIndex: 'description',
    width: '30%'
  },
  {
    title: 'operation',
    dataIndex: 'operation'
  }
]

const dataSource: Ref<ThirdPartySiteVerification[]> = ref([])
const loading = ref(false)
const get = async () => {
  try {
    loading.value = true
    const res: any = await GetThirdPartySiteVerification()
    dataSource.value = res.data.data.list || []
  } catch (e) {
    console.log(e)
  } finally {
    loading.value = false
  }
}
get()

const onDelete = async (id: string) => {
  const res: any = await DeleteThirdPartySiteVerification(id)
  if (res && res.data.code === 0) {
    message.success('删除成功！')
    formState.key = ''
    formState.value = ''
    formState.description = ''
    await get()
  } else {
    message.error('Failed:', res)
  }
}

const formState = reactive<ThirdPartySiteVerificationRequest>({
  key: '',
  value: '',
  description: ''
})
const onFinish = async (req: ThirdPartySiteVerificationRequest) => {
  const res: any = await AddThirdPartySiteVerification(req)
  if (res && res.data.code === 0) {
    message.success('添加成功！')
    formState.key = ''
    formState.value = ''
    formState.description = ''
    await get()
  } else {
    message.error('Failed:', res)
  }
}

const onFinishFailed = (errorInfo: any) => {
  message.error('Failed:', errorInfo)
}
</script>
<style scoped>
.editable-cell {
  position: relative;

  .editable-cell-input-wrapper,
  .editable-cell-text-wrapper {
    padding-right: 24px;
  }

  .editable-cell-text-wrapper {
    padding: 5px 24px 5px 5px;
  }

  .editable-cell-icon,
  .editable-cell-icon-check {
    position: absolute;
    right: 0;
    width: 20px;
    cursor: pointer;
  }

  .editable-cell-icon {
    margin-top: 4px;
    display: none;
  }

  .editable-cell-icon-check {
    line-height: 28px;
  }

  .editable-cell-icon:hover,
  .editable-cell-icon-check:hover {
    color: #108ee9;
  }

  .editable-add-btn {
    margin-bottom: 8px;
  }
}

.editable-cell:hover .editable-cell-icon {
  display: inline-block;
}
</style>
