<!--
 * 优诚配运 - 快速下单（近期购买 + 我的收藏）
 * 设计规范：自然清新
-->
<template>
    <pageWrapper>
        <!-- 顶部 Tab -->
        <view class="header-tabs">
            <view class="tab-bar">
                <view
                    class="tab-item"
                    :class="{ active: currentTab === 0 }"
                    @click="switchTab(0)"
                >
                    <text>近期购买</text>
                    <view class="tab-line" v-if="currentTab === 0"></view>
                </view>
                <view
                    class="tab-item"
                    :class="{ active: currentTab === 1 }"
                    @click="switchTab(1)"
                >
                    <text>我的收藏</text>
                    <view class="tab-line" v-if="currentTab === 1"></view>
                </view>
            </view>
        </view>

        <!-- 搜索栏 -->
        <view class="search-bar">
            <view class="search-box">
                <u-icon name="search" size="28rpx" color="#999999"></u-icon>
                <input
                    class="search-input"
                    v-model="keyword"
                    type="text"
                    confirm-type="search"
                    @confirm="searchConfirm"
                    placeholder="搜索商品名称"
                    placeholder-class="placeholder"
                />
            </view>
        </view>

        <!-- 未登录 -->
        <view class="login-tip" v-if="!token">
            <view class="login-icon">
                <u-icon name="account" size="80rpx" color="#CCCCCC"></u-icon>
            </view>
            <view class="login-text">登录后查看购买记录和收藏</view>
            <view class="login-btn" @click="showLogin">去登录</view>
        </view>

        <!-- 内容区 -->
        <view class="content-area" v-else>
            <!-- 近期购买 -->
            <scroll-view
                v-if="currentTab === 0 && payGoodslist.length > 0"
                class="goods-scroll"
                scroll-y
                :style="{ height: scrollViewHeight + 'px' }"
                refresher-enabled
                :refresher-triggered="triggered"
                @refresherrefresh="onRefresh"
                @scrolltolower="scrollTolower"
            >
                <GoodsList
                    :vertical="true"
                    :lists="payGoodslist"
                    price-type="￥"
                    :is-audit="isAudit"
                    :is-add-cart="true"
                    imgWidth="160rpx"
                    imgHeight="160rpx"
                    :show-pay-count="true"
                    @updateCart="updatePayGoodsCart"
                ></GoodsList>

                <view class="load-more">
                    <u-loadmore :status="payLoadMore" nomore-text="没有更多了" />
                </view>
            </scroll-view>

            <!-- 我的收藏 -->
            <scroll-view
                v-if="currentTab === 1 && favoritesList.length > 0"
                class="goods-scroll"
                scroll-y
                :style="{ height: scrollViewHeight + 'px' }"
                refresher-enabled
                :refresher-triggered="triggered"
                @refresherrefresh="onRefresh"
                @scrolltolower="scrollTolower"
            >
                <GoodsList
                    :vertical="true"
                    :lists="favoritesList"
                    price-type="￥"
                    :is-audit="isAudit"
                    :is-add-cart="true"
                    imgWidth="160rpx"
                    imgHeight="160rpx"
                    @onGoodsLongClick="onGoodsLongClick"
                    @updateCart="updateFavGoodsCart"
                ></GoodsList>

                <view class="load-more">
                    <u-loadmore :status="favLoadMore" nomore-text="没有更多了" />
                </view>
            </scroll-view>

            <!-- 空状态 -->
            <view class="empty-wrap" v-else>
                <view class="empty-icon">
                    <u-icon :name="currentTab === 0 ? 'clock' : 'heart'" size="120rpx" color="#CCCCCC"></u-icon>
                </view>
                <view class="empty-text">{{ currentTab === 0 ? '暂无近期购买记录' : '暂无收藏商品' }}</view>
                <view class="empty-btn" @click="goShopping">去选购</view>
            </view>
        </view>

        <!-- 删除收藏确认 -->
        <u-modal
            :show="showDeleteFav"
            :showCancelButton="true"
            title="取消收藏"
            content="确定取消收藏该商品吗？"
            @confirm="confirmDeleteFav"
            @cancel="showDeleteFav = false"
        ></u-modal>

        <!-- 登录弹窗 -->
        <loginPop :show="showLoginDialog" @close="hideLogin" @success="loginSuccess"></loginPop>

        <!-- 底部导航 -->
        <Tabbar :tabsId="2"></Tabbar>

        <u-toast ref="toast" style="z-index: 9999"></u-toast>
    </pageWrapper>
</template>

