<template>
  <a-card title="分类列表">
    <template #extra>
      <div class="flex gap-x-3">
        <a-tooltip title="刷新数据">
          <a-button
            shape="circle"
            :icon="h(ReloadOutlined)"
            :loading="loading"
            @click="getCategories"
          />
        </a-tooltip>
      </div>
    </template>
    <div>
      <a-button @click="openCreateCategory" class="mb-5">新增分类</a-button>
      <a-modal
        v-model:open="visible"
        title="新增分类"
        ok-text="提交"
        cancel-text="取消"
        @ok="addCategory"
        @cancel="resetCategoryForm"
      >
        <TaxonomyCreateForm
          ref="formRef"
          v-model:name="formState.name"
          v-model:route="formState.route"
          v-model:description="formState.description"
          v-model:enabled="formState.enabled"
          v-model:show-in-nav="formState.show_in_nav"
          type="category"
        />
      </a-modal>
    </div>
    <div>
      <a-spin :spinning="loading">
        <a-table :columns="columns" :data-source="data" :pagination="pagination" @change="change">
          <template #bodyCell="{ column, text, record }">
            <template v-if="column.dataIndex === 'description'">
              <div>
                <a-textarea
                  v-if="editableData[record.id]"
                  v-model:value="
                    editableData[record.id][column.dataIndex as keyof UpdateCategoryRequest]
                  "
                  style="margin: -5px 0"
                />
                <template v-else>
                  {{ text }}
                </template>
              </div>
            </template>

            <template v-if="column.dataIndex === 'created_at'">
              {{ dayjs.unix(text).format('YYYY-MM-DD HH:mm:ss') }}
            </template>

            <template v-if="column.dataIndex === 'updated_at'">
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
                  @confirm="deleteCategory(record.id, record.post_count)"
                >
                  <a>删除</a>
                </a-popconfirm>
              </div>
            </template>
          </template>
        </a-table>
      </a-spin>
    </div>
  </a-card>
</template>
<script lang="ts" setup>
import originalAxios from 'axios'
import { ref, reactive, type UnwrapRef, computed, h } from 'vue'
import type { PageRequest } from '@/interfaces/Common'
import {
  AddCategory,
  type CategoryRequest,
  ChangeCategoryEnabled,
  ChangeCategoryShowInNav,
  DeleteCategory,
  GetCategories,
  type ICategory,
  UpdateCategory,
  type UpdateCategoryRequest
} from '@/interfaces/Category'
import { message } from 'ant-design-vue'
import { cloneDeep } from 'lodash-es'
import dayjs from 'dayjs'
import { ReloadOutlined } from '@ant-design/icons-vue'
import TaxonomyCreateForm from '@/components/form/TaxonomyCreateForm.vue'

document.title = '分类列表 - 后台管理'

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
    title: '文章数量',
    key: 'post_count',
    dataIndex: 'post_count'
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
    key: 'created_at',
    dataIndex: 'created_at'
  },
  {
    title: '最后一次修改的时间',
    key: 'updated_at',
    dataIndex: 'updated_at'
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
  sortField: 'created_at',
  sortOrder: 'desc'
} as PageRequest)

const total = ref(0)

const pagination = computed(() => ({
  total: total.value,
  current: pageReq.value.pageNo,
  pageSize: pageReq.value.pageSize
}))

const loading = ref(false)

const getCategories = async () => {
  try {
    loading.value = true
    const response: any = await GetCategories(pageReq.value)
    data.value = response.data.data?.list || []
    total.value = response.data.data?.totalCount || 0
  } catch (error) {
    console.log(error)
  } finally {
    loading.value = false
  }
}

getCategories()

// 添加分类
const formRef = ref<InstanceType<typeof TaxonomyCreateForm>>()
const visible = ref(false)
const formState = reactive<CategoryRequest>({
  name: '',
  route: '',
  description: '',
  show_in_nav: false,
  enabled: true
})

const resetCategoryForm = () => {
  formState.name = ''
  formState.route = ''
  formState.description = ''
  formState.show_in_nav = false
  formState.enabled = true
  formRef.value?.resetRouteTouched()
  formRef.value?.clearValidate()
}

const openCreateCategory = () => {
  resetCategoryForm()
  visible.value = true
}

const addCategory = () => {
  if (formRef.value) {
    formRef.value
      .validateFields()
      .then(async () => {
        try {
          const response: any = await AddCategory(formState)
          if (response.data.code !== 0) {
            message.error(response.data.message)
            return
          }
          message.success('添加成功')
          visible.value = false
          resetCategoryForm()
          await getCategories()
        } catch (error) {
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
    const response: any = await ChangeCategoryEnabled(record.id, record.enabled)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('修改成功')
  } catch (error) {
    console.log(error)
  }
  await getCategories()
}

const changeCategoryNav = async (record: ICategory) => {
  try {
    const response: any = await ChangeCategoryShowInNav(record.id, record.show_in_nav)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('设置成功')
  } catch (error) {
    console.log(error)
  }
  await getCategories()
}

// 删除
const deleteCategory = async (id: string, postCount: number) => {
  try {
    if (postCount > 0) {
      message.warn('该分类下有文章，不能删除')
      return
    }
    const response: any = await DeleteCategory(id)
    if (response.data.code !== 0) {
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
        if (error.response.data.status === 404) {
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
const editableData: UnwrapRef<Record<string, UpdateCategoryRequest>> = reactive({})
const edit = (id: string) => {
  editableData[id] = cloneDeep(data.value.filter((item) => id === item.id)[0])
}

const save = async (id: string) => {
  const editableDatum = editableData[id]
  try {
    const response: any = await UpdateCategory(id, editableDatum)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('更新成功')
    delete editableData[id]
    await getCategories()
  } catch (error) {
    console.log(error)
  }
}
const cancel = (key: string) => {
  delete editableData[key]
}

const change = (pg: any) => {
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
