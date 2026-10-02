<template>
    <pageWrapper>
        <scroll-view scroll-y class="page-scroll">
            <view class="section-card" v-if="user && user.auditStatus === 0">
                <view class="address-row" v-if="shipmentType === '0'" @click="addressShow">
                    <view class="address-icon">
                        <u-icon name="map" color="#22A84F" size="40rpx"></u-icon>
                    </view>
                    <view class="address-content" v-if="addressId > 0">
                        <view class="address-main">
                            <text class="address-name">{{ address.name }}</text>
                            <text class="address-sex">{{ address.sex === 1 ? '先生' : '女士' }}</text>
                            <text class="address-phone">{{ address.mobile }}</text>
                        </view>
                        <view class="address-detail">{{ address.title }}{{ address.detail }}</view>
                    </view>
                    <view class="address-empty" v-else>
                        <text>请选择收货地址</text>
                    </view>
                    <u-icon name="arrow-right" color="#CCCCCC" size="32rpx"></u-icon>
                </view>
                <view class="pickup-row" v-if="shipmentType === '1'">
                    <view class="pickup-icon">
                        <u-icon name="bag" color="#22A84F" size="40rpx"></u-icon>
                    </view>
                    <view class="pickup-info">
                        <text class="pickup-label">到店自提</text>
                        <text class="pickup-tip">请到店出示订单编号取货</text>
                    </view>
                </view>
                <view class="shipment-tabs">
                    <view class="tab-item" :class="{ active: shipmentType === '0' }" @click="shipmentTypeChange('0')">
                        <u-icon name="car" size="32rpx" :color="shipmentType === '0' ? '#22A84F' : '#999999'"></u-icon>
                        <text>配送</text>
                    </view>
                    <view class="tab-item" :class="{ active: shipmentType === '1' }" @click="shipmentTypeChange('1')">
                        <u-icon name="bag" size="32rpx" :color="shipmentType === '1' ? '#22A84F' : '#999999'"></u-icon>
                        <text>自提</text>
                    </view>
                </view>
            </view>
            <view class="section-title">订单详情</view>
            <view class="section-card goods-card">
                <view class="goods-item" v-for="cart in list" :key="cart.ID">
                    <image class="goods-img" :src="getGoodsImage(cart)" mode="aspectFill"></image>
                    <view class="goods-info">
                        <text class="goods-name">{{ cart.goods && cart.goods.name }}</text>
                        <text class="goods-spec">规格：{{ cart.goods && cart.goods.unit }}</text>
                        <view class="goods-bottom">
                            <text class="goods-price">
                                <text class="sym">¥</text>{{ getGoodsPrice(cart) }}
                            </text>
                            <text class="goods-num">x{{ cart.num }}</text>
                        </view>
                    </view>
                </view>
            </view>
            <view class="section-card remark-card">
                <text class="remark-label">备注</text>
                <input class="remark-input" v-model="remark" placeholder="选填，可备注特殊需求" />
            </view>
            <view class="bottom-safe"></view>
        </scroll-view>
        <view class="bottom-bar">
            <view class="total-info">
                <text class="total-label">合计：</text>
                <text class="total-price">
                    <text class="sym">¥</text>{{ total }}
                </text>
            </view>
            <view class="submit-btn" @click="submit">
                <text>提交订单</text>
            </view>
        </view>
        <addressPop :show="showLoginDialog" @close="addressClose" :addressId="addressId" @checked="addressChecked"></addressPop>
        <u-modal :show="showSettlmentUnpaid" showCancelButton closeOnClickOverlay @confirm="callPhone"
            @cancel="settlmentUnpaidCancel" confirmText="联系商家" cancelText="稍后处理" title="未结算订单提醒">
            <view class="modal-body">
                <view class="modal-main">您有{{ preSettlmentInfo.month }}月未结算的订单需要处理</view>
                <view class="modal-info">
                    <view>共 <text class="hl">{{ preSettlmentInfo.monthUnpaid }}</text> 个订单未结算</view>
                    <view>金额 <text class="hl">¥{{ preSettlmentInfo.settlementUnpaid }}</text></view>
                </view>
            </view>
        </u-modal>
        <u-toast ref="toast" style="z-index: 9999"></u-toast>
    </pageWrapper>
</template>

<script>
import { getCheckedCartList } from "@/api/cart"
import config from '@/config/config.js'
import { getToken, getRole, getSettlmentInfo, getUser } from '@/store/storage.js'
import { getDefaultAddressInfo } from '@/api/address'
import { createOrder } from '@/api/order'
import { getGoodsInfo } from '@/api/goods.js'
import addressPop from '@/components/addressPop/addressPop.vue'

