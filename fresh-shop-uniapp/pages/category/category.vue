<!--
 * 优诚配运 - 分类页
 * 设计规范：自然清新
-->
<template>
    <pageWrapper>
        <!-- 顶部搜索 -->
        <view class="header-bar">
            <view class="search-box" @click="searchClick">
                <u-icon name="search" color="#999999" size="32rpx"></u-icon>
                <text class="search-txt">搜索商品</text>
            </view>
        </view>

        <!-- 主体：左侧分类 + 右侧商品 -->
        <view class="main-content">
            <!-- 左侧分类 -->
            <scroll-view class="left-panel" scroll-y :scroll-top="leftScroll">
                <view
                    class="category-item"
                    :class="{ active: currentCategoryId === item.ID }"
                    v-for="item in categoryList"
                    :key="item.ID"
                    @click="selectCategory(item.ID)"
                >
                    <text class="category-name">{{ item.title }}</text>
                    <view class="category-indicator" v-if="currentCategoryId === item.ID"></view>
                </view>
            </scroll-view>

            <!-- 右侧商品 -->
            <view class="right-panel">
                <!-- 品牌横向滚动 -->
                <view class="brand-section" v-if="brandList.length > 1">
                    <scroll-view scroll-x enable-flex class="brand-scroll">
                        <view
                            class="brand-chip"
                            :class="{ active: currentBrandId === 0 }"
                            @click="selectBrand(0)"
                        >全部</view>
                        <view
                            class="brand-chip"
                            :class="{ active: currentBrandId === b.ID }"
                            v-for="b in brandList"
                            :key="b.ID"
                            @click="selectBrand(b.ID)"
                        >
                            <u--image width="40rpx" height="40rpx" :src="b.logo" radius="8rpx"></u--image>
                            <text class="brand-name">{{ b.name }}</text>
                        </view>
                    </scroll-view>
                </view>

                <!-- 商品列表 -->
                <scroll-view
                    class="goods-scroll"
                    scroll-y
                    :style="{ height: scrollHeight + 'px' }"
                    refresher-enabled
                    :refresher-triggered="isRefreshing"
                    @refresherrefresh="onRefresh"
                    @scrolltolower="onScrollLower"
                    :scroll-anchoring="true"
                >
                    <GoodsList
                        :vertical="true"
                        :lists="goodsArr"
                        price-type="¥"
                        :is-audit="isAudit"
                    ></GoodsList>

                    <view class="load-more" v-if="goodsArr.length > 0">
                        <u-loadmore
                            :status="loadMore"
                            loading-text="加载中..."
                            loadmore-text="上拉加载更多"
                            nomore-text="没有更多了"
                        />
                    </view>

                    <u-empty
                        v-if="goodsArr.length === 0 && !isLoading"
                        :style="{ height: scrollHeight * 0.6 + 'px' }"
                        text="暂无商品"
                        mode="data"
                        icon="http://cdn.uviewui.com/uview/empty/data.png"
                    ></u-empty>
                </scroll-view>
            </view>
        </view>

        <!-- 登录悬浮 -->
        <loginSuspend :show="loginSuspendShow" @success="loginSuccess"></loginSuspend>

        <!-- 底部导航 -->
        <Tabbar :tabsId="1" />

        <u-toast style="z-index:9998;" ref="toast"></u-toast>
    </pageWrapper>
</template>

