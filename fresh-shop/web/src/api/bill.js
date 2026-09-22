import service from '@/utils/request'

// @Tags Bill
// @Summary 获取账单列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query { page: 1, pageSize: 10 } true "分页参数"
// @Success 200 {string} string "{"success":true,"data":{"list":[],"total":0},"msg":"查询成功"}"
// @Router /business/bill/list [get]
export const getBillList = (params) => {
  return service({
    url: '/business/bill/list',
    method: 'get',
    params
  })
}

// @Tags Bill
// @Summary 生成账单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body { year: 2026, month: 9 } true "账单年月"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"账单生成成功"}"
// @Router /business/bill/generate [post]
export const generateBill = (data) => {
  return service({
    url: '/business/bill/generate',
    method: 'post',
    data
  })
}

// @Tags Bill
// @Summary 获取账单详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /business/bill/:id [get]
export const getBillDetail = (id) => {
  return service({
    url: `/business/bill/${id}`,
    method: 'get'
  })
}

// @Tags Bill
// @Summary 更新账单状态
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"状态更新成功"}"
// @Router /business/bill/:id/status [put]
export const updateBillStatus = (id, status) => {
  return service({
    url: `/business/bill/${id}/status`,
    method: 'put',
    data: { status }
  })
}
