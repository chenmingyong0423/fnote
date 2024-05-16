<template>
  <a-card title="轮播图列表">
    <template #extra>
      <div class="flex gap-x-3">
        <a-tooltip title="刷新数据">
          <a-button
            shape="circle"
            :icon="h(ReloadOutlined)"
            :loading="loading"
            @click="getCarousel"
          />
        </a-tooltip>
      </div>
    </template>
    <a-button @click="visible = true" class="mb-5">新增轮播图</a-button>
    <a-modal
      v-model:open="visible"
      title="新增轮播图"
      ok-text="提交"
      cancel-text="取消"
      @ok="addCarousel"
    >
      <a-modal v-model:open="showPosts" title="请选择文章">
        <a-input-search
          v-model:value="req.keyword"
          placeholder="请输入关键字"
          class="mb-2 w-200px"
          @search="searchPost"
          @pressEnter="searchPost"
          allow-clear
        />
        <a-table
          :columns="postColumns"
          :data-source="posts"
          :pagination="pagination"
          @change="change"
          bordered
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'id' || column.key === 'title'">
              <span>{{ record[column.key] }}</span>
            </template>
            <template v-else-if="column.dataIndex === 'operation'">
              <div>
                <a-button @click="choosePost(record)">选择</a-button>
              </div>
            </template>
          </template>
        </a-table>
      </a-modal>
      <a-form ref="formRef" :model="formState" layout="vertical" name="form_in_modal">
        <a-form-item
          name="id"
          label="文章索引"
          :rules="[{ required: true, message: '请选择文章' }]"
        >
          <a-input v-model:value="formState.id" disabled class="w-30%" />
          <a-button @click="showPostList" class="ml-1 w-20%">选择文章</a-button>
        </a-form-item>
        <a-form-item
          name="cover_img"
          label="封面"
          :rules="[{ required: true, message: '请选择封面' }]"
        >
          <StaticUpload
            :image-url="formState.cover_img"
            @update:imageUrl="(value) => (formState.cover_img = value)"
            :authorization="userStore.token"
          />
        </a-form-item>
        <a-form-item name="title" label="标题" :rules="[{ required: true, message: '请输入标题' }]">
          <a-input v-model:value="formState.title" />
        </a-form-item>
        <a-form-item name="summary" label="摘要">
          <a-textarea v-model:value="formState.summary" />
        </a-form-item>
        <a-form-item name="color" label="选择与封面搭配的字体颜色" tooltip="默认黑色">
          <div class="flex gap-x-1">
            <span
              class="w-8 h-8 inline-block black_border"
              :style="{ backgroundColor: formState.color }"
            ></span>
            <a-button @click="showColorModal = true">选择字体颜色</a-button>
          </div>
        </a-form-item>
        <a-form-item label="是否显示" name="show" class="collection-create-form_last-form-item">
          <a-radio-group v-model:value="formState.show">
            <a-radio :value="true">true</a-radio>
            <a-radio :value="false">false</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>
    <a-modal
      width="60%"
      class="h-100 relative ml-auto"
      v-model:open="showColorModal"
      title="请选择字体颜色"
      :cancel-button-props="{ size: 0 }"
      @ok="showColorModal = false"
    >
      <div
        class="relative w-full h-85 slide-up item group flex cursor-pointer ease-linear duration-100 mb-5"
      >
        <img class="w-full h-full" :src="serverHost + formState.cover_img" :alt="formState.title" />
        <div
          class="w-90% flex flex-col flex-center text-center absolute top-50% left-50% translate--50% translate--50%"
        >
          <div class="text-10 font-bold" :style="{ color: formState.color || '#000' }">
            {{ formState.title }}
          </div>
          <div class="text-8" :style="{ color: formState.color || '#000' }">
            {{ formState.summary }}
          </div>
        </div>
      </div>
      <a-input v-model:value="formState.color" type="color" />
      <template #footer>
        <a-button key="submit" type="primary" @click="showColorModal = false">确定</a-button>
      </template>
    </a-modal>
    <a-modal
      width="60%"
      class="h-100 relative ml-auto"
      v-model:open="showColorModal4Edit"
      title="请选择字体颜色"
      :cancel-button-props="{ size: 0 }"
      @ok="showColorModal4Edit = false"
    >
      <div
        class="relative w-full h-85 slide-up item group flex cursor-pointer ease-linear duration-100 mb-5"
      >
        <img
          class="w-full h-full"
          :src="serverHost + editableData[editId]['cover_img']"
          :alt="editableData[editId]['title']"
        />
        <div
          class="w-90% flex flex-col flex-center text-center absolute top-50% left-50% translate--50% translate--50%"
        >
          <div
            class="text-10 font-bold"
            :style="{ color: editableData[editId]['color'] || '#000' }"
          >
            {{ editableData[editId]['title'] }}
          </div>
          <div class="text-8" :style="{ color: editableData[editId]['color'] || '#000' }">
            {{ editableData[editId]['summary'] }}
          </div>
        </div>
      </div>
      <a-input v-model:value="editableData[editId]['color']" type="color" />
      <template #footer>
        <a-button key="submit" type="primary" @click="showColorModal4Edit = false">确定</a-button>
      </template>
    </a-modal>
    <a-spin :spinning="loading">
      <a-table :columns="columns" :data-source="data">
        <template #bodyCell="{ column, text, record }">
          <template v-if="column.key === 'id'">
            <a :href="baseHost + '/posts/' + record.id" target="_blank">{{
              `${baseHost}/posts/${record.id}`
            }}</a>
          </template>
          <template v-if="column.key === 'cover_img'">
            <StaticUpload
              :image-url="record.cover_img"
              @update:imageUrl="(value) => (editableData[record.id].cover_img = value)"
              :authorization="userStore.token"
              v-if="editableData[record.id]"
            />
            <a-image :width="200" :src="serverHost + record.cover_img" v-else />
          </template>
          <template v-if="column.dataIndex === 'title'">
            <a-input
              v-model:value="editableData[record.id][column.dataIndex as keyof CarouselRequest]"
              v-if="editableData[record.id]"
            />
            <span v-else>{{ text }}</span>
          </template>
          <template v-if="column.dataIndex === 'color'">
            <div class="flex gap-x-1" v-if="editableData[record.id]">
              <span
                class="w-8 h-8 inline-block"
                :style="{
                  backgroundColor: editableData[record.id][
                    column.dataIndex as keyof CarouselRequest
                  ] as string
                }"
              ></span>
              <a-button @click="showColorModal4Edit = true">选择字体颜色</a-button>
            </div>
            <div class="w-full h-3 black_border" :style="{ backgroundColor: text }" v-else></div>
          </template>
          <template v-if="column.dataIndex === 'summary'">
            <a-textarea
              v-model:value="editableData[record.id][column.dataIndex as keyof CarouselRequest]"
              v-if="editableData[record.id]"
            />
            <span v-else>{{ text }}</span>
          </template>
          <template v-if="column.dataIndex === 'created_at'">
            {{ dayjs.unix(text).format('YYYY-MM-DD HH:mm:ss') }}
          </template>
          <template v-if="column.dataIndex === 'updated_at'">
            {{ dayjs.unix(text).format('YYYY-MM-DD HH:mm:ss') }}
          </template>
          <template v-if="column.key === 'show'">
            <a-switch
              v-model:checked="record.show"
              @change="changeShowStatus(record.id, record.show)"
            />
          </template>

          <template v-else-if="column.dataIndex === 'operation'">
            <div class="editable-row-operations">
              <span v-if="editableData[record.id]" class="flex gap-x-1">
                <a-typography-link @click="save(record.id)">保存</a-typography-link>
                <a-popconfirm title="确定取消？" @confirm="cancel(record.id)">
                  <a>取消</a>
                </a-popconfirm>
              </span>
              <span v-else class="flex gap-x-1">
                <a @click="edit(record.id)">编辑</a>
                <a-popconfirm
                  v-if="data.length"
                  title="确认删除？"
                  @confirm="deleteCarousel(record.id)"
                >
                  <a>删除</a>
                </a-popconfirm>
              </span>
            </div>
          </template>
        </template>
      </a-table>
    </a-spin> </a-card
  >>
