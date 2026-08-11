import { createRouter, createWebHistory } from 'vue-router'
import { pinia, useAppStore } from '../store/index.js'

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
  {
    path: '/chat/:id',
    name: 'ContactChat',
    component: () => import('../views/chat/contact/ContactChat.vue')
  },
  {
    path: '/chat/sessionList',
    name: 'SessionList',
    component: () => import('../views/chat/session/SessionList.vue')
  },
  {
    path: '/manager',
    name: 'Manager',
    component: () => import('../views/manager/Manager.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
});

router.beforeEach((to, from, next) => {
  const store = useAppStore(pinia)
  if (!store.userInfo.uuid || !store.token) {
    if (to.path === '/login' || to.path === '/register') {
      next()
      return
    }
    next('/login')
  } else {
    next()
  }
})

export default router
