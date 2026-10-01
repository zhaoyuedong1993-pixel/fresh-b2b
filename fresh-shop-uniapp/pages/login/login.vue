<!--
 * 优诚配运 - 登录页
 * 设计规范：自然清新
-->
<template>
    <view class="login-page">
        <!-- 顶部品牌区 -->
        <view class="brand-section">
            <view class="brand-icon">
                <u-icon name="shopping-cart-fill" size="80" color="#FFFFFF"></u-icon>
            </view>
            <text class="brand-name">优诚配运</text>
            <text class="brand-desc">新鲜食材 企业配送</text>
        </view>

        <!-- 登录表单 -->
        <view class="form-section">
            <view class="form-card">
                <!-- 手机号 -->
                <view class="form-item">
                    <view class="form-label">
                        <u-icon name="phone" size="36rpx" color="#22A84F"></u-icon>
                        <text>手机号</text>
                    </view>
                    <input class="form-input" v-model="phone" type="number" maxlength="11"
                        placeholder="请输入手机号" placeholder-class="placeholder" />
                </view>

                <!-- 密码 -->
                <view class="form-item">
                    <view class="form-label">
                        <u-icon name="lock" size="36rpx" color="#22A84F"></u-icon>
                        <text>密码</text>
                    </view>
                    <input class="form-input" v-model="password" type="password" placeholder="请输入密码"
                        placeholder-class="placeholder" />
                </view>

                <!-- 登录按钮 -->
                <view class="login-btn" :class="{ loading }" @click="handleLogin">
                    <text v-if="!loading">登录</text>
                    <text v-else>登录中...</text>
                </view>

                <!-- 注册链接 -->
                <view class="register-link" @click="goRegister">
                    <text>还没有账号？</text>
                    <text class="link-text">立即注册</text>
                </view>
            </view>
        </view>
    </view>
</template>

<script>
    import { loginByPhone } from '@/api/login.js'
    import { setToken, setExpires, setUser, setOpenId } from '@/store/storage.js'

    export default {
        data() {
            return {
                phone: '',
                password: '',
                loading: false
            }
        },
        methods: {
            async handleLogin() {
                if (!this.phone) {
                    uni.showToast({ title: '请输入手机号', icon: 'none' })
                    return
                }
                if (!/^1[3-9]\d{9}$/.test(this.phone)) {
                    uni.showToast({ title: '请输入合法手机号', icon: 'none' })
                    return
                }
                if (!this.password) {
                    uni.showToast({ title: '请输入密码', icon: 'none' })
                    return
                }

                this.loading = true
                try {
                    const res = await loginByPhone({
                        phone: this.phone,
                        password: this.password
                    })

                    if (res.code === 0) {
                        setToken(res.data.token)
                        setExpires(res.data.expiresAt)
                        setUser(res.data.user)
                        setOpenId(res.data.user.openId || '')
                        uni.showToast({ title: '登录成功', icon: 'success' })
                        setTimeout(() => {
                            uni.switchTab({ url: '/pages/index/index' })
                        }, 1500)
                    } else {
                        uni.showToast({ title: res.msg || '登录失败', icon: 'none' })
                    }
                } catch (e) {
                    uni.showToast({ title: '登录异常，请重试', icon: 'none' })
                } finally {
                    this.loading = false
                }
            },
            goRegister() {
                uni.navigateTo({ url: '/pages/login/register' })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .login-page {
        min-height: 100vh;
        background: linear-gradient(180deg, #22A84F 0%, #E8F8EC 60%, #F5F7F4 100%);
        display: flex;
        flex-direction: column;
    }

    /* 品牌区 */
    .brand-section {
        padding: 120rpx 0 80rpx;
        text-align: center;
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 16rpx;
    }

    .brand-icon {
        width: 160rpx;
        height: 160rpx;
        background: rgba(255, 255, 255, 0.2);
        border-radius: 40rpx;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-bottom: 8rpx;
    }

    .brand-name {
        font-size: 48rpx;
        font-weight: 700;
        color: #FFFFFF;
        letter-spacing: 4rpx;
    }

    .brand-desc {
        font-size: 28rpx;
        color: rgba(255, 255, 255, 0.8);
    }

    /* 表单区 */
    .form-section {
        flex: 1;
        padding: 0 40rpx;
    }

    .form-card {
        background: #FFFFFF;
        border-radius: 24rpx;
        padding: 48rpx 40rpx;
        box-shadow: 0 8rpx 40rpx rgba(34, 168, 79, 0.15);
    }

    .form-item {
        margin-bottom: 36rpx;
    }

    .form-label {
        display: flex;
        align-items: center;
        gap: 12rpx;
        margin-bottom: 16rpx;

        text {
            font-size: 28rpx;
            color: #1A1A1A;
            font-weight: 500;
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

    .login-btn {
        height: 96rpx;
        background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
        border-radius: 48rpx;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-top: 24rpx;
        box-shadow: 0 4rpx 24rpx rgba(34, 168, 79, 0.3);

        text {
            font-size: 32rpx;
            font-weight: 600;
            color: #FFFFFF;
        }

        &.loading {
            opacity: 0.7;
        }

        &:active {
            transform: scale(0.98);
        }
    }

    .register-link {
        text-align: center;
        margin-top: 36rpx;
        font-size: 26rpx;
        color: #666666;

        .link-text {
            color: #22A84F;
            font-weight: 600;
            margin-left: 8rpx;
        }
    }
</style>
