import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/register/index')
  },
  {
    path: '/home',
    name: 'UserHome',
    component: () => import('@/views/userhome'),
    redirect: '/user-home/staff',
    children: [
      {
        path: '/user-home/staff',
        name: 'Staff',
        component: () => import('@/views/staff/index')
      },
      {
        path: '/user-home/visualization',
        name: 'Visualization',
        component: () => import('@/views/userhome/component/visualization')
      }
    ]
  },
  {
    path: '/admin-home',
    name: 'AdminHome',
    component: () => import('@/views/adminhome/index'),
    redirect: '/admin-home/staff',
    children: [
      {
        path: '/admin-home/staff',
        name: 'AdminStaff',
        component: () => import('@/views/adminhome/component/staff')
      },
      {
        path: '/admin-home/revise',
        name: 'Revise',
        component: () => import('@/views/userhome/component/revise')
      },
      {
        path: '/admin-home/user',
        name: 'User',
        component: () => import('@/views/adminhome/component/user')
      },
      {
        path: '/admin-home/ip',
        name: 'Ip',
        component: () => import('@/views/adminhome/component/ip')
      },
      {
        path: '/admin-home/track',
        name: 'Track',
        component: () => import('@/views/adminhome/component/track')
      },
      {
        path: '/admin-home/record',
        name: 'Record',
        component: () => import('@/views/adminhome/component/record')
      }
    ]
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

router.beforeEach((to, from, next) => {
  if (to.path === '/login') {
    next()
  } else if (to.path === '/register') {
    next()
  } else {
    const token = localStorage.getItem('Authorization')
    if (token === null || token === '') {
      next('/login')
    } else {
      next()
    }
  }
})

export default router
