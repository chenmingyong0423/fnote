import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'index',
      component: () => import('../views/IndexView.vue')
    },
    {
      path: '/post/list',
      name: 'post',
      component: () => import('../views/post/PostListView.vue')
    },
    {
      path: '/post',
      name: 'add-post',
      component: () => import('../views/post/AddPostView.vue')
    },
    {
      path: '/drafts/:id',
      name: 'edit-post',
      component: () => import('../views/post/UpdatePostView.vue')
    },
    {
      path: '/comment',
      name: 'comment',
      component: () => import('../views/comment/CommentView.vue')
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
