import request from '@/utils/request'

const installApi = {
  InitDb: '/install/init-db',
}

export function initDb (parameter) {
  return request({
    url: installApi.InitDb,
    method: 'post',
    data: parameter
  })
}

