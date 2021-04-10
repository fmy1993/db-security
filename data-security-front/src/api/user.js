import request from '@/util/request'

export const getSalt = () => {
  return request({
    url: '/salt',
    method: 'GET'
  })
}

export const login = data => {
  return request({
    url: '/login',
    method: 'POST',
    data
  })
}

export const register = data => {
  return request({
    url: '/register',
    method: 'POST',
    data
  })
}

export const getCaptcha = () => {
  return request({
    url: '/captcha',
    method: 'GET'
  })
}

export const logout = () => {
  return request({
    url: '/logout',
    method: 'POST'
  })
}

export const revise = data => {
  return request({
    url: '/revise',
    method: 'PUT',
    data
  })
}

export const getUserInfo = () => {
  return request({
    url: '/user/info',
    method: 'GET'
  })
}

export const getAllUsers = () => {
  return request({
    url: '/all_users',
    method: 'GET'
  })
}

export const freezeUser = (userId) => {
  return request({
    url: `/freeze/${userId}`,
    method: 'POST'
  })
}

export const freeUser = (userId) => {
  return request({
    url: `/free/${userId}`,
    method: 'POST'
  })
}
