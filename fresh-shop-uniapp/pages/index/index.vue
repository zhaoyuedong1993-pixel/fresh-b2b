<!--
 * @Author: dalefeng
 * @Date: 2023-03-23 15:52:23
 * @LastEditors:
 * @LastEditTime: 2026-09-23
-->
<template>
	<pageWrapper>
		<!-- 顶部搜索区域 -->
		<view class="header">
			<view class="header-content">
				<view class="logo">
					<text class="logo-text">Fresh B2B</text>
				</view>
				<view class="search-box" @click="searchClick">
					<u-icon name="search" color="#999" size="18"></u-icon>
					<text class="search-placeholder">搜索商品</text>
				</view>
			</view>
		</view>

		<!-- 分类导航 -->
		<view class="category-nav">
			<view
				v-for="(cat, index) in category"
				:key="cat.ID"
				:class="['category-item', { active: currentCategory === cat.ID }]"
				@click="selectCategory(cat.ID)"
			>
				<text>{{ cat.title }}</text>
			</view>
		</view>

		<!-- 商品列表 -->
		<scroll-view scroll-y="true" class="goods-scroll" refresher-enabled="true" :refresher-triggered="refreshing" @refresherrefresh="onRefresh" @scrolltolower="loadMore">
			<view class="goods-list">
				<view
					v-for="item in goodsList"
					:key="item.id"
					class="goods-item"
					@click="toGoodsInfo(item)"
				>
					<image class="goods-image" :src="item.imgUrl || item.goodsImg" mode="aspectFill"></image>
					<view class="goods-info">
						<text class="goods-name">{{ item.name }}</text>
						<view class="goods-bottom">
							<text class="goods-price">¥{{ item.price }}</text>
							<view class="add-btn" @click.stop="addToCart(item)">
								<u-icon name="plus" color="#fff" size="16"></u-icon>
							</view>
						</view>
					</view>
				</view>
			</view>

			<view class="load-more" v-if="goodsList.length > 0">
				<u-loadmore :status="loadMoreStatus" />
			</view>

			<view class="empty" v-if="goodsList.length === 0 && !loading">
				<u-empty text="暂无商品" icon="https://cdn.uviewui.com/uview/empty/data.png"></u-empty>
			</view>
		</scroll-view>

		<!-- 底部导航 -->
		<Tabbar :tabsId="0" />

		<!-- 购物车悬浮按钮 -->
		<view class="cart-float" @click="goCart">
			<u-badge :value="cartCount" :overflow-count="99" absolute>
				<u-icon name="shopping-cart" size="28" color="#fff"></u-icon>
			</u-badge>
		</view>
	</pageWrapper>
</template>

<script>
import request from '@/utils/request.js'

export default {
	data() {
		return {
			category: [],
			currentCategory: null,
			goodsList: [],
			page: 1,
			pageSize: 20,
			loadMoreStatus: 'loadmore',
			refreshing: false,
			loading: false,
			cartCount: 0
		}
	},
	onLoad() {
		this.loadCategory()
		this.loadCartCount()
	},
	onShow() {
		this.loadCartCount()
	},
	methods: {
		// 加载分类
		async loadCategory() {
			const res = await request({
				url: '/category/getCategoryList',
				method: 'GET'
			})
			if (res.code === 0) {
				this.category = res.data.list || []
				// 默认选中第一个分类
				if (this.category.length > 0) {
					this.currentCategory = this.category[0].ID
					this.loadGoods()
				}
			}
		},

		// 选择分类
		selectCategory(catId) {
			this.currentCategory = catId
			this.page = 1
			this.goodsList = []
			this.loadGoods()
		},

		// 加载商品
		async loadGoods() {
			if (this.loading) return
			this.loading = true

			const res = await request({
				url: '/goods/getGoodsList',
				method: 'GET',
				data: {
					categoryId: this.currentCategory,
					page: this.page,
					pageSize: this.pageSize
				}
			})

			this.loading = false
			this.refreshing = false

			if (res.code === 0) {
				const list = res.data.list || []
				if (this.page === 1) {
					this.goodsList = list
				} else {
					this.goodsList = [...this.goodsList, ...list]
				}
				this.loadMoreStatus = list.length < this.pageSize ? 'nomore' : 'loadmore'
			}
		},

		// 下拉刷新
		onRefresh() {
			this.refreshing = true
			this.page = 1
			this.loadGoods()
		},

		// 加载更多
		loadMore() {
			if (this.loadMoreStatus === 'nomore') return
			this.page++
			this.loadGoods()
		},

		// 搜索
		searchClick() {
			uni.navigateTo({
				url: '/pages/search/search'
			})
		},

		// 商品详情
		toGoodsInfo(item) {
			uni.navigateTo({
				url: `/pages/goods/detail?id=${item.id}`
			})
		},

		// 加入购物车
		async addToCart(item) {
			await request({
				url: '/cart/createCart',
				method: 'POST',
				data: {
					goodsId: item.id,
					quantity: 1
				}
			})
			this.loadCartCount()
			uni.showToast({ title: '已加入购物车', icon: 'success' })
		},

		// 购物车数量
		async loadCartCount() {
			const res = await request({
				url: '/cart/getCartList',
				method: 'GET'
			})
			if (res.code === 0) {
				const list = res.data.list || []
				this.cartCount = list.reduce((sum, item) => sum + item.quantity, 0)
			}
		},

		// 去购物车
		goCart() {
			uni.switchTab({
				url: '/pages/cart/cart'
			})
		}
	}
}
</script>

