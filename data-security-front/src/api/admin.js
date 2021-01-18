import request from '@/util/request'

export const adminPatient = page => {
  return request({
    url: '/patient/admin/patient?page=' + page.toString(),
    method: 'POST'
  })
}

export const adminSearch = (page, pattern) => {
  return request({
    url: '/patient/admin/search?page=' + page.toString() + '&pattern=' + pattern,
    method: 'GET'
  })
}

export const adminAddPatient = data => {
  return request({
    url: '/patient/admin/add',
    method: 'POST',
    data
  })
}

export const adminDeletePatient = data => {
  return request({
    url: '/patient/admin/delete',
    method: 'POST',
    data
  })
}

export const adminUpdatePatient = data => {
  return request({
    url: '/patient/admin/update',
    method: 'POST',
    data: data
  })
}

export const adminGetAllUsers = () => {
  return request({
    url: '/user/admin/users',
    method: 'POST'
  })
}

export const adminFreezeUser = user => {
  return request({
    url: '/user/admin/freeze',
    method: 'POST',
    data: user
  })
}

export const adminFreeUser = user => {
  return request({
    url: '/user/admin/free',
    method: 'POST',
    data: user
  })
}

export const adminGetAllIps = () => {
  return request({
    url: '/ip/admin/ips',
    method: 'POST'
  })
}

export const adminFreeIp = ip => {
  return request({
    url: '/ip/admin/free',
    method: 'POST',
    data: ip
  })
}

export const adminAddIp = ip => {
  return request({
    url: '/ip/admin/add',
    method: 'POST',
    data: {
      ip: ip
    }
  })
}

export const adminGetAllIis = () => {
  return request({
    url: '/ii/admin/iis',
    method: 'POST'
  })
}

export const adminDp = () => {
  return request({
    url: '/patient/admin/dp',
    method: 'POST',
    async: false
  })
}

export const adminTrack = data => {
  return request({
    url: '/patient/admin/track',
    method: 'POST',
    data
  })
}
