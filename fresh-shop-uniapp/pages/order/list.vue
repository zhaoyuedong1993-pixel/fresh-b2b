<!--
 * 优诚配运 - 订单列表
 * 设计规范：自然清新
-->
<template>
    <pageWrapper>
        <!-- 顶部 Tab -->
        <view class="header-tabs">
            <view class="tab-bar">
                <view
                    class="tab-item"
                    :class="{ active: currentTab === index }"
                    v-for="(tab, index) in tabs"
                    :key="tab.status"
                    @click="switchTab(index)"
                >
                    <text>{{ tab.name }}</text>
                    <view class="tab-line" v-if="currentTab === index"></view>
                </view>
            </view>
        </view>

        <!-- 订单列表 -->
        <scroll-view class="order-scroll" scroll-y :style="{ height: scrollHeight + 'px' }"
            refresher-enabled :refresher-triggered="isRefreshing" @refresherrefresh="onRefresh"
            @scrolltolower="onScrollLower">
            <view class="order-list" v-if="orderList.length > 0">
                <orderList :status="currentStatus" />
            </view>

            <!-- 空状态 -->
            <view class="empty-wrap" v-else>
                <view class="empty-icon">
                    <u-icon name="order" size="120rpx" color="#CCCCCC"></u-icon>
                </view>
                <view class="empty-text">暂无相关订单</view>
                <view class="empty-btn" @click="goShopping">去选购</view>
            </view>
        </scroll-view>

        <u-toast ref="toast" style="z-index: 9999"></u-toast>
    </pageWrapper>
</template>

<script>
    import orderList from "@/components/orderList/orderList"

    export default {
        components: { orderList },
        data() {
            return {
                tabs: [
                    { name: '全部', status: null },
                    { name: '未付款', status: 0 },
                    { name: '备货中', status: 1 },
                    { name: '配送中', status: 2 },
                    { name: '已完成', status: 3 }
                ],
                currentTab: 0,
                currentStatus: null,
                orderList: [],
                scrollHeight: 600,
                isRefreshing: false
            }
        },
        onLoad(options) {
            if (options.status !== undefined && options.status !== 'null') {
                const s = parseInt(options.status)
                const idx = this.tabs.findIndex(t => t.status === s)
                if (idx >= 0) {
                    this.currentTab = idx
                    this.currentStatus = s
                }
            }
        },
        onReady() {
            uni.getSystemInfo({
                success: (res) => {
                    this.scrollHeight = res.windowHeight - 90
                }
            })
        },
        methods: {
            switchTab(index) {
                this.currentTab = index
                this.currentStatus = this.tabs[index].status
            },
            async onRefresh() {
                this.isRefreshing = true
                // orderList 组件会自己处理刷新
                setTimeout(() => {
                    this.isRefreshing = false
                    this.$message(this.$refs.toast).success('刷新成功')
                }, 1000)
            },
            onScrollLower() { },
            goShopping() {
                uni.switchTab({ url: '/pages/index/index' })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .header-tabs {
        height: 90rpx;
        background: #FFFFFF;
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        z-index: 100;
        box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.06);
    }

    .tab-bar {
        display: flex;
        height: 90rpx;
    }

    .tab-item {
        flex: 1;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        position: relative;
        font-size: 28rpx;
        color: #999999;
        transition: all 0.2s;

        &.active {
            color: #22A84F;
            font-weight: 600;
        }
    }

    .tab-line {
        position: absolute;
        bottom: 0;
        left: 50%;
        transform: translateX(-50%);
        width: 48rpx;
        height: 6rpx;
        background: linear-gradient(90deg, #22A84F, #1A9A45);
        border-radius: 3rpx;
    }

    .order-scroll {
        padding-top: 90rpx;
        background: #F5F7F4;
    }

    .order-list {
        padding: 24rpx;
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

    .empty-btn {
        background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
        color: #FFFFFF;
        font-size: 28rpx;
        padding: 20rpx 60rpx;
        border-radius: 40rpx;
        box-shadow: 0 4rpx 16rpx rgba(34, 168, 79, 0.3);
    }
</style>