<style scoped>
.header {
	background: linear-gradient(135deg, #4CAF50, #81C784);
	padding: 20rpx 30rpx;
	padding-top: calc(20rpx + env(safe-area-inset-top));
}

.header-content {
	display: flex;
	align-items: center;
	justify-content: space-between;
}

.logo-text {
	font-size: 36rpx;
	font-weight: bold;
	color: #fff;
}

.search-box {
	flex: 1;
	margin-left: 30rpx;
	background: rgba(255, 255, 255, 0.9);
	border-radius: 30rpx;
	padding: 16rpx 24rpx;
	display: flex;
	align-items: center;
}

.search-placeholder {
	color: #999;
	font-size: 28rpx;
	margin-left: 10rpx;
}

.category-nav {
	display: flex;
	background: #fff;
	padding: 20rpx 0;
	border-bottom: 1rpx solid #f0f0f0;
}

.category-item {
	flex: 1;
	text-align: center;
	font-size: 30rpx;
	color: #666;
	position: relative;
	padding: 10rpx 0;
}

.category-item.active {
	color: #4CAF50;
	font-weight: bold;
}

.category-item.active::after {
	content: '';
	position: absolute;
	bottom: 0;
	left: 50%;
	transform: translateX(-50%);
	width: 60rpx;
	height: 4rpx;
	background: #4CAF50;
	border-radius: 2rpx;
}

.goods-scroll {
	height: calc(100vh - 200rpx - env(safe-area-inset-bottom) - 100rpx);
}

.goods-list {
	display: flex;
	flex-wrap: wrap;
	padding: 20rpx;
}

.goods-item {
	width: calc(50% - 10rpx);
	background: #fff;
	border-radius: 16rpx;
	margin-bottom: 20rpx;
	overflow: hidden;
}

.goods-item:nth-child(odd) {
	margin-right: 20rpx;
}

.goods-image {
	width: 100%;
	height: 320rpx;
	background: #f5f5f5;
}

.goods-info {
	padding: 20rpx;
}

.goods-name {
	font-size: 28rpx;
	color: #333;
	display: -webkit-box;
	-webkit-line-clamp: 2;
	-webkit-box-orient: vertical;
	overflow: hidden;
	line-height: 1.4;
}

.goods-bottom {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-top: 16rpx;
}

.goods-price {
	color: #4CAF50;
	font-size: 32rpx;
	font-weight: bold;
}

.add-btn {
	width: 48rpx;
	height: 48rpx;
	background: #4CAF50;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
}

.load-more {
	padding: 30rpx;
	text-align: center;
}

.cart-float {
	position: fixed;
	right: 30rpx;
	bottom: 200rpx;
	width: 100rpx;
	height: 100rpx;
	background: #4CAF50;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	box-shadow: 0 8rpx 20rpx rgba(76, 175, 80, 0.4);
}

.empty {
	padding: 100rpx 0;
}
</style>
