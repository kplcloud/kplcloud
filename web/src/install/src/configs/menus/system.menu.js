//https://pictogrammers.github.io/@mdi/font/5.4.55/ icon
export default [
  {
    icon: 'mdi-card-bulleted-settings-outline', key: 'menu.systemPage', text: 'System Pages', regex: /^\/system/,
    items: [
      { icon: 'mdi-file-settings-outline', key: 'menu.systemSettingPage', text: 'System Setting Pages', link: '/system/settings' },
      { icon: 'mdi-account-details-outline', key: 'menu.systemUserPage', text: 'System User Pages', link: '/system/users' },
      { icon: 'mdi-account-group-outline', key: 'menu.systemRolePage', text: 'System Role Pages', link: '/system/roles' },
      { icon: 'mdi-group', key: 'menu.systemNamespacePage', text: 'System Namespace Pages', link: '/system/namespaces' },
      { icon: 'mdi-api', key: 'menu.systemPermissionPage', text: 'System Permission Pages', link: '/system/permissions' }
    ]
  }
]
