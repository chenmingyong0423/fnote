<template>
  <div>
    <div>
      <a-button type="primary" @click="visible = true">新增分类</a-button>
      <a-modal
        v-model:open="visible"
        title="新增分类"
        ok-text="提交"
        cancel-text="取消"
        @ok="addCategory"
      >
        <a-form ref="formRef" :model="formState" layout="vertical" name="form_in_modal">
          <a-form-item
            name="name"
            label="名称"
            :rules="[{ required: true, message: '请输入分类名称' }]"
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
          <a-form-item name="description" label="描述">
            <a-textarea v-model:value="formState.description" />
          </a-form-item>
          <a-form-item
            label="是否显示在导航栏上"
            name="show_in_nav"
            class="collection-create-form_last-form-item"
          >
            <a-radio-group v-model:value="formState.show_in_nav">
              <a-radio :value="true">true</a-radio>
              <a-radio :value="false">false</a-radio>
            </a-radio-group>
          </a-form-item>
          <a-form-item
            label="是否启用"
            name="enabled"
            class="collection-create-form_last-form-item"
          >
            <a-radio-group v-model:value="formState.enabled">
              <a-radio :value="true">true</a-radio>
              <a-radio :value="false">false</a-radio>
            </a-radio-group>
          </a-form-item>
        </a-form>
      </a-modal>
    </div>
    <div>
      <a-table :columns="columns" :data-source="data" :pagination="pagination" @change="change">
        <template #bodyCell="{ column, text, record }">
          <template v-if="column.dataIndex === 'description'">
            <div>
              <a-textarea
                v-if="editableData[record.id]"
                v-model:value="editableData[record.id][column.dataIndex]"
                style="margin: -5px 0"
              />
              <template v-else>
                {{ text }}
              </template>
            </div>
          </template>

          <template v-if="column.dataIndex === 'create_time'">
            {{ dayjs.unix(text).format('YYYY-MM-DD HH:mm:ss') }}
          </template>

          <template v-if="column.dataIndex === 'update_time'">
            {{ dayjs.unix(text).format('YYYY-MM-DD HH:mm:ss') }}
          </template>

          <template v-if="column.key === 'enabled'">
            <a-switch v-model:checked="record.enabled" @change="changeCategoryEnabled(record)" />
          </template>

          <template v-if="column.key === 'show_in_nav'">
            <a-switch v-model:checked="record.show_in_nav" @change="changeCategoryNav(record)" />
          </template>

          <template v-else-if="column.dataIndex === 'operation'">
            <div class="editable-row-operations">
              <span v-if="editableData[record.id]">
                <a-typography-link @click="save(record.id)">保存</a-typography-link>
                <a-popconfirm title="确定取消？" @confirm="cancel(record.id)">
                  <a>取消</a>
                </a-popconfirm>
              </span>
              <span v-else>
                <a @click="edit(record.id)">编辑</a>
              </span>

              <a-popconfirm
                v-if="data.length"
                title="确认删除？"
                @confirm="deleteCategory(record.id)"
              >
                <a>删除</a>
              </a-popconfirm>
            </div>
          </template>
        </template>
      </a-table>
    </div>
  </div>
</template>
<script lang="ts" setup>
import axios from '@/http/axios'
import originalAxios from 'axios'
import { ref, reactive, toRaw, type UnwrapRef, computed } from 'vue'
import type { FormInstance } from 'ant-design-vue'
import type { IBaseResponse, IPageData, IResponse, PageRequest } from '@/interfaces/Common'
import type { CategoryRequest, ICategory } from '@/interfaces/Category'
import { message } from 'ant-design-vue'
import { cloneDeep } from 'lodash-es'
import dayjs from 'dayjs'
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
    title: '描述',
    key: 'description',
    dataIndex: 'description'
  },
  {
    title: '状态',
    key: 'enabled',
    dataIndex: 'enabled'
  },
  {
    title: '导航栏显示',
    key: 'show_in_nav',
    dataIndex: 'show_in_nav'
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

const data = ref<ICategory[]>([])

const pageReq = ref<PageRequest>({
  pageNo: 1,
  pageSize: 5,
  sortField: 'create_time',
  sortOrder: 'desc'
} as PageRequest)

const total = ref(0)

const pagination = computed(() => ({
  total: total.value,
  current: pageReq.value.pageNo,
  pageSize: pageReq.value.pageSize
}))

const getCategories = async () => {
  try {
    const response = await axios.get<IResponse<IPageData<ICategory>>>('/admin/categories', {
      params: pageReq.value
    })
    data.value = response.data.data?.list || []
    total.value = response.data.data?.totalCount || 0
  } catch (error) {
    console.log(error)
  }
}

getCategories()

// 添加分类
const formRef = ref<FormInstance>()
const visible = ref(false)
const formState = reactive<CategoryRequest>({
  name: '',
  route: '',
  description: '',
  show_in_nav: false,
  enabled: false
})

const addCategory = () => {
  if (formRef.value) {
    formRef.value
      .validateFields()
      .then(async (values) => {
        try {
          // 提交 body 参数 values
          const response = await axios.post<IBaseResponse>('/admin/categories', formState)
          if (response.data.code !== 200) {
            message.error(response.data.message)
            return
          }
          message.success('添加成功')
          visible.value = false
          if (formRef.value) {
            formRef.value.resetFields()
          }
          await getCategories()
        } catch (error) {
          console.log(error)
          if (originalAxios.isAxiosError(error)) {
            // 这是一个由 axios 抛出的错误
            if (error.response) {
              if (error.response.status === 409) {
                message.error('分类名称或路由重复')
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

const changeCategoryEnabled = async (record: ICategory) => {
  try {
    const response = await axios.put<IBaseResponse>(`/admin/categories/enabled/${record.id}`, {
      enabled: record.enabled
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
  await getCategories()
}

const changeCategoryNav = async (record: ICategory) => {
  try {
    const response = await axios.put<IBaseResponse>(`/admin/categories/navigation/${record.id}`, {
      show_in_nav: record.show_in_nav
    })
    if (response.data.code !== 200) {
      message.error(response.data.message)
      return
    }
    message.success('设置成功')
  } catch (error) {
    console.log(error)
    message.error('设置失败')
  }
  await getCategories()
}

// 删除
const deleteCategory = async (id: string) => {
  try {
    // 提交 body 参数 values
    const response = await axios.delete<IBaseResponse>(`/admin/categories/${id}`)
    if (response.data.code !== 200) {
      message.error(response.data.message)
      return
    }
    message.success('删除成功')
    await getCategories()
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

// 编辑
const editableData: UnwrapRef<Record<string, ICategory>> = reactive({})
const edit = (id: string) => {
  editableData[id] = cloneDeep(data.value.filter((item) => id === item.id)[0])
}

const save = async (id: string) => {
  const editableDatum = editableData[id]
  try {
    // 提交 body 参数 values
    const response = await axios.put<IBaseResponse>(
      `/admin/categories/${editableDatum.id}`,
      editableDatum
    )
    if (response.data.code !== 200) {
      message.error(response.data.message)
      return
    }
    message.success('更新成功')
    delete editableData[id]
    await getCategories()
  } catch (error) {
    console.log(error)
    message.error('更新失败')
  }
}
const cancel = (key: string) => {
  delete editableData[key]
}

const change = (pg, filters, sorter, { currentDataSource }) => {
  pageReq.value.pageNo = pg.current
  pageReq.value.pageSize = pg.pageSize
  getCategories()
}
</script>

<style scoped>
.collection-create-form_last-form-item {
  margin-bottom: 0;
}

.editable-row-operations a {
  margin-right: 8px;
}
</style>
