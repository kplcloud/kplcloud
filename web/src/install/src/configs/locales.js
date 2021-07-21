import en from '../translations/en'
import zh from '../translations/zh'

const supported = ['en', 'zh']
let locale = 'zh'

try {
  // get browser default language
  const { 0: browserLang } = navigator.language.split('-')

  if (supported.includes(browserLang)) locale = browserLang
} catch (e) {
  // console.log(e)
}

export default {
  // current locale
  locale,

  // when translation is not available fallback to that locale
  fallbackLocale: 'zh',

  // availabled locales for user selection
  availableLocales: [{
    code: 'en',
    flag: 'us',
    label: 'English',
    messages: en
  }, {
    code: 'zh',
    flag: 'cn',
    label: '中文',
    messages: zh
  }]
}
