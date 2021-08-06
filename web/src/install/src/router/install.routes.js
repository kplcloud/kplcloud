export default [{
  path: '/install',
  redirect: '/install/db-step'
}, {
  path: '/install/:step',
  name: 'install-page',
  component: () => import('@/pages/install/HomePage.vue'),
  meta: {
    layout: 'error'
  }
}]
