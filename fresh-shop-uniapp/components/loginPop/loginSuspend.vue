<!--
 * 优诚配运 - 登录悬浮提示
 * 设计规范：自然清新
-->
<template>
	<view class="login-suspend" v-if="show">
		<view class="suspend-bar" @click="showLogin">
			<view class="bar-content">
				<view class="bar-icon">
					<u-icon name="account" size="36rpx" color="#FFFFFF"></u-icon>
				</view>
				<text class="bar-text">登录查看专属价格及下单订货</text>
			</view>
			<view class="bar-btn">
				<text>立即登录</text>
			</view>
		</view>
	</view>
</template>

<script>
	import loginPop from '@/components/loginPop/loginPop.vue'

	export default {
		name: 'loginSuspend',
		components: {
			loginPop
		},
		props: {
			show: {
				type: Boolean,
				default: false
			}
		},
		data() {
			return {
				showLoginDialog: false
			}
		},
		methods: {
			showLogin() {
				this.showLoginDialog = true
			},
			hideLogin() {
				this.showLoginDialog = false
			},
			loginSuccess() {
				this.hideLogin()
				this.$emit('success')
			}
		}
	}
</script>

<style lang="scss" scoped>
	.login-suspend {
		position: fixed;
		bottom: 0;
		left: 0;
		right: 0;
		z-index: 999;
		padding-bottom: env(safe-area-inset-bottom);
	}

	.suspend-bar {
		height: 100rpx;
		background: linear-gradient(135deg, rgba(0, 0, 0, 0.8) 0%, rgba(0, 0, 0, 0.7) 100%);
		backdrop-filter: blur(10px);
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0 32rpx;
	}

	.bar-content {
		display: flex;
		align-items: center;
		gap: 16rpx;
	}

	.bar-icon {
		width: 64rpx;
		height: 64rpx;
		background: rgba(255, 255, 255, 0.2);
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.bar-text {
		font-size: 26rpx;
		color: #FFFFFF;
	}

	.bar-btn {
		height: 64rpx;
		padding: 0 32rpx;
		background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
		border-radius: 32rpx;
		display: flex;
		align-items: center;
		justify-content: center;

		text {
			font-size: 26rpx;
			font-weight: 600;
			color: #FFFFFF;
		}
	}
</style>
