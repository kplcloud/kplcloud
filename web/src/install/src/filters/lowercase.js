import Vue from 'vue'

Vue.filter('lowercase', (value) => {
  if (!value) return ''

  return value.toString().toLowerCase()
})
