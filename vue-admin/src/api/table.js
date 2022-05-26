import request from '@/utils/request'

export function getList(params) {
  return request({
    url: '/api/i2p/metadata',
    method: 'get',
    params
  })
}

export function getExperList(params) {
  return request({
    url: '/api/i2p/exper/list',
    method: 'get',
    params
  })
}
export function getFileContent(params) {
  return request({
    url: '/api/i2p/exper/filecontent',
    method: 'get',
    params
  })
}

export function getRouterInfo(params) {
  return request({
    url: '/api/i2p/metadata_content',
    method: 'get',
    params
  })
}
export function getRouterList(params) {
  return request({
    url: '/api/i2p/exper/getFileList',
    method: 'get',
    params
  })
}
export function getDockerState(params) {
  return request({
    url: '/api/i2p/exper/getDockerState',
    method: 'get',
    params
  })
}

export function getFileChanged() {
  return request({
    url: '/api/i2p/exper/getFileChange',
    method: 'get',
  })
}

export function createNewLab(params) {
  return request({
    // url:'api/i2p/Exper/create',
    url: '/api/i2p/saveluafile',
    method: 'post',
    data: {
      file_name: params['filename'],
      file_content: params['filecontent'],
      exper_type: params['experPoint'],
      exper_time: params['experTime'],
      date1: params['date1'],
      date2: params['date2'],
      isPublic: params['isPublic'],
      start_time: params['startTime'],
      add_time: params['addTime'],
      desc: params['desc'],
      user: params['user'],
    }
  })
}
