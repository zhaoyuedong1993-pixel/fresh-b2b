<!--
 * 优诚配运 - 注册页
 * 设计规范：自然清新
-->
<template>
    <view class="register-page">
        <!-- 顶部品牌区 -->
        <view class="brand-section">
            <view class="brand-icon">
                <u-icon name="shopping-cart-fill" size="80" color="#FFFFFF"></u-icon>
            </view>
            <text class="brand-name">优诚配运</text>
            <text class="brand-desc">新鲜食材 企业配送</text>
        </view>

        <!-- 注册表单 -->
        <view class="form-section">
            <view class="form-card">
                <!-- 申请类型切换 -->
                <view class="type-selector">
                    <view
                        class="type-item"
                        :class="{ active: applyType === 'platform' }"
                        @click="applyType = 'platform'"
                    >
                        <text>申请入驻</text>
                    </view>
                    <view
                        class="type-item"
                        :class="{ active: applyType === 'join' }"
                        @click="applyType = 'join'"
                    >
                        <text>加入公司</text>
                    </view>
                </view>

                <!-- 申请入驻表单 -->
                <view class="form-content" v-if="applyType === 'platform'">
                    <view class="form-item">
                        <view class="form-label">
                            <u-icon name="home" size="36rpx" color="#22A84F"></u-icon>
                            <text>公司名称</text>
                        </view>
                        <input
                            class="form-input"
                            v-model="formData.name"
                            placeholder="请输入公司名称"
                            placeholder-class="placeholder"
                        />
                    </view>

                    <view class="form-item">
                        <view class="form-label">
                            <u-icon name="map" size="36rpx" color="#22A84F"></u-icon>
                            <text>公司地址</text>
                        </view>
                        <input
                            class="form-input"
                            v-model="formData.address"
                            placeholder="请输入公司地址"
                            placeholder-class="placeholder"
                        />
                    </view>

                    <view class="form-item">
                        <view class="form-label">
                            <u-icon name="account" size="36rpx" color="#22A84F"></u-icon>
                            <text>联系人</text>
                        </view>
                        <input
                            class="form-input"
                            v-model="formData.contactName"
                            placeholder="请输入联系人姓名"
                            placeholder-class="placeholder"
                        />
                    </view>

                    <view class="form-item">
                        <view class="form-label">
                            <u-icon name="phone" size="36rpx" color="#22A84F"></u-icon>
                            <text>手机号</text>
                        </view>
                        <input
                            class="form-input"
                            v-model="formData.phone"
                            type="number"
                            maxlength="11"
                            placeholder="请输入手机号"
                            placeholder-class="placeholder"
                        />
                    </view>

                    <view class="tip-text">
                        <u-icon name="info-circle" size="28rpx" color="#999999"></u-icon>
                        <text>提交后请等待平台管理员审核</text>
                    </view>
                </view>

                <!-- 加入公司表单 -->
                <view class="form-content" v-else>
                    <view class="form-item">
                        <view class="form-label">
                            <u-icon name="tags" size="36rpx" color="#22A84F"></u-icon>
                            <text>邀请码</text>
                        </view>
                        <input
                            class="form-input"
                            v-model="formData.invitationCode"
                            placeholder="请输入公司邀请码"
                            placeholder-class="placeholder"
                        />
                    </view>

                    <!-- 公司信息预览 -->
                    <view class="company-preview" v-if="companyInfo">
                        <view class="company-icon">
                            <u-icon name="building" size="40rpx" color="#22A84F"></u-icon>
                        </view>
                        <view class="company-info">
                            <text class="company-name">{{ companyInfo.name }}</text>
                            <text class="company-contact">联系人：{{ companyInfo.contact }}</text>
                        </view>
                        <u-icon name="checkmark-circle-fill" size="36rpx" color="#22A84F"></u-icon>
                    </view>

                    <view class="form-item">
                        <view class="form-label">
                            <u-icon name="phone" size="36rpx" color="#22A84F"></u-icon>
                            <text>手机号</text>
                        </view>
                        <input
                            class="form-input"
                            v-model="formData.phone"
                            type="number"
                            maxlength="11"
                            placeholder="请输入手机号"
                            placeholder-class="placeholder"
                        />
                    </view>

                    <view class="tip-text">
                        <u-icon name="info-circle" size="28rpx" color="#999999"></u-icon>
                        <text>提交后请等待公司管理员审核</text>
                    </view>
                </view>

                <!-- 注册按钮 -->
                <view class="register-btn" :class="{ loading }" @click="handleRegister">
                    <text v-if="!loading">提交申请</text>
                    <text v-else>提交中...</text>
                </view>

                <!-- 返回登录 -->
                <view class="back-login" @click="goLogin">
                    <text>已有账号？</text>
                    <text class="link">立即登录</text>
                </view>
            </view>
        </view>
    </view>
</template>

