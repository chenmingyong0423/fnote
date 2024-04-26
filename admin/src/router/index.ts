import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { isInit } from '@/interfaces/Config'
import { message } from 'ant-design-vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: () => {
        return '/home/dashboard/traffic-stats'
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
      name: 'home',
      component: () => import('../views/HomeView.vue'),
      children: [
        {
          path: '/home/dashboard/traffic-stats',
          name: 'traffic stats',
          component: () => import('../views/dashboard/TrafficStatsView.vue')
        },
        {
          path: '/home/dashboard/content-stats',
          name: 'content stats',
          component: () => import('../views/dashboard/ContentStatsView.vue')
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
          path: '/home/post/draft/list',
          name: 'post-draft',
          component: () => import('../views/post/PostDraftListView.vue')
        },
        {
          path: '/home/post/draft/:id',
          name: 'edit-post',
          component: () => import('../views/post/EditPostDraftView.vue')
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

let flag = true
router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore()
  if (flag) {
    await isInit()
      .then((res) => {
        if (res.data.code === 0) {
          userStore.initialization = res.data.data.initStatus
        }
      })
      .catch((err) => {
        console.error(err)
      })
    flag = false
  }
  if (!userStore.initialization) {
    if (to.name !== 'init') {
      message.warn('网站未初始化，请初始化').then((r) => r)
      next({ name: 'init' })
    } else {
      next()
    }
    return
  }

  // 检查用户是否登录，示例逻辑
  if (to.name !== 'login' && !userStore.isLoggedIn) {
    // 如果用户未登录，重定向到登录页面
    next({ name: 'login' })
  } else if (to.name == 'init') {
    message.warn('网站已经初始化').then((r) => r)
    next({ name: 'home' })
  } else {
    // 如果用户已登录，或者访问的是登录页面，正常导航
    next()
  }
})

export default router
