import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/post',
      name: 'post',
      component: () => import('../views/PostListView.vue')
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
    }
  ]
})

export default router
