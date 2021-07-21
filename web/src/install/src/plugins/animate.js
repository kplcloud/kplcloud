import Vue from 'vue'

/**
 * Animate html elements with the help of the library Animate.css
 * https://animate.style/
 *
 * @param {*} node
 * @param {String} animationName
 * @param {Function} callback when the animation ends
 */
Vue.prototype.$animate = function (node, animationName, callback) {
  node.classList.add('animate__animated', `animate__${animationName}`)

  function handleAnimationEnd() {
    node.classList.remove('animate__animated', `animate__${animationName}`)
    node.removeEventListener('animationend', handleAnimationEnd)

    // eslint-disable-next-line callback-return
    if (typeof callback === 'function') callback()
  }

  node.addEventListener('animationend', handleAnimationEnd)
}
