import request from '@/utils/request'

// 获取公司列表
export function getCompanyList(params) {
  return request({
    url: '/company/list',
    method: 'get',
    params
  })
}

// 获取公司详情
export function getCompanyDetail(id) {
  return request({
    url: '/company/detail',
    method: 'get',
    params: { id }
  })
}

// 创建公司
export function createCompany(data) {
  return request({
    url: '/company/create',
    method: 'post',
    data
  })
}

// 更新公司
export function updateCompany(data) {
  return request({
    url: '/company/update',
    method: 'post',
    data
  })
}

// 删除公司
export function deleteCompany(id) {
  return request({
    url: '/company/delete',
    method: 'post',
    data: { id }
  })
}

// 获取待审核入驻申请列表
export function getPendingCompanyList(params) {
  return request({
    url: '/company/pendingList',
    method: 'get',
    params
  })
}

// 审核入驻申请
export function auditCompany(data) {
  return request({
    url: '/company/audit',
    method: 'post',
    data
  })
}
