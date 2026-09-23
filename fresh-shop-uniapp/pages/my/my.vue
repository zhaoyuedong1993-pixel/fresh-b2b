<!--
 * @Author: dalefeng
 * @Date: 2023-03-23 14:30:12
 * @LastEditors:
 * @LastEditTime: 2026-09-23
-->
<template>
	<pageWrapper>
		<!-- 顶部用户信息 -->
		<view class="user-header">
			<view class="user-info" v-if="isLogin">
				<image class="avatar" :src="userInfo.headerImg || 'https://minio.kl.do/picture/images/avatar/face.png'" mode="aspectFill"></image>
				<view class="info">
					<text class="nickname">{{ userInfo.nickName || userInfo.userName }}</text>
					<text class="company">{{ companyName }}</text>
				</view>
			</view>
			<view class="user-info" v-else @click="toLogin">
				<image class="avatar" src="https://minio.kl.do/picture/images/avatar/face.png" mode="aspectFill"></image>
				<view class="info">
					<text class="nickname">点击登录</text>
				</view>
			</view>
		</view>

		<!-- 功能菜单 -->
		<view class="menu-section">
			<view class="menu-item" @click="toBill" v-if="isLogin">
				<view class="menu-left">
					<u-icon name="account-fill" size="22" color="#4CAF50"></u-icon>
					<text class="menu-text">我的账单</text>
				</view>
				<u-icon name="arrow-right" size="16" color="#ccc"></u-icon>
			</view>

			<view class="menu-item" @click="toOrder" v-if="isLogin">
				<view class="menu-left">
					<u-icon name="order" size="22" color="#FF9800"></u-icon>
					<text class="menu-text">我的订单</text>
				</view>
				<u-icon name="arrow-right" size="16" color="#ccc"></u-icon>
			</view>

			<view class="menu-item" @click="toAddress" v-if="isLogin">
				<view class="menu-left">
					<u-icon name="map" size="22" color="#2196F3"></u-icon>
					<text class="menu-text">收货地址</text>
				</view>
				<u-icon name="arrow-right" size="16" color="#ccc"></u-icon>
			</view>

			<view class="menu-item" @click="toFavorites" v-if="isLogin">
				<view class="menu-left">
					<u-icon name="star" size="22" color="#E91E63"></u-icon>
					<text class="menu-text">我的收藏</text>
				</view>
				<u-icon name="arrow-right" size="16" color="#ccc"></u-icon>
			</view>
		</view>

		<!-- 退出登录 -->
		<view class="logout-btn" v-if="isLogin" @click="logout">
			<text>退出登录</text>
		</view>

		<!-- 底部导航 -->
		<Tabbar :tabsId="3" />
	</pageWrapper>
</template>

<script>
import { mapState, mapMutations } from 'vuex'
import request from '@/utils/request.js'

export default {
	data() {
		return {
			companyName: ''
		}
	},
	computed: {
		...mapState(['userInfo', 'isLogin'])
	},
	onShow() {
		if (this.isLogin) {
			this.loadCompanyInfo()
		}
	},
	methods: {
		...mapMutations(['logout']),

		toLogin() {
			uni.navigateTo({ url: '/pages/login/login' })
		},

		toBill() {
			uni.navigateTo({ url: '/pages/bill/list' })
		},

		toOrder() {
			uni.navigateTo({ url: '/pages/order/list' })
		},

		toAddress() {
			uni.navigateTo({ url: '/pages/address/address' })
		},

		toFavorites() {
			uni.navigateTo({ url: '/pages/goods/pointGoods' })
		},

		async loadCompanyInfo() {
			// 从本地存储获取公司名称
			const company = uni.getStorageSync('companyInfo')
			if (company) {
				this.companyName = company.name || ''
			}
		},

		logout() {
			uni.showModal({
				title: '提示',
				content: '确定要退出登录吗？',
				success: (res) => {
					if (res.confirm) {
						this.clearLogin()
					}
				}
			})
		},

		clearLogin() {
			uni.removeStorageSync('token')
			uni.removeStorageSync('userInfo')
			uni.removeStorageSync('companyInfo')
			this.$store.commit('logout')
			uni.reLaunch({
				url: '/pages/my/my'
			})
		}
	}
}
</script>

<style scoped>
.user-header {
	background: linear-gradient(135deg, #4CAF50, #81C784);
	padding: 60rpx 40rpx 80rpx;
}

.user-info {
	display: flex;
	align-items: center;
}

.avatar {
	width: 120rpx;
	height: 120rpx;
	border-radius: 60rpx;
	border: 4rpx solid rgba(255, 255, 255, 0.5);
}

.info {
	margin-left: 30rpx;
}

.nickname {
	display: block;
	font-size: 36rpx;
	font-weight: bold;
	color: #fff;
}

.company {
	display: block;
	font-size: 26rpx;
	color: rgba(255, 255, 255, 0.8);
	margin-top: 10rpx;
}

.menu-section {
	background: #fff;
	margin: -40rpx 30rpx 30rpx;
	border-radius: 20rpx;
	padding: 20rpx 0;
	box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.08);
}

.menu-item {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 30rpx 40rpx;
	border-bottom: 1rpx solid #f5f5f5;
}

.menu-item:last-child {
	border-bottom: none;
}

.menu-left {
	display: flex;
	align-items: center;
}

.menu-text {
	font-size: 30rpx;
	color: #333;
	margin-left: 20rpx;
}

.logout-btn {
	margin: 40rpx 30rpx;
	background: #fff;
	border-radius: 20rpx;
	padding: 30rpx;
	text-align: center;
	font-size: 30rpx;
	color: #f56c6c;
	box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.08);
}
</style>