<script>
    import Tabbar from '@/components/tabbar/tabbar.vue'
    import loginPop from '@/components/loginPop/loginPop.vue'
    import { getRecentlyPurchasedGoodsList, getRecentlyPurchasedGoodsListLoading } from '@/api/order.js'
    import { getFavoritesListPage, getFavoritesListPageLoding, favorites } from '@/api/favorites.js'
    import { deleteCartByIds } from '@/api/cart.js'
    import GoodsList from '@/components/goodsList/goodsList.vue'
    import { getUser, getToken, setUser } from '@/store/storage.js'
    import { getUserAuditStatus } from '@/api/user.js'
    import config from '@/config/config.js'

    export default {
        components: {
            Tabbar,
            loginPop,
            GoodsList
        },
        data() {
            return {
                token: '',
                isAudit: false,
                currentTab: 0,
                keyword: '',
                scrollViewHeight: 600,
                triggered: false,
                payGoodslist: [],
                favoritesList: [],
                payPage: { page: 1, pageSize: 10, total: 0, isMore: true },
                favPage: { page: 1, pageSize: 10, total: 0, isMore: true },
                payLoadMore: 'loadmore',
                favLoadMore: 'loadmore',
                showLoginDialog: false,
                showDeleteFav: false,
                currentDeleteFavId: null
            }
        },
        onLoad() {
            this.token = getToken()
            if (this.token) {
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
            }
            this.calculateHeight()
        },
        onShow() {
            if (this.token) {
                this.loadData()
            }
        },
        methods: {
            calculateHeight() {
                uni.getSystemInfo({
                    success: (res) => {
                        this.scrollViewHeight = res.windowHeight - 180
                    }
                })
            },
            switchTab(index) {
                this.currentTab = index
                if (this.keyword) this.keyword = ''
                if (index === 0 && this.payGoodslist.length === 0) {
                    this.getRecentlyGoodsList(0)
                } else if (index === 1 && this.favoritesList.length === 0) {
                    this.getFavoritesGoodsList(0)
                }
            },
            loadData() {
                if (this.currentTab === 0) {
                    this.getRecentlyGoodsList(0)
                } else {
                    this.getFavoritesGoodsList(0)
                }
            },
            async getRecentlyGoodsList(type = 0) {
                if (type === 0) {
                    this.payPage = { page: 1, pageSize: 10, total: 0, isMore: true }
                    this.payLoadMore = 'loadmore'
                }

                const data = { page: this.payPage.page, pageSize: this.payPage.pageSize }
                if (this.keyword) data.name = this.keyword

                try {
                    const res = type === 0
                        ? await getRecentlyPurchasedGoodsListLoading(data)
                        : await getRecentlyPurchasedGoodsList(data)

                    if (res.code !== 0) return false

                    const list = res.data?.list || []
                    list.forEach(item => {
                        if (item.images?.[0]?.url && item.images[0].url.slice(0, 4) !== 'http') {
                            item.images[0].url = config.baseUrl + '/' + item.images[0].url
                        }
                        if (item.weight) {
                            item.weight = item.weight > 1000 ? item.weight / 1000 + 'kg' : item.weight + 'g'
                        }
                        if (item.cartNum == null) item.cartNum = 0
                    })

                    if (type === 0) {
                        this.payGoodslist = list
                    } else {
                        this.payGoodslist = [...this.payGoodslist, ...list]
                    }

                    this.payPage.page++
                    this.payPage.total = res.data?.total || 0
                    this.payPage.isMore = this.payPage.page * this.payPage.pageSize < this.payPage.total
                    this.payLoadMore = this.payPage.isMore ? 'loadmore' : 'nomore'
                    return true
                } catch (e) {
                    return false
                }
            },
            async getFavoritesGoodsList(type = 0) {
                if (type === 0) {
                    this.favPage = { page: 1, pageSize: 10, total: 0, isMore: true }
                    this.favLoadMore = 'loadmore'
                }

                const data = { page: this.favPage.page, pageSize: this.favPage.pageSize }
                if (this.keyword) data.name = this.keyword

                try {
                    const res = type === 0
                        ? await getFavoritesListPageLoding(data)
                        : await getFavoritesListPage(data)

                    if (res.code !== 0) return false

                    const list = res.data?.list || []
                    list.forEach(item => {
                        if (item.images?.[0]?.url && item.images[0].url.slice(0, 4) !== 'http') {
                            item.images[0].url = config.baseUrl + '/' + item.images[0].url
                        }
                        if (item.weight) {
                            item.weight = item.weight > 1000 ? item.weight / 1000 + 'kg' : item.weight + 'g'
                        }
                        if (item.cartNum == null) item.cartNum = 0
                    })

                    if (type === 0) {
                        this.favoritesList = list
                    } else {
                        this.favoritesList = [...this.favoritesList, ...list]
                    }

                    this.favPage.page++
                    this.favPage.total = res.data?.total || 0
                    this.favPage.isMore = this.favPage.page * this.favPage.pageSize < this.favPage.total
                    this.favLoadMore = this.favPage.isMore ? 'loadmore' : 'nomore'
                    return true
                } catch (e) {
                    return false
                }
            },
            async updatePayGoodsCart(index, cardId, num) {
                num = parseInt(num)
                let originNum = this.payGoodslist[index].cartNum || 0
                if (originNum > 0 && num < this.payGoodslist[index].minCount) {
                    this.$message(this.$refs.toast).error(`商品最低购买${this.payGoodslist[index].minCount}件`)
                    num = 0
                }
                if (originNum === 0 && num < this.payGoodslist[index].minCount) {
                    num = this.payGoodslist[index].minCount
                }
                if (originNum !== 0 && num === 0) {
                    await deleteCartByIds({ ids: [cardId] })
                    this.payGoodslist[index].cartNum = 0
                    return
                }
                this.payGoodslist[index].cartNum = num
                await this.addCartReq(this.payGoodslist[index], index, num, 'pay')
            },
            async updateFavGoodsCart(index, cardId, num) {
                num = parseInt(num)
                let originNum = this.favoritesList[index].cartNum || 0
                if (originNum > 0 && num < this.favoritesList[index].minCount) {
                    this.$message(this.$refs.toast).error(`商品最低购买${this.favoritesList[index].minCount}件`)
                    num = 0
                }
                if (originNum === 0 && num < this.favoritesList[index].minCount) {
                    num = this.favoritesList[index].minCount
                }
                if (originNum !== 0 && num === 0) {
                    await deleteCartByIds({ ids: [cardId] })
                    this.favoritesList[index].cartNum = 0
                    return
                }
                this.favoritesList[index].cartNum = num
                await this.addCartReq(this.favoritesList[index], index, num, 'fav')
            },
            async addCartReq(goodsInfo, index, num, type) {
                const res = await uni.request({
                    url: `${config.baseUrl}/cart/addCart`,
                    method: 'POST',
                    data: { goodsId: goodsInfo.ID, specType: 0, num: num }
                })
                if (res.data?.code !== 0) {
                    this.$message(this.$refs.toast).error('操作失败')
                    if (type === 'pay') {
                        this.payGoodslist[index].cartNum = this.payGoodslist[index].cartNum || 0
                    } else {
                        this.favoritesList[index].cartNum = this.favoritesList[index].cartNum || 0
                    }
                }
            },
            onGoodsLongClick(goods) {
                this.currentDeleteFavId = goods.ID
                this.showDeleteFav = true
            },
            async confirmDeleteFav() {
                if (!this.currentDeleteFavId) return
                try {
                    const res = await favorites({ goodsId: this.currentDeleteFavId })
                    if (res.code === 0) {
                        this.$message(this.$refs.toast).success('已取消收藏')
                        this.getFavoritesGoodsList(0)
                    }
                } catch (e) {
                    this.$message(this.$refs.toast).error('操作失败')
                } finally {
                    this.showDeleteFav = false
                    this.currentDeleteFavId = null
                }
            },
            searchConfirm() {
                this.loadData()
            },
            async onRefresh() {
                this.triggered = true
                const b = await this.loadData()
                this.$message(this.$refs.toast).success(b ? '刷新成功' : '刷新失败')
                this.triggered = false
            },
            async scrollTolower() {
                if (this.currentTab === 0) {
                    if (this.payLoadMore === 'loading' || !this.payPage.isMore) return
                    this.payLoadMore = 'loading'
                    await this.getRecentlyGoodsList(1)
                    this.payLoadMore = this.payPage.isMore ? 'loadmore' : 'nomore'
                } else {
                    if (this.favLoadMore === 'loading' || !this.favPage.isMore) return
                    this.favLoadMore = 'loading'
                    await this.getFavoritesGoodsList(1)
                    this.favLoadMore = this.favPage.isMore ? 'loadmore' : 'nomore'
                }
            },
            showLogin() {
                this.showLoginDialog = true
            },
            hideLogin() {
                this.showLoginDialog = false
            },
            async loginSuccess() {
                this.hideLogin()
                this.token = getToken()
                this.loadData()
                this.$message(this.$refs.toast).success('登录成功')
            },
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
        position: sticky;
        top: 0;
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

    .search-bar {
        padding: 20rpx 24rpx;
        background: #FFFFFF;
    }

    .search-box {
        height: 64rpx;
        background: #F5F7F4;
        border-radius: 32rpx;
        display: flex;
        align-items: center;
        padding: 0 24rpx;
        gap: 12rpx;
    }

    .search-input {
        flex: 1;
        height: 100%;
        font-size: 26rpx;
        color: #1A1A1A;
    }

    .placeholder {
        color: #CCCCCC;
        font-size: 26rpx;
    }

    .login-tip {
        display: flex;
        flex-direction: column;
        align-items: center;
        padding-top: 200rpx;
        gap: 24rpx;
    }

    .login-icon {
        width: 160rpx;
        height: 160rpx;
        background: #F5F7F4;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .login-text {
        font-size: 28rpx;
        color: #999999;
    }

    .login-btn {
        background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
        color: #FFFFFF;
        font-size: 28rpx;
        padding: 20rpx 60rpx;
        border-radius: 40rpx;
        box-shadow: 0 4rpx 16rpx rgba(34, 168, 79, 0.3);
    }

    .content-area {
        padding: 24rpx;
    }

    .goods-scroll {
        padding-bottom: 24rpx;
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
</style>
