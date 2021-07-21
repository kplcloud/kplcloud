import Vue from 'vue'
import Vuex from 'vuex'

// Global vuex
import AppModule from './app'
import user from './user'
// import actions from  './app/actions'

Vue.use(Vuex)

/**
 * Main Vuex Store
 */
const store = new Vuex.Store({
  modules: {
    app: AppModule,
    user: user
  }
})

export default store
