import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/post',
      name: 'post',
      component: () => import('../views/post/PostListView.vue')
    },
    {
      path: '/post/edit',
      name: 'post-edit',
      component: () => import('../views/post/PostEditView.vue')
    },
    {
      path: '/category',
      name: 'category',
      component: () => import('../views/CategoryListView.vue')
    },
    {
      path: '/tag',
      name: 'tag',
      component: () => import('../views/TagListView.vue')
    },
    {
      path: '/friend',
      name: 'friend',
      component: () => import('@/views/friend/FriendView.vue')
    },
    {
      path: '/setting',
      name: 'setting',
      component: () => import('@/views/setting/SettingView.vue')
    }
  ]
})

export default router
