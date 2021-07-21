import Vue from 'vue'
import VueI18n from 'vue-i18n'

import config from '../configs'

const { locale, availableLocales, fallbackLocale } = config.locales

/**
 * Vue Translations
 * https://kazupon.github.io/vue-i18n/
 */
Vue.use(VueI18n)

const messages = {}

availableLocales.forEach((l) => { messages[l.code] = l.messages })

export const i18n = new VueI18n({
  locale,
  fallbackLocale,
  messages
})

i18n.locales = availableLocales

export default i18n
