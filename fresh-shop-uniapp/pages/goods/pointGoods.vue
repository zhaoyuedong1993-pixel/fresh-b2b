<!--
 * 优诚配运 - 积分商品
 * 设计规范：自然清新
-->
<template>
    <pageWrapper>
        <!-- 顶部搜索栏 -->
        <view class="search-header">
            <view class="search-box">
                <u-icon name="search" size="32rpx" color="#999999"></u-icon>
                <input
                    class="search-input"
                    v-model="keyword"
                    type="text"
                    confirm-type="search"
                    @confirm="searchConfirm"
                    placeholder="搜索积分商品"
                    placeholder-class="placeholder"
                />
                <u-icon v-if="keyword" name="close-circle-fill" size="32rpx" color="#CCCCCC" @click="clearKeyword"></u-icon>
            </view>
        </view>

        <!-- 积分余额 -->
        <view class="point-banner">
            <view class="point-info">
                <u-icon name="star-fill" size="40rpx" color="#FFD700"></u-icon>
                <view class="point-text">
                    <text class="point-label">我的积分</text>
                    <text class="point-value">{{ pointAmount }} <text class="point-unit">积分</text></text>
                </view>
            </view>
        </view>

        <!-- 商品列表 -->
        <view class="goods-content">
            <scroll-view
                v-if="goodsArr.length > 0"
                class="goods-scroll"
                scroll-y
                :style="{ height: scrollViewHeight + 'px' }"
                :scroll-top="scrollTop"
                refresher-enabled
                :refresher-triggered="triggered"
                @refresherrefresh="onRefresh"
                @scrolltolower="scrollTolower"
                @scroll="onScroll"
            >
                <GoodsList :lists="goodsArr" :isPoint="true" :is-audit="isAudit"></GoodsList>

                <view class="load-more">
                    <u-loadmore
                        :status="loadMore"
                        loading-text="努力加载中..."
                        loadmore-text="上拉加载更多"
                        nomore-text="没有更多了"
                    />
                </view>
            </scroll-view>

            <!-- 空状态 -->
            <view class="empty-wrap" v-else-if="!loading">
                <view class="empty-icon">
                    <u-icon name="star" size="120rpx" color="#CCCCCC"></u-icon>
                </view>
                <view class="empty-text">暂无积分商品</view>
                <view class="empty-btn" @click="clearKeyword">清空搜索</view>
            </view>

            <!-- 加载中 -->
            <view class="loading-wrap" v-if="loading">
                <u-loading-icon size="48rpx" color="#22A84F"></u-loading-icon>
            </view>
        </view>

        <u-back-top :scroll-top="scrollTop" @click="toTop"></u-back-top>

        <u-toast ref="toast" style="z-index: 9998;"></u-toast>
    </pageWrapper>
</template>

