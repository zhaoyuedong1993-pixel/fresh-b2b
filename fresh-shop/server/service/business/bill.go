package business

import (
	"time"

	"fresh-shop/server/global"
)

// BillService 账单服务
type BillService struct{}

// GetBillList 获取月度账单列表（按月聚合已完成订单）
func (s *BillService) GetBillList(companyId uint, page, pageSize int) ([]map[string]interface{}, int64) {
	var total int64
	global.DB.Table("shop_order").
		Where("company_id = ? AND status IN (2,3) AND status_cancel = 0 AND status_refund = 0", companyId).
		Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	var list []map[string]interface{}
	offset := (page - 1) * pageSize
	global.DB.Table("shop_order").
		Select(`DATE_FORMAT(created_at, '%Y-%m') as period,
			COUNT(*) as order_count,
			COALESCE(SUM(total), 0) as total_amount,
			MIN(settlement_type) as settlement_type`).
		Where("company_id = ? AND status IN (2,3) AND status_cancel = 0 AND status_refund = 0", companyId).
		Group("DATE_FORMAT(created_at, '%Y-%m')").
		Order("period DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&list)

	return list, total
}

// GetBillDetailByPeriod 根据账期获取账单详情
func (s *BillService) GetBillDetailByPeriod(companyId uint, period string) (map[string]interface{}, error) {
	startDate := period + "-01"
	t, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}
	endDate := t.AddDate(0, 1, 0).Format("2006-01-02")

	// 统计
	var result struct {
		Count          int
		Total          float64
		SettlementType int
	}
	global.DB.Table("shop_order").
		Select("COUNT(*) as count, COALESCE(SUM(total), 0) as total, MIN(settlement_type) as settlement_type").
		Where("company_id = ? AND status IN (2,3) AND status_cancel = 0 AND status_refund = 0 AND created_at >= ? AND created_at < ?",
			companyId, startDate, endDate).
		Scan(&result)

	// 订单明细
	var orders []map[string]interface{}
	global.DB.Table("shop_order").
		Select("id, order_sn, total, status, created_at, shipment_name, shipment_mobile").
		Where("company_id = ? AND status IN (2,3) AND status_cancel = 0 AND status_refund = 0 AND created_at >= ? AND created_at < ?",
			companyId, startDate, endDate).
		Order("created_at DESC").
		Find(&orders)

	return map[string]interface{}{
		"period":        period,
		"order_count":   result.Count,
		"total_amount":  result.Total,
		"settlement_type": result.SettlementType,
		"status":        0, // 0=未结算 1=已结算
		"orders":        orders,
	}, nil
}

// UpdateBillStatusByPeriod 按账期更新结算状态
func (s *BillService) UpdateBillStatusByPeriod(companyId uint, period string, status int) error {
	startDate := period + "-01"
	t, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return err
	}
	endDate := t.AddDate(0, 1, 0).Format("2006-01-02")
	return global.DB.Table("shop_order").
		Where("company_id = ? AND status IN (2,3) AND created_at >= ? AND created_at < ?",
			companyId, startDate, endDate).
		Update("settlement_type", status+1).Error // 1=月结未结 2=月结已结
}
