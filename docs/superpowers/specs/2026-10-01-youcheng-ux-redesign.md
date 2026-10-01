# 优诚配运小程序 UI 设计规范

> **Goal:** 重构「优诚配运」小程序全链路 UI，统一设计语言，传达「自然清新」品牌调性。
> **品牌名:** 优诚配运（原「启运冻品」全部替换）
> **风格:** 自然清新 — 新鲜、可靠、专业

---

## 设计系统

### 色彩规范

| 角色 | 色值 | 用途 |
|------|------|------|
| 主色 | `#22A84F` | 按钮、Tab bar、强调元素 |
| 主色浅 | `#E8F8EC` | 背景、标签 |
| 主色深 | `#1A7A38` | 按钮按下态 |
| 点缀色 | `#F97316` | 价格、重要提示 |
| 背景 | `#F5F7F4` | 页面底色 |
| 卡片 | `#FFFFFF` | 卡片、弹层 |
| 正文 | `#1A1A1A` | 主要文字 |
| 次要文字 | `#666666` | 辅助说明 |
| 占位文字 | `#999999` | 输入框 placeholder |
| 分割线 | `#EEEEEE` | 列表分割、边框 |
| 危险色 | `#EF4444` | 删除、错误 |
| 成功色 | `#22C55E` | 成功提示 |

### 字体规范

| 层级 | 大小 | 字重 | 用途 |
|------|------|------|------|
| 页面标题 | 36rpx | 600 | 页面大标题 |
| 卡片标题 | 32rpx | 600 | 模块标题 |
| 正文 | 28rpx | 400 | 主要内容 |
| 辅助文字 | 24rpx | 400 | 标签、说明 |
| 小字 | 22rpx | 400 | 次要信息 |

### 间距规范（基于 8rpx 网格）

- `xs`: 8rpx
- `sm`: 16rpx
- `md`: 24rpx
- `lg`: 32rpx
- `xl`: 48rpx

### 圆角规范

- 小组件（标签、徽章）: `8rpx`
- 卡片、按钮: `16rpx`
- 大弹层: `24rpx`
- 头像、图标容器: `50%`

### 阴影规范

```scss
// 卡片阴影
box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.06);

// 弹层阴影
box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.12);
```

### 动效规范

- 页面切换: `300ms ease`
- 按钮反馈: `150ms`
- 列表加载: `fade` + `slide-up`, `400ms`
- 卡片入场: `stagger 50ms`

---

## 页面重构清单

### 批次 1: 首页 + 分类 + 商品列表
- `pages/index/index.vue` — 首页
- `pages/category/category.vue` — 分类页
- `components/goodsList/goodsList.vue` — 商品卡片组件

### 批次 2: 交易流程
- `pages/goods/detail.vue` — 商品详情
- `pages/cart/cart.vue` — 购物车
- `pages/order/submit.vue` — 提交订单
- `components/shopCart/shopCart.vue` — 购物车浮层

### 批次 3: 订单管理
- `pages/order/list.vue` — 订单列表
- `pages/order/detail.vue` — 订单详情
- `components/orderList/orderList.vue` — 订单卡片组件
- `pages/address/address.vue` — 地址管理
- `pages/address/addressForm.vue` — 地址表单

### 批次 4: 账户体系
- `pages/my/my.vue` — 个人中心
- `pages/my/memberInfo.vue` — 会员信息
- `pages/login/login.vue` — 登录页
- `pages/login/register.vue` — 注册页
- `pages/bill/list.vue` — 账单列表
- `pages/bill/detail.vue` — 账单详情

### 共享组件
- `components/tabbar/tabbar.vue` — 底部导航
- `components/pageWrapper/pageWrapper.vue` — 页面包装
- `components/loginPop/loginPop.vue` — 登录弹层
- `components/loginSuspend/loginSuspend.vue` — 登录悬浮
- `components/addressPop/addressPop.vue` — 地址选择
- `components/searchPage/searchPage.vue` — 搜索页
- `components/filterDropdown/filterDropdown.vue` — 筛选下拉

---

## 共享改动

### uni.scss 设计 token
所有颜色、间距、圆角统一为 SCSS 变量。

### 通用工具类（可选）
废弃散乱的 `king-*` 类，统一为 Tailwind-like 或纯 SCSS。

### manifest.json
- 应用名改为「优诚配运」
- 移除所有「启运冻品」相关字样

### pages.json
- 检查所有页面标题是否包含旧品牌名

### 各页面顶部导航
统一样式，搜索框 + 品牌感

### Tab Bar
- 图标用线性风格
- 主色高亮

---

## 实施顺序

1. 先统一 `uni.scss` 设计 token（所有页面共享）
2. 批次 1 改完验收
3. 批次 2 改完验收
4. 批次 3 改完验收
5. 批次 4 改完验收
6. 全局检查：品牌名替换、manifest 更新
