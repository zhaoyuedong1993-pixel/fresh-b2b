<!--
 * 注册页面 - 支持申请入驻平台和加入已有公司
-->
<template>
	<view class="register-page">
		<!-- 顶部背景 -->
		<view class="header-bg">
			<view class="logo-area">
				<view class="logo-icon">
					<u-icon name="shopping-cart-fill" size="60" color="#fff"></u-icon>
				</view>
				<text class="app-name">Fresh B2B</text>
				<text class="app-desc">申请加入我们</text>
			</view>
		</view>

		<!-- 注册表单 -->
		<view class="register-form">
			<!-- 申请类型切换 -->
			<view class="type-selector">
				<view class="type-item" :class="{ active: applyType === 'platform' }" @click="applyType = 'platform'">
					<text>申请入驻平台</text>
				</view>
				<view class="type-item" :class="{ active: applyType === 'join' }" @click="applyType = 'join'">
					<text>加入已有公司</text>
				</view>
			</view>

			<!-- 申请入驻平台表单 -->
			<view v-if="applyType === 'platform'">
				<view class="form-item">
					<view class="form-label">
						<u-icon name="home" size="20" color="#4CAF50"></u-icon>
						<text>公司名称</text>
					</view>
					<input class="form-input" v-model="formData.name" placeholder="请输入公司名称" />
				</view>

				<view class="form-item">
					<view class="form-label">
						<u-icon name="map" size="20" color="#4CAF50"></u-icon>
						<text>公司地址</text>
					</view>
					<input class="form-input" v-model="formData.address" placeholder="请输入公司地址" />
				</view>

				<view class="form-item">
					<view class="form-label">
						<u-icon name="account" size="20" color="#4CAF50"></u-icon>
						<text>联系人</text>
					</view>
					<input class="form-input" v-model="formData.contactName" placeholder="请输入联系人姓名" />
				</view>

				<view class="form-item">
					<view class="form-label">
						<u-icon name="phone" size="20" color="#4CAF50"></u-icon>
						<text>手机号</text>
					</view>
					<input class="form-input" v-model="formData.phone" type="number" maxlength="11" placeholder="请输入手机号" />
				</view>

				<view class="tip-text">
					提交后请等待平台管理员审核
				</view>
			</view>

			<!-- 加入已有公司表单 -->
			<view v-else>
				<view class="form-item">
					<view class="form-label">
						<u-icon name="tags" size="20" color="#4CAF50"></u-icon>
						<text>邀请码</text>
					</view>
					<input class="form-input" v-model="formData.invitationCode" placeholder="请输入公司邀请码" />
				</view>

				<!-- 公司信息预览 -->
				<view v-if="companyInfo" class="company-preview">
					<view class="company-name">{{ companyInfo.name }}</view>
					<view class="company-contact">联系人：{{ companyInfo.contact }}</view>
				</view>

				<view class="form-item">
					<view class="form-label">
						<u-icon name="phone" size="20" color="#4CAF50"></u-icon>
						<text>手机号</text>
					</view>
					<input class="form-input" v-model="formData.phone" type="number" maxlength="11" placeholder="请输入手机号" />
				</view>

				<view class="tip-text">
					提交后请等待公司管理员审核
				</view>
			</view>

			<button class="register-btn" :loading="loading" @click="handleRegister">
				<text v-if="!loading">提交申请</text>
			</button>

			<view class="back-login" @click="goLogin">
				已有账号？<text class="link">立即登录</text>
			</view>
		</view>
	</view>
</template>

<script>
import { registerCompany, getCompanyByInviteCode, joinCompany } from '@/api/login.js'

export default {
	data() {
		return {
			applyType: 'platform', // platform: 申请入驻平台, join: 加入已有公司
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
		// 检查邀请码
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
		// 提交申请
		async handleRegister() {
			// 验证手机号
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
					// 申请入驻平台
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
					// 加入已有公司
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
					uni.showToast({ title: res.msg || '提交成功', icon: 'success' })
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

<style scoped>
.register-page {
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

.register-form {
	margin: -60rpx 30rpx 0;
	background: #fff;
	border-radius: 20rpx;
	padding: 40rpx;
	box-shadow: 0 8rpx 30rpx rgba(0, 0, 0, 0.1);
}

.type-selector {
	display: flex;
	background: #f5f5f5;
	border-radius: 16rpx;
	padding: 6rpx;
	margin-bottom: 30rpx;
}

.type-item {
	flex: 1;
	text-align: center;
	padding: 16rpx 0;
	border-radius: 12rpx;
	font-size: 28rpx;
	color: #666;
	transition: all 0.3s;
}

.type-item.active {
	background: #4CAF50;
	color: #fff;
	font-weight: 500;
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

.company-preview {
	background: #e8f5e9;
	border-radius: 12rpx;
	padding: 20rpx;
	margin-bottom: 30rpx;
}

.company-name {
	font-size: 32rpx;
	font-weight: bold;
	color: #2e7d32;
	margin-bottom: 8rpx;
}

.company-contact {
	font-size: 26rpx;
	color: #666;
}

.tip-text {
	font-size: 24rpx;
	color: #999;
	margin-bottom: 30rpx;
	text-align: center;
}

.register-btn {
	margin-top: 40rpx;
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

.register-btn::after {
	border: none;
}

.back-login {
	text-align: center;
	margin-top: 40rpx;
	font-size: 28rpx;
	color: #666;
}

.link {
	color: #4CAF50;
}
</style>
