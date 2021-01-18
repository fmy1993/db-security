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
    path: '/user-home',
    name: 'UserHome',
    component: () => import('@/views/user/home/index'),
    redirect: '/user-home/patient',
    children: [
      {
        path: '/user-home/patient',
        name: 'Patient',
        component: () => import('@/views/user/home/component/patient')
      },
      {
        path: '/user-home/revise',
        name: 'Revise',
        component: () => import('@/views/user/home/component/revise')
      },
      {
        path: '/user-home/visualization',
        name: 'Visualization',
        component: () => import('@/views/user/home/component/visualization')
      }
    ]
  },
  {
    path: '/admin-home',
    name: 'AdminHome',
    component: () => import('@/views/admin/index'),
    redirect: '/admin-home/patient',
    children: [
      {
        path: '/admin-home/patient',
        name: 'AdminPatient',
        component: () => import('@/views/admin/component/patient')
      },
      {
        path: '/admin-home/revise',
        name: 'Revise',
        component: () => import('@/views/user/home/component/revise')
      },
      {
        path: '/admin-home/user',
        name: 'User',
        component: () => import('@/views/admin/component/user')
      },
      {
        path: '/admin-home/ip',
        name: 'Ip',
        component: () => import('@/views/admin/component/ip')
      },
      {
        path: '/admin-home/implant',
        name: 'Implant',
        component: () => import('@/views/admin/component/implant')
      },
      {
        path: '/admin-home/track',
        name: 'Track',
        component: () => import('@/views/admin/component/track')
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
