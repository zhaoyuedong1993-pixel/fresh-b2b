import service from '@/utils/request'

// @Tags Bill
// @Summary 获取账单列表（按月聚合）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{"list":[],"total":0},"msg":"查询成功"}"
// @Router /order/bill/list [get]
export const getBillList = (params) => {
  return service({
    url: '/order/bill/list',
    method: 'get',
    params
  })
}

// @Tags Bill
// @Summary 获取账单详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /order/bill/detail [get]
export const getBillDetail = (params) => {
  return service({
    url: '/order/bill/detail',
    method: 'get',
    params
  })
}

// @Tags Bill
// @Summary 更新账单状态
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"状态更新成功"}"
// @Router /order/bill/status [put]
export const updateBillStatus = (data) => {
  return service({
    url: '/order/bill/status',
    method: 'put',
    data
  })
}
