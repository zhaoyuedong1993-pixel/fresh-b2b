<!--
 * 优诚配运 - 商品详情
 * 设计规范：自然清新
-->
<template>
    <pageWrapper>
        <!-- 商品图片轮播 -->
        <view class="goods-swiper" v-if="goods.images && goods.images.length > 0">
            <swiper
                class="swiper"
                :indicator-dots="true"
                :autoplay="false"
                :current="currentImgIndex"
                @change="onSwiperChange"
                indicator-color="rgba(255,255,255,0.5)"
                indicator-active-color="#FFFFFF"
            >
                <swiper-item v-for="(img, index) in goods.images" :key="index" @click="previewImages(index)">
                    <image class="goods-image" :src="img.url" mode="aspectFill"></image>
                </swiper-item>
            </swiper>
            <view class="img-indicator">
                <text>{{ currentImgIndex + 1 }}/{{ goods.images.length }}</text>
            </view>
        </view>

        <!-- 商品基本信息 -->
        <view class="goods-info">
            <view class="price-row">
                <view class="price-box">
                    <text class="price-symbol">¥</text>
                    <text class="price-value">{{ formatPrice(displayPrice) }}</text>
                    <text class="price-unit">/{{ goods.unit || '份' }}</text>
                </view>
                <view class="origin-tag" v-if="goods.origin">
                    <u-icon name="position" size="24rpx" color="#22A84F"></u-icon>
                    <text>{{ goods.origin }}</text>
                </view>
            </view>

            <view class="title-row">
                <view class="brand-tag" v-if="goods.brand && goods.brand.name">
                    {{ goods.brand.name }}
                </view>
                <text class="goods-name">{{ goods.name }}</text>
            </view>

            <view class="goods-weight" v-if="goods.weight">
                <text>{{ goods.weight }}</text>
            </view>

            <!-- 库存状态 -->
            <view class="stock-row">
                <text class="stock-label">库存：</text>
                <text class="stock-value" :class="{ 'out': goods.store <= 0 }">
                    {{ goods.store > 100 ? '充足' : goods.store > 0 ? `剩余 ${goods.store}` : '缺货' }}
                </text>
            </view>
        </view>

        <!-- 商品详情 -->
        <view class="goods-detail">
            <view class="detail-header">
                <u-icon name="file-text" size="32rpx" color="#22A84F"></u-icon>
                <text>商品详情</text>
            </view>
            <view class="detail-content" v-if="goods.desc && goods.desc.details">
                <u-parse :content="goods.desc.details" :tagStyle="tagStyle"></u-parse>
            </view>
            <view class="detail-empty" v-else>
                <text>暂无商品详情</text>
            </view>
        </view>

        <!-- 底部操作栏 -->
        <view class="bottom-bar">
            <view class="action-icons">
                <view class="action-icon" @click="toggleFavorite">
                    <u-icon :name="goods.isFavorite ? 'heart-fill' : 'heart'"
                        :color="goods.isFavorite ? '#EF4444' : '#666666'" size="44rpx"></u-icon>
                    <text>收藏</text>
                </view>
                <view class="action-icon" @click="goCart">
                    <u-icon name="shopping-cart" color="#666666" size="44rpx"></u-icon>
                    <text>购物车</text>
                    <view class="cart-badge" v-if="goods.cartTotalNum > 0">
                        <text>{{ goods.cartTotalNum > 99 ? '99+' : goods.cartTotalNum }}</text>
                    </view>
                </view>
            </view>

            <view class="action-buttons">
                <!-- 积分商品 -->
                <template v-if="goods.goodsArea === 1">
                    <view
                        class="action-btn exchange"
                        :class="{ disabled: pointAmount < goods.costPrice || goods.store <= 0 }"
                        @click="exchangePointGoods"
                    >
                        <text v-if="goods.store <= 0">暂时无货</text>
                        <text v-else-if="pointAmount < goods.costPrice">积分不足</text>
                        <text v-else>{{ goods.costPrice }} 积分兑换</text>
                    </view>
                </template>

                <!-- 普通商品 -->
                <template v-else>
                    <view class="quantity-control" v-if="goods.cartNum > 0">
                        <view class="qty-btn minus" @click="updateCartNum(-1)">
                            <u-icon name="minus" size="28rpx" color="#22A84F"></u-icon>
                        </view>
                        <text class="qty-num">{{ goods.cartNum }}</text>
                        <view class="qty-btn plus" @click="updateCartNum(1)">
                            <u-icon name="plus" size="28rpx" color="#FFFFFF"></u-icon>
                        </view>
                    </view>

                    <view
                        class="action-btn add-cart"
                        :class="{ disabled: goods.store <= 0 }"
                        @click="handleAddCart"
                    >
                        <text v-if="goods.store <= 0">暂时无货</text>
                        <text v-else-if="goods.minCount > 1">最低 {{ goods.minCount }} 件起购</text>
                        <text v-else>加入购物车</text>
                    </view>

                    <view
                        class="action-btn buy-now"
                        :class="{ disabled: goods.store <= 0 }"
                        @click="handleBuyNow"
                    >
                        <text>立即下单</text>
                    </view>
                </template>
            </view>
        </view>

        <!-- 底部占位 -->
        <view class="bottom-placeholder"></view>

        <!-- 登录悬浮框 -->
        <loginSuspend :show="loginSuspendShow" @success="loginSuccess"></loginSuspend>

        <u-toast ref="toast" style="z-index: 9998;"></u-toast>
    </pageWrapper>