</template>

<script lang="ts" setup>
import dayjs from 'dayjs'
import { computed, h, reactive, ref, type UnwrapRef } from 'vue'
import { type FormInstance, message } from 'ant-design-vue'
import {
  AddCarousel,
  type CarouselRequest,
  type CarouselVO,
  ChangeCarouselShowStatus,
  DeleteCarousel,
  GetCarousel,
  UpdateCarousel
} from '@/interfaces/Config'
import { cloneDeep } from 'lodash-es'
import originalAxios from 'axios'
import StaticUpload from '@/components/upload/StaticUpload.vue'
import { useUserStore } from '@/stores/user'
import { GetPost, type IPost, type PageRequest } from '@/interfaces/Post'
import { ReloadOutlined } from '@ant-design/icons-vue'

document.title = '轮播图配置 - 后台管理'
const userStore = useUserStore()
const columns = [
  {
    title: '封面',
    dataIndex: 'cover_img',
    key: 'cover_img'
  },
  {
    title: 'url',
    dataIndex: 'id',
    key: 'id'
  },
  {
    title: '标题',
    dataIndex: 'title',
    key: 'title'
  },
  {
    title: '摘要',
    dataIndex: 'summary',
    key: 'summary'
  },
  {
    title: '字体颜色',
    key: 'color',
    dataIndex: 'color'
  },
  {
    title: '是否显示',
    key: 'show',
    dataIndex: 'show'
  },
  {
    title: '发布时间',
    key: 'created_at',
    dataIndex: 'created_at'
  },
  {
    title: '更新时间',
    key: 'updated_at',
    dataIndex: 'updated_at'
  },
  {
    title: 'operation',
    dataIndex: 'operation'
  }
]

