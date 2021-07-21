import Vue from 'vue'

Vue.filter('uppercase', (value) => {
  if (!value) return ''

  return value.toString().toUpperCase()
})
