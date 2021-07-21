import Vue from 'vue'

/**
 * Copy to clipboard the text passed
 * @param {String} text string to copy
 * @param {String} toastText message to appear on the toast notification
 */
Vue.prototype.$clipboard = function (text, toastText = 'Copied to clipboard!') {
  const el = document.createElement('textarea')

  el.value = text
  document.body.appendChild(el)
  el.select()
  document.execCommand('copy')
  document.body.removeChild(el)

  if (this.$store && this.$store.dispatch) this.$store.dispatch('app/showToast', toastText)
}
