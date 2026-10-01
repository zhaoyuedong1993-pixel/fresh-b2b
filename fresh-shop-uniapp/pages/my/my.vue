<!--
 * 优诚配运 - 个人中心
 * 设计规范：自然清新
-->
<template>
    <pageWrapper>
        <!-- 头部用户信息 -->
        <view class="profile-header">
            <!-- 已登录 -->
            <view class="user-card" v-if="token && user" @click="toMemberInfo">
                <view class="user-avatar">
                    <image :src="user.headerImg || '/static/my/face.png'" mode="aspectFill"></image>
                </view>
                <view class="user-info">
                    <text class="user-name">{{ user.nickName || '用户' }}</text>
                    <view class="user-tags">
                        <view class="tag" v-if="user.auditStatus === 1">
                            <text>{{ role.authorityName ? role.authorityName.replace('客户', '') : '' }}</text>
                        </view>
                        <view class="tag warn" v-else-if="[2, 3].includes(user.auditStatus)">
                            <text>等待审核</text>
                        </view>
                        <view class="tag warn" v-else-if="user.auditStatus === 4">
                            <text>审核不通过</text>
                        </view>
                        <view class="tag" v-else>
                            <text>未填写信息</text>
                        </view>
                        <view class="tag point">
                            <text>积分 {{ point }}</text>
                        </view>
                    </view>
                </view>
                <u-icon name="arrow-right" color="#CCCCCC" size="32rpx"></u-icon>
            </view>

            <!-- 未登录 -->
            <view class="user-card not-login" v-else @click="showLogin">
                <view class="user-avatar">
                    <image src="/static/my/face.png" mode="aspectFill"></image>
                </view>
                <view class="user-info">
                    <text class="user-name login-tip">点击登录</text>
                </view>
                <u-icon name="arrow-right" color="#CCCCCC" size="32rpx"></u-icon>
            </view>

            <!-- 月结信息 -->
            <view class="monthly-info" v-if="role.authorityId === 1001 && orderStatusCount.monthUnpaid > 0">
                <view class="monthly-item">
                    <text class="label">{{ orderStatusCount.month }}月未结</text>
                    <text class="value">{{ orderStatusCount.monthUnpaid }} 单</text>
                </view>
                <view class="divider"></view>
                <view class="monthly-item">
                    <text class="label">未结金额</text>
                    <text class="value primary">¥{{ orderStatusCount.settlementUnpaid }}</text>
                </view>
            </view>
        </view>

        <!-- 订单入口 -->
        <view class="order-section">
            <view class="section-header">
                <text class="section-title">我的订单</text>
                <view class="more-btn" @click="toOrderList(null)">
                    <text>全部订单</text>
                    <u-icon name="arrow-right" color="#CCCCCC" size="28rpx"></u-icon>
                </view>
            </view>
            <view class="order-tabs">
                <view class="order-tab" v-for="item in orderTabs" :key="item.status" @click="toOrderList(item.status)">
                    <view class="tab-icon">
                        <u-icon :name="item.icon" :color="item.color" size="48rpx"></u-icon>
                        <view class="badge" v-if="orderStatusCount[item.badge] > 0">
                            <text>{{ orderStatusCount[item.badge] }}</text>
                        </view>
                    </view>
                    <text class="tab-name">{{ item.name }}</text>
                </view>
            </view>
        </view>

        <!-- 菜单列表 -->
        <view class="menu-section">
            <view class="menu-item" v-for="item in menuList" :key="item.name" @click="handleMenu(item)">
                <view class="menu-left">
                    <u-icon :name="item.icon" color="#666666" size="40rpx"></u-icon>
                    <text class="menu-name">{{ item.name }}</text>
                </view>
                <u-icon name="arrow-right" color="#CCCCCC" size="32rpx"></u-icon>
            </view>
        </view>

        <!-- 退出登录按钮 -->
        <view class="logout-btn" v-if="token" @click="logout">
            <text>退出登录</text>
        </view>

        <loginPop :show="showLoginDialog" @close="hideLogin" @success="loginSuccess" />
        <Tabbar :tabsId="4" />
        <u-toast ref="toast" style="z-index: 9999"></u-toast>
    </pageWrapper>
