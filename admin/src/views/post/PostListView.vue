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
    <a-modal
      v-model:open="coverEditorVisible"
      title="编辑文章封面"
      ok-text="保存"
      cancel-text="取消"
      :confirm-loading="coverSaving"
      @ok="updateCover"
      @cancel="closeCoverEditor"
    >
      <a-form layout="vertical">
        <a-form-item label="文章">
          <a-input :value="currentPost?.title" disabled />
        </a-form-item>
        <a-form-item label="封面" required>
          <StaticUpload
            :image-url="coverImage"
            @update:imageUrl="(value) => (coverImage = value)"
            :authorization="userStore.token"
          />
        </a-form-item>
      </a-form>
    </a-modal>
    <a-spin :spinning="loading">
      <a-table
        :columns="columns"
        :data-source="posts"
        :pagination="pagination"
        :scroll="{ x: 1570 }"
        row-key="id"
        @change="change"
        bordered
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'id'">
            <div class="post-link-cell">
              <a :href="getPostUrl(record.id)" target="_blank" rel="noopener noreferrer">查看</a>
              <a-button type="link" size="small" @click="copyPostLink(record.id)">
                复制链接
              </a-button>
            </div>
          </template>
          <template v-if="column.key === 'cover_img'">
            <a-image
              :width="88"
              :height="56"
              :src="serverHost + record.cover_img"
              class="cover-img"
            />
          </template>
          <template v-else-if="column.key === 'summary'">
            <a-tooltip :title="record.summary">
              <div class="summary-text">{{ record.summary }}</div>
            </a-tooltip>
          </template>
          <template v-else-if="column.key === 'word_count'">
            <span>{{ formatPostStats(record.word_count) }}</span>
          </template>
          <template v-else-if="column.key === 'categories'">
            <span class="taxonomy-tags">
              <a-tag
                v-for="category in record.categories"
                :key="category.id"
                :color="category.name.length > 5 ? 'geekblue' : 'green'"
              >
                {{ category.name }}
              </a-tag>
            </span>
          </template>
          <template v-else-if="column.key === 'tags'">
            <span class="taxonomy-tags">
              <a-tag
                v-for="tag in record.tags"
                :key="tag.id"
                :color="tag.name.length > 5 ? 'geekblue' : 'green'"
              >
                {{ tag.name }}
              </a-tag>
            </span>
          </template>
          <template v-if="column.key === 'is_displayed'">
            <a-switch
              :checked="record.is_displayed"
              :loading="isDisplayUpdating(record.id)"
              @change="changeDisplayStatus(record, Boolean($event))"
            />
          </template>
          <template v-if="column.key === 'is_comment_allowed'">
            <a-switch
              :checked="record.is_comment_allowed"
              :loading="isCommentUpdating(record.id)"
              @change="changeCommentAllowedStatus(record, Boolean($event))"
            />
          </template>
          <template v-else-if="column.key === 'created_at' || column.key === 'updated_at'">
            <span>{{ dayjs.unix(record[column.key]).format('YYYY-MM-DD HH:mm:ss') }}</span>
          </template>
          <template v-else-if="column.dataIndex === 'operation'">
            <div class="post-actions">
              <span>
                <a-button
                  type="link"
                  size="small"
                  :loading="isContentCopying(record.id)"
                  @click="copyPostContent(record.id)"
                >
                  复制正文
                </a-button>
              </span>
              <span>
                <a-button
                  type="link"
                  size="small"
                  :loading="isWechatCopying(record.id)"
                  @click="copyPostWechatContent(record.id)"
                >
                  复制公众号格式
                </a-button>
              </span>
              <span>
                <a @click="router.push(`/home/post/draft/${record.id}`)">编辑</a>
              </span>
              <span>
                <a @click="openCoverEditor(record)">编辑封面</a>
              </span>
              <a-popconfirm v-if="posts.length" title="确认删除？" @confirm="deletePost(record)">
                <a-button type="link" size="small" danger :loading="isDeleting(record.id)">
                  删除
                </a-button>
              </a-popconfirm>
            </div>
          </template>
        </template>
      </a-table>
    </a-spin>
  </a-card>
