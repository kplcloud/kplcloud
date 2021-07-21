import Vue from 'vue'

// For full framework
import Vuetify from 'vuetify/lib/framework'
// For a-la-carte components - https://vuetifyjs.com/en/customization/a-la-carte/
// import Vuetify from 'vuetify/lib'

import * as directives from 'vuetify/lib/directives'
import i18n from './vue-i18n'
import config from '../configs'

/**
 * Vuetify Plugin
 * Main components library
 *
 * https://vuetifyjs.com/
 *
 */
Vue.use(Vuetify, {
  directives
})

export default new Vuetify({
  rtl: config.theme.isRTL,
  theme: {
    dark: config.theme.globalTheme === 'dark',
    options: {
      customProperties: true
    },
    themes: {
      dark: config.theme.dark,
      light: config.theme.light
    }
  },
  lang: {
    current: config.locales.locale,
    // To use Vuetify own translations without Vue-i18n comment the next line
    t: (key, ...params) => i18n.t(key, params)
  }
})
