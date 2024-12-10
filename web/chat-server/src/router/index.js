import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: { name: 'Login' }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/access/Login.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/access/Register.vue')
  },
  {
    path: '/chat/owninfo',
    name: 'OwnInfo',
    component: () => import('../views/chat/user/OwnInfo.vue')
  },
  {
    path: '/chat/contactlist',
    name: 'ContactList',
    component: () => import('../views/chat/contact/ContactList.vue')
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  next()
})

export default router
