import Vue from 'vue'
import moment from 'moment-timezone'
import store from '../store'

Vue.filter('formatDate', (value, filterFormat) => {
  const { zone, format } = store.state.app.time

  if (value) {
    return moment(value).tz(zone).format(filterFormat || format || 'lll')
  }

  return ''
})