<script>
    import config from '@/config/config.js'
    import { getGoodsPageList, getGoodsPageListLoading } from '@/api/goods.js'
    import { getAccountInfo } from '@/api/account.js'
    import GoodsList from '@/components/goodsList/goodsList.vue'
    import { getUser, setUser } from '@/store/storage.js'
    import { getUserAuditStatus } from '@/api/user.js'

    export default {
        components: {
            GoodsList
        },
        data() {
            return {
                keyword: '',
                pointAmount: 0,
                goodsArr: [],
                scrollViewHeight: 600,
                scrollTop: 0,
                page: {
                    page: 1,
                    pageSize: 10,
                    total: 0,
                    isMore: true
                },
                triggered: false,
                loadMore: 'loadmore',
                loading: false,
                isAudit: false,
                scrollTimer: null
            }
        },
        onLoad() {
            const user = getUser()
            if (user) {
                if (user.auditStatus === 1) {
                    this.isAudit = true
                } else {
                    getUserAuditStatus().then(res => {
                        if (res.data?.auditStatus === 1) {
                            this.isAudit = true
                            user.auditStatus = 1
                            setUser(user)
                        }
                    })
                }
            }

            this.loadPointAmount()
            this.getGoodsListData(0)
            this.calculateHeight()
        },
        methods: {
            calculateHeight() {
                uni.getSystemInfo({
                    success: (res) => {
                        this.scrollViewHeight = res.windowHeight - 180
                    }
                })
            },
            async loadPointAmount() {
                try {
                    const res = await getAccountInfo(2)
                    if (res.code === 0) {
                        this.pointAmount = res.data?.account?.amount || 0
                    }
                } catch (e) { }
            },
            async getGoodsListData(type = 0) {
                if (type === 0) {
                    this.page.page = 1
                    this.page.isMore = true
                    this.loadMore = 'loadmore'
                }

                this.loading = type === 0

                const data = {
                    goodsArea: 1,
                    page: this.page.page,
                    pageSize: this.page.pageSize
                }

                if (this.keyword) data.name = this.keyword

                try {
                    const res = type === 0
                        ? await getGoodsPageListLoading(data)
                        : await getGoodsPageList(data)

                    if (res.code !== 0) return false

                    const list = res.data?.list || []
                    list.forEach(item => {
                        if (item.images?.[0]?.url && item.images[0].url.slice(0, 4) !== 'http') {
                            item.images[0].url = config.baseUrl + '/' + item.images[0].url
                        }
                    })

                    this.page.total = res.data?.total || 0
                    if (this.page.page * this.page.pageSize >= this.page.total) {
                        this.page.isMore = false
                        this.loadMore = 'nomore'
                    }

                    if (type === 0) {
                        this.goodsArr = list
                    } else {
                        this.goodsArr = [...this.goodsArr, ...list]
                    }

                    this.page.page++
                    return true
                } catch (e) {
                    return false
                } finally {
                    this.loading = false
                }
            },
            async onRefresh() {
                this.triggered = true
                await this.loadPointAmount()
                const b = await this.getGoodsListData(0)
                this.$message(this.$refs.toast).success(b ? '刷新成功' : '刷新失败')
                this.triggered = false
            },
            async scrollTolower() {
                if (this.loadMore === 'loading' || !this.page.isMore) return
                this.loadMore = 'loading'
                await this.getGoodsListData(1)
                this.loadMore = this.page.isMore ? 'loadmore' : 'nomore'
            },
            searchConfirm() {
                this.getGoodsListData(0)
            },
            clearKeyword() {
                this.keyword = ''
                this.getGoodsListData(0)
            },
            onScroll(e) {
                if (this.scrollTimer) clearTimeout(this.scrollTimer)
                this.scrollTimer = setTimeout(() => {
                    this.scrollTop = e.detail.scrollTop
                }, 100)
            },
            toTop() {
                this.scrollTop = 0
            }
        }
    }
</script>

<style lang="scss" scoped>
    .search-header {
        display: flex;
        align-items: center;
        padding: 20rpx 24rpx;
        background: #FFFFFF;
        position: sticky;
        top: 0;
        z-index: 100;
        box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);
    }

    .search-box {
        flex: 1;
        height: 72rpx;
        background: #F5F7F4;
        border-radius: 36rpx;
        display: flex;
        align-items: center;
        padding: 0 24rpx;
        gap: 12rpx;
    }

    .search-input {
        flex: 1;
        height: 100%;
        font-size: 28rpx;
        color: #1A1A1A;
    }

    .placeholder {
        color: #CCCCCC;
        font-size: 28rpx;
    }

    .point-banner {
        background: linear-gradient(135deg, #F97316 0%, #EA580C 100%);
        padding: 24rpx 32rpx;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .point-info {
        display: flex;
        align-items: center;
        gap: 16rpx;
    }

    .point-text {
        display: flex;
        flex-direction: column;
    }

    .point-label {
        font-size: 24rpx;
        color: rgba(255, 255, 255, 0.8);
    }

    .point-value {
        font-size: 48rpx;
        font-weight: 700;
        color: #FFFFFF;
    }

    .point-unit {
        font-size: 24rpx;
        font-weight: 400;
    }

    .goods-scroll {
        padding: 24rpx;
    }

    .load-more {
        padding: 24rpx 0;
    }

    .empty-wrap {
        display: flex;
        flex-direction: column;
        align-items: center;
        padding-top: 160rpx;
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
    }

    .loading-wrap {
        display: flex;
        justify-content: center;
        padding: 48rpx;
    }
</style>
