<template>
    <pageWrapper>
        <!-- 账单概览 -->
        <view class="bill-overview">
            <view class="overview-item">
                <view class="label">账期</view>
                <view class="value">{{ billDetail.period || '-' }}</view>
            </view>
            <view class="overview-item">
                <view class="label">订单数</view>
                <view class="value">{{ billDetail.order_count || 0 }} 笔</view>
            </view>
            <view class="overview-item highlight">
                <view class="label">总金额</view>
                <view class="value">¥{{ formatAmount(billDetail.total_amount) }}</view>
            </view>
            <view class="overview-item">
                <view class="label">状态</view>
                <view :class="['value', 'status', billDetail.status === 1 ? 'settled' : 'unsettled']">
                    {{ billDetail.status === 1 ? '已结算' : '待结算' }}
                </view>
            </view>
        </view>

        <!-- 订单明细 -->
        <view class="section">
            <view class="section-title">订单明细</view>
            <view class="order-list" v-if="billDetail.orders && billDetail.orders.length > 0">
                <view class="order-item" v-for="order in billDetail.orders" :key="order.id">
                    <view class="order-info">
                        <view class="order-sn">{{ order.order_sn }}</view>
                        <view class="order-date">{{ formatDate(order.created_at) }}</view>
                    </view>
                    <view class="order-amount">¥{{ formatAmount(order.total) }}</view>
                </view>
            </view>
            <view class="empty-orders" v-else>
                <text>暂无订单记录</text>
            </view>
        </view>
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
                }
            } catch (e) {
                console.error(e)
                uni.showToast({ title: '加载失败', icon: 'none' })
            } finally {
                uni.hideLoading()
            }
        },

        formatAmount(amount) {
            return (amount || 0).toFixed(2)
        },

        formatDate(dateStr) {
            if (!dateStr) return '-'
            const date = new Date(dateStr)
            return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
        }
    }
}
</script>

<style lang="scss" scoped>
.bill-overview {
    background: linear-gradient(135deg, #4CAF50, #8BC34A);
    padding: 30rpx;
    color: #fff;
    display: flex;
    flex-wrap: wrap;
}

.overview-item {
    width: 50%;
    padding: 20rpx 0;
    text-align: center;

    .label {
        font-size: 24rpx;
        opacity: 0.9;
    }

    .value {
        font-size: 32rpx;
        font-weight: bold;
        margin-top: 10rpx;
    }

    &.highlight .value {
        font-size: 40rpx;
    }

    .status.settled {
        color: #67c23a;
        background: #fff;
        padding: 4rpx 16rpx;
        border-radius: 20rpx;
    }

    .status.unsettled {
        color: #fff;
        background: rgba(255, 255, 255, 0.2);
        padding: 4rpx 16rpx;
        border-radius: 20rpx;
    }
}

.section {
    padding: 30rpx;
}

.section-title {
    font-size: 32rpx;
    font-weight: bold;
    color: #333;
    margin-bottom: 20rpx;
    padding-left: 10rpx;
    border-left: 6rpx solid #4CAF50;
}

.order-list {
    background: #fff;
    border-radius: 16rpx;
    overflow: hidden;
}

.order-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 30rpx;
    border-bottom: 1rpx solid #f0f0f0;

    &:last-child {
        border-bottom: none;
    }
}

.order-info {
    .order-sn {
        font-size: 28rpx;
        color: #333;
    }

    .order-date {
        font-size: 24rpx;
        color: #999;
        margin-top: 8rpx;
    }
}

.order-amount {
    font-size: 32rpx;
    font-weight: bold;
    color: #4CAF50;
}

.empty-orders {
    text-align: center;
    padding: 60rpx 0;
    color: #999;
    background: #fff;
    border-radius: 16rpx;
}
</style>
