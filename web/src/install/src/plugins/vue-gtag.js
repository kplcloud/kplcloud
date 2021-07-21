import Vue from 'vue'
import VueGtag from 'vue-gtag'
import router from '../router'
import config from '../configs'

const { gaId } = config.analytics

if (gaId) {
  Vue.use(VueGtag, {
    config: { id: gaId }
  }, router)
}
