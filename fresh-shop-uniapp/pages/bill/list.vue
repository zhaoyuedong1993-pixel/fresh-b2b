<!--
 * 优诚配运 - 账单列表
 * 设计规范：自然清新
-->
<template>
    <pageWrapper>
        <!-- 顶部 -->
        <view class="page-header">
            <text class="page-title">我的账单</text>
            <text class="page-desc">月度账单明细</text>
        </view>

        <!-- 账单列表 -->
        <view class="bill-list" v-if="billList.length > 0">
            <view class="bill-card" v-for="item in billList" :key="item.period" @click="goDetail(item)">
                <view class="bill-left">
                    <text class="period">{{ item.period }}</text>
                    <text class="count">共 {{ item.order_count }} 笔订单</text>
                </view>
                <view class="bill-right">
                    <text class="amount">
                        <text class="sym">¥</text>{{ formatAmount(item.total_amount) }}
                    </text>
                    <view class="status" :class="item.settlement_type === 2 ? 'settled' : 'pending'">
                        {{ item.settlement_type === 2 ? '已结算' : '待结算' }}
                    </view>
                </view>
                <u-icon name="arrow-right" color="#CCCCCC" size="32rpx"></u-icon>
            </view>
        </view>

        <!-- 空状态 -->
        <view class="empty-wrap" v-else-if="!loading">
            <view class="empty-icon">
                <u-icon name="red-packet" size="120rpx" color="#CCCCCC"></u-icon>
            </view>
            <view class="empty-text">暂无账单记录</view>
        </view>

        <!-- 加载中 -->
        <view class="loading-wrap" v-if="loading">
            <text>加载中...</text>
        </view>

        <u-toast ref="toast" style="z-index:9999"></u-toast>
    </pageWrapper>
</template>

<script>
    import { getBillList } from '@/api/business/bill.js'
    import { getToken } from '@/store/storage.js'

    export default {
        data() {
            return {
                billList: [],
                loading: false,
                page: 1,
                pageSize: 20,
                hasMore: true
            }
        },
        onLoad() {
            this.loadBillList()
        },
        onReachBottom() {
            if (this.hasMore) {
                this.loadBillList(true)
            }
        },
        methods: {
            async loadBillList(append = false) {
                const token = getToken()
                if (!token) {
                    uni.showToast({ title: '请先登录', icon: 'none' })
                    return
                }

                if (append) {
                    this.page++
                } else {
                    this.page = 1
                }

                this.loading = true
                try {
                    const res = await getBillList({
                        page: this.page,
                        pageSize: this.pageSize
                    })

                    if (res.code === 0) {
                        const list = res.data?.list || []
                        if (append) {
                            this.billList = [...this.billList, ...list]
                        } else {
                            this.billList = list
                        }
                        this.hasMore = list.length >= this.pageSize
                    }
                } catch (e) {
                    uni.showToast({ title: '加载失败', icon: 'none' })
                } finally {
                    this.loading = false
                }
            },
            formatAmount(amount) {
                return (Number(amount) || 0).toFixed(2)
            },
            goDetail(item) {
                uni.navigateTo({
                    url: `/pages/bill/detail?period=${item.period}`
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .page-header {
        background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
        padding: 48rpx 32rpx 32rpx;
    }

    .page-title {
        display: block;
        font-size: 48rpx;
        font-weight: 700;
        color: #FFFFFF;
    }

    .page-desc {
        display: block;
        font-size: 26rpx;
        color: rgba(255, 255, 255, 0.7);
        margin-top: 8rpx;
    }

    .bill-list {
        padding: 24rpx;
    }

    .bill-card {
        background: #FFFFFF;
        border-radius: 20rpx;
        padding: 28rpx 24rpx;
        margin-bottom: 20rpx;
        display: flex;
        align-items: center;
        box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.06);

        &:active {
            background: #FAFAFA;
        }
    }

    .bill-left {
        flex: 1;
    }

    .period {
        display: block;
        font-size: 32rpx;
        font-weight: 600;
        color: #1A1A1A;
        margin-bottom: 8rpx;
    }

    .count {
        display: block;
        font-size: 24rpx;
        color: #999999;
    }

    .bill-right {
        text-align: right;
        margin-right: 16rpx;
    }

    .amount {
        display: block;
        font-size: 36rpx;
        font-weight: 700;
        color: #F97316;
        margin-bottom: 8rpx;
    }

    .sym {
        font-size: 24rpx;
        font-weight: 600;
    }

    .status {
        font-size: 24rpx;
        border-radius: 20rpx;
        padding: 4rpx 16rpx;
        display: inline-block;

        &.settled {
            background: #E8F8EC;
            color: #22A84F;
        }

        &.pending {
            background: #FEF3E2;
            color: #F97316;
        }
    }

    .empty-wrap {
        display: flex;
        flex-direction: column;
        align-items: center;
        padding-top: 200rpx;
        gap: 24rpx;
    }

    .empty-text {
        font-size: 28rpx;
        color: #999999;
    }

    .loading-wrap {
        text-align: center;
        padding: 32rpx;
        color: #999999;
        font-size: 26rpx;
    }
</style>
