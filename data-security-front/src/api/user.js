import request from '@/util/request'

export const getSalt = () => {
  return request({
    url: '/user/salt',
    method: 'GET'
  })
}

export const login = data => {
  return request({
    url: '/user/login',
    method: 'POST',
    data
  })
}

export const register = data => {
  return request({
    url: '/user/register',
    method: 'POST',
    data
  })
}

export const getCaptcha = () => {
  return request({
    url: '/user/captcha',
    method: 'GET'
  })
}

export const logout = () => {
  return request({
    url: '/user/logout',
    method: 'POST'
  })
}

export const revise = data => {
  return request({
    url: '/user/revise',
    method: 'POST',
    data
  })
}
