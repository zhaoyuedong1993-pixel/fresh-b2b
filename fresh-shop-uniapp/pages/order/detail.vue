<!--
 * 优诚配运 - 订单详情
 * 设计规范：自然清新
-->
<template>
    <pageWrapper>
        <!-- 顶部状态栏 -->
        <view class="status-header" :class="'status-' + order.status">
            <view class="status-icon">
                <u-icon :name="statusIcon" size="48" color="#FFFFFF"></u-icon>
            </view>
            <view class="status-text">{{ statusText }}</view>
        </view>

        <!-- 收货地址 -->
        <view class="card address-card" v-if="order.address">
            <view class="address-icon">
                <u-icon name="map" size="36" color="#22A84F"></u-icon>
            </view>
            <view class="address-info">
                <view class="address-user">
                    <text class="name">{{ order.address.userName }}</text>
                    <text class="phone">{{ order.address.userPhone }}</text>
                </view>
                <view class="address-detail">
                    {{ order.address.detailAddress || order.address.address }}
                </view>
            </view>
        </view>

        <!-- 商品列表 -->
        <view class="card goods-card">
            <view class="card-title">
                <text>商品清单</text>
            </view>
            <view class="goods-list">
                <view class="goods-item" v-for="item in order.goodsList" :key="item.id">
                    <image class="goods-img" :src="item.image" mode="aspectFill"></image>
                    <view class="goods-info">
                        <text class="goods-name">{{ item.name }}</text>
                        <view class="goods-spec" v-if="item.spec">
                            <text>{{ item.spec }}</text>
                        </view>
                    </view>
                    <view class="goods-right">
                        <text class="goods-price">¥{{ formatPrice(item.price) }}</text>
                        <text class="goods-num">x{{ item.num }}</text>
                    </view>
                </view>
            </view>
        </view>

        <!-- 订单信息 -->
        <view class="card info-card">
            <view class="info-row">
                <text class="info-label">订单编号</text>
                <view class="info-value">
                    <text>{{ order.orderSn || order.order_sn || '-' }}</text>
                    <text class="copy-btn" @click="copyOrderSn">复制</text>
                </view>
            </view>
            <view class="info-row">
                <text class="info-label">下单时间</text>
                <text class="info-value">{{ formatDate(order.createdAt || order.created_at) }}</text>
            </view>
            <view class="info-row" v-if="order.deliveryTime">
                <text class="info-label">发货时间</text>
                <text class="info-value">{{ formatDate(order.deliveryTime) }}</text>
            </view>
            <view class="info-row" v-if="order.receiveTime">
                <text class="info-label">收货时间</text>
                <text class="info-value">{{ formatDate(order.receiveTime) }}</text>
            </view>
            <view class="info-row">
                <text class="info-label">配送方式</text>
                <text class="info-value">{{ order.shippingType === 0 ? '配送到家' : '自提' }}</text>
            </view>
            <view class="info-row">
                <text class="info-label">支付方式</text>
                <text class="info-value">{{ paymentText }}</text>
            </view>
            <view class="info-row">
                <text class="info-label">订单备注</text>
                <text class="info-value remark">{{ order.remark || '无' }}</text>
            </view>
        </view>

        <!-- 价格汇总 -->
        <view class="card price-card">
            <view class="price-row">
                <text class="price-label">商品金额</text>
                <text class="price-value">¥{{ formatPrice(order.goodsAmount || order.goods_amount) }}</text>
            </view>
            <view class="price-row">
                <text class="price-label">配送费</text>
                <text class="price-value">¥{{ formatPrice(order.freightAmount || order.freight_amount || 0) }}</text>
            </view>
            <view class="price-row total">
                <text class="price-label">实付金额</text>
                <text class="price-value primary">¥{{ formatPrice(order.totalAmount || order.total_amount) }}</text>
            </view>
        </view>

        <!-- 底部操作栏 -->
        <view class="bottom-bar" v-if="showBottomBar">
            <view class="bottom-left" v-if="order.status === 0">
                <view class="total-info">
                    <text class="total-label">实付</text>
                    <text class="total-price">¥{{ formatPrice(order.totalAmount || order.total_amount) }}</text>
                </view>
            </view>
            <view class="bottom-right">
                <view class="action-btn cancel" v-if="order.status === 0" @click="cancelOrder">
                    取消订单
                </view>
                <view class="action-btn pay" v-if="order.status === 0" @click="payOrder">
                    立即付款
                </view>
                <view class="action-btn confirm" v-if="order.status === 2" @click="confirmReceive">
                    确认收货
                </view>
                <view class="action-btn delete" v-if="order.status === 3 || order.status === 4" @click="deleteOrder">
                    删除订单
                </view>
            </view>
        </view>

        <!-- 底部占位 -->
        <view class="bottom-placeholder" v-if="showBottomBar"></view>

        <u-toast ref="toast" style="z-index: 9999"></u-toast>
    </pageWrapper>
</template>