<script>
    import { registerCompany, getCompanyByInviteCode, joinCompany } from '@/api/login.js'

    export default {
        data() {
            return {
                applyType: 'platform',
                formData: {
                    name: '',
                    address: '',
                    contactName: '',
                    phone: '',
                    invitationCode: ''
                },
                companyInfo: null,
                loading: false
            }
        },
        watch: {
            'formData.invitationCode': {
                handler(val) {
                    if (val && val.length >= 6) {
                        this.checkInviteCode()
                    } else {
                        this.companyInfo = null
                    }
                }
            }
        },
        methods: {
            async checkInviteCode() {
                try {
                    const res = await getCompanyByInviteCode(this.formData.invitationCode)
                    if (res.code === 0) {
                        this.companyInfo = res.data
                    } else {
                        this.companyInfo = null
                    }
                } catch (e) {
                    this.companyInfo = null
                }
            },
            async handleRegister() {
                if (!this.formData.phone) {
                    uni.showToast({ title: '请输入手机号', icon: 'none' })
                    return
                }
                if (!/^1[3-9]\d{9}$/.test(this.formData.phone)) {
                    uni.showToast({ title: '请输入合法手机号', icon: 'none' })
                    return
                }

                this.loading = true

                try {
                    let res
                    if (this.applyType === 'platform') {
                        if (!this.formData.name) {
                            uni.showToast({ title: '请输入公司名称', icon: 'none' })
                            this.loading = false
                            return
                        }
                        res = await registerCompany({
                            name: this.formData.name,
                            address: this.formData.address,
                            contactName: this.formData.contactName,
                            phone: this.formData.phone,
                            password: ''
                        })
                    } else {
                        if (!this.formData.invitationCode) {
                            uni.showToast({ title: '请输入邀请码', icon: 'none' })
                            this.loading = false
                            return
                        }
                        if (!this.companyInfo) {
                            uni.showToast({ title: '请输入正确的邀请码', icon: 'none' })
                            this.loading = false
                            return
                        }
                        res = await joinCompany({
                            invitationCode: this.formData.invitationCode,
                            phone: this.formData.phone
                        })
                    }

                    if (res.code === 0) {
                        uni.showToast({ title: '提交成功', icon: 'success' })
                        setTimeout(() => {
                            uni.redirectTo({ url: '/pages/login/login' })
                        }, 1500)
                    } else {
                        uni.showToast({ title: res.msg || '提交失败', icon: 'none' })
                    }
                } catch (e) {
                    uni.showToast({ title: '提交失败', icon: 'none' })
                } finally {
                    this.loading = false
                }
            },
            goLogin() {
                uni.redirectTo({ url: '/pages/login/login' })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .register-page {
        min-height: 100vh;
        background: linear-gradient(180deg, #22A84F 0%, #E8F8EC 60%, #F5F7F4 100%);
    }

    .brand-section {
        padding: 100rpx 0 60rpx;
        text-align: center;
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 16rpx;
    }

    .brand-icon {
        width: 140rpx;
        height: 140rpx;
        background: rgba(255, 255, 255, 0.2);
        border-radius: 36rpx;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-bottom: 8rpx;
    }

    .brand-name {
        font-size: 44rpx;
        font-weight: 700;
        color: #FFFFFF;
        letter-spacing: 4rpx;
    }

    .brand-desc {
        font-size: 26rpx;
        color: rgba(255, 255, 255, 0.8);
    }

    .form-section {
        padding: 0 40rpx;
    }

    .form-card {
        background: #FFFFFF;
        border-radius: 24rpx;
        padding: 40rpx 36rpx;
        box-shadow: 0 8rpx 40rpx rgba(34, 168, 79, 0.15);
    }

    .type-selector {
        display: flex;
        background: #F5F7F4;
        border-radius: 16rpx;
        padding: 6rpx;
        margin-bottom: 36rpx;
    }

    .type-item {
        flex: 1;
        text-align: center;
        padding: 16rpx 0;
        border-radius: 12rpx;
        font-size: 28rpx;
        color: #666666;
        transition: all 0.2s;

        &.active {
            background: #FFFFFF;
            color: #22A84F;
            font-weight: 600;
            box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.08);
        }
    }

    .form-content {
        margin-bottom: 24rpx;
    }

    .form-item {
        margin-bottom: 28rpx;
    }

    .form-label {
        display: flex;
        align-items: center;
        gap: 12rpx;
        margin-bottom: 12rpx;

        text {
            font-size: 28rpx;
            color: #666666;
        }
    }

    .form-input {
        height: 88rpx;
        background: #F5F7F4;
        border-radius: 16rpx;
        padding: 0 24rpx;
        font-size: 30rpx;
        color: #1A1A1A;
    }

    .placeholder {
        color: #CCCCCC;
        font-size: 28rpx;
    }

    .company-preview {
        display: flex;
        align-items: center;
        gap: 16rpx;
        padding: 24rpx;
        background: #E8F8EC;
        border-radius: 16rpx;
        margin-bottom: 28rpx;
    }

    .company-icon {
        width: 72rpx;
        height: 72rpx;
        background: #FFFFFF;
        border-radius: 16rpx;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .company-info {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 4rpx;
    }

    .company-name {
        font-size: 30rpx;
        font-weight: 600;
        color: #1A1A1A;
    }

    .company-contact {
        font-size: 24rpx;
        color: #666666;
    }

    .tip-text {
        display: flex;
        align-items: center;
        gap: 8rpx;
        justify-content: center;
        margin-bottom: 24rpx;

        text {
            font-size: 24rpx;
            color: #999999;
        }
    }

    .register-btn {
        height: 96rpx;
        background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
        border-radius: 48rpx;
        display: flex;
        align-items: center;
        justify-content: center;
        box-shadow: 0 4rpx 24rpx rgba(34, 168, 79, 0.3);

        text {
            font-size: 32rpx;
            font-weight: 600;
            color: #FFFFFF;
        }

        &.loading {
            opacity: 0.7;
        }
    }

    .back-login {
        text-align: center;
        margin-top: 36rpx;
        font-size: 26rpx;
        color: #666666;

        .link {
            color: #22A84F;
            font-weight: 600;
            margin-left: 8rpx;
        }
    }
</style>
