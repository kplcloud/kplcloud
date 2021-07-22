export default [{
  path: '/install',
  name: 'install-page',
  component: () => import('@/pages/install/HomePage.vue'),
  meta: {
    layout: 'error'
  }
}]
