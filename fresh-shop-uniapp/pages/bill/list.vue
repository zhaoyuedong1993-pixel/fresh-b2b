<template>
    <pageWrapper>
        <!-- 头部 -->
        <view class="header">
            <view class="title">我的账单</view>
            <view class="subtitle">月度账单明细</view>
        </view>

        <!-- 账单列表 -->
        <view class="bill-list" v-if="billList.length > 0">
            <view
                class="bill-item"
                v-for="item in billList"
                :key="item.id"
                @click="goDetail(item)"
            >
                <view class="bill-info">
                    <view class="period">{{ item.period }}</view>
                    <view class="count">共 {{ item.order_count }} 笔订单</view>
                </view>
                <view class="bill-amount">
                    <view class="money">¥{{ formatAmount(item.total_amount) }}</view>
                    <view :class="['status', item.status === 1 ? 'settled' : 'unsettled']">
                        {{ item.status === 1 ? '已结算' : '待结算' }}
                    </view>
                </view>
                <view class="arrow">›</view>
            </view>
        </view>

        <!-- 空状态 -->
        <view class="empty" v-else>
            <image class="empty-icon" src="../../../static/nopicture.jpg" mode="aspectFit"></image>
            <view class="empty-text">暂无账单记录</view>
        </view>

        <!-- 加载状态 -->
        <view class="loading" v-if="loading">
            <text>加载中...</text>
        </view>
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
        if (this.hasMore && !this.loading) {
            this.page++
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

            this.loading = true
            try {
                const res = await getBillList({
                    page: this.page,
                    pageSize: this.pageSize
                })

                if (res.code === 0) {
                    const list = res.data.list || []
                    if (append) {
                        this.billList = [...this.billList, ...list]
                    } else {
                        this.billList = list
                    }
                    this.hasMore = list.length >= this.pageSize
                }
            } catch (e) {
                console.error(e)
                uni.showToast({ title: '加载失败', icon: 'none' })
            } finally {
                this.loading = false
            }
        },

        formatAmount(amount) {
            return (amount || 0).toFixed(2)
        },

        goDetail(item) {
            uni.navigateTo({
                url: `/pages/bill/detail?id=${item.id}&period=${item.period}`
            })
        }
    }
}
</script>

<style lang="scss" scoped>
.header {
    background: linear-gradient(135deg, #4CAF50, #8BC34A);
    padding: 40rpx 30rpx;
    color: #fff;

    .title {
        font-size: 40rpx;
        font-weight: bold;
    }

    .subtitle {
        font-size: 24rpx;
        margin-top: 10rpx;
        opacity: 0.9;
    }
}

.bill-list {
    padding: 20rpx;
}

.bill-item {
    background: #fff;
    border-radius: 16rpx;
    padding: 30rpx;
    margin-bottom: 20rpx;
    display: flex;
    align-items: center;
    box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.05);
}

.bill-info {
    flex: 1;

    .period {
        font-size: 32rpx;
        font-weight: bold;
        color: #333;
    }

    .count {
        font-size: 24rpx;
        color: #999;
        margin-top: 8rpx;
    }
}

.bill-amount {
    text-align: right;
    margin-right: 20rpx;

    .money {
        font-size: 36rpx;
        font-weight: bold;
        color: #4CAF50;
    }

    .status {
        font-size: 24rpx;
        margin-top: 8rpx;
    }

    .status.settled {
        color: #67c23a;
    }

    .status.unsettled {
        color: #e6a23c;
    }
}

.arrow {
    font-size: 40rpx;
    color: #ccc;
}

.empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 200rpx 0;

    .empty-icon {
        width: 200rpx;
        height: 200rpx;
        opacity: 0.5;
    }

    .empty-text {
        font-size: 28rpx;
        color: #999;
        margin-top: 20rpx;
    }
}

.loading {
    text-align: center;
    padding: 30rpx;
    color: #999;
    font-size: 24rpx;
}
</style>
