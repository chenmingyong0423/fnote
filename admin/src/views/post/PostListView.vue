<template>
  <a-button type="primary" @click="router.push('/post/edit')">发布文章</a-button>
  <a-table :columns="columns" :data-source="posts" :pagination="pagination" @change="change">
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
        {{`https://chenmingyong.cn/posts/${record.id}`}}
      </template>
      <template v-if="column.key === 'cover_img'">
        <a-image :width="200" :src="record.cover_img" />
      </template>
      <template v-else-if="column.key === 'categories'">
        <span>
          <a-tag
            v-for="category in record.categories"
            :key="category"
            :color="category === 'loser' ? 'volcano' : category.length > 5 ? 'geekblue' : 'green'"
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
      <template v-else-if="column.key === 'create_time' || column.key === 'update_time'">
        <span>{{ dayjs.unix(record[column.key]).format('YYYY-MM-DD HH:mm:ss') }}</span>
      </template>
      <template v-else-if="column.dataIndex === 'operation'">
        <a-popconfirm v-if="data.length" title="确认删除？" @confirm="deletePost(record)">
          <a>删除</a>
        </a-popconfirm>
      </template>
    </template>
  </a-table>
</template>
<script lang="ts" setup>
import { SmileOutlined } from '@ant-design/icons-vue'
import axios from '@/http/axios'
import { computed, ref } from 'vue'
import type { IPost, PageRequest } from '@/interfaces/Post'
import type { IBaseResponse, IPageData, IResponse } from '@/interfaces/Common'
import router from '@/router'
import dayjs from 'dayjs'
import { message } from 'ant-design-vue'
import { template } from 'lodash-es'

const columns = [
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
    dataIndex: 'categories'
  },
  {
    title: '标签',
    key: 'tags',
    dataIndex: 'tags'
  },
  {
    title: '发布时间',
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

const data = [
  {
    key: '1',
    coverImg: 'John Brown',
    title: 32,
    summary: 'New York No. 1 Lake Park',
    categories: ['nice', 'developer'],
    tags: ['nice', 'developer'],
    createTime: 121123123312132,
    updateTime: 121123123312132
  },
  {
    key: '2',
    coverImg: 'John Brown',
    title: 32,
    summary: 'New York No. 1 Lake Park',
    categories: ['nice', 'developer'],
    tags: ['nice', 'developer'],
    createTime: 121123123312132,
    updateTime: 121123123312132
  },
  {
    key: '3',
    coverImg: 'John Brown',
    title: 32,
    summary: 'New York No. 1 Lake Park',
    categories: ['nice', 'developer'],
    tags: ['nice', 'developer'],
    createTime: 121123123312132,
    updateTime: 121123123312132
  },
  {
    key: '4',
    coverImg: 'John Brown',
    title: 32,
    summary: 'New York No. 1 Lake Park',
    categories: ['nice', 'developer'],
    tags: ['nice', 'developer'],
    createTime: 121123123312132,
    updateTime: 121123123312132
  },
  {
    key: '5',
    coverImg: 'John Brown',
    title: 32,
    summary: 'New York No. 1 Lake Park',
    categories: ['nice', 'developer'],
    tags: ['nice', 'developer'],
    createTime: 121123123312132,
    updateTime: 121123123312132
  },
  {
    key: '6',
    coverImg: 'John Brown',
    title: 32,
    summary: 'New York No. 1 Lake Park',
    categories: ['nice', 'developer'],
    tags: ['nice', 'developer'],
    createTime: 121123123312132,
    updateTime: 121123123312132
  }
]

const req = ref<PageRequest>({
  pageNo: 1,
  pageSize: 5,
  sortField: 'create_time',
  sortOrder: 'desc'
} as PageRequest)

const posts = ref<IPost[]>([])

const total = ref(0)

const pagination = computed(() => ({
  total: total.value,
  current: req.value.pageNo,
  pageSize: req.value.pageSize
}))

const change = (pg, filters, sorter, { currentDataSource }) => {
  req.value.pageNo = pg.current
  req.value.pageSize = pg.pageSize
  getPosts()
}

const getPosts = async () => {
  try {
    const response = await axios.get<IResponse<IPageData<IPost>>>('/admin/posts', {
      params: req.value
    })
    posts.value = response.data.data?.list || []
    total.value = response.data.data?.totalCount || 0
  } catch (error) {
    console.log(error)
  }
}

const deletePost = async (record: IPost) => {
  try {
    console.log(record)
    const response = await axios.delete<IBaseResponse>(`/admin/posts/${record.id}`)
    if (response.data.code !== 200) {
      message.error(response.data.message)
      return
    }
    message.success('删除成功')
    await getPosts()
  } catch (error) {
    console.log(error)
    message.error('删除失败')
  }
}

getPosts()
</script>
