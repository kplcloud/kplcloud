export default [{
  path: '/show',
  name: 'show-page',
  component: () => import(/* webpackChunkName: "auth-signin" */ '@/pages/show/HomePage.vue'),
  meta: {
    layout: 'error'
  }
}]
