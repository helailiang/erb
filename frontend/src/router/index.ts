import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { getToken } from '@/utils/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/Login.vue'),
    meta: {
      title: '登录',
      requiresAuth: false,
    },
  },
  {
    path: '/',
    component: () => import('@/layouts/DefaultLayout.vue'),
    redirect: '/dashboard',
    meta: {
      requiresAuth: true,
    },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/Dashboard.vue'),
        meta: {
          title: '仪表盘',
        },
      },
      {
        path: 'users',
        name: 'UserList',
        component: () => import('@/views/user/UserList.vue'),
        meta: {
          title: '用户管理',
        },
      },
      {
        path: 'reports',
        name: 'ReportList',
        component: () => import('@/views/report/ReportList.vue'),
        meta: {
          title: '日报管理',
        },
      },
      {
        path: 'reports/create',
        name: 'ReportCreate',
        component: () => import('@/views/report/ReportCreate.vue'),
        meta: {
          title: '写日报',
        },
      },
      {
        path: 'reports/:id',
        name: 'ReportDetail',
        component: () => import('@/views/report/ReportDetail.vue'),
        meta: {
          title: '日报详情',
        },
      },
      {
        path: 'reports/:id/edit',
        name: 'ReportEdit',
        component: () => import('@/views/report/ReportEdit.vue'),
        meta: {
          title: '编辑日报',
        },
      },
      {
        path: 'statistics',
        name: 'Statistics',
        component: () => import('@/views/statistics/Statistics.vue'),
        meta: {
          title: '统计分析',
        },
      },
      {
        path: 'roles',
        name: 'RoleList',
        component: () => import('@/views/role/RoleList.vue'),
        meta: {
          title: '角色管理',
        },
      },
      {
        path: 'roles',
        name: 'RoleList',
        component: () => import('@/views/role/RoleList.vue'),
        meta: {
          title: '角色管理',
        },
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  const token = getToken()

  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/')
  } else {
    next()
  }
})

export default router

