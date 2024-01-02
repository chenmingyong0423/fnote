<template>
  <div>
    <div>
      <a-button type="primary" @click="visible = true">新增标签</a-button>
      <a-modal
        v-model:open="visible"
        title="新增分类"
        ok-text="提交"
        cancel-text="取消"
        @ok="addTag"
      >
        <a-form ref="formRef" :model="formState" layout="vertical" name="form_in_modal">
          <a-form-item
            name="name"
            label="名称"
            :rules="[{ required: true, message: '请输入标签名称' }]"
          >
            <a-input v-model:value="formState.name" />
          </a-form-item>
          <a-form-item
            name="route"
            label="前端路由"
            :rules="[{ required: true, message: '请输入前端路由' }]"
          >
            <a-input v-model:value="formState.route" />
          </a-form-item>
          <a-form-item
            label="是否启用"
            name="disabled"
            class="collection-create-form_last-form-item"
          >
            <a-radio-group v-model:value="formState.disabled">
              <a-radio :value="true">true</a-radio>
              <a-radio :value="false">false</a-radio>
            </a-radio-group>
          </a-form-item>
        </a-form>
      </a-modal>
    </div>
    <div>
      <a-table :columns="columns" :data-source="data">
        <template #bodyCell="{ column, text, record }">
          <template v-if="column.dataIndex === 'create_time'">
            {{ dayjs.unix(text).format('YYYY-MM-DD HH:mm:ss') }}
          </template>

          <template v-if="column.dataIndex === 'update_time'">
            {{ dayjs.unix(text).format('YYYY-MM-DD HH:mm:ss') }}
          </template>

          <template v-if="column.key === 'disabled'">
            <a-switch v-model:checked="record.disabled" @change="changeTagDisabled(record)" />
          </template>

          <template v-if="column.key === 'show_in_nav'">
            <a-switch v-model:checked="record.show_in_nav" @change="changeTagNav(record)" />
          </template>

          <template v-else-if="column.dataIndex === 'operation'">
            <a-popconfirm v-if="data.length" title="确认删除？" @confirm="deleteTag(record.id)">
              <a>删除</a>
            </a-popconfirm>
          </template>
        </template>
      </a-table>
    </div>
  </div>
</template>
<script lang="ts" setup>
import axios from '@/http/axios'
import originalAxios from 'axios'
import { ref, reactive, toRaw, type UnwrapRef } from 'vue'
import type { FormInstance } from 'ant-design-vue'
import type { IBaseResponse, IPageData, IResponse, PageRequest } from '@/interfaces/Common'
import { message } from 'ant-design-vue'
import { cloneDeep } from 'lodash-es'
import dayjs from 'dayjs'
import type { Tag, TagRequest } from '@/interfaces/Tag'
const columns = [
  {
    title: '名称',
    dataIndex: 'name',
    key: 'name'
  },
  {
    title: '路由',
    dataIndex: 'route',
    key: 'route'
  },
  {
    title: '状态',
    key: 'disabled',
    dataIndex: 'disabled'
  },
  {
    title: '创建时间',
    key: 'create_time',
    dataIndex: 'create_time'
  },
  {
    title: '最后一次修改的时间',
    key: 'update_time',
    dataIndex: 'update_time'
  },
  {
    title: 'operation',
    dataIndex: 'operation'
  }
]

const data = ref<Tag[]>([])

const pageReq = ref<PageRequest>({
  pageNo: 1,
  pageSize: 10,
  sortField: 'create_time',
  sortOrder: 'desc'
} as PageRequest)

const getTags = async () => {
  try {
    const response = await axios.get<IResponse<IPageData<Tag>>>('/admin/tags', {
      params: pageReq.value
    })
    data.value = response.data.data?.list || []
  } catch (error) {
    console.log(error)
  }
}

getTags()

// 添加分类
const formRef = ref<FormInstance>()
const visible = ref(false)
const formState = reactive<TagRequest>({
  name: '',
  route: '',
  disabled: false
})

const addTag = () => {
  if (formRef.value) {
    formRef.value
      .validateFields()
      .then(async (values) => {
        try {
          // 提交 body 参数 values
          const response = await axios.post<IBaseResponse>('/admin/tags', formState)
          if (response.data.code !== 200) {
            message.error(response.data.message)
            return
          }
          message.success('添加成功')
          visible.value = false
          if (formRef.value) {
            formRef.value.resetFields()
          }
          await getTags()
        } catch (error) {
          console.log(error)
          if (originalAxios.isAxiosError(error)) {
            // 这是一个由 axios 抛出的错误
            if (error.response) {
              if (error.response.status === 409) {
                message.error('标签名称或路由重复')
                return
              }
            } else if (error.request) {
              // 请求已发出，但没有收到响应
              console.log('No response received:', error.request)
            } else {
              // 在设置请求时触发了一个错误
              console.log('Error Message:', error.message)
            }
          }
          message.error('添加失败')
        }
      })
      .catch((info) => {
        console.log('Validate Failed:', info)
        message.warning('请检查表单是否填写正确')
      })
  }
}

const changeTagDisabled = async (record: Tag) => {
  try {
    const response = await axios.put<IBaseResponse>(`/admin/tags/disabled/${record.id}`, {
      disabled: record.disabled
    })
    if (response.data.code !== 200) {
      message.error(response.data.message)
      return
    }
    message.success('修改成功')
  } catch (error) {
    console.log(error)
    message.error('修改失败')
  }
  await getTags()
}

// 删除
const deleteTag = async (id: string) => {
  try {
    // 提交 body 参数 values
    const response = await axios.delete<IBaseResponse>(`/admin/tags/${id}`)
    if (response.data.code !== 200) {
      message.error(response.data.message)
      return
    }
    message.success('删除成功')
    await getTags()
  } catch (error) {
    console.log(error)
    if (originalAxios.isAxiosError(error)) {
      // 这是一个由 axios 抛出的错误
      if (error.response) {
        if (error.response.status === 404) {
          message.error('id 不存在')
          return
        }
      } else if (error.request) {
        // 请求已发出，但没有收到响应
        console.log('No response received:', error.request)
      } else {
        // 在设置请求时触发了一个错误
        console.log('Error Message:', error.message)
      }
    }
    message.error('删除失败')
  }
}
</script>

<style scoped>
.collection-create-form_last-form-item {
  margin-bottom: 0;
}
</style>
