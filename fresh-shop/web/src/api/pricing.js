import service from '@/utils/request'

// @Tags Pricing
// @Summary 批量更新采购价
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body { updates: { [goodsId]: costPrice } } true "采购价更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"采购价更新成功"}"
// @Router /business/pricing/batch-cost-price [post]
export const batchUpdateCostPrice = (data) => {
  return service({
    url: '/business/pricing/batch-cost-price',
    method: 'post',
    data
  })
}

// @Tags Pricing
// @Summary 应用价格（根据采购价和加价比例计算商户价）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"价格已生效"}"
// @Router /business/pricing/apply [post]
export const applyPricing = () => {
  return service({
    url: '/business/pricing/apply',
    method: 'post'
  })
}

// @Tags Pricing
// @Summary 获取当前公司的加价比例
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{"markupRate":10},"msg":"查询成功"}"
// @Router /business/pricing/markup-rate [get]
export const getMarkupRate = () => {
  return service({
    url: '/business/pricing/markup-rate',
    method: 'get'
  })
}

// @Tags Pricing
// @Summary AI解析文案并批量更新采购价
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body { text: string } true "待解析的文本"
// @Success 200 {string} string "{"success":true,"data":{"total":5,"updated":3,"items":[]},"msg":"解析成功"}"
// @Router /business/pricing/batch-parse [post]
export const batchParseText = (data) => {
  return service({
    url: '/business/pricing/batch-parse',
    method: 'post',
    data
  })
}