</template>

<script>
    import Tabbar from '@/components/tabbar/tabbar.vue'
    import loginPop from '@/components/loginPop/loginPop.vue'
    import { getUserInfo, setSelfInfo } from "@/api/user"
    import { getOrderStatusCount } from "@/api/order"
    import { getCompanyById } from "@/api/login"
    import { getUser, getToken, setUser, setToken, setOpenId, getRole, setRole, setSettlmentInfo } from '@/store/storage.js'
    import { parseDateStr } from "@/utils/date"
    import config from '@/config/config.js'

    export default {
        components: { Tabbar, loginPop },
        data() {
            return {
                user: {},
                role: {},
                company: null,
                point: 0,
                orderStatusCount: {},
                token: '',
                showLoginDialog: false,
                orderTabs: [
                    { name: '待付款', icon: 'checkbox-mark', color: '#F97316', status: 0, badge: 'unpaid' },
                    { name: '备货中', icon: 'clock', color: '#22A84F', status: 1, badge: 'delivered' },
                    { name: '配送中', icon: 'car', color: '#3B82F6', status: 2, badge: 'shipped' },
                    { name: '已完成', icon: 'checkmark-circle', color: '#666666', status: 3, badge: 'success' }
                ],
                menuList: [
                    { name: '我的账单', icon: 'red-packet', action: 'bill' },
                    { name: '会员信息', icon: 'account', action: 'member' },
                    { name: '积分兑换', icon: 'star', action: 'point' },
                    { name: '收货地址', icon: 'map', action: 'address' },
                    { name: '清除缓存', icon: 'trash', action: 'clear' },
                    { name: '联系我们', icon: 'phone', action: 'contact' }
                ]
            }
        },
        onLoad() {
            this.init()
        },
        onShow() {
            this.user = getUser()
            this.token = getToken()
        },
        methods: {
            async init() {
                this.token = getToken()
                if (!this.token) return

                const res = await getUserInfo()
                if (res.code === 401) {
                    this.token = ''
                    this.user = ''
                    return
                }
                this.user = res.data.userInfo
                this.role = res.data.userInfo.authority
                this.point = res.data.point
                setUser(this.user)
                setRole(this.role)

                // 月结统计
                const date = new Date()
                date.setMonth(date.getMonth() - 1)
                const preRes = await getOrderStatusCount({ settlementMonth: parseDateStr(date.toString()) })
                if (preRes.data?.monthUnpaid > 0) {
                    this.orderStatusCount = preRes.data
                    setSettlmentInfo(this.orderStatusCount)
                } else {
                    const curRes = await getOrderStatusCount()
                    this.orderStatusCount = curRes.data || {}
                }
            },
            toMemberInfo() {
                uni.navigateTo({ url: '/pages/my/memberInfo' })
            },
            toOrderList(status) {
                uni.navigateTo({ url: '/pages/order/list?status=' + (status ?? 'null') })
            },
            handleMenu(item) {
                switch (item.action) {
                    case 'bill': uni.navigateTo({ url: '/pages/bill/list' }); break
                    case 'member': uni.navigateTo({ url: '/pages/my/memberInfo' }); break
                    case 'point': uni.navigateTo({ url: '/pages/goods/pointGoods' }); break
                    case 'address': uni.navigateTo({ url: '/pages/address/address' }); break
                    case 'clear':
                        uni.showModal({
                            title: '清除缓存', content: '是否确认清除？', success: ({ confirm }) => {
                                if (confirm) {
                                    uni.clearStorageSync()
                                    this.$message(this.$refs.toast).success('清除成功')
                                }
                            }
                        }); break
                    case 'contact':
                        uni.makePhoneCall({ phoneNumber: config.phone }); break
                }
            },
            showLogin() {
                this.showLoginDialog = true
            },
            hideLogin() {
                this.showLoginDialog = false
            },
            async loginSuccess(u) {
                this.hideLogin()
                this.token = getToken()
                this.user = u
                await this.init()
                this.$message(this.$refs.toast).success('登录成功')
            },
            logout() {
                uni.showModal({
                    title: '退出登录', content: '是否确认退出？', success: ({ confirm }) => {
                        if (confirm) {
                            this.token = ''
                            this.user = ''
                            setUser('')
                            setToken('')
                            setOpenId('')
                            this.$message(this.$refs.toast).success('已退出')
                        }
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    /* 头部用户卡片 */
    .profile-header {
        background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
        padding: 48rpx 32rpx 32rpx;
    }

    .user-card {
        display: flex;
        align-items: center;
        gap: 24rpx;
        background: #FFFFFF;
        border-radius: 24rpx;
        padding: 32rpx;
        box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.1);

        &.not-login {
            .login-tip {
                color: #999999;
            }
        }
    }

    .user-avatar {
        width: 120rpx;
        height: 120rpx;
        border-radius: 60rpx;
        overflow: hidden;
        flex-shrink: 0;

        image {
            width: 100%;
            height: 100%;
        }
    }

    .user-info {
        flex: 1;
    }

    .user-name {
        font-size: 36rpx;
        font-weight: 600;
        color: #1A1A1A;
        display: block;
        margin-bottom: 16rpx;
    }

    .user-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 12rpx;
    }

    .tag {
        background: #E8F8EC;
        color: #22A84F;
        font-size: 22rpx;
        padding: 6rpx 16rpx;
        border-radius: 20rpx;

        &.warn {
            background: #FEF3E2;
            color: #F97316;
        }

        &.point {
            background: #FEF3E2;
            color: #F97316;
        }
    }

    /* 月结信息 */
    .monthly-info {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 48rpx;
        margin-top: 24rpx;
        background: rgba(255, 255, 255, 0.15);
        border-radius: 16rpx;
        padding: 20rpx;
    }

    .monthly-item {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 8rpx;

        .label {
            font-size: 24rpx;
            color: rgba(255, 255, 255, 0.7);
        }

        .value {
            font-size: 32rpx;
            font-weight: 600;
            color: #FFFFFF;

            &.primary {
                color: #FFD700;
            }
        }
    }

    .divider {
        width: 1rpx;
        height: 48rpx;
        background: rgba(255, 255, 255, 0.3);
    }

    /* 订单区块 */
    .order-section {
        margin: -40rpx 24rpx 24rpx;
        background: #FFFFFF;
        border-radius: 24rpx;
        box-shadow: 0 4rpx 24rpx rgba(0, 0, 0, 0.08);
        padding: 32rpx;
    }

    .section-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 28rpx;
    }

    .section-title {
        font-size: 32rpx;
        font-weight: 600;
        color: #1A1A1A;
    }

    .more-btn {
        display: flex;
        align-items: center;
        gap: 4rpx;
        font-size: 26rpx;
        color: #999999;
    }

    .order-tabs {
        display: flex;
        justify-content: space-around;
    }

    .order-tab {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 12rpx;
    }

    .tab-icon {
        position: relative;
        width: 80rpx;
        height: 80rpx;
        background: #F5F7F4;
        border-radius: 20rpx;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .badge {
        position: absolute;
        top: -8rpx;
        right: -8rpx;
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

    .tab-name {
        font-size: 24rpx;
        color: #666666;
    }

    /* 菜单区块 */
    .menu-section {
        margin: 0 24rpx 24rpx;
        background: #FFFFFF;
        border-radius: 24rpx;
        box-shadow: 0 4rpx 24rpx rgba(0, 0, 0, 0.08);
        overflow: hidden;
    }

    .menu-item {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 32rpx;
        border-bottom: 1rpx solid #EEEEEE;

        &:last-child {
            border-bottom: none;
        }
    }

    .menu-left {
        display: flex;
        align-items: center;
        gap: 16rpx;
    }

    .menu-name {
        font-size: 28rpx;
        color: #1A1A1A;
    }

    /* 退出按钮 */
    .logout-btn {
        margin: 40rpx 24rpx;
        height: 96rpx;
        background: #FFFFFF;
        border-radius: 48rpx;
        display: flex;
        align-items: center;
        justify-content: center;
        box-shadow: 0 4rpx 24rpx rgba(0, 0, 0, 0.08);

        text {
            font-size: 30rpx;
            color: #EF4444;
        }
    }
</style>
