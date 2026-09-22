import request from "@/utils/request"

// 获取账单列表
export const getBillList = (data) => {
    return request({
        url: `/business/bill/list`,
        method: 'GET',
        loading: true,
        toLogin: true,
        data
    })
}

// 获取账单详情
export const getBillDetail = (id) => {
    return request({
        url: `/business/bill/${id}`,
        method: 'GET',
        loading: true,
        toLogin: true
    })
}