<script>
    import { getOrderDetail } from '@/api/order.js'
    import { cancelOrder as cancelOrderApi, confirmOrder, deleteOrder as deleteOrderApi } from '@/api/order.js'
    import { getToken } from '@/store/storage.js'

    export default {
        data() {
            return {
                orderId: 0,
                order: {}
            }
        },
        computed: {
            statusText() {
                const map = { 0: '待付款', 1: '备货中', 2: '配送中', 3: '已完成', 4: '已取消' }
                return map[this.order.status] || '未知状态'
            },
            statusIcon() {
                const map = { 0: 'rmb-circle', 1: 'clock', 2: 'car', 3: 'checkmark-circle', 4: 'close-circle' }
                return map[this.order.status] || 'question'
            },
            paymentText() {
                const map = { 0: '微信支付', 1: '支付宝', 2: '余额支付', 5: '月结' }
                return map[this.order.payment] || '微信支付'
            },
            showBottomBar() {
                return [0, 2, 3, 4].includes(this.order.status)
            }
        },
        onLoad(options) {
            if (options.id) {
                this.orderId = parseInt(options.id)
            }
            this.loadOrderDetail()
        },
        methods: {
            async loadOrderDetail() {
                const token = getToken()
                if (!token) {
                    uni.showToast({ title: '请先登录', icon: 'none' })
                    return
                }

                uni.showLoading({ title: '加载中...' })
                try {
                    const res = await getOrderDetail({ ID: this.orderId })
                    if (res.code === 0) {
                        this.order = res.data || {}
                    }
                } catch (e) {
                    uni.showToast({ title: '加载失败', icon: 'none' })
                } finally {
                    uni.hideLoading()
                }
            },
            formatPrice(price) {
                return (Number(price) || 0).toFixed(2)
            },
            formatDate(dateStr) {
                if (!dateStr) return '-'
                const date = new Date(dateStr)
                return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
            },
            copyOrderSn() {
                const sn = this.order.orderSn || this.order.order_sn || ''
                uni.setClipboardData({
                    data: sn,
                    success: () => {
                        uni.showToast({ title: '已复制', icon: 'success' })
                    }
                })
            },
            cancelOrder() {
                uni.showModal({
                    title: '取消订单',
                    content: '确定要取消该订单吗？',
                    success: async (res) => {
                        if (res.confirm) {
                            try {
                                const result = await cancelOrderApi({ ID: this.orderId })
                                if (result.code === 0) {
                                    this.$message(this.$refs.toast).success('订单已取消')
                                    this.loadOrderDetail()
                                }
                            } catch (e) {
                                this.$message(this.$refs.toast).error('操作失败')
                            }
                        }
                    }
                })
            },
            payOrder() {
                uni.navigateTo({
                    url: `/pages/order/submit?orderId=${this.orderId}`
                })
            },
            confirmReceive() {
                uni.showModal({
                    title: '确认收货',
                    content: '确认已收到货物吗？',
                    success: async (res) => {
                        if (res.confirm) {
                            try {
                                const result = await confirmOrder({ ID: this.orderId })
                                if (result.code === 0) {
                                    this.$message(this.$refs.toast).success('已确认收货')
                                    this.loadOrderDetail()
                                } else {
                                    this.$message(this.$refs.toast).error(result.msg || '操作失败')
                                }
                            } catch (e) {
                                this.$message(this.$refs.toast).error('操作失败')
                            }
                        }
                    }
                })
            },
            deleteOrder() {
                uni.showModal({
                    title: '删除订单',
                    content: '确定要删除该订单吗？',
                    success: async (res) => {
                        if (res.confirm) {
                            try {
                                const result = await deleteOrderApi({ ID: this.orderId })
                                if (result.code === 0) {
                                    this.$message(this.$refs.toast).success('订单已删除')
                                    setTimeout(() => {
                                        uni.navigateBack()
                                    }, 1500)
                                }
                            } catch (e) {
                                this.$message(this.$refs.toast).error('操作失败')
                            }
                        }
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .status-header {
        background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
        padding: 48rpx 32rpx;
        display: flex;
        align-items: center;
        gap: 20rpx;

        &.status-0 {
            background: linear-gradient(135deg, #F97316 0%, #EA580C 100%);
        }

        &.status-1 {
            background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
        }

        &.status-2 {
            background: linear-gradient(135deg, #3B82F6 0%, #2563EB 100%);
        }

        &.status-3 {
            background: linear-gradient(135deg, #666666 0%, #555555 100%);
        }

        &.status-4 {
            background: linear-gradient(135deg, #999999 0%, #888888 100%);
        }
    }

    .status-text {
        font-size: 36rpx;
        font-weight: 600;
        color: #FFFFFF;
    }

    .card {
        margin: 24rpx;
        background: #FFFFFF;
        border-radius: 24rpx;
        box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.06);
        overflow: hidden;
    }

    .address-card {
        display: flex;
        align-items: flex-start;
        padding: 28rpx;
        gap: 20rpx;
    }

    .address-icon {
        flex-shrink: 0;
        width: 64rpx;
        height: 64rpx;
        background: #E8F8EC;
        border-radius: 16rpx;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .address-info {
        flex: 1;
    }

    .address-user {
        display: flex;
        align-items: center;
        gap: 16rpx;
        margin-bottom: 12rpx;

        .name {
            font-size: 30rpx;
            font-weight: 600;
            color: #1A1A1A;
        }

        .phone {
            font-size: 28rpx;
            color: #666666;
        }
    }

    .address-detail {
        font-size: 26rpx;
        color: #666666;
        line-height: 1.5;
    }

    .goods-card {
        padding: 0;
    }

    .card-title {
        padding: 28rpx 28rpx 20rpx;
        border-bottom: 1rpx solid #EEEEEE;

        text {
            font-size: 30rpx;
            font-weight: 600;
            color: #1A1A1A;
        }
    }

    .goods-list {
        padding: 0 28rpx;
    }

    .goods-item {
        display: flex;
        align-items: center;
        padding: 24rpx 0;
        border-bottom: 1rpx solid #F5F5F5;

        &:last-child {
            border-bottom: none;
        }
    }

    .goods-img {
        width: 140rpx;
        height: 140rpx;
        border-radius: 16rpx;
        background: #F5F7F4;
        flex-shrink: 0;
    }

    .goods-info {
        flex: 1;
        padding: 0 20rpx;
    }

    .goods-name {
        font-size: 28rpx;
        color: #1A1A1A;
        display: block;
        margin-bottom: 8rpx;
    }

    .goods-spec {
        font-size: 24rpx;
        color: #999999;
    }

    .goods-right {
        text-align: right;
    }

    .goods-price {
        display: block;
        font-size: 28rpx;
        font-weight: 600;
        color: #F97316;
        margin-bottom: 8rpx;
    }

    .goods-num {
        display: block;
        font-size: 24rpx;
        color: #999999;
    }

    .info-card {
        padding: 0;
    }

    .info-row {
        display: flex;
        align-items: flex-start;
        justify-content: space-between;
        padding: 28rpx;
        border-bottom: 1rpx solid #F5F5F5;

        &:last-child {
            border-bottom: none;
        }
    }

    .info-label {
        font-size: 28rpx;
        color: #999999;
        flex-shrink: 0;
    }

    .info-value {
        font-size: 28rpx;
        color: #1A1A1A;
        text-align: right;
        max-width: 65%;
        display: flex;
        align-items: center;
        gap: 12rpx;

        &.remark {
            max-width: none;
            color: #666666;
        }
    }

    .copy-btn {
        font-size: 22rpx;
        color: #22A84F;
        background: #E8F8EC;
        padding: 4rpx 12rpx;
        border-radius: 8rpx;
    }

    .price-card {
        padding: 28rpx;
    }

    .price-row {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 16rpx;

        &:last-child {
            margin-bottom: 0;
        }

        &.total {
            margin-top: 16rpx;
            padding-top: 16rpx;
            border-top: 1rpx dashed #EEEEEE;
        }
    }

    .price-label {
        font-size: 28rpx;
        color: #666666;
    }

    .price-value {
        font-size: 28rpx;
        color: #1A1A1A;

        &.primary {
            font-size: 36rpx;
            font-weight: 700;
            color: #F97316;
        }
    }

    .bottom-bar {
        position: fixed;
        bottom: 0;
        left: 0;
        right: 0;
        height: 120rpx;
        background: #FFFFFF;
        box-shadow: 0 -2rpx 20rpx rgba(0, 0, 0, 0.06);
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 0 32rpx;
        padding-bottom: env(safe-area-inset-bottom);
        z-index: 100;
    }

    .bottom-left {
        display: flex;
        align-items: center;
    }

    .total-info {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
    }

    .total-label {
        font-size: 24rpx;
        color: #999999;
    }

    .total-price {
        font-size: 36rpx;
        font-weight: 700;
        color: #F97316;
    }

    .bottom-right {
        display: flex;
        align-items: center;
        gap: 20rpx;
    }

    .action-btn {
        height: 72rpx;
        padding: 0 36rpx;
        border-radius: 36rpx;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 28rpx;
        font-weight: 500;

        &.cancel {
            background: #F5F5F5;
            color: #666666;
        }

        &.pay {
            background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
            color: #FFFFFF;
            box-shadow: 0 4rpx 16rpx rgba(34, 168, 79, 0.3);
        }

        &.confirm {
            background: linear-gradient(135deg, #3B82F6 0%, #2563EB 100%);
            color: #FFFFFF;
            box-shadow: 0 4rpx 16rpx rgba(59, 130, 246, 0.3);
        }

        &.delete {
            background: #FEE2E2;
            color: #EF4444;
        }
    }

    .bottom-placeholder {
        height: 140rpx;
    }
</style>