<script>
    import Tabbar from '@/components/tabbar/tabbar.vue'
    import GoodsList from '@/components/goodsList/goodsList.vue'
    import loginSuspend from '@/components/loginPop/loginSuspend.vue'
    import { getCategoryListAll } from '@/api/category.js'
    import { getBrandListByCategoryId } from '@/api/brand.js'
    import { getGoodsPageList } from '@/api/goods.js'
    import config from '@/config/config.js'
    import { getUser, getToken, setUser } from '@/store/storage.js'
    import { getUserAuditStatus } from '@/api/user'

    export default {
        components: { Tabbar, GoodsList, loginSuspend },
        data() {
            return {
                isAudit: false,
                isRefreshing: false,
                isLoading: false,
                loginSuspendShow: false,
                scrollHeight: 600,
                leftScroll: 0,
                categoryList: [{ ID: 0, title: '全部' }],
                brandList: [],
                goodsArr: [],
                page: { page: 1, pageSize: 12, total: 0, isMore: true },
                loadMore: 'loadmore',
                currentCategoryId: 0,
                currentBrandId: 0
            }
        },
        onLoad() {
            let user = getUser()
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
            const t = getToken()
            if (!t) {
                this.loginSuspendShow = true
            }
            this.init()
        },
        onReady() {
            uni.getSystemInfo({
                success: (res) => {
                    this.scrollHeight = res.windowHeight - 110
                }
            })
        },
        methods: {
            async init() {
                const res = await getCategoryListAll()
                if (res.code !== 0) return
                this.categoryList = [{ ID: 0, title: '全部' }, ...(res.data || [])]
                await this.loadBrandList()
                await this.loadGoodsList(false)
            },
            async loadBrandList() {
                const res = await getBrandListByCategoryId({ categoryId: this.currentCategoryId })
                if (res.code !== 0 || !res.data) {
                    this.brandList = []
                    return
                }
                res.data.forEach(item => {
                    if (item.logo && item.logo.slice(0, 4) !== 'http') {
                        item.logo = config.baseUrl + '/' + item.logo
                    }
                })
                this.brandList = res.data
            },
            async loadGoodsList(append = false) {
                const data = {}
                if (this.currentCategoryId > 0) {
                    data.categoryId = this.currentCategoryId
                }
                if (this.currentBrandId > 0) {
                    data.brandId = this.currentBrandId
                }
                if (!append) {
                    this.page.page = 1
                    this.page.isMore = true
                    this.loadMore = 'loadmore'
                }
                data.page = this.page.page
                data.pageSize = this.page.pageSize

                this.isLoading = true
                const res = append ? await getGoodsPageList(data) : await this.$refs.toast ?
                    new Promise(resolve => {
                        getGoodsPageList(data).then(resolve)
                    }) : await getGoodsPageList(data)
                if (res.code !== 0) {
                    this.isLoading = false
                    return
                }

                res.data.list?.forEach(item => {
                    if (item.images?.[0] && item.images[0].url?.slice(0, 4) !== 'http') {
                        item.images[0].url = config.baseUrl + '/' + item.images[0].url
                    }
                })

                this.page.total = res.data.total || 0
                if (append) {
                    this.goodsArr = [...this.goodsArr, ...(res.data.list || [])]
                } else {
                    this.goodsArr = res.data.list || []
                }

                if (this.page.page * this.page.pageSize >= this.page.total) {
                    this.page.isMore = false
                    this.loadMore = 'nomore'
                }
                this.page.page++
                this.isLoading = false
            },
            selectCategory(id) {
                if (this.currentCategoryId === id) return
                this.currentCategoryId = id
                this.currentBrandId = 0
                this.loadBrandList()
                this.loadGoodsList(false)
            },
            selectBrand(id) {
                this.currentBrandId = id
                this.loadGoodsList(false)
            },
            searchClick() {
                uni.navigateTo({ url: '/pages/search/search' })
            },
            async onRefresh() {
                this.isRefreshing = true
                await this.loadGoodsList(false)
                this.isRefreshing = false
                this.$message(this.$refs.toast).success("刷新成功")
            },
            async onScrollLower() {
                if (this.loadMore === 'loading' || !this.page.isMore) return
                this.loadMore = 'loading'
                await this.loadGoodsList(true)
                this.loadMore = this.page.isMore ? 'loadmore' : 'nomore'
            },
            loginSuccess() {
                this.loginSuspendShow = false
            }
        }
    }
</script>

<style lang="scss" scoped>
    /* 顶部搜索 */
    .header-bar {
        height: 100rpx;
        padding: 16rpx 24rpx;
        background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
        display: flex;
        align-items: center;
    }

    .search-box {
        display: flex;
        align-items: center;
        height: 72rpx;
        background: #FFFFFF;
        border-radius: 36rpx;
        padding: 0 28rpx;
        flex: 1;
        box-shadow: 0 4rpx 16rpx rgba(34, 168, 79, 0.2);
    }

    .search-txt {
        margin-left: 12rpx;
        font-size: 26rpx;
        color: #999999;
    }

    /* 主体 */
    .main-content {
        display: flex;
        height: calc(100vh - 180rpx);
        background: #F5F7F4;
    }

    /* 左侧分类 */
    .left-panel {
        width: 180rpx;
        background: #FFFFFF;
        border-radius: 0 24rpx 24rpx 0;
        padding: 16rpx 0;
    }

    .category-item {
        display: flex;
        align-items: center;
        justify-content: center;
        height: 100rpx;
        position: relative;
        transition: all 0.2s;
    }

    .category-name {
        font-size: 26rpx;
        color: #666666;
        text-align: center;
        padding: 0 16rpx;
        transition: all 0.2s;
    }

    .category-item.active {
        background: #F5F7F4;
    }

    .category-item.active .category-name {
        color: #22A84F;
        font-weight: 600;
    }

    .category-indicator {
        position: absolute;
        left: 0;
        top: 50%;
        transform: translateY(-50%);
        width: 6rpx;
        height: 48rpx;
        background: linear-gradient(180deg, #22A84F, #1A9A45);
        border-radius: 0 3rpx 3rpx 0;
    }

    /* 右侧商品 */
    .right-panel {
        flex: 1;
        margin-left: 16rpx;
        padding: 16rpx 0;
        display: flex;
        flex-direction: column;
    }

    /* 品牌横向滚动 */
    .brand-section {
        background: #FFFFFF;
        border-radius: 16rpx;
        padding: 16rpx;
        margin-bottom: 16rpx;
        box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);
    }

    .brand-scroll {
        white-space: nowrap;
        display: flex;
        align-items: center;
    }

    .brand-chip {
        display: inline-flex;
        align-items: center;
        gap: 8rpx;
        padding: 10rpx 20rpx;
        background: #F5F7F4;
        border-radius: 32rpx;
        margin-right: 12rpx;
        flex-shrink: 0;
        transition: all 0.2s;
        border: 2rpx solid transparent;
    }

    .brand-chip.active {
        background: #E8F8EC;
        border-color: #22A84F;
    }

    .brand-name {
        font-size: 24rpx;
        color: #666666;
        white-space: nowrap;
    }

    .brand-chip.active .brand-name {
        color: #22A84F;
        font-weight: 500;
    }

    /* 商品滚动区 */
    .goods-scroll {
        flex: 1;
        background: #FFFFFF;
        border-radius: 16rpx;
        box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);
    }

    .load-more {
        padding: 20rpx 0;
    }
</style>
