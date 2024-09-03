<template>
  <a-card title="文章列表">
    <template #extra>
      <div class="flex gap-x-3">
        <a-tooltip title="刷新数据">
          <a-button shape="circle" :icon="h(ReloadOutlined)" :loading="loading" @click="getPosts" />
        </a-tooltip>
      </div>
    </template>
    <div>
      <div>
        <a-button type="primary" @click="router.push('/home/post')" class="mb-3">发布文章</a-button>
        <a-input-search
          v-model:value="req.keyword"
          placeholder="请输入关键字"
          style="width: 200px"
          @search="searchPost"
          @pressEnter="searchPost"
          allow-clear
          class="float-right"
        />
      </div>
    </div>
    <a-spin :spinning="loading">
      <a-table
        :columns="columns"
        :data-source="posts"
        :pagination="pagination"
        @change="change"
        bordered
      >
        <template #headerCell="{ column }">
          <template v-if="column.key === 'name'">
            <span>
              <smile-outlined />
              Name
            </span>
          </template>
        </template>

        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'id'">
            <a :href="baseHost + '/posts/' + record.id" target="_blank">{{
              `${baseHost}/posts/${record.id}`
            }}</a>
          </template>
          <template v-if="column.key === 'cover_img'">
            <a-image :width="200" :src="serverHost + record.cover_img" />
          </template>
          <template v-else-if="column.key === 'categories'">
            <span>
              <a-tag
                v-for="category in record.categories"
                :key="category"
                :color="
                  category === 'loser' ? 'volcano' : category.length > 5 ? 'geekblue' : 'green'
                "
              >
                {{ category.name.toUpperCase() }}
              </a-tag>
            </span>
          </template>
          <template v-else-if="column.key === 'tags'">
            <span>
              <a-tag
                v-for="tag in record.tags"
                :key="tag"
                :color="tag === 'loser' ? 'volcano' : tag.length > 5 ? 'geekblue' : 'green'"
              >
                {{ tag.name.toUpperCase() }}
              </a-tag>
            </span>
          </template>
          <template v-if="column.key === 'is_displayed'">
            <a-switch
              v-model:checked="record.is_displayed"
              @change="changeDisplayStatus(record.id, record.is_displayed)"
            />
          </template>
          <template v-if="column.key === 'is_comment_allowed'">
            <a-switch
              v-model:checked="record.is_comment_allowed"
              @change="changeCommentAllowedStatus(record.id, record.is_comment_allowed)"
            />
          </template>
          <template v-else-if="column.key === 'created_at' || column.key === 'updated_at'">
            <span>{{ dayjs.unix(record[column.key]).format('YYYY-MM-DD HH:mm:ss') }}</span>
          </template>
          <template v-else-if="column.dataIndex === 'operation'">
            <div class="flex gap-x-1">
              <span>
                <a @click="copyContent(record.id)">复制</a>
              </span>
              <span>
                <a @click="router.push(`/home/post/draft/${record.id}`)">编辑</a>
              </span>
              <a-popconfirm v-if="posts.length" title="确认删除？" @confirm="deletePost(record)">
                <a>删除</a>
              </a-popconfirm>
            </div>
          </template>
        </template>
      </a-table>
    </a-spin>
  </a-card>
</template>
<script lang="ts" setup>
import { ReloadOutlined, SmileOutlined } from '@ant-design/icons-vue'
import { computed, h, ref } from 'vue'
import {
  ChangeCommentAllowedStatus,
  ChangePostDisplayStatus,
  DeletePost,
  GetPost,
  GetPostById,
  GetPostBySug,
  type IPost,
  type PageRequest
} from '@/interfaces/Post'
import router from '@/router'
import dayjs from 'dayjs'
import { message } from 'ant-design-vue'
import type { TableColumnType, TableProps } from 'ant-design-vue'
import { GetSelectedCategories, type SelectCategory } from '@/interfaces/Category'
import { GetSelectedTags, type SelectTag } from '@/interfaces/Tag'

document.title = '文章列表 - 后台管理'

