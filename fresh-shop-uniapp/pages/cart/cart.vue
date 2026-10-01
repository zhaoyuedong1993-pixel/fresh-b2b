<!--
 * 优诚配运 - 购物车
 * 设计规范：自然清新
-->
<template>
    <pageWrapper>
        <!-- 未登录 -->
        <view class="not-login" v-if="!token">
            <view class="empty-card">
                <view class="empty-icon">
                    <u-icon name="shopping-cart" size="120rpx" color="#CCCCCC"></u-icon>
                </view>
                <view class="empty-title">购物车是空的</view>
                <view class="empty-desc">登录后查看您的购物车商品</view>
                <view class="empty-btn" @click="showLogin">立即登录</view>
            </view>
        </view>

        <!-- 已登录 -->
        <view class="cart-wrap" v-else>
            <!-- 页面标题 -->
            <view class="page-header">
                <text class="page-title">购物车</text>
                <view class="header-right" v-if="list.length > 0" @click="toggleEdit">
                    <text>{{ isEdit ? '完成' : '编辑' }}</text>
                </view>
            </view>

            <!-- 购物车内容 -->
            <scroll-view class="cart-scroll" scroll-y :style="{ height: scrollHeight + 'px' }" refresher-enabled
                :refresher-triggered="isRefreshing" @refresherrefresh="onRefresh">
                <shopCart :list="list" :isEdit="isEdit" @onRefresh="onRefresh" @delect="delectCart"
                    @update="updateCart" @accounts="accounts" @deleteCart="deleteCartByIndex" />
                <view class="bottom-safe" v-if="list.length === 0">
                    <view class="empty-card">
                        <view class="empty-icon">
                            <u-icon name="shopping-cart" size="120rpx" color="#CCCCCC"></u-icon>
                        </view>
                        <view class="empty-title">购物车是空的</view>
                        <view class="empty-desc">快去选购心仪商品吧</view>
                        <view class="empty-btn" @click="goHome">去选购</view>
                    </view>
                </view>
            </scroll-view>
        </view>

        <Tabbar :tabsId="3" />
        <loginPop :show="showLoginDialog" @close="hideLogin" @success="loginSuccess" />
        <u-toast ref="toast" style="z-index: 9999"></u-toast>
    </pageWrapper>
</template>

<script>
    import Tabbar from '@/components/tabbar/tabbar.vue'
    import shopCart from '@/components/shopCart/shopCart.vue'
    import loginPop from '@/components/loginPop/loginPop.vue'
    import { getToken } from '@/store/storage.js'
    import { getCartList } from '@/api/cart'
    import config from '@/config/config.js'

    export default {
        components: { Tabbar, shopCart, loginPop },
        data() {
            return {
                token: '',
                list: [],
                showLoginDialog: false,
                scrollHeight: 600,
                isRefreshing: false,
                isEdit: false
            }
        },
        onLoad() {
            this.token = getToken()
        },
        onShow() {
            if (this.token) {
                this.getCartListData()
            }
        },
        mounted() {
            uni.getSystemInfo({
                success: (res) => {
                    this.scrollHeight = res.windowHeight - 120
                }
            })
        },
        methods: {
            async getCartListData() {
                const res = await getCartList(this.$refs.toast)
                if (res.code === 401) {
                    this.token = ''
                    return
                }
                res.data.list?.forEach(item => {
                    if (item.goods?.images?.[0] && item.goods.images[0].url?.slice(0, 4) !== 'http') {
                        item.goods.images[0].url = config.baseUrl + '/' + item.goods.images[0].url
                    }
                })
                this.list = res.data.list || []
            },
            updateCart(list) {
                this.list = list
            },
            deleteCartByIndex(index) {
                this.list.splice(index, 1)
            },
            toggleEdit() {
                this.isEdit = !this.isEdit
            },
            async onRefresh() {
                this.isRefreshing = true
                await this.getCartListData()
                this.isRefreshing = false
                this.$message(this.$refs.toast).success('刷新成功')
            },
            delectCart(e) { },
            accounts(e) { },
            showLogin() {
                this.showLoginDialog = true
            },
            hideLogin() {
                this.showLoginDialog = false
            },
            loginSuccess() {
                this.hideLogin()
                this.token = getToken()
                this.$message(this.$refs.toast).success('登录成功')
                this.getCartListData()
            },
            goHome() {
                uni.switchTab({ url: '/pages/index/index' })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .not-login {
        min-height: 100vh;
        background: #F5F7F4;
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 40rpx;
    }

    .empty-card {
        background: #FFFFFF;
        border-radius: 24rpx;
        padding: 60rpx 40rpx;
        text-align: center;
        box-shadow: 0 4rpx 24rpx rgba(0, 0, 0, 0.06);
        width: 100%;
    }

    .empty-icon {
        margin-bottom: 24rpx;
    }

    .empty-title {
        font-size: 32rpx;
        font-weight: 600;
        color: #1A1A1A;
        margin-bottom: 12rpx;
    }

    .empty-desc {
        font-size: 26rpx;
        color: #999999;
        margin-bottom: 40rpx;
    }

    .empty-btn {
        display: inline-block;
        background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
        color: #FFFFFF;
        font-size: 28rpx;
        padding: 20rpx 60rpx;
        border-radius: 40rpx;
        box-shadow: 0 4rpx 16rpx rgba(34, 168, 79, 0.3);
    }

    .cart-wrap {
        height: 100vh;
        display: flex;
        flex-direction: column;
        background: #F5F7F4;
    }

    .page-header {
        height: 100rpx;
        padding: 0 32rpx;
        background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
        display: flex;
        align-items: center;
        justify-content: space-between;
        flex-shrink: 0;
    }

    .page-title {
        font-size: 36rpx;
        font-weight: 600;
        color: #FFFFFF;
    }

    .header-right text {
        font-size: 28rpx;
        color: #FFFFFF;
        opacity: 0.9;
    }

    .cart-scroll {
        flex: 1;
    }

    .bottom-safe {
        padding: 40rpx 24rpx;
    }
</style>
