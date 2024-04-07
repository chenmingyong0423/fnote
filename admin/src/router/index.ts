import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { isInit } from '@/interfaces/Config'
import { login } from '@/interfaces/User'
import { message } from 'ant-design-vue'

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
      path: '/init',
      name: 'init',
      component: () => import('../views/initialization/InitView.vue')
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
        },
        {
          path: '/home/backup',
          name: 'backup',
          component: () => import('@/views/backup/BackupView.vue')
        }
      ]
    }
  ]
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore();
  if (!userStore.isInit) {
    if (to.name !== 'init') {
      next({ name: 'init' });
    } else {
      next();
    }
  }

  // 检查用户是否登录，示例逻辑
  if (to.name !== 'login' && !userStore.isLoggedIn) {
    // 如果用户未登录，重定向到登录页面
    next({ name: 'login' });
  } else {
    // 如果用户已登录，或者访问的是登录页面，正常导航
    next();
  }
})

export default router