export default {
    components: { addressPop },
    data() {
        return {
            token: '',
            role: {},
            user: {},
            list: [],
            pointGoodsId: 0,
            addressId: 0,
            address: {},
            total: 0,
            showSettlmentUnpaid: false,
            relationPhone: '',
            showLoginDialog: false,
            shipmentType: '0',
            remark: '',
            showPointPay: false,
            pointAmount: 0,
            preSettlmentInfo: {}
        }
    },
    onLoad(options) {
        this.relationPhone = config.phone
        if (options.pointGoodsId) {
            this.pointGoodsId = parseInt(options.pointGoodsId)
        }
        this.role = getRole()
        this.preSettlmentInfo = getSettlmentInfo()
        if (this.preSettlmentInfo && this.preSettlmentInfo.monthUnpaid > 0) {
            this.showSettlmentUnpaid = true
        }
    },
    mounted() {
        this.token = getToken()
        this.user = getUser()
        if (this.user && this.user.auditStatus === 1) {
            this.shipmentType = '1'
        }
        if (!this.token) {
            uni.showToast({ title: '请先登录', icon: 'none' })
            setTimeout(() => {
                uni.redirectTo({ url: '/pages/my/my' })
            }, 1500)
            return
        }
        if (this.pointGoodsId > 0) {
            this.getPointGoodsData()
        } else {
            this.getCartListData()
        }
        this.getAddressInfo()
    },
    methods: {
        getGoodsImage(cart) {
            if (cart.goods && cart.goods.images && cart.goods.images[0]) {
                let url = cart.goods.images[0].url
                if (url.slice(0, 4) !== 'http') {
                    url = config.baseUrl + '/' + url
                }
                return url
            }
            return '/static/nopicture.jpg'
        },
        getGoodsPrice(cart) {
            if (!cart.goods) return '0.00'
            const price = cart.goods.price > 0 && cart.goods.price < cart.goods.costPrice
                ? cart.goods.price : cart.goods.costPrice || 0
            return price.toFixed(2)
        },
        async submit() {
            const user = getUser()
            if (user && user.auditStatus === 0) {
                uni.showToast({ title: '您的账号正在审核中，暂无法下单', icon: 'none' })
                return
            }
            const data = {
                remarks: this.remark,
                addressId: this.addressId,
                shipmentType: parseInt(this.shipmentType)
            }
            if (this.pointGoodsId) {
                data.pointGoodsId = this.pointGoodsId
            }
            const res = await createOrder(data)
            if (res.code !== 0) {
                uni.showToast({ title: res.msg || '提交失败', icon: 'none' })
                return
            }
            if (this.pointGoodsId > 0) {
                uni.showToast({ title: '兑换成功', icon: 'success' })
                uni.redirectTo({ url: '/pages/order/detail?id=' + res.data.order.ID })
                return
            }
            uni.showToast({ title: '订单已提交', icon: 'success' })
            uni.redirectTo({ url: '/pages/order/detail?id=' + res.data.order.ID })
        },
        addressChecked(info) {
            this.address = info
            this.addressId = info.ID
            this.addressClose()
        },
        addressShow() {
            this.showLoginDialog = true
        },
        addressClose() {
            this.showLoginDialog = false
        },
        shipmentTypeChange(type) {
            this.shipmentType = type
            if (type === '1') {
                this.address = {}
                this.addressId = 0
            } else {
                this.getAddressInfo()
            }
        },
        async getPointGoodsData() {
            const res = await getGoodsInfo({ ID: this.pointGoodsId })
            if (res.code !== 0) return
            const g = res.data.regoods
            if (g.images && g.images[0]) {
                g.images.forEach((item, index) => {
                    if (item.url && item.url.slice(0, 4) !== 'http') {
                        g.images[index].url = config.baseUrl + '/' + item.url
                    }
                })
            }
            this.list.push({ goodsId: g.ID, specType: 0, num: 1, goods: g })
            this.total = (g.costPrice || 0).toFixed(2)
        },
        async getCartListData() {
            const res = await getCheckedCartList()
            if (res.code !== 0) {
                uni.showToast({ title: res.msg || '加载失败', icon: 'none' })
                return
            }
            if (!res.data.list || res.data.list.length === 0) {
                uni.showToast({ title: '无商品数据', icon: 'none' })
                setTimeout(() => {
                    uni.redirectTo({ url: '/pages/cart/cart' })
                }, 1500)
                return
            }
            let total = 0
            res.data.list.forEach(item => {
                if (item.goods && item.goods.images && item.goods.images[0]) {
                    let url = item.goods.images[0].url
                    if (url.slice(0, 4) !== 'http') {
                        item.goods.images[0].url = config.baseUrl + '/' + url
                    }
                }
                const price = item.goods && item.goods.price > 0 && item.goods.price < item.goods.costPrice
                    ? item.goods.price : (item.goods ? item.goods.costPrice : 0) || 0
                total += price * item.num
            })
            this.total = total.toFixed(2)
            this.list = res.data.list
        },
        getAddressInfo() {
            getDefaultAddressInfo().then(res => {
                if (res.data) {
                    this.address = res.data
                    this.addressId = res.data.ID
                }
            })
        },
        callPhone() {
            uni.makePhoneCall({ phoneNumber: this.relationPhone })
        },
        settlmentUnpaidCancel() {
            this.showSettlmentUnpaid = false
            uni.navigateBack({ delta: 1 })
        }
    }
}
</script>

