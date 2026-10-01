<!--
 * 优诚配运 - 登录弹窗
 * 设计规范：自然清新
-->
<template>
	<u-popup :show="show" @close="close" mode="bottom" round="24rpx" :closeable="true" :closeIconPos="'top-left'">
		<view class="login-popup">
			<view class="popup-header">
				<text class="popup-title">登录</text>
				<text class="popup-subtitle">登录后享受专属价格和便捷下单</text>
			</view>

			<view class="popup-content">
				<!-- 微信登录 -->
				<button class="wx-login-btn" open-type="getPhoneNumber" @getphonenumber="getPhoneNumber">
					<view class="wx-icon">
						<u-icon name="weixin-fill" size="40rpx" color="#FFFFFF"></u-icon>
					</view>
					<text>微信用户一键登录</text>
				</button>

				<!-- 分隔线 -->
				<view class="divider">
					<view class="divider-line"></view>
					<text class="divider-text">其他方式</text>
					<view class="divider-line"></view>
				</view>

				<!-- 手机号登录 -->
				<view class="phone-login-btn" @click="goPhoneLogin">
					<u-icon name="phone" size="36rpx" color="#22A84F"></u-icon>
					<text>手机号登录/注册</text>
				</view>
			</view>

			<view class="popup-footer">
				<text class="agreement-text">登录即表示同意</text>
				<text class="agreement-link">《用户服务协议》</text>
				<text class="agreement-text">和</text>
				<text class="agreement-link">《隐私政策》</text>
			</view>
		</view>

		<u-toast ref="toast" style="z-index: 9999;"></u-toast>
	</u-popup>
</template>

<script>
	import { setToken, setExpires, setUser, setOpenId } from '@/store/storage.js'
	import { getWeChatOpenIdByCode, wxLogin } from '@/api/login.js'

	export default {
		name: 'loginPop',
		data() {
			return {
				openId: '',
				sessionKey: ''
			}
		},
		props: {
			show: {
				type: Boolean,
				default: false
			},
			closeable: {
				type: Boolean,
				default: true
			}
		},
		watch: {
			show(val) {
				if (val) {
					this.initLogin()
				}
			}
		},
		methods: {
			async initLogin() {
				try {
					const loginRes = await this.loginGetJsCode()
					if (loginRes.code && loginRes.errMsg === 'login:ok') {
						const openIdRes = await getWeChatOpenIdByCode({ js_code: loginRes.code })
						if (openIdRes.data?.errcode === 0) {
							this.openId = openIdRes.data.openid
							this.sessionKey = openIdRes.data.session_key
						}
					}
				} catch (e) {
					console.error('Login init error:', e)
				}
			},
			loginGetJsCode() {
				return new Promise((resolve, reject) => {
					uni.login({
						success: resolve,
						fail: reject
					})
				})
			},
			getPhoneNumber(e) {
				if (e?.detail?.errMsg === 'getPhoneNumber:ok') {
					this.handleWxLogin(e.detail.encryptedData, e.detail.iv)
				} else {
					this.$message(this.$refs.toast).error('获取手机号失败')
				}
			},
			async handleWxLogin(encryptedData, iv) {
				if (!this.sessionKey) {
					this.$message(this.$refs.toast).warning('请稍后重试')
					return
				}

				try {
					const res = await wxLogin({
						encryptedData,
						iv,
						sessionKey: this.sessionKey,
						openid: this.openId
					})

					if (res.code === 0) {
						setToken(res.data.token)
						setExpires(res.data.expiresAt)
						setOpenId(res.data.user.openId || '')
						setUser(res.data.user)

						this.$message(this.$refs.toast).success('登录成功')
						this.$emit('success', res.data.user)

						if (res.data.user.auditStatus !== 1) {
							uni.navigateTo({ url: '/pages/my/memberInfo' })
						}
					} else {
						this.$message(this.$refs.toast).error(res.msg || '登录失败')
					}
				} catch (e) {
					this.$message(this.$refs.toast).error('登录失败')
				}
			},
			close() {
				this.$emit('close')
			},
			goPhoneLogin() {
				this.close()
				uni.navigateTo({ url: '/pages/login/login' })
			}
		}
	}
</script>

<style lang="scss" scoped>
	.login-popup {
		padding: 48rpx 40rpx;
		padding-bottom: calc(env(safe-area-inset-bottom) + 40rpx);
	}

	.popup-header {
		text-align: center;
		margin-bottom: 48rpx;
	}

	.popup-title {
		display: block;
		font-size: 40rpx;
		font-weight: 700;
		color: #1A1A1A;
		margin-bottom: 12rpx;
	}

	.popup-subtitle {
		display: block;
		font-size: 26rpx;
		color: #999999;
	}

	.popup-content {
		margin-bottom: 40rpx;
	}

	.wx-login-btn {
		width: 100%;
		height: 96rpx;
		background: linear-gradient(135deg, #07C160 0%, #06AD56 100%);
		border-radius: 48rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 16rpx;
		border: none;
		box-shadow: 0 4rpx 20rpx rgba(7, 193, 96, 0.3);

		text {
			font-size: 32rpx;
			font-weight: 600;
			color: #FFFFFF;
		}

		&::after {
			border: none;
		}
	}

	.wx-icon {
		width: 56rpx;
		height: 56rpx;
		background: rgba(255, 255, 255, 0.2);
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.divider {
		display: flex;
		align-items: center;
		margin: 40rpx 0;
		gap: 20rpx;
	}

	.divider-line {
		flex: 1;
		height: 1rpx;
		background: #EEEEEE;
	}

	.divider-text {
		font-size: 24rpx;
		color: #999999;
		flex-shrink: 0;
	}

	.phone-login-btn {
		width: 100%;
		height: 96rpx;
		background: #FFFFFF;
		border: 2rpx solid #22A84F;
		border-radius: 48rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 12rpx;

		text {
			font-size: 32rpx;
			font-weight: 600;
			color: #22A84F;
		}
	}

	.popup-footer {
		text-align: center;
		font-size: 22rpx;
		color: #999999;
	}

	.agreement-text {
		color: #999999;
	}

	.agreement-link {
		color: #22A84F;
	}
</style>
