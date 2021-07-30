import request from '@/utils/request'

const installApi = {
  InitDb: '/install/init-db',
  InitPlatform: '/install/init-platform',
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