</template>
<script lang="ts" setup>
import { ReloadOutlined } from '@ant-design/icons-vue'
import { computed, h, ref } from 'vue'
import {
  ChangeCommentAllowedStatus,
  ChangePostDisplayStatus,
  DeletePost,
  GetPost,
  GetPostById,
  UpdatePostCover,
  type IPost,
  type PageRequest
} from '@/interfaces/Post'
import router from '@/router'
import dayjs from 'dayjs'
import { message } from 'ant-design-vue'
import type { TableColumnType, TableProps } from 'ant-design-vue'
import { GetSelectedCategories, type SelectCategory } from '@/interfaces/Category'
import { GetSelectedTags, type SelectTag } from '@/interfaces/Tag'
import StaticUpload from '@/components/upload/StaticUpload.vue'
import { useUserStore } from '@/stores/user'
import { toErrorMessage } from '@/utils/error'
import { copyMarkdownAsWechat } from '@/utils/wechat'

document.title = '文章列表 - 后台管理'

const userStore = useUserStore()

const showSorterTooltip = ref('点击升序排序')
const columns = computed<TableColumnType[]>(() => {
  return [
    {
      title: '封面',
      dataIndex: 'cover_img',
      key: 'cover_img',
      width: 120
    },
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      width: 220,
      ellipsis: true
    },
    {
      title: '链接',
      dataIndex: 'id',
      key: 'id',
      width: 150
    },
    {
      title: '摘要',
      dataIndex: 'summary',
      key: 'summary',
      width: 260
    },
    {
      title: '字数',
      dataIndex: 'word_count',
      key: 'word_count',
      width: 150
    },
    {
      title: '分类',
      key: 'categories',
      dataIndex: 'categories',
      width: 180,
      filters: categories.value
    },
    {
      title: '标签',
      key: 'tags',
      dataIndex: 'tags',
      width: 180,
      filters: tags.value
    },
    {
      title: '显示',
      key: 'is_displayed',
      dataIndex: 'is_displayed',
      width: 90
    },
    {
      title: '评论',
      key: 'is_comment_allowed',
      dataIndex: 'is_comment_allowed',
      width: 90
    },
    {
      title: '发布时间',
      key: 'created_at',
      dataIndex: 'created_at',
      sorter: (p1: IPost, p2: IPost) => p1.created_at - p2.created_at,
      defaultSortOrder: 'descend',
      sortDirections: ['descend', 'ascend'],
      showSorterTooltip: { title: showSorterTooltip.value },
      width: 170
    },
    {
      title: '更新时间',
      key: 'updated_at',
      dataIndex: 'updated_at',
      width: 170
    },
    {
      title: '操作',
      dataIndex: 'operation',
      width: 310,
      fixed: 'right'
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
const coverEditorVisible = ref(false)
const coverSaving = ref(false)
const currentPost = ref<IPost | null>(null)
const coverImage = ref('')
const displayUpdatingIds = ref<Set<string>>(new Set())
const commentUpdatingIds = ref<Set<string>>(new Set())
const deletingIds = ref<Set<string>>(new Set())
const contentCopyingIds = ref<Set<string>>(new Set())
const wechatCopyingIds = ref<Set<string>>(new Set())

const setPending = (source: typeof displayUpdatingIds, id: string, pending: boolean) => {
  const next = new Set(source.value)
  if (pending) {
    next.add(id)
  } else {
    next.delete(id)
  }
  source.value = next
}

const isDisplayUpdating = (id: string) => displayUpdatingIds.value.has(id)

const isCommentUpdating = (id: string) => commentUpdatingIds.value.has(id)

const isDeleting = (id: string) => deletingIds.value.has(id)

const isContentCopying = (id: string) => contentCopyingIds.value.has(id)

const isWechatCopying = (id: string) => wechatCopyingIds.value.has(id)

const formatPostStats = (wordCount?: number) => {
  if (!wordCount) {
    return '-'
  }
  return `${wordCount} 字 / 约 ${Math.max(1, Math.ceil(wordCount / 400))} 分钟`
}

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
  if (isDeleting(record.id)) {
    return
  }
  try {
    setPending(deletingIds, record.id, true)
    const response: any = await DeletePost(record.id)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('删除成功')
    await getPosts()
  } catch (error) {
    message.error(toErrorMessage(error, '删除失败'))
  } finally {
    setPending(deletingIds, record.id, false)
  }
}

getPosts()

const changeDisplayStatus = async (record: IPost, isDisplayed: boolean) => {
  if (isDisplayUpdating(record.id)) {
    return
  }
  const previous = record.is_displayed
  try {
    record.is_displayed = isDisplayed
    setPending(displayUpdatingIds, record.id, true)
    const response: any = await ChangePostDisplayStatus(record.id, isDisplayed)
    if (response.data.code !== 0) {
      record.is_displayed = previous
      message.error(response.data.message)
      return
    }
    message.success('更新成功')
    await getPosts()
  } catch (error) {
    record.is_displayed = previous
    message.error(toErrorMessage(error, '更新显示状态失败'))
  } finally {
    setPending(displayUpdatingIds, record.id, false)
  }
}

const changeCommentAllowedStatus = async (record: IPost, isCommentAllowed: boolean) => {
  if (isCommentUpdating(record.id)) {
    return
  }
  const previous = record.is_comment_allowed
  try {
    record.is_comment_allowed = isCommentAllowed
    setPending(commentUpdatingIds, record.id, true)
    const response: any = await ChangeCommentAllowedStatus(record.id, isCommentAllowed)
    if (response.data.code !== 0) {
      record.is_comment_allowed = previous
      message.error(response.data.message)
      return
    }
    message.success('更新成功')
    await getPosts()
  } catch (error) {
    record.is_comment_allowed = previous
    message.error(toErrorMessage(error, '更新评论状态失败'))
  } finally {
    setPending(commentUpdatingIds, record.id, false)
  }
}

const searchPost = () => {
  req.value.pageNo = 1
  getPosts()
}

const openCoverEditor = (record: IPost) => {
  currentPost.value = record
  coverImage.value = record.cover_img
  coverEditorVisible.value = true
}

const closeCoverEditor = () => {
  currentPost.value = null
  coverImage.value = ''
}

const updateCover = async () => {
  if (!currentPost.value) {
    return
  }

  if (!coverImage.value) {
    message.warning('请选择封面')
    return
  }

  try {
    coverSaving.value = true
    const response: any = await UpdatePostCover(currentPost.value.id, coverImage.value)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('更新成功')
    coverEditorVisible.value = false
    closeCoverEditor()
    await getPosts()
  } catch (error) {
    console.log(error)
  } finally {
    coverSaving.value = false
  }
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

const getPostUrl = (id: string) => {
  return `${baseHost}/posts/${id}`
}

const copyPostLink = async (id: string) => {
  await navigator.clipboard.writeText(getPostUrl(id))
  message.success('链接复制成功')
}

const copyPostContent = async (id: string) => {
  if (isContentCopying(id)) {
    return
  }
  try {
    setPending(contentCopyingIds, id, true)
    const response = await GetPostById(id)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    const content = response.data.data?.content || ''
    await navigator.clipboard.writeText(content)
    message.success('正文复制成功')
  } catch (error) {
    message.error(toErrorMessage(error, '正文复制失败'))
  } finally {
    setPending(contentCopyingIds, id, false)
  }
}

const copyPostWechatContent = async (id: string) => {
  if (isWechatCopying(id)) {
    return
  }
  try {
    setPending(wechatCopyingIds, id, true)
    const response = await GetPostById(id)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    const post = response.data.data
    const content = post?.content || ''
    if (!content.trim()) {
      message.warning('正文为空，无法复制公众号格式')
      return
    }
    await copyMarkdownAsWechat(content, { author: post?.author })
    message.success('公众号格式已复制')
  } catch (error) {
    message.error(toErrorMessage(error, '公众号格式复制失败'))
  } finally {
    setPending(wechatCopyingIds, id, false)
  }
}
</script>

<style scoped>
.cover-img {
  object-fit: cover;
  border-radius: 4px;
}

.summary-text {
  display: -webkit-box;
  overflow: hidden;
  color: var(--app-text);
  line-height: 1.5;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

.taxonomy-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.post-link-cell,
.post-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  white-space: nowrap;
}

.post-link-cell :deep(.ant-btn) {
  padding: 0;
}

.post-actions :deep(.ant-btn) {
  padding: 0;
}
</style>
