import request from '@/utils/request'

const installApi = {
  InitDb: '/install/init-db',
  InitPlatform: '/install/init-platform',
  InitLogo: '/install/init-logo',
  InitCors: '/install/init-cors',
  InitRedis: '/install/init-redis',
}

export function initDb (parameter) {
  return request({
    url: installApi.InitDb,
    method: 'post',
    data: parameter
  })
}

export function initPlatform (parameter) {
  return request({
    url: installApi.InitPlatform,
    method: 'post',
    data: parameter
  })
}

export function initLogo (parameter) {
  return request({
    url: installApi.InitLogo,
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data: parameter
  })
}

export function initCors (parameter) {
  return request({
    url: installApi.InitCors,
    method: 'post',
    data: parameter
  })
}

export function initRedis (parameter) {
  return request({
    url: installApi.InitRedis,
    method: 'post',
    data: parameter
  })
}