const showSorterTooltip = ref('点击升序排序')
const columns = computed<TableColumnType[]>(() => {
  return [
    {
      title: '封面',
      dataIndex: 'cover_img',
      key: 'cover_img'
    },
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title'
    },
    {
      title: 'url',
      dataIndex: 'id',
      key: 'id'
    },
    {
      title: '摘要',
      dataIndex: 'summary',
      key: 'summary'
    },
    {
      title: '分类',
      key: 'categories',
      dataIndex: 'categories',
      filters: categories.value
    },
    {
      title: '标签',
      key: 'tags',
      dataIndex: 'tags',
      filters: tags.value
    },
    {
      title: '是否显示',
      key: 'is_displayed',
      dataIndex: 'is_displayed'
    },
    {
      title: '是否允许评论',
      key: 'is_comment_allowed',
      dataIndex: 'is_comment_allowed'
    },
    {
      title: '发布时间',
      key: 'created_at',
      dataIndex: 'created_at',
      sorter: (p1: IPost, p2: IPost) => p1.created_at - p2.created_at,
      defaultSortOrder: 'descend',
      sortDirections: ['descend', 'ascend'],
      showSorterTooltip: { title: showSorterTooltip.value }
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
})

const serverHost = import.meta.env.VITE_API_HOST
const baseHost = import.meta.env.VITE_BASE_HOST
const req = ref<PageRequest>({
  pageNo: 1,
  pageSize: 5,
  sortField: 'created_at',
  sortOrder: 'DESC',
  keyword: ''
} as PageRequest)

const posts = ref<IPost[]>([])

const total = ref(0)

const pagination = computed(() => ({
  total: total.value,
  current: req.value.pageNo,
  pageSize: req.value.pageSize
}))

const change: TableProps<IPost>['onChange'] = (pagination, filters, sorter: any) => {
  req.value.pageNo = <number>pagination.current
  req.value.pageSize = <number>pagination.pageSize
  req.value.sortField = sorter.field
  switch (sorter.order) {
    case 'ascend':
      req.value.sortOrder = 'ASC'
      showSorterTooltip.value = '点击默认排序'
      break
    case 'descend':
      req.value.sortOrder = 'DESC'
      showSorterTooltip.value = '点击升序排序'
      break
    default:
      req.value.sortOrder = 'DESC'
      showSorterTooltip.value = '点击降序排序'
  }
  req.value.category_filter = filters.categories as string[]
  req.value.tag_filter = filters.tags as string[]
  console.log(req.value)
  getPosts()
}

const loading = ref(false)

const getPosts = async () => {
  try {
    loading.value = true
    const response = await GetPost(req.value)
    posts.value = response.data.data?.list || []
    total.value = response.data.data?.totalCount || 0
  } catch (error) {
    console.log(error)
  } finally {
    loading.value = false
  }
}

const deletePost = async (record: IPost) => {
  try {
    console.log(record)
    const response: any = await DeletePost(record.id)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('删除成功')
    await getPosts()
  } catch (error) {
    console.log(error)
  }
}

getPosts()

const changeDisplayStatus = async (id: string, is_displayed: boolean) => {
  try {
    const response: any = await ChangePostDisplayStatus(id, is_displayed)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('更新成功')
    await getPosts()
  } catch (error) {
    console.log(error)
  }
}

const changeCommentAllowedStatus = async (id: string, is_comment_allowed: boolean) => {
  try {
    const response: any = await ChangeCommentAllowedStatus(id, is_comment_allowed)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('更新成功')
    await getPosts()
  } catch (error) {
    console.log(error)
  }
}

const searchPost = () => {
  getPosts()
}

interface Filter {
  text: string
  value: string
}

const categories = ref<Filter[]>([])

const getCategories = async () => {
  try {
    const response = await GetSelectedCategories()
    response.data.data?.list.forEach((item: SelectCategory) => {
      categories.value?.push({
        text: item.label,
        value: item.value
      })
    })
  } catch (error) {
    console.log(error)
  }
}
getCategories()

const tags = ref<Filter[]>([])
const getTags = async () => {
  try {
    const response = await GetSelectedTags()
    response.data.data?.list.forEach((item: SelectTag) => {
      tags.value?.push({
        text: item.label,
        value: item.value
      })
    })
  } catch (error) {
    console.log(error)
  }
}
getTags()

const copyContent = async (id: string) => {
  const response = await GetPostById(id)
  if (response.data.code !== 0) {
    message.error(response.data.message)
    return
  }
  const content = response.data.data?.content
  await navigator.clipboard.writeText(content)
  message.success('复制成功')
}
</script>
