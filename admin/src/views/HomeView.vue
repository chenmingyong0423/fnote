<template>
  <a-layout>
    <a-layout-sider v-model:collapsed="collapsed" :trigger="null" collapsible>
      <div class="logo" />
      <a-menu
        v-model:openKeys="state.openKeys"
        v-model:selectedKeys="state.selectedKeys"
        mode="inline"
        theme="dark"
        :inline-collapsed="state.collapsed"
        :items="items"
        @click="itemClick"
      ></a-menu>
    </a-layout-sider>
    <a-layout>
      <a-layout-header style="background: #fff; padding: 0">
        <menu-unfold-outlined
          v-if="collapsed"
          class="trigger"
          @click="() => (collapsed = !collapsed)"
        />
        <menu-fold-outlined v-else class="trigger" @click="() => (collapsed = !collapsed)" />
      </a-layout-header>
      <a-layout-content class="p-5" :style="{ minHeight: '780px' }">
        <RouterView />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>
<script lang="ts" setup>
import { ref } from 'vue'
const collapsed = ref<boolean>(false)
import router from '@/router'

import { reactive, watch, h } from 'vue'
import { MenuFoldOutlined, MenuUnfoldOutlined, PieChartOutlined } from '@ant-design/icons-vue'
import { useRoute } from 'vue-router'
const route = useRoute()
const path = route.path

const state = reactive({
  collapsed: false,
  selectedKeys: [path],
  openKeys: ['dashboard']
})

const items = reactive([
  {
    key: 'dashboard',
    icon: () => h(PieChartOutlined),
    label: '仪表盘',
    title: '仪表盘',
    children: [
      {
        key: '/home/dashboard/traffic-stats',
        label: '流量统计',
        title: '流量统计'
      },
      {
        key: '/home/dashboard/content-stats',
        label: '内容发布统计',
        title: '内容发布统计'
      }
    ]
  },
  {
    key: 'sub post',
    icon: () => h(PieChartOutlined),
    label: '文章管理',
    title: '文章管理',
    children: [
      {
        key: '/home/post/list',
        label: '文章列表',
        title: '文章列表'
      },
      {
        key: '/home/post/draft/list',
        label: '草稿箱',
        title: '草稿箱'
      }
    ]
  },
  {
    key: 'sub comment',
    icon: () => h(PieChartOutlined),
    label: '评论管理',
    title: '评论管理',
    children: [
      {
        key: '/home/comment',
        label: '评论列表',
        title: '评论列表'
      }
    ]
  },
  {
    key: 'sub category',
    icon: () => h(PieChartOutlined),
    label: '分类管理',
    title: '分类管理',
    children: [
      {
        key: '/home/category',
        label: '分类列表',
        title: '分类列表'
      }
    ]
  },
  {
    key: 'sub tag',
    icon: () => h(PieChartOutlined),
    label: '标签管理',
    title: '标签管理',
    children: [
      {
        key: '/home/tag',
        label: '标签列表',
        title: '标签列表'
      }
    ]
  },
  {
    key: 'sub friend',
    icon: () => h(PieChartOutlined),
    label: '友链管理',
    title: '友链管理',
    children: [
      {
        key: '/home/friend',
        label: '友链列表',
        title: '友链列表'
      }
    ]
  },
  {
    key: 'sub blog',
    icon: () => h(PieChartOutlined),
    label: '系统',
    title: '系统',
    children: [
      {
        key: '/home/setting',
        label: '博客设置',
        title: '博客设置'
      },
      {
        key: '/home/backup',
        label: '备份',
        title: '备份'
      }
    ]
  }
])

watch(
  () => state.openKeys,
  (_val, oldVal) => {
    state.preOpenKeys = oldVal
  }
)

const itemClick = (item: any) => {
  router.push(item.key)
}
</script>
<style scoped>
.trigger {
  font-size: 18px;
  line-height: 64px;
  padding: 0 24px;
  cursor: pointer;
  transition: color 0.3s;
}

.trigger:hover {
  color: #1890ff;
}

.logo {
  height: 32px;
  background: rgba(255, 255, 255, 0.3);
  margin: 16px;
}
</style>
