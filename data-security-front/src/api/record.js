import request from '@/util/request'

export const getAllRecords = () => {
  return request({
    url: '/download_record',
    method: 'GET'
  })
}
