# Fresh-B2B 生鲜食材订购系统

[![Go](https://img.shields.io/badge/Go-1.18+-blue.svg)](https://golang.org/)
[![Vue](https://img.shields.io/badge/Vue.js-3.x-green.svg)](https://vuejs.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)

**Fresh-B2B** 是一款面向餐饮企业的 B2B 生鲜食材订购平台，基于 Go (Gin) 和 Vue 3 构建。

## ✨ 核心功能

- **多租户 SaaS** - 支持多个公司独立运营，数据完全隔离
- **AI 智能调价** - 采购价 × 加价比例 = 商户价，一键批量生效
- **月度账单** - 自动生成月度账单，线下结算
- **订单管理** - 待确认 → 已确认 → 已完成

## 🏗️ 系统架构

```
┌─────────────────────────────────────────────┐
│              微信小程序（商户端）              │
│         商品浏览 | 购物车 | 下单 | 账单        │
└────────────────────┬────────────────────────┘
                     │
┌────────────────────▼────────────────────────┐
│           SaaS 管理后台（公司运营）            │
│      商品管理 | AI调价 | 订单确认 | 账单       │
└────────────────────┬────────────────────────┘
                     │
┌────────────────────▼────────────────────────┐
│              Go (Gin) API 服务                │
│     多租户中间件 | AI调价 | 账单服务          │
└─────────────────────────────────────────────┘
```

## 📂 目录结构

```
fresh-b2b/
├── fresh-shop/server/        # Go 后端 API
├── fresh-shop/web/          # Vue 3 管理后台
├── fresh-shop-uniapp/       # uni-app 微信小程序
└── docs/                    # 设计文档和计划
```

## 🚀 快速启动

### 环境准备

- Go >= 1.18
- Node.js >= 16.0
- MySQL >= 5.7
- Redis

### 后端

```bash
cd fresh-b2b/fresh-shop/server
cp config.yaml.example config.yaml
# 修改 config.yaml 中的数据库配置

# 执行数据库迁移
mysql -u root -p fresh-shop < sql/migration/001_multi_tenant.sql

go mod tidy
go run main.go
```

### 管理后台

```bash
cd fresh-b2b/fresh-shop/web
pnpm install
pnpm serve
# 访问 http://localhost:8080
```

### 微信小程序

使用 HBuilder X 导入 `fresh-shop-uniapp` 项目，修改 `config/config.js` 中的 API 地址。

## 📄 开源许可

本项目基于 [MIT](/LICENSE) 许可。
