-- ============================================
-- 生鲜食材订购系统 - 多租户改造迁移脚本
-- 执行前请备份数据库！
-- ============================================

-- 1. 创建公司表
CREATE TABLE IF NOT EXISTS `sys_company` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '公司名称',
  `contact` varchar(50) COMMENT '联系人',
  `phone` varchar(20) COMMENT '联系电话',
  `address` varchar(200) COMMENT '地址',
  `markup_rate` decimal(5,2) DEFAULT 10.00 COMMENT '加价比例(%)',
  `cutoff_time` varchar(10) DEFAULT '22:00' COMMENT '截单时间',
  `status` tinyint(1) DEFAULT 1 COMMENT '状态: 0禁用 1启用',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='公司表';

-- 2. 用户表增加公司字段
ALTER TABLE `sys_users` ADD COLUMN `company_id` bigint(20) DEFAULT NULL COMMENT '公司ID' AFTER `header_img`;
ALTER TABLE `sys_users` ADD COLUMN `user_type` tinyint(1) DEFAULT 3 COMMENT '用户类型: 1超管 2公司管理员 3商户' AFTER `company_id`;
CREATE INDEX idx_company_id ON sys_users(company_id);

-- 3. 分类表增加公司字段
ALTER TABLE `shop_category` ADD COLUMN `company_id` bigint(20) DEFAULT NULL COMMENT '公司ID';
CREATE INDEX idx_category_company_id ON shop_category(company_id);

-- 4. 商品表增加公司字段和采购价
ALTER TABLE `shop_goods` ADD COLUMN `company_id` bigint(20) DEFAULT NULL COMMENT '公司ID';
ALTER TABLE `shop_goods` ADD COLUMN `cost_price` decimal(10,2) DEFAULT NULL COMMENT '采购价';
CREATE INDEX idx_goods_company_id ON shop_goods(company_id);

-- 5. 购物车表增加公司字段
ALTER TABLE `shop_cart` ADD COLUMN `company_id` bigint(20) DEFAULT NULL COMMENT '公司ID';
CREATE INDEX idx_cart_company_id ON shop_cart(company_id);

-- 6. 收藏表增加公司字段
ALTER TABLE `shop_favorites` ADD COLUMN `company_id` bigint(20) DEFAULT NULL COMMENT '公司ID';
CREATE INDEX idx_favorites_company_id ON shop_favorites(company_id);

-- 7. 订单表增加公司字段、订单状态、配送日期
ALTER TABLE `shop_order` ADD COLUMN `company_id` bigint(20) DEFAULT NULL COMMENT '公司ID';
ALTER TABLE `shop_order` ADD COLUMN `order_status` tinyint(1) DEFAULT 0 COMMENT '订单状态: 0待确认 1已确认 2已完成 3已取消';
ALTER TABLE `shop_order` ADD COLUMN `delivery_date` date DEFAULT NULL COMMENT '配送日期';
CREATE INDEX idx_order_company_id ON shop_order(company_id);
CREATE INDEX idx_order_status ON shop_order(order_status);

-- 8. 收货地址表增加公司字段
ALTER TABLE `shop_user_address` ADD COLUMN `company_id` bigint(20) DEFAULT NULL COMMENT '公司ID';
CREATE INDEX idx_address_company_id ON shop_user_address(company_id);

-- 9. 账单表
CREATE TABLE IF NOT EXISTS `shop_bill` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `company_id` bigint(20) NOT NULL COMMENT '公司ID',
  `bill_no` varchar(50) NOT NULL COMMENT '账单编号',
  `period` varchar(20) NOT NULL COMMENT '账期(2026-09)',
  `order_count` int DEFAULT 0 COMMENT '订单数量',
  `total_amount` decimal(12,2) DEFAULT 0 COMMENT '总金额',
  `status` tinyint(1) DEFAULT 0 COMMENT '状态: 0待结算 1已结算',
  `remark` varchar(500) COMMENT '备注',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_company_period` (`company_id`, `period`),
  INDEX `idx_bill_company_id` (`company_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='账单表';

-- 10. 插入平台公司（超管所属）
INSERT INTO `sys_company` (`id`, `name`, `markup_rate`, `cutoff_time`, `status`, `created_at`) VALUES (1, '平台管理', 10.00, '22:00', 1, NOW());

-- 11. 给现有数据设置 company_id（假设超管的 authority_id = 888）
UPDATE `sys_users` SET `company_id` = 1, `user_type` = 1 WHERE `authority_id` = 888;
UPDATE `sys_users` SET `company_id` = 1, `user_type` = 3 WHERE `company_id` IS NULL AND `authority_id` != 888;

UPDATE `shop_category` SET `company_id` = 1 WHERE `company_id` IS NULL;
UPDATE `shop_goods` SET `company_id` = 1 WHERE `company_id` IS NULL;
UPDATE `shop_cart` SET `company_id` = 1 WHERE `company_id` IS NULL;
UPDATE `shop_favorites` SET `company_id` = 1 WHERE `company_id` IS NULL;
UPDATE `shop_order` SET `company_id` = 1 WHERE `company_id` IS NULL;
UPDATE `shop_user_address` SET `company_id` = 1 WHERE `company_id` IS NULL;

-- ============================================
-- 迁移完成！
-- ============================================
