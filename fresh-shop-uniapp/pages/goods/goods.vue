<!--
 * 优诚配运 - 商品列表
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
                    placeholder="搜索商品名称"
                    placeholder-class="placeholder"
                />
                <u-icon v-if="keyword" name="close-circle-fill" size="32rpx" color="#CCCCCC" @click="clearKeyword"></u-icon>
            </view>
            <view class="filter-btn" @click="showFilter = !showFilter">
                <u-icon name="filter" size="36rpx" :color="hasFilter ? '#22A84F' : '#666666'"></u-icon>
            </view>
        </view>

        <!-- 筛选面板 -->
        <view class="filter-panel" v-if="showFilter">
            <view class="filter-section">
                <view class="filter-label">排序</view>
                <view class="filter-tags">
                    <view
                        class="filter-tag"
                        :class="{ active: sortType === 0 }"
                        @click="setSortType(0)"
                    >智能排序</view>
                    <view
                        class="filter-tag"
                        :class="{ active: sortType === 1 }"
                        @click="setSortType(1)"
                    >最新上架</view>
                    <view
                        class="filter-tag"
                        :class="{ active: sortType === 2 }"
                        @click="setSortType(2)"
                    >销量优先</view>
                </view>
            </view>

            <view class="filter-section" v-if="categoryList.length > 0">
                <view class="filter-label">分类</view>
                <view class="filter-tags">
                    <view
                        class="filter-tag"
                        :class="{ active: categoryId === 0 }"
                        @click="setCategoryId(0)"
                    >全部</view>
                    <view
                        class="filter-tag"
                        :class="{ active: categoryId === item.ID }"
                        v-for="item in categoryList"
                        :key="item.ID"
                        @click="setCategoryId(item.ID)"
                    >{{ item.title }}</view>
                </view>
            </view>

            <view class="filter-section" v-if="brandList.length > 0">
                <view class="filter-label">品牌</view>
                <view class="filter-tags">
                    <view
                        class="filter-tag"
                        :class="{ active: brandId === 0 }"
                        @click="setBrandId(0)"
                    >全部</view>
                    <view
                        class="filter-tag"
                        :class="{ active: brandId === item.ID }"
                        v-for="item in brandList"
                        :key="item.ID"
                        @click="setBrandId(item.ID)"
                    >{{ item.name }}</view>
                </view>
            </view>

            <view class="filter-actions">
                <view class="filter-reset" @click="resetFilter">重置</view>
                <view class="filter-confirm" @click="confirmFilter">确定</view>
            </view>
        </view>

        <!-- 商品列表 -->
        <view class="goods-content" :class="{ 'with-filter': showFilter }">
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
                <GoodsList :lists="goodsArr" price-type="￥" :is-audit="isAudit"></GoodsList>

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
                    <u-icon name="shopping-cart" size="120rpx" color="#CCCCCC"></u-icon>
                </view>
                <view class="empty-text">暂无相关商品</view>
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
    import { getBrandListAll } from '@/api/brand.js'
    import { getCategoryListAll } from '@/api/category.js'
    import { getGoodsPageList, getGoodsPageListLoading } from '@/api/goods.js'
    import GoodsList from '@/components/goodsList/goodsList.vue'
    import { getUser } from '@/store/storage.js'
    import { setUser } from '@/store/storage.js'
    import { getUserAuditStatus } from '@/api/user.js'

    export default {
        components: {
            GoodsList
        },
        data() {
            return {
                keyword: '',
                sortType: 0,
                categoryId: 0,
                brandId: 0,
                categoryList: [],
                brandList: [],
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
                showFilter: false,
                isAudit: false,
                scrollTimer: null
            }
        },
        computed: {
            hasFilter() {
                return this.categoryId !== 0 || this.brandId !== 0 || this.sortType !== 0
            }
        },
        onLoad(options) {
            if (options.keyword) this.keyword = options.keyword
            if (options.categoryId) this.categoryId = parseInt(options.categoryId)
            if (options.brandId) this.brandId = parseInt(options.brandId)
            if (options.sortType) this.sortType = parseInt(options.sortType)

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

            this.init()
            this.calculateHeight()
        },
        onShow() {
            // 返回时刷新
            if (this.goodsArr.length > 0) {
                this.onRefresh()
            }
        },
        methods: {
            async init() {
                this.loading = true
                try {
                    const [categoryRes, brandRes] = await Promise.all([
                        getCategoryListAll(),
                        getBrandListAll()
                    ])
                    this.categoryList = categoryRes.data || []
                    this.brandList = brandRes.data || []
                } catch (e) { }

                await this.getGoodsListData(0)
                this.loading = false
            },
            calculateHeight() {
                uni.getSystemInfo({
                    success: (res) => {
                        this.scrollViewHeight = res.windowHeight - 110
                    }
                })
            },
            async getGoodsListData(type = 0) {
                if (type === 0) {
                    this.page.page = 1
                    this.page.isMore = true
                    this.loadMore = 'loadmore'
                }

                const data = {
                    goodsArea: 0,
                    page: this.page.page,
                    pageSize: this.page.pageSize
                }

                if (this.keyword) data.name = this.keyword
                if (this.categoryId) data.categoryId = this.categoryId
                if (this.brandId) data.brandId = this.brandId
                if (this.sortType) data.sortType = this.sortType

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
                }
            },
            async onRefresh() {
                this.triggered = true
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
                this.showFilter = false
                this.getGoodsListData(0)
            },
            clearKeyword() {
                this.keyword = ''
                this.getGoodsListData(0)
            },
            setSortType(val) {
                this.sortType = val
            },
            setCategoryId(val) {
                this.categoryId = val
            },
            setBrandId(val) {
                this.brandId = val
            },
            resetFilter() {
                this.sortType = 0
                this.categoryId = 0
                this.brandId = 0
            },
            confirmFilter() {
                this.showFilter = false
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
        gap: 16rpx;
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

    .filter-btn {
        width: 72rpx;
        height: 72rpx;
        background: #F5F7F4;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .filter-panel {
        background: #FFFFFF;
        padding: 24rpx;
        border-bottom: 1rpx solid #EEEEEE;
    }

    .filter-section {
        margin-bottom: 24rpx;

        &:last-of-type {
            margin-bottom: 0;
        }
    }

    .filter-label {
        font-size: 26rpx;
        color: #999999;
        margin-bottom: 16rpx;
    }

    .filter-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 16rpx;
    }

    .filter-tag {
        padding: 12rpx 28rpx;
        background: #F5F7F4;
        border-radius: 32rpx;
        font-size: 26rpx;
        color: #666666;
        border: 2rpx solid transparent;

        &.active {
            background: #E8F8EC;
            color: #22A84F;
            border-color: #22A84F;
        }
    }

    .filter-actions {
        display: flex;
        gap: 20rpx;
        margin-top: 24rpx;
    }

    .filter-reset {
        flex: 1;
        height: 80rpx;
        background: #F5F5F5;
        border-radius: 40rpx;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 28rpx;
        color: #666666;
    }

    .filter-confirm {
        flex: 2;
        height: 80rpx;
        background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
        border-radius: 40rpx;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 28rpx;
        font-weight: 600;
        color: #FFFFFF;
    }

    .goods-content {
        &.with-filter {
            padding-top: 0;
        }
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
    }

    .loading-wrap {
        display: flex;
        justify-content: center;
        padding: 48rpx;
    }
</style>
