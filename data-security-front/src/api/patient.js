import request from '@/util/request'

export const download = () => {
  return request({
    url: '/patient/download',
    method: 'GET',
    responseType: 'blob'
  })
}

export const patient = page => {
  return request({
    url: '/patient/data?page=' + page.toString(),
    method: 'POST'
  })
}

export const search = (page, pattern) => {
  return request({
    url: '/patient/search?page=' + page.toString() + '&pattern=' + pattern,
    method: 'GET'
  })
}

export const analysis = data => {
  return request({
    url: '/patient/analysis',
    method: 'post',
    data
  })
}
