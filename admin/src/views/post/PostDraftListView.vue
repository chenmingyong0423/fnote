<template>
  <a-card title="草稿箱">
    <template #extra>
      <div class="flex gap-x-3">
        <a-tooltip title="刷新数据">
          <a-button
            shape="circle"
            :icon="h(ReloadOutlined)"
            :loading="loading"
            @click="getPostDrafts"
          />
        </a-tooltip>
      </div>
    </template>
    <a-spin :spinning="loading">
      <a-list
        class="demo-loadmore-list bg-white"
        item-layout="horizontal"
        :pagination="pagination"
        :data-source="listData"
      >
        <template #renderItem="{ item }">
          <a-list-item>
            <template #actions>
              <a key="list-loadmore-edit" @click="router.push(`/home/post/draft/${item.id}`)"
                >编辑</a
              >
              <a-popconfirm title="确认删除？" @confirm="deletePostDraft(item.id)">
                <a>删除</a>
              </a-popconfirm>
            </template>
            <a-skeleton avatar :title="false" :loading="!!item.loading" active>
              <a-list-item-meta
                :description="dayjs.unix(item.created_at).format('YYYY-MM-DD HH:mm:ss')"
              >
                <template #title>
                  {{ item.title }}
                </template>
              </a-list-item-meta>
            </a-skeleton>
          </a-list-item>
        </template>
      </a-list>
    </a-spin>
  </a-card>
</template>
<script lang="ts" setup>
import { h, ref } from 'vue'
import { DeletePostDraftById, GetPostDraft, type PageRequest } from '@/interfaces/Post'
import router from '@/router'
import dayjs from 'dayjs'
import { message } from 'ant-design-vue'
import type { PostDraftBrief } from '@/interfaces/PostDraft'
import { ReloadOutlined } from '@ant-design/icons-vue'

document.title = '草稿箱 - 后台管理'

const req = ref<PageRequest>({
  pageNo: 1,
  pageSize: 5,
  sortField: 'created_at',
  sortOrder: 'DESC'
} as PageRequest)

const listData = ref<PostDraftBrief[]>([])

const pagination = {
  onChange: (page: number) => {
    req.value.pageNo = page
    getPostDrafts()
  },
  pageSize: 5,
  total: 0
}
const loading = ref(false)
const getPostDrafts = async () => {
  try {
    loading.value = true
    const response = await GetPostDraft(req.value)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    listData.value = response.data.data?.list || []
    pagination.total = response.data.data?.totalCount || 0
  } catch (error) {
    console.log(error)
  } finally {
    loading.value = false
  }
}

const deletePostDraft = async (id: string) => {
  try {
    const response: any = await DeletePostDraftById(id)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('删除成功')
    await getPostDrafts()
  } catch (error) {
    console.log(error)
  }
}

getPostDrafts()
</script>

<style scoped>
.demo-loadmore-list {
  min-height: 350px;
}
</style>
