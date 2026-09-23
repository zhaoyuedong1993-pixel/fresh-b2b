package business

import (
	"fmt"
	"time"

	"fresh-shop/server/global"
)

// BillService 账单服务
type BillService struct{}

// GenerateBill 生成月度账单
func (s *BillService) GenerateBill(companyId uint, year, month int) error {
	period := fmt.Sprintf("%d-%02d", year, month)

	// 检查是否已存在
	var count int64
	global.DB.Table("shop_bill").
		Where("company_id = ? AND period = ?", companyId, period).
		Count(&count)
	if count > 0 {
		return nil // 已存在，跳过
	}

	// 统计该账期的已完成订单
	var result struct {
		Count int
		Total float64
	}

	startDate := fmt.Sprintf("%d-%02d-01", year, month)
	endDate := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)

	global.DB.Table("shop_order").
		Select("COUNT(*) as count, COALESCE(SUM(`total`), 0) as total").
		Where("company_id = ? AND order_status = 2 AND created_at >= ? AND created_at < ?",
			companyId, startDate, endDate).
		Scan(&result)

	// 生成账单编号
	billNo := fmt.Sprintf("BILL%s%02d", period, companyId)

	// 创建账单
	bill := map[string]interface{}{
		"company_id":   companyId,
		"bill_no":      billNo,
		"period":       period,
		"order_count":  result.Count,
		"total_amount": result.Total,
		"status":       0,
	}

	return global.DB.Table("shop_bill").Create(&bill).Error
}

// GetBillList 获取账单列表
func (s *BillService) GetBillList(companyId uint, page, pageSize int) ([]map[string]interface{}, int64) {
	var total int64
	global.DB.Table("shop_bill").
		Where("company_id = ?", companyId).
		Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	var list []map[string]interface{}
	global.DB.Table("shop_bill").
		Where("company_id = ?", companyId).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list)

	return list, total
}

// GetBillDetail 获取账单详情
func (s *BillService) GetBillDetail(billId uint) (map[string]interface{}, error) {
	var bill map[string]interface{}
	if err := global.DB.Table("shop_bill").Where("id = ?", billId).First(&bill).Error; err != nil {
		return nil, err
	}

	// 获取关联订单
	var orders []map[string]interface{}
	period := bill["period"].(string)
	global.DB.Table("shop_order").
		Select("id, order_sn, total, order_status, created_at").
		Where("company_id = ? AND DATE_FORMAT(created_at, '%Y-%m') = ?",
			bill["company_id"], period).
		Order("created_at DESC").
		Find(&orders)

	bill["orders"] = orders
	return bill, nil
}

// UpdateBillStatus 更新账单状态
func (s *BillService) UpdateBillStatus(billId uint, status int) error {
	return global.DB.Table("shop_bill").
		Where("id = ?", billId).
		Update("status", status).Error
}
