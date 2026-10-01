import request from "@/utils/request"

// 获取账单列表（按月聚合）
export const getBillList = (data) => {
    return request({
        url: `/order/bill/list`,
        method: 'GET',
        loading: true,
        toLogin: true,
        data
    })
}

// 获取账单详情
export const getBillDetail = (data) => {
    return request({
        url: `/order/bill/detail`,
        method: 'GET',
        loading: true,
        toLogin: true,
        data
    })
}
