<template>
  <a-table :columns="columns" :data-source="data">
    <template #headerCell="{ column }">
      <template v-if="column.key === 'name'">
        <span>
          <smile-outlined />
          Name
        </span>
      </template>
    </template>

    <template #bodyCell="{ column, record }">
      <template v-if="column.key === 'coverImg'">
        <a-image
          :width="200"
          src="https://zos.alipayobjects.com/rmsportal/jkjgkEfvpUPVyRjUImniVslZfWPnJuuZ.png"
        />
      </template>
      <template v-else-if="column.key === 'categories'">
        <span>
          <a-tag
            v-for="category in record.categories"
            :key="category"
            :color="category === 'loser' ? 'volcano' : category.length > 5 ? 'geekblue' : 'green'"
          >
            {{ category.toUpperCase() }}
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
            {{ tag.toUpperCase() }}
          </a-tag>
        </span>
      </template>
    </template>
  </a-table>
</template>
<script lang="ts" setup>
import { SmileOutlined, DownOutlined } from '@ant-design/icons-vue'
import axios from '@/http/axios'
import { ref } from 'vue'
import type { IPost, PageRequest } from '@/interfaces/Post'
import type { IResponse } from '@/interfaces/Common'

const columns = [
  {
    title: '封面',
    dataIndex: 'coverImg',
    key: 'coverImg'
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
    key: 'createTime',
    dataIndex: 'createTime'
  },
  {
    title: '最后一次修改的时间',
    key: 'updateTime',
    dataIndex: 'updateTime'
  },
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
    updateTime: 121123123312132,
  },
  {
    key: '2',
    coverImg: 'John Brown',
    title: 32,
    summary: 'New York No. 1 Lake Park',
    categories: ['nice', 'developer'],
    tags: ['nice', 'developer'],
    createTime: 121123123312132,
    updateTime: 121123123312132,
  },
  {
    key: '3',
    coverImg: 'John Brown',
    title: 32,
    summary: 'New York No. 1 Lake Park',
    categories: ['nice', 'developer'],
    tags: ['nice', 'developer'],
    createTime: 121123123312132,
    updateTime: 121123123312132,
  },
  {
    key: '4',
    coverImg: 'John Brown',
    title: 32,
    summary: 'New York No. 1 Lake Park',
    categories: ['nice', 'developer'],
    tags: ['nice', 'developer'],
    createTime: 121123123312132,
    updateTime: 121123123312132,
  },
  {
    key: '5',
    coverImg: 'John Brown',
    title: 32,
    summary: 'New York No. 1 Lake Park',
    categories: ['nice', 'developer'],
    tags: ['nice', 'developer'],
    createTime: 121123123312132,
    updateTime: 121123123312132,
  },{
    key: '6',
    coverImg: 'John Brown',
    title: 32,
    summary: 'New York No. 1 Lake Park',
    categories: ['nice', 'developer'],
    tags: ['nice', 'developer'],
    createTime: 121123123312132,
    updateTime: 121123123312132,
  },

]

const req = ref<PageRequest>({
  pageNo: 1,
  pageSize: 10,
  sortField: "create_time",
  sortOrder: "desc",
} as PageRequest)

const posts = ref<IPost[]>([]);

const getPosts = async () => {
  try {
    const response = await axios.get<IResponse<IPost>>('/admin/posts', {
      params: req.value
    });
    posts.value = response.data.data?.list || [];
    console.log(posts.value)
  } catch (error) {
    console.log(error);
  }
};

getPosts();
</script>

