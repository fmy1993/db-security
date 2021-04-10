import request from '@/util/request'

export const getAllIp = () => {
  return request({
    url: '/all_ip',
    method: 'GET'
  })
}

export const freeIp = (ipId) => {
  return request({
    url: `/free_ip/${ipId}`,
    method: 'DELETE'
  })
}

export const freezeIp = (data) => {
  return request({
    url: '/freeze_ip',
    method: 'POST',
    data
  })
}
