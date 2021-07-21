import request from '@/utils/request'

const showApi = {
  Info: '/s',
}

export function showInfo(params) {
  return request({
    url: `${showApi.Info}/${params.code}/info`,
    method: 'get',
    params
  })
}
