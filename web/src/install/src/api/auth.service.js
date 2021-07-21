import request from '@/utils/request'

const authApi = {
  Login: '/admin/auth/login',
  Logout: '/admin/auth/login'
}

/**
 * login func
 * parameter: {
 *     username: '',
 *     password: '',
 *     remember_me: true,
 *     captcha: '12345'
 * }
 * @param parameter
 * @returns {*}
 */
export function login(parameter) {
  return request({
    url: authApi.Login,
    method: 'post',
    data: parameter
  })
}

export function logout(parameter) {
  return request({
    url: authApi.Login,
    method: 'post',
    data: parameter
  })
}