const formRef = ref<FormInstance>()
const visible = ref(false)
const formState = reactive<CarouselRequest>({
  id: '',
  title: '',
  summary: '',
  cover_img: '',
  color: '#000',
  show: true
})

const data = ref<CarouselVO[]>([])
const serverHost = import.meta.env.VITE_API_HOST
const baseHost = import.meta.env.VITE_BASE_HOST

const loading = ref(false)

const getCarousel = async () => {
  try {
    loading.value = true
    const response: any = await GetCarousel()
    data.value = response.data.data?.list || []
  } catch (error) {
    console.log(error)
  } finally {
    loading.value = false
  }
}

getCarousel()

// 编辑
const editableData: UnwrapRef<Record<string, CarouselRequest>> = reactive({})
const edit = (id: string) => {
  editableData[id] = cloneDeep(data.value.filter((item) => id === item.id)[0])
  editId.value = id
}

const save = async (id: string) => {
  const editableDatum = editableData[id]
  try {
    const response: any = await UpdateCarousel(id, editableDatum)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('更新成功')
    delete editableData[id]
    editId.value = ''
    await getCarousel()
  } catch (error) {
    console.log(error)
  }
}
const cancel = (key: string) => {
  delete editableData[key]
  editId.value = ''
}

const deleteCarousel = async (id: string) => {
  try {
    const response: any = await DeleteCarousel(id)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('删除成功')
    await getCarousel()
  } catch (error) {
    console.log(error)
  }
}

const addCarousel = () => {
  if (formRef.value) {
    formRef.value
      .validateFields()
      .then(async () => {
        try {
          const response: any = await AddCarousel(formState)
          if (response.data.code !== 0) {
            message.error(response.data.message)
            return
          }
          message.success('添加成功')
          visible.value = false
          if (formRef.value) {
            formRef.value.resetFields()
            formState.color = '#000'
          }
          await getCarousel()
        } catch (error) {
          if (originalAxios.isAxiosError(error)) {
            // 这是一个由 axios 抛出的错误
            if (error.response) {
              if (error.response.status === 409) {
                message.error('轮播图已存在')
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

const changeShowStatus = async (id: string, show: boolean) => {
  const response: any = await ChangeCarouselShowStatus(id, show)
  if (response.data.code !== 0) {
    message.error(response.data.message)
    return
  }
  message.success('更新成功')
  await getCarousel()
}

const showPosts = ref(false)
const posts = ref<IPost[]>([])
const postColumns = [
  {
    title: 'id',
    dataIndex: 'id',
    key: 'id'
  },
  {
    title: '标题',
    dataIndex: 'title',
    key: 'title'
  },
  {
    title: 'operation',
    dataIndex: 'operation'
  }
]

const showPostList = async () => {
  showPosts.value = true
  await getPosts()
}

const total = ref(0)

const pagination = computed(() => ({
  total: total.value,
  current: req.value.pageNo,
  pageSize: req.value.pageSize
}))
const req = ref<PageRequest>({
  pageNo: 1,
  pageSize: 5,
  sortField: 'create_time',
  sortOrder: 'desc'
} as PageRequest)
const choosePost = (record: IPost) => {
  formState.id = record.id
  formState.title = record.title
  formState.cover_img = record.cover_img
  showPosts.value = false
}
const change = (pg: any) => {
  req.value.pageNo = pg.current
  req.value.pageSize = pg.pageSize
  getPosts()
}
const getPosts = async () => {
  try {
    const response = await GetPost(req.value)
    posts.value = response.data.data?.list || []
    total.value = response.data.data?.totalCount || 0
  } catch (error) {
    console.log(error)
  }
}
const searchPost = () => {
  getPosts()
}

const showColorModal = ref(false)
const showColorModal4Edit = ref(false)
const editId = ref('')
</script>
