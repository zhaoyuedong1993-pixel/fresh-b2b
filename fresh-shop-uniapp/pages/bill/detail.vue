<!--
 * 优诚配运 - 账单详情
 * 设计规范：自然清新
-->
<template>
    <pageWrapper>
        <!-- 顶部概览 -->
        <view class="bill-header">
            <view class="header-top">
                <view class="back-btn" @click="goBack">
                    <u-icon name="arrow-left" size="36rpx" color="#FFFFFF"></u-icon>
                </view>
                <text class="header-title">{{ period || '账单详情' }}</text>
            </view>
            <view class="bill-overview">
                <view class="overview-item">
                    <text class="overview-label">订单数</text>
                    <text class="overview-value">{{ billDetail.order_count || 0 }} <text class="unit">笔</text></text>
                </view>
                <view class="divider"></view>
                <view class="overview-item highlight">
                    <text class="overview-label">账单金额</text>
                    <text class="overview-value primary">
                        ¥{{ formatAmount(billDetail.total_amount) }}
                    </text>
                </view>
                <view class="divider"></view>
                <view class="overview-item">
                    <text class="overview-label">状态</text>
                    <view class="status-tag" :class="billDetail.status === 2 ? 'settled' : 'pending'">
                        {{ billDetail.status === 2 ? '已结算' : '待结算' }}
                    </view>
                </view>
            </view>
        </view>

        <!-- 订单明细 -->
        <view class="section">
            <view class="section-header">
                <view class="section-title">
                    <u-icon name="order" size="32rpx" color="#22A84F"></u-icon>
                    <text>订单明细</text>
                </view>
                <text class="section-count">{{ billDetail.orders?.length || 0 }} 笔</text>
            </view>

            <view class="order-list" v-if="billDetail.orders && billDetail.orders.length > 0">
                <view
                    class="order-card"
                    v-for="order in billDetail.orders"
                    :key="order.id"
                    @click="goOrderDetail(order)"
                >
                    <view class="order-info">
                        <view class="order-sn">{{ order.order_sn || order.orderSn }}</view>
                        <view class="order-date">{{ formatDate(order.created_at || order.createdAt) }}</view>
                    </view>
                    <view class="order-amount">
                        <text class="amount">¥{{ formatAmount(order.total) }}</text>
                        <u-icon name="arrow-right" size="28rpx" color="#CCCCCC"></u-icon>
                    </view>
                </view>
            </view>

            <view class="empty-orders" v-else>
                <u-icon name="order" size="80rpx" color="#CCCCCC"></u-icon>
                <text>暂无订单记录</text>
            </view>
        </view>

        <u-toast ref="toast" style="z-index: 9999"></u-toast>
    </pageWrapper>
</template>

<script>
    import { getBillDetail } from '@/api/business/bill.js'

    export default {
        data() {
            return {
                billId: null,
                period: '',
                billDetail: {}
            }
        },
        onLoad(options) {
            if (options.id) {
                this.billId = options.id
            }
            if (options.period) {
                this.period = options.period
            }
            this.loadBillDetail()
        },
        methods: {
            async loadBillDetail() {
                if (!this.billId) return

                uni.showLoading({ title: '加载中...' })
                try {
                    const res = await getBillDetail(this.billId)
                    if (res.code === 0) {
                        this.billDetail = res.data || {}
                        if (!this.period && this.billDetail.period) {
                            this.period = this.billDetail.period
                        }
                    }
                } catch (e) {
                    uni.showToast({ title: '加载失败', icon: 'none' })
                } finally {
                    uni.hideLoading()
                }
            },
            formatAmount(amount) {
                return (Number(amount) || 0).toFixed(2)
            },
            formatDate(dateStr) {
                if (!dateStr) return '-'
                const date = new Date(dateStr)
                return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
            },
            goBack() {
                uni.navigateBack()
            },
            goOrderDetail(order) {
                const id = order.id || order.ID
                if (id) {
                    uni.navigateTo({ url: `/pages/order/detail?id=${id}` })
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .bill-header {
        background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
        padding: 0 0 32rpx;
    }

    .header-top {
        display: flex;
        align-items: center;
        gap: 16rpx;
        padding: 24rpx 24rpx 16rpx;
    }

    .back-btn {
        width: 56rpx;
        height: 56rpx;
        background: rgba(255, 255, 255, 0.2);
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .header-title {
        font-size: 32rpx;
        font-weight: 600;
        color: #FFFFFF;
    }

    .bill-overview {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 32rpx 24rpx;
        margin: 0 24rpx;
        background: rgba(255, 255, 255, 0.15);
        border-radius: 24rpx;
        backdrop-filter: blur(10px);
    }

    .overview-item {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 8rpx;
        flex: 1;

        &.highlight {
            flex: 1.2;
        }
    }

    .overview-label {
        font-size: 24rpx;
        color: rgba(255, 255, 255, 0.7);
    }

    .overview-value {
        font-size: 32rpx;
        font-weight: 600;
        color: #FFFFFF;

        &.primary {
            font-size: 40rpx;
            color: #FFD700;
        }

        .unit {
            font-size: 22rpx;
            font-weight: 400;
        }
    }

    .status-tag {
        padding: 6rpx 20rpx;
        border-radius: 20rpx;
        font-size: 24rpx;

        &.settled {
            background: #E8F8EC;
            color: #22A84F;
        }

        &.pending {
            background: rgba(255, 255, 255, 0.2);
            color: #FFFFFF;
        }
    }

    .divider {
        width: 1rpx;
        height: 48rpx;
        background: rgba(255, 255, 255, 0.2);
    }

    .section {
        padding: 24rpx;
    }

    .section-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 20rpx;
    }

    .section-title {
        display: flex;
        align-items: center;
        gap: 8rpx;

        text {
            font-size: 30rpx;
            font-weight: 600;
            color: #1A1A1A;
        }
    }

    .section-count {
        font-size: 26rpx;
        color: #999999;
    }

    .order-list {
        background: #FFFFFF;
        border-radius: 24rpx;
        overflow: hidden;
        box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.06);
    }

    .order-card {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 28rpx;
        border-bottom: 1rpx solid #F0F0F0;

        &:last-child {
            border-bottom: none;
        }

        &:active {
            background: #FAFAFA;
        }
    }

    .order-info {
        flex: 1;
    }

    .order-sn {
        font-size: 28rpx;
        font-weight: 500;
        color: #1A1A1A;
        margin-bottom: 8rpx;
    }

    .order-date {
        font-size: 24rpx;
        color: #999999;
    }

    .order-amount {
        display: flex;
        align-items: center;
        gap: 8rpx;
    }

    .amount {
        font-size: 30rpx;
        font-weight: 600;
        color: #F97316;
    }

    .empty-orders {
        display: flex;
        flex-direction: column;
        align-items: center;
        padding: 80rpx 0;
        background: #FFFFFF;
        border-radius: 24rpx;

        text {
            font-size: 28rpx;
            color: #999999;
            margin-top: 20rpx;
        }
    }
</style>
