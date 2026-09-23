<!--
 * 登录页面
-->
<template>
	<view class="login-page">
		<!-- 顶部背景 -->
		<view class="header-bg">
			<view class="logo-area">
				<view class="logo-icon">
					<u-icon name="shopping-cart-fill" size="60" color="#fff"></u-icon>
				</view>
				<text class="app-name">Fresh B2B</text>
				<text class="app-desc">生鲜食材采购平台</text>
			</view>
		</view>

		<!-- 登录表单 -->
		<view class="login-form">
			<view class="form-item">
				<view class="form-label">
					<u-icon name="account" size="20" color="#4CAF50"></u-icon>
					<text>用户名</text>
				</view>
				<input
					class="form-input"
					v-model="username"
					placeholder="请输入用户名"
					placeholder-class="placeholder"
				/>
			</view>

			<view class="form-item">
				<view class="form-label">
					<u-icon name="lock" size="20" color="#4CAF50"></u-icon>
					<text>密码</text>
				</view>
				<input
					class="form-input"
					v-model="password"
					type="password"
					placeholder="请输入密码"
					placeholder-class="placeholder"
				/>
			</view>

			<button class="login-btn" :loading="loading" @click="handleLogin">
				<text v-if="!loading">登录</text>
			</button>
		</view>
	</view>
</template>

<script>
import request from '@/utils/request.js'

export default {
	data() {
		return {
			username: '',
			password: '',
			loading: false
		}
	},
	methods: {
		async handleLogin() {
			if (!this.username) {
				uni.showToast({ title: '请输入用户名', icon: 'none' })
				return
			}
			if (!this.password) {
				uni.showToast({ title: '请输入密码', icon: 'none' })
				return
			}

			this.loading = true

			try {
				const res = await request({
					url: '/base/login',
					method: 'POST',
					data: {
						username: this.username,
						password: this.password
					}
				})

				if (res.code === 0) {
					// 保存登录信息
					uni.setStorageSync('token', res.data.token)
					uni.setStorageSync('userInfo', res.data.user)

					// 更新 Vuex 状态
					this.$store.commit('login', res.data.user)

					uni.showToast({ title: '登录成功', icon: 'success' })

					// 跳转到首页
					setTimeout(() => {
						uni.switchTab({
							url: '/pages/index/index'
						})
					}, 1500)
				} else {
					uni.showToast({ title: res.msg || '登录失败', icon: 'none' })
				}
			} catch (e) {
				uni.showToast({ title: '登录失败', icon: 'none' })
			} finally {
				this.loading = false
			}
		}
	}
}
</script>

<style scoped>
.login-page {
	min-height: 100vh;
	background: #f5f5f5;
}

.header-bg {
	background: linear-gradient(135deg, #4CAF50, #81C784);
	padding: 100rpx 0 120rpx;
	border-radius: 0 0 60rpx 60rpx;
}

.logo-area {
	display: flex;
	flex-direction: column;
	align-items: center;
}

.logo-icon {
	width: 160rpx;
	height: 160rpx;
	background: rgba(255, 255, 255, 0.2);
	border-radius: 40rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	margin-bottom: 30rpx;
}

.app-name {
	font-size: 48rpx;
	font-weight: bold;
	color: #fff;
	margin-bottom: 10rpx;
}

.app-desc {
	font-size: 28rpx;
	color: rgba(255, 255, 255, 0.8);
}

.login-form {
	margin: -60rpx 30rpx 0;
	background: #fff;
	border-radius: 20rpx;
	padding: 40rpx;
	box-shadow: 0 8rpx 30rpx rgba(0, 0, 0, 0.1);
}

.form-item {
	margin-bottom: 30rpx;
}

.form-label {
	display: flex;
	align-items: center;
	margin-bottom: 16rpx;
}

.form-label text {
	margin-left: 12rpx;
	font-size: 28rpx;
	color: #333;
}

.form-input {
	height: 88rpx;
	background: #f5f5f5;
	border-radius: 16rpx;
	padding: 0 24rpx;
	font-size: 28rpx;
}

.placeholder {
	color: #999;
}

.login-btn {
	margin-top: 60rpx;
	height: 96rpx;
	background: linear-gradient(135deg, #4CAF50, #81C784);
	border-radius: 48rpx;
	color: #fff;
	font-size: 32rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	border: none;
}

.login-btn::after {
	border: none;
}
</style>