<style lang="scss" scoped>
.page-scroll {
    height: calc(100vh - 120rpx);
    background: #F5F7F4;
}
.section-title {
    font-size: 28rpx;
    font-weight: 600;
    color: #1A1A1A;
    padding: 24rpx 24rpx 12rpx;
}
.section-card {
    margin: 0 24rpx 20rpx;
    background: #FFFFFF;
    border-radius: 20rpx;
}
.address-row, .pickup-row {
    display: flex;
    align-items: center;
    padding: 28rpx 24rpx;
    gap: 20rpx;
}
.address-icon, .pickup-icon {
    width: 72rpx;
    height: 72rpx;
    background: #E8F8EC;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
}
.address-content { flex: 1; }
.address-main {
    display: flex;
    align-items: center;
    gap: 8rpx;
    margin-bottom: 8rpx;
}
.address-name { font-size: 30rpx; font-weight: 600; color: #1A1A1A; }
.address-sex, .address-phone { font-size: 26rpx; color: #666666; }
.address-detail { font-size: 24rpx; color: #999999; }
.address-empty { flex: 1; font-size: 28rpx; color: #999999; }
.pickup-info { display: flex; flex-direction: column; gap: 8rpx; }
.pickup-label { font-size: 30rpx; font-weight: 600; color: #1A1A1A; }
.pickup-tip { font-size: 24rpx; color: #999999; }
.shipment-tabs { display: flex; border-top: 1rpx solid #EEEEEE; }
.tab-item {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8rpx;
    padding: 24rpx 0;
    font-size: 28rpx;
    color: #999999;
    border-top: 4rpx solid transparent;
}
.tab-item.active { color: #22A84F; border-top-color: #22A84F; }
.goods-card { padding: 0; }
.goods-item {
    display: flex;
    padding: 24rpx;
    gap: 20rpx;
    border-bottom: 1rpx solid #EEEEEE;
}
.goods-item:last-child { border-bottom: none; }
.goods-img { width: 160rpx; height: 160rpx; border-radius: 16rpx; flex-shrink: 0; }
.goods-info { flex: 1; display: flex; flex-direction: column; justify-content: space-between; }
.goods-name { font-size: 28rpx; font-weight: 500; color: #1A1A1A; }
.goods-spec { font-size: 24rpx; color: #999999; margin-top: 8rpx; }
.goods-bottom { display: flex; align-items: center; justify-content: space-between; margin-top: 8rpx; }
.goods-price { color: #F97316; font-size: 32rpx; font-weight: 700; }
.sym { font-size: 24rpx; font-weight: 600; }
.goods-num { font-size: 26rpx; color: #666666; }
.remark-card { display: flex; align-items: center; padding: 24rpx; gap: 16rpx; }
.remark-label { font-size: 28rpx; color: #1A1A1A; font-weight: 500; flex-shrink: 0; }
.remark-input { flex: 1; font-size: 28rpx; color: #1A1A1A; }
.bottom-safe { height: 140rpx; }
.bottom-bar {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    height: 120rpx;
    background: #FFFFFF;
    display: flex;
    align-items: center;
    padding: 0 24rpx;
    box-shadow: 0 -2rpx 12rpx rgba(0, 0, 0, 0.06);
    z-index: 100;
}
.total-info { flex: 1; display: flex; align-items: baseline; }
.total-label { font-size: 28rpx; color: #666666; }
.total-price { font-size: 40rpx; font-weight: 700; color: #F97316; }
.submit-btn {
    background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
    color: #FFFFFF;
    font-size: 30rpx;
    font-weight: 600;
    padding: 24rpx 60rpx;
    border-radius: 40rpx;
}
.modal-body { padding: 32rpx; }
.modal-main {
    font-size: 32rpx;
    font-weight: 600;
    color: #1A1A1A;
    text-align: center;
    margin-bottom: 24rpx;
}
.modal-info {
    background: #F5F7F4;
    border-radius: 12rpx;
    padding: 24rpx;
    font-size: 28rpx;
    color: #666666;
    line-height: 1.8;
}
.hl { color: #22A84F; font-weight: 600; margin: 0 6rpx; }
</style>
