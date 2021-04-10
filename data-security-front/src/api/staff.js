import request from '@/util/request'

export const download = () => {
  return request({
    url: '/download/staff',
    method: 'GET',
    responseType: 'blob'
  })
}

export const getStaff = data => {
  return request({
    url: '/staff',
    method: 'POST',
    data
  })
}

export const getOriStaff = data => {
  return request({
    url: '/ori_staff',
    method: 'POST',
    data
  })
}

export const addStaff = (data) => {
  return request({
    url: '/add_staff',
    method: 'POST',
    data
  })
}

export const deleteStaff = staffId => {
  return request({
    url: `/staff/${staffId}`,
    method: 'DELETE'
  })
}

export const updateStaff = (staff, staffId) => {
  return request({
    url: `/staff/${staffId}`,
    method: 'PUT',
    data: staff
  })
}

export const differentialPrivacy = () => {
  return request({
    url: '/dp',
    method: 'POST'
  })
}

export const getAnalysis = data => {
  return request({
    url: '/analysis',
    method: 'POST',
    data
  })
}

export const getMostIdcard = () => {
  return request({
    url: '/idcard/most',
    method: 'GET'
  })
}
