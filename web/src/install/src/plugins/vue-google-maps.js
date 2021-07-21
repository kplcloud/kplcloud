import Vue from 'vue'
import * as VueGoogleMaps from 'vue2-google-maps'

import config from '../configs'

const { key } = config.maps

/**
 * Vue Google Maps plugin
 * https://github.com/Jeson-gk/vue2-google-maps
 */
Vue.use(VueGoogleMaps, {
  load: {
    // REPLACE key on configs/maps.js
    key,
    libraries: 'places' // This is required if you use the Autocomplete plugin
  },
  installComponents: true
})