</template>

<script>
    import { getToken, getUser, setUser } from '@/store/storage.js'
    import loginSuspend from '@/components/loginPop/loginSuspend.vue'
    import { getGoodsInfo } from '@/api/goods.js'
    import { favorites } from '@/api/favorites.js'
    import { addCart, selectGoodsSingeChecked } from '@/api/cart.js'
    import { getAccountInfo } from '@/api/account.js'
    import { getUserAuditStatus } from '@/api/user.js'
    import config from '@/config/config.js'

    export default {
        components: {
            loginSuspend
        },
        data() {
            return {
                id: 0,
                goods: {},
                token: '',
                isAudit: false,
                loginSuspendShow: false,
                currentImgIndex: 0,
                pointAmount: 0,
                tagStyle: {
                    img: 'width: 100%; vertical-align: bottom;'
                }
            }
        },
        computed: {
            displayPrice() {
                if (this.isAudit) {
                    if (this.goods.goodsArea === 1) {
                        return this.goods.costPrice
                    }
                    return this.goods.price > 0 && this.goods.price < this.goods.costPrice
                        ? this.goods.price
                        : this.goods.costPrice
                }
                return this.goods.costPrice || 0
            }
        },
        onLoad(options) {
            if (!options.id) {
                uni.showToast({ title: '参数错误', icon: 'error' })
                setTimeout(() => uni.navigateBack(), 1500)
                return
            }

            this.id = parseInt(options.id)
            this.token = getToken()

            if (!this.token) {
                this.loginSuspendShow = true
            }

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

            this.loadGoods()
        },
        methods: {
            async loadGoods() {
                uni.showLoading({ title: '加载中...' })
                try {
                    const res = await getGoodsInfo({ ID: this.id })
                    if (res.code !== 0) return

                    const regoods = res.data.regoods || {}

                    // 处理图片URL
                    if (regoods.images && regoods.images.length > 0) {
                        regoods.images.forEach((item, index) => {
                            if (item.url && item.url.slice(0, 4) !== 'http') {
                                regoods.images[index].url = config.baseUrl + '/' + item.url
                            }
                        })
                    } else {
                        regoods.images = [{ url: '/static/nopicture.jpg' }]
                    }

                    // 处理重量
                    if (regoods.weight) {
                        regoods.weight = regoods.weight > 1000
                            ? (regoods.weight / 1000) + 'kg'
                            : regoods.weight + 'g'
                    }

                    this.goods = regoods

                    // 积分商品获取积分余额
                    if (regoods.goodsArea === 1) {
                        const accRes = await getAccountInfo(2)
                        if (accRes.code === 0) {
                            this.pointAmount = accRes.data?.account?.amount || 0
                        }
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
            onSwiperChange(e) {
                this.currentImgIndex = e.detail.current
            },
            previewImages(index) {
                const urls = this.goods.images.map(img => img.url)
                uni.previewImage({
                    current: index,
                    urls: urls
                })
            },
            async toggleFavorite() {
                if (!this.token) {
                    this.loginSuspendShow = true
                    return
                }

                try {
                    const res = await favorites({ goodsId: this.id })
                    if (res.code === 0) {
                        this.goods.isFavorite = !this.goods.isFavorite
                        this.$message(this.$refs.toast).success(
                            this.goods.isFavorite ? '收藏成功' : '取消收藏'
                        )
                    }
                } catch (e) {
                    this.$message(this.$refs.toast).error('操作失败')
                }
            },
            goCart() {
                uni.navigateTo({ url: '/pages/cart/cart' })
            },
            handleAddCart() {
                if (!this.isAudit) {
                    this.$message(this.$refs.toast).warning('审核通过后才可以下单！').then(() => {
                        uni.navigateTo({ url: '/pages/my/memberInfo' })
                    })
                    return
                }
                if (this.goods.store <= 0) return

                let num = this.goods.minCount > 0 ? this.goods.minCount : 1
                this.addCartReq(num)
            },
            async updateCartNum(delta) {
                let newNum = (this.goods.cartNum || 0) + delta
                if (newNum < 0) return
                if (newNum < this.goods.minCount && newNum !== 0) {
                    newNum = this.goods.minCount
                }
                if (newNum === 0) {
                    newNum = 0
                }
                await this.addCartReq(newNum)
            },
            async addCartReq(num) {
                const res = await addCart({
                    goodsId: this.id,
                    specType: 0,
                    num: num
                })
                if (res.code === 0) {
                    this.goods.cartNum = num
                    this.goods.cartTotalNum = (this.goods.cartTotalNum || 0) + (num - (this.goods.cartNum || 0))
                    this.$message(this.$refs.toast).success('添加成功')
                }
            },
            async handleBuyNow() {
                if (!this.isAudit) {
                    this.$message(this.$refs.toast).warning('审核通过后才可以下单！')
                    return
                }
                if (this.goods.store <= 0) return

                const num = this.goods.cartNum > 0 ? this.goods.cartNum : 1
                const res = await selectGoodsSingeChecked({
                    goodsId: this.id,
                    specType: 0,
                    num: num
                })
                if (res.code === 0) {
                    uni.navigateTo({ url: '/pages/order/submit' })
                }
            },
            exchangePointGoods() {
                if (!this.isAudit) {
                    this.$message(this.$refs.toast).warning('审核通过后才可以兑换！')
                    return
                }
                if (this.pointAmount < this.goods.costPrice) {
                    this.$message(this.$refs.toast).error('积分不足')
                    return
                }
                if (this.goods.store <= 0) {
                    this.$message(this.$refs.toast).error('暂时无货')
                    return
                }
                uni.navigateTo({ url: `/pages/order/submit?pointGoodsId=${this.id}` })
            },
            loginSuccess() {
                this.loginSuspendShow = false
                this.token = getToken()
                this.loadGoods()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .goods-swiper {
        position: relative;
        width: 100%;
        height: 600rpx;
        background: #F5F7F4;
    }

    .swiper {
        width: 100%;
        height: 100%;
    }

    .goods-image {
        width: 100%;
        height: 100%;
    }

    .img-indicator {
        position: absolute;
        bottom: 24rpx;
        right: 24rpx;
        background: rgba(0, 0, 0, 0.4);
        padding: 8rpx 20rpx;
        border-radius: 24rpx;

        text {
            font-size: 24rpx;
            color: #FFFFFF;
        }
    }

    .goods-info {
        background: #FFFFFF;
        padding: 28rpx;
        margin-bottom: 16rpx;
    }

    .price-row {
        display: flex;
        align-items: baseline;
        justify-content: space-between;
        margin-bottom: 16rpx;
    }

    .price-box {
        display: flex;
        align-items: baseline;
    }

    .price-symbol {
        font-size: 28rpx;
        font-weight: 600;
        color: #F97316;
    }

    .price-value {
        font-size: 52rpx;
        font-weight: 700;
        color: #F97316;
        line-height: 1;
    }

    .price-unit {
        font-size: 26rpx;
        color: #999999;
        margin-left: 8rpx;
    }

    .origin-tag {
        display: flex;
        align-items: center;
        gap: 6rpx;
        background: #E8F8EC;
        padding: 8rpx 16rpx;
        border-radius: 20rpx;

        text {
            font-size: 22rpx;
            color: #22A84F;
        }
    }

    .title-row {
        display: flex;
        align-items: center;
        flex-wrap: wrap;
        gap: 12rpx;
        margin-bottom: 12rpx;
    }

    .brand-tag {
        background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
        color: #FFFFFF;
        font-size: 22rpx;
        padding: 4rpx 12rpx;
        border-radius: 8rpx;
    }

    .goods-name {
        font-size: 32rpx;
        font-weight: 600;
        color: #1A1A1A;
        line-height: 1.4;
    }

    .goods-weight {
        font-size: 24rpx;
        color: #999999;
        margin-bottom: 12rpx;
    }

    .stock-row {
        display: flex;
        align-items: center;
        font-size: 26rpx;
    }

    .stock-label {
        color: #999999;
    }

    .stock-value {
        color: #22A84F;

        &.out {
            color: #EF4444;
        }
    }

    .goods-detail {
        background: #FFFFFF;
        margin-bottom: 16rpx;
    }

    .detail-header {
        display: flex;
        align-items: center;
        gap: 12rpx;
        padding: 28rpx;
        border-bottom: 1rpx solid #EEEEEE;

        text {
            font-size: 30rpx;
            font-weight: 600;
            color: #1A1A1A;
        }
    }

    .detail-content {
        padding: 24rpx;
    }

    .detail-empty {
        padding: 48rpx;
        text-align: center;
        color: #999999;
        font-size: 28rpx;
    }

    .bottom-bar {
        position: fixed;
        bottom: 0;
        left: 0;
        right: 0;
        background: #FFFFFF;
        box-shadow: 0 -2rpx 20rpx rgba(0, 0, 0, 0.06);
        display: flex;
        align-items: center;
        padding: 16rpx 24rpx;
        padding-bottom: calc(env(safe-area-inset-bottom) + 16rpx);
        z-index: 100;
    }

    .action-icons {
        display: flex;
        gap: 32rpx;
    }

    .action-icon {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 4rpx;
        position: relative;

        text {
            font-size: 20rpx;
            color: #666666;
        }
    }

    .cart-badge {
        position: absolute;
        top: -8rpx;
        right: -12rpx;
        min-width: 32rpx;
        height: 32rpx;
        background: #EF4444;
        border-radius: 16rpx;
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 0 6rpx;

        text {
            font-size: 20rpx;
            color: #FFFFFF;
        }
    }

    .action-buttons {
        flex: 1;
        display: flex;
        align-items: center;
        justify-content: flex-end;
        gap: 16rpx;
        margin-left: 24rpx;
    }

    .quantity-control {
        display: flex;
        align-items: center;
        gap: 16rpx;
    }

    .qty-btn {
        width: 56rpx;
        height: 56rpx;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;

        &.minus {
            background: #E8F8EC;
            border: 2rpx solid #22A84F;
        }

        &.plus {
            background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
        }
    }

    .qty-num {
        font-size: 32rpx;
        font-weight: 600;
        color: #1A1A1A;
        min-width: 48rpx;
        text-align: center;
    }

    .action-btn {
        height: 72rpx;
        padding: 0 32rpx;
        border-radius: 36rpx;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 28rpx;
        font-weight: 600;

        &.disabled {
            opacity: 0.5;
        }

        &.add-cart {
            background: #FEF3E2;
            color: #F97316;
        }

        &.buy-now {
            background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
            color: #FFFFFF;
            box-shadow: 0 4rpx 16rpx rgba(34, 168, 79, 0.3);
        }

        &.exchange {
            background: linear-gradient(135deg, #F97316 0%, #EA580C 100%);
            color: #FFFFFF;
            box-shadow: 0 4rpx 16rpx rgba(249, 115, 22, 0.3);
        }
    }

    .bottom-placeholder {
        height: 140rpx;
    }
</style>
