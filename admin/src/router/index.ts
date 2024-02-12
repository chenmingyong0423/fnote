import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: () => {
        return '/home/index'
      }
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/login/LoginView.vue')
    },
    {
      path: '/home',
      name: 'base',
      component: () => import('../views/HomeView.vue'),
      children: [
        {
          path: '/home/index',
          name: 'home',
          component: () => import('../views/IndexView.vue')
        },
        {
          path: '/home/post/list',
          name: 'post',
          component: () => import('../views/post/PostListView.vue')
        },
        {
          path: '/home/post',
          name: 'add-post',
          component: () => import('../views/post/AddPostView.vue')
        },
        {
          path: '/home/drafts/:id',
          name: 'edit-post',
          component: () => import('../views/post/UpdatePostView.vue')
        },
        {
          path: '/home/comment',
          name: 'comment',
          component: () => import('../views/comment/CommentView.vue')
        },
        {
          path: '/home/category',
          name: 'category',
          component: () => import('../views/CategoryListView.vue')
        },
        {
          path: '/home/tag',
          name: 'tag',
          component: () => import('../views/TagListView.vue')
        },
        {
          path: '/home/friend',
          name: 'friend',
          component: () => import('@/views/friend/FriendView.vue')
        },
        {
          path: '/home/setting',
          name: 'setting',
          component: () => import('@/views/setting/SettingView.vue')
        }
      ]
    }
  ]
})

export default router
