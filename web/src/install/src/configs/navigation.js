// import menuPages from './menus/pages.menu'
export default {
  // main navigation - side menu
  menu: [{
    text: '',
    key: '',
    items: [
      { icon: 'mdi-view-dashboard-outline', key: 'menu.dashboard', text: 'Show', link: '/show' }
    ]
  }],

  // footer links
  footer: [{
    text: 'Docs',
    key: 'menu.docs',
    href: '/',
    target: '_self'
  }]
}
