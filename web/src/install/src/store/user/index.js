import storage from 'store'
// import { login, logout } from '@/api/auth.service'
import { ACCESS_TOKEN } from '@/store/mutation-types'
// import { welcome } from '@/utils/util'

const user = {
  namespaced: true,
  state: {
    token: '',
    name: '',
    welcome: '',
    avatar: '',
    roles: [],
    info: {}
  },

  mutations: {
    SET_TOKEN: (state, token) => {
      state.token = token
    },
    SET_NAME: (state, { name, welcome }) => {
      state.name = name
      state.welcome = welcome
    },
    SET_AVATAR: (state, avatar) => {
      state.avatar = avatar
    },
    SET_ROLES: (state, roles) => {
      state.roles = roles
    },
    SET_INFO: (state, info) => {
      state.info = info
    }
  },

  actions: {
    // 登录
    Login({ commit }, userInfo) {
      return new Promise((resolve, reject) => {
        console.log(userInfo, reject)
        // login(userInfo).then(response => {
        //   if(!response) {
            resolve()
        //     return
        //   }
        //   const result = response.data
        //   if (result.token === '') {
        //     resolve(result)
        //     return
        //   }
        //   storage.set(ACCESS_TOKEN, result.token, 7 * 24 * 60 * 60 * 1000)
        //   commit('SET_TOKEN', result.token)
        //   commit('SET_NAME', result.username)
        //   resolve(result)
        // }).catch(error => {
        //   reject(error)
        // })
      })
    },

    // 获取用户信息
    // GetInfo ({ commit }) {
    //   return new Promise((resolve, reject) => {
    //     getInfo().then(response => {
    //       const result = response.data
    //
    //       if (result.role && result.role.permissions.length > 0) {
    //         const role = result.role
    //         role.permissions = result.role.permissions
    //         role.permissions.map(per => {
    //           if (per.actionEntitySet != null && per.actionEntitySet.length > 0) {
    //             const action = per.actionEntitySet.map(action => { return action.action })
    //             per.actionList = action
    //           }
    //         })
    //         role.permissionList = role.permissions.map(permission => { return permission.permissionId })
    //         commit('SET_ROLES', result.role)
    //         commit('SET_INFO', result)
    //       } else {
    //         reject(new Error('getInfo: roles must be a non-null array !'))
    //       }
    //
    //       commit('SET_NAME', { name: result.username, welcome: welcome() })
    //       commit('SET_AVATAR', result.avatar)
    //
    //       resolve(response)
    //     }).catch(error => {
    //       reject(error)
    //     })
    //   })
    // },

    // 登出
    Logout ({ commit }) {
      return new Promise((resolve) => {
        // logout(state.token).then(() => {
          commit('SET_TOKEN', '')
          commit('SET_ROLES', [])
          commit('SET_INFO', {})
          commit('SET_NAME', {})
          commit('SET_AVATAR', '')
          storage.remove(ACCESS_TOKEN)
          resolve()
        // }).catch(() => {
        //   resolve()
        // }).finally(() => {
        // })
      })
    }

  }
}

export default user
