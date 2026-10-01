/*
 * @Author: dalefeng
 * @Date: 2023-04-14 20:40:42
 * @LastEditors: dalefeng
 * @LastEditTime: 2023-04-18 17:57:36
 */
import request from "@/utils/request"

export const getWeChatOpenIdByCode = (data) => {
  return request({
    url: '/wechat/code2Session',
    data,
  })
}

export const wxLogin = (data) => {
  return request({
    url: '/base/loginWx',
    method: 'POST',
    loading: true,
    data,
  })
}

// 小程序手机号+密码登录
export const loginByPhone = (data) => {
  return request({
    url: '/base/loginByPhone',
    method: 'POST',
    loading: true,
    data,
  })
}

// 小程序注册
export const register = (data) => {
  return request({
    url: '/base/register',
    method: 'POST',
    loading: true,
    data,
  })
}

// 申请入驻平台
export const registerCompany = (data) => {
  return request({
    url: '/base/registerCompany',
    method: 'POST',
    loading: true,
    data,
  })
}

// 通过邀请码获取公司信息
export const getCompanyByInviteCode = (code) => {
  return request({
    url: '/base/getCompanyByInviteCode',
    method: 'GET',
    params: { code },
  })
}

// 申请加入公司
export const joinCompany = (data) => {
  return request({
    url: '/base/joinCompany',
    method: 'POST',
    loading: true,
    data,
  })
}

// 通过公司ID获取公司信息
export const getCompanyById = (id) => {
  return request({
    url: '/company/detail',
    method: 'GET',
    params: { id },
  })
}
