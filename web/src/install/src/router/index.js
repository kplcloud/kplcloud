import Vue from 'vue'
import Router from 'vue-router'
import NProgress from 'nprogress' // progress bar
// import { setDocumentTitle, domTitle } from '@/utils/domUtil'
import storage from 'store'
import { ACCESS_TOKEN } from '@/store/mutation-types'
// import {=} from 'vuetify'

// Routes
import ShowRouters from './show.routes'
// import store from '../store'

Vue.use(Router)

export const routes = [{
  path: '/',
  redirect: '/show'
},
  ...ShowRouters,
  {
    path: '/blank',
    name: 'blank',
    component: () => import(/* webpackChunkName: "blank" */ '@/pages/BlankPage.vue')
  },
  {
    path: '*',
    name: 'error',
    component: () => import(/* webpackChunkName: "error" */ '@/pages/error/NotFoundPage.vue'),
    meta: {
      layout: 'error'
    }
  }]

const router = new Router({
  mode: 'hash',
  base: process.env.BASE_URL || '/',
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) return savedPosition

    return { x: 0, y: 0 }
  },
  routes
})

NProgress.configure({ showSpinner: false }) // NProgress Configuration

const allowList = ['auth-signin', 'auth-signin-locked', 'error', 'login', 'show', 'show-page'] // no redirect allowList
const loginRoutePath = '/auth/signin'
const defaultRoutePath = '/dashboard/analytics'

/**
 * Before each route update
 */
router.beforeEach((to, from, next) => {
  NProgress.start() // start progress bar
  // to.meta && (typeof to.meta.title !== 'undefined' && setDocumentTitle(to.meta.title `- ${domTitle}`))
  if (storage.get(ACCESS_TOKEN)) {
    if (to.path === loginRoutePath) {
      const { redirect } = to.query
      if (redirect !== '') {
        next({ path: redirect })
      } else {
        next({ path: defaultRoutePath })
      }
      NProgress.done()
    } else {
      // const redirect = decodeURIComponent(from.query.redirect || to.path)
      // console.log(from.query.redirect, to.path, redirect)
      // if (to.path === redirect) {
      //   // set the replace: true so the navigation will not leave a history record
      //   next({ ...to, replace: true })
      // } else {
      //   // 跳转到目的路由
      //   next({ path: redirect })
      // }
      // NProgress.done()

      // check login user.roles is null
      // if (store.getters.roles.length === 0) {
      //   // request login userInfo
      //   store
      //     .dispatch('GetInfo')
      //     .then(res => {
      //       const roles = res.data && res.data.role
      //       // generate dynamic router
      //       store.dispatch('GenerateRoutes', { roles }).then(() => {
      //         // 根据roles权限生成可访问的路由表
      //         // 动态添加可访问路由表
      //         // router.addRoutes(store.getters.addRouters)
      //         // 请求带有 redirect 重定向时，登录自动重定向到该地址
      //         const redirect = decodeURIComponent(from.query.redirect || to.path)
      //         if (to.path === redirect) {
      //           // set the replace: true so the navigation will not leave a history record
      //           next({ ...to, replace: true })
      //         } else {
      //           // 跳转到目的路由
      //           next({ path: redirect })
      //         }
      //       })
      //     })
      //     .catch(() => {
      //       // notification.error({
      //       //   message: '错误',
      //       //   description: '请求用户信息失败，请重试'
      //       // })
      //       // 失败时，获取用户信息失败时，调用登出，来清空历史保留信息
      //       store.dispatch('Logout').then(() => {
      //         next({ path: loginRoutePath, query: { redirect: to.fullPath } })
      //       })
      //     })
      // } else {
      //   next()
      // }
    }
  } else {
    if (allowList.includes(to.name)) {
      // 在免登录名单，直接进入
      next()
    } else {
      next({ path: loginRoutePath, query: { redirect: to.fullPath } })
      NProgress.done() // if current page is login will not trigger afterEach hook, so manually handle it
    }
  }
  return next()
})

/**
 * After each route update
 */
router.afterEach(() => {
  NProgress.done() // finish progress bar
})

export default router
