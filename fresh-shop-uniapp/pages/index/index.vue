<!--
 * 优诚配运 - 首页
 * 设计规范：自然清新
-->
<template>
	<pageWrapper>
		<!-- 顶部搜索区 -->
		<view class="home-header">
			<view class="search-bar" @click="searchClick">
				<u-icon name="search" color="#999999" size="32rpx"></u-icon>
				<text class="search-placeholder">搜索商品</text>
			</view>
		</view>

		<!-- 页面内容 -->
		<scroll-view scroll-y class="home-scroll" :style="{ height: scrollHeight + 'px' }" refresher-enabled
			:refresher-triggered="isRefreshing" @refresherrefresh="onRefresh" @scrolltolower="onScrollLower"
			:scroll-anchoring="true">
			<!-- 轮播图 -->
			<view class="banner-wrap">
				<u-swiper :list="banner" keyName="imgUrl" indicator indicatorMode="line" :height="320" circular
					bgColor="transparent" :autoplay="true" :interval="3000" @click="clickBanner" radius="16rpx">
				</u-swiper>
			</view>

			<!-- 分类网格 -->
			<view class="category-card">
				<view class="category-grid">
					<view class="category-item" v-for="c in category" :key="c.ID" @click="toGoodsByCategory(c.ID)">
						<view class="category-icon-wrap">
							<u--image width="96rpx" height="96rpx" :src="c.imgUrl" shape="circle"></u--image>
						</view>
						<text class="category-title">{{ c.title }}</text>
					</view>
				</view>
			</view>

			<!-- 商品 Tab -->
			<view class="goods-tabs-wrap">
				<view class="goods-tabs">
					<view class="tab-item" :class="{ active: goodsTabsId === 0 }" @click="changeGoodsTabs(0)">
						<text class="tab-text">热销商品</text>
						<view class="tab-indicator" v-if="goodsTabsId === 0"></view>
					</view>
					<view class="tab-item" :class="{ active: goodsTabsId === 1 }" @click="changeGoodsTabs(1)">
						<text class="tab-text">新品上市</text>
						<view class="tab-indicator" v-if="goodsTabsId === 1"></view>
					</view>
				</view>
			</view>

			<!-- 商品列表 -->
			<view class="goods-list-wrap">
				<swiper :current="goodsTabsId" @change="onChangeGoodsTabs" :style="{ height: listSwiperHeight + 'px' }">
					<!-- 热销商品 -->
					<swiper-item>
						<GoodsList :lists="goodsHotArr" price-type="¥" @onGoods="toGoodsInfo" :is-audit="isAudit"></GoodsList>
						<view class="load-more" v-if="goodsHotArr.length > 0">
							<u-loadmore :status="hotLoadMore" loading-text="加载中..." loadmore-text="上拉加载更多" nomore-text="没有更多了" />
						</view>
						<view class="empty-tip" v-if="goodsHotArr.length === 0 && !isLoading">
							<u-empty text="暂无热销商品" :icon="emptyIcon"></u-empty>
						</view>
					</swiper-item>
					<!-- 新品上市 -->
					<swiper-item>
						<GoodsList :lists="goodsNewArr" price-type="¥" :is-audit="isAudit"></GoodsList>
						<view class="load-more" v-if="goodsNewArr.length > 0">
							<u-loadmore :status="newLoadMore" loading-text="加载中..." loadmore-text="上拉加载更多" nomore-text="没有更多了" />
						</view>
						<view class="empty-tip" v-if="goodsNewArr.length === 0 && !isLoading">
							<u-empty text="暂无新品" :icon="emptyIcon"></u-empty>
						</view>
					</swiper-item>
				</swiper>
			</view>
		</scroll-view>

		<!-- 登录悬浮 -->
		<loginSuspend :show="loginSuspendShow" @success="loginSuccess"></loginSuspend>

		<!-- 底部导航 -->
		<Tabbar :tabsId="0" />

		<!-- 未结算提醒弹窗 -->
		<u-modal :show="showSettlmentUnpaid" showCancelButton closeOnClickOverlay @confirm="callPhone"
			@cancel="showSettlmentUnpaid = false" @close="showSettlmentUnpaid = false" confirmText="联系商家"
			cancelText="稍后处理" title="未结算订单提醒">
			<view class="modal-content">
				<view class="modal-main">您有{{ preOrderStatus.month }}月未结算的订单需要处理</view>
				<view class="modal-info">
					<view class="info-row">共有 <text class="highlight">{{ preOrderStatus.monthUnpaid }}</text> 个订单未结算</view>
					<view class="info-row">结算金额 <text class="highlight">¥{{ preOrderStatus.settlementUnpaid }}</text></view>
				</view>
				<view class="modal-tip">为确保您的正常使用，请尽快处理</view>
			</view>
		</u-modal>

		<u-toast style="z-index:9998;" ref="toast"></u-toast>
	</pageWrapper>
</template>

<script>
	import Tabbar from '@/components/tabbar/tabbar.vue'
	import loginSuspend from '@/components/loginPop/loginSuspend.vue'
	import GoodsList from '@/components/goodsList/goodsList.vue'
	import config from '@/config/config.js'
	import { getBannerList } from '@/api/banner.js'
	import { getHomeCategoryList } from '@/api/category.js'
	import { getGoodsPageList } from '@/api/goods.js'
	import { getToken, getUser, setUser, setRole, setSettlmentInfo } from '@/store/storage.js'
	import { getUserAuditStatus, getUserInfo } from "@/api/user";
	import { getOrderStatusCount } from "@/api/order";
	import { parseDateStr } from "@/utils/date";

	export default {
		components: { Tabbar, GoodsList, loginSuspend },
		data() {
			return {
				isAudit: false,
				isRefreshing: false,
				isLoading: false,
				loginSuspendShow: false,
				showSettlmentUnpaid: false,
				relationPhone: '',
				goodsTabsId: 0,
				scrollHeight: 1000,
				listSwiperHeight: 600,
				hotLoadMore: 'loadmore',
				newLoadMore: 'loadmore',
				hotPage: { page: 1, pageSize: 12, total: 0, isMore: true },
				newPage: { page: 1, pageSize: 12, total: 0, isMore: true },
				banner: [],
				category: [],
				goodsHotArr: [],
				goodsNewArr: [],
				preOrderStatus: {},
				emptyIcon: 'http://cdn.uviewui.com/uview/empty/data.png'
			}
		},
		onLoad() {
			let user = getUser()
			if (user) {
				if (user.auditStatus === 1) {
					this.isAudit = true
				} else {
					getUserAuditStatus().then(res => {
						if (res.data?.auditStatus === 1) {
							this.isAudit = true
							user.auditStatus = 1
							setUser(user)
						}
					})
				}
			}
			this.relationPhone = config.phone
			this.getBanner()
			this.getHomeCategory()
			this.getGoodsListData(0)
			this.getGoodsListData(1)

			const t = getToken()
			if (!t) {
				this.loginSuspendShow = true
				return
			}

			this.getUserInfoData()
			const date = new Date()
			date.setMonth(date.getMonth() - 1)
			this.getOrderStatusCountInfo(date)
		},
		onShow() {
			const t = getToken()
			if (t && this.loginSuspendShow) {
				this.loginSuspendShow = false
				this.getUserInfoData()
			}
		},
		mounted() {
			uni.getSystemInfo({
				success: (res) => {
					this.scrollHeight = res.windowHeight - 120
					this.listSwiperHeight = res.windowHeight - 380
				}
			})
		},
		methods: {
			async getUserInfoData() {
				const t = getToken()
				if (!t) return
				const res = await getUserInfo()
				if (res.code === 0) {
					setUser(res.data.userInfo)
					setRole(res.data.userInfo.authority)
				}
			},
			getOrderStatusCountInfo(date) {
				const settlementMonth = parseDateStr(date.toString())
				getOrderStatusCount({ settlementMonth }).then((res) => {
					if (res.data?.monthUnpaid > 0) {
						setSettlmentInfo(res.data)
						this.preOrderStatus = res.data
						this.showSettlmentUnpaid = true
					}
				})
			},
			onShareAppMessage() {
				return {
					title: '优诚配运 - 新鲜食材配送',
					path: '/pages/index/index',
					imageUrl: '/static/logo.png'
				}
			},
			onShareTimeline() {
				return {
					title: '优诚配运 - 新鲜食材配送',
					imageUrl: '/static/logo.png'
				}
			},
			searchClick() {
				uni.navigateTo({ url: '/pages/search/search' })
			},
			getBanner() {
				getBannerList().then(res => {
					if (res.code !== 0) return
					res.data.list.forEach(item => {
						if (item.imgUrl && item.imgUrl.slice(0, 4) !== 'http') {
							item.imgUrl = config.baseUrl + "/" + item.imgUrl
						}
					})
					this.banner = res.data.list || []
				})
			},
			getHomeCategory() {
				getHomeCategoryList().then(res => {
					if (res.code !== 0) return
					res.data.list.forEach(item => {
						if (item.imgUrl && item.imgUrl.slice(0, 4) !== 'http') {
							item.imgUrl = config.baseUrl + "/" + item.imgUrl
						}
					})
					this.category = res.data.list || []
				})
			},
			changeGoodsTabs(id) {
				this.goodsTabsId = id
			},
			onChangeGoodsTabs(e) {
				this.goodsTabsId = e.detail.current
			},
			async getGoodsListData(tabId, append = false) {
				const data = { goodsArea: 0 }
				if (!append) {
					this.hotPage.page = 1
					this.newPage.page = 1
					this.hotPage.isMore = true
					this.newPage.isMore = true
					this.hotLoadMore = 'loadmore'
					this.newLoadMore = 'loadmore'
				}
				if (tabId === 0) {
					data.isHot = 1
					data.page = this.hotPage.page
					data.pageSize = this.hotPage.pageSize
				} else if (tabId === 1) {
					data.isNew = 1
					data.page = this.newPage.page
					data.pageSize = this.newPage.pageSize
				} else return

				this.isLoading = true
				const res = await getGoodsPageList(data)
				this.isLoading = false
				if (res.code !== 0) return

				res.data.list?.forEach(item => {
					if (item.images?.[0] && item.images[0].url?.slice(0, 4) !== 'http') {
						item.images[0].url = config.baseUrl + "/" + item.images[0].url
					}
				})

				if (tabId === 0) {
					this.hotPage.total = res.data.total || 0
					if (append) {
						this.goodsHotArr = [...this.goodsHotArr, ...(res.data.list || [])]
					} else {
						this.goodsHotArr = res.data.list || []
					}
					if (this.hotPage.page * this.hotPage.pageSize >= this.hotPage.total) {
						this.hotPage.isMore = false
						this.hotLoadMore = 'nomore'
					}
					this.hotPage.page++
				} else {
					this.newPage.total = res.data.total || 0
					if (append) {
						this.goodsNewArr = [...this.goodsNewArr, ...(res.data.list || [])]
					} else {
						this.goodsNewArr = res.data.list || []
					}
					if (this.newPage.page * this.newPage.pageSize >= this.newPage.total) {
						this.newPage.isMore = false
						this.newLoadMore = 'nomore'
					}
					this.newPage.page++
				}
			},
			async onRefresh() {
				this.isRefreshing = true
				await this.getGoodsListData(this.goodsTabsId, false)
				this.isRefreshing = false
				this.$message(this.$refs.toast).success("刷新成功")
			},
			async onScrollLower() {
				if (this.goodsTabsId === 0) {
					if (this.hotLoadMore === 'loading' || !this.hotPage.isMore) return
					this.hotLoadMore = 'loading'
					await this.getGoodsListData(0, true)
					this.hotLoadMore = this.hotPage.isMore ? 'loadmore' : 'nomore'
				} else {
					if (this.newLoadMore === 'loading' || !this.newPage.isMore) return
					this.newLoadMore = 'loading'
					await this.getGoodsListData(1, true)
					this.newLoadMore = this.newPage.isMore ? 'loadmore' : 'nomore'
				}
			},
			toGoodsByCategory(categoryId) {
				uni.navigateTo({ url: `/pages/goods/goods?categoryId=${categoryId}` })
			},
			toGoodsInfo(goods) {
				uni.navigateTo({ url: `/pages/goods/goodsInfo?id=${goods.ID}` })
			},
			loginSuccess() {
				this.loginSuspendShow = false
				this.getUserInfoData()
			},
			clickBanner(index) {
				const b = this.banner[index]
				if (b?.type === 1 && b.toPath) {
					uni.navigateTo({ url: b.toPath })
				}
			},
			callPhone() {
				uni.makePhoneCall({ phoneNumber: this.relationPhone })
			}
		}
	}
</script>

<style lang="scss" scoped>
	/* 顶部搜索区 */
	.home-header {
		padding: 16rpx 24rpx 12rpx;
		background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
	}

	.search-bar {
		display: flex;
		align-items: center;
		height: 72rpx;
		background: #FFFFFF;
		border-radius: 36rpx;
		padding: 0 28rpx;
		box-shadow: 0 4rpx 16rpx rgba(34, 168, 79, 0.2);
	}

	.search-placeholder {
		margin-left: 12rpx;
		font-size: 26rpx;
		color: #999999;
	}

	/* 页面滚动区 */
	.home-scroll {
		background: #F5F7F4;
	}

	/* 轮播图 */
	.banner-wrap {
		padding: 20rpx 24rpx 0;
	}

	/* 分类网格 */
	.category-card {
		margin: 20rpx 24rpx;
		background: #FFFFFF;
		border-radius: 20rpx;
		padding: 28rpx 16rpx;
		box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
	}

	.category-grid {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 16rpx;
	}

	.category-item {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 12rpx;
	}

	.category-icon-wrap {
		width: 96rpx;
		height: 96rpx;
		background: #F5F7F4;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		overflow: hidden;
	}

	.category-title {
		font-size: 24rpx;
		color: #1A1A1A;
		text-align: center;
	}

	/* 商品 Tab */
	.goods-tabs-wrap {
		padding: 0 24rpx;
	}

	.goods-tabs {
		display: flex;
		background: #FFFFFF;
		border-radius: 16rpx;
		padding: 0 8rpx;
		box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);
	}

	.tab-item {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 24rpx 0 16rpx;
		position: relative;
	}

	.tab-text {
		font-size: 28rpx;
		font-weight: 500;
		color: #999999;
		transition: all 0.2s;
	}

	.tab-item.active .tab-text {
		color: #22A84F;
		font-weight: 600;
	}

	.tab-indicator {
		position: absolute;
		bottom: 0;
		left: 50%;
		transform: translateX(-50%);
		width: 48rpx;
		height: 6rpx;
		background: linear-gradient(90deg, #22A84F, #1A9A45);
		border-radius: 3rpx;
	}

	/* 商品列表 */
	.goods-list-wrap {
		padding: 20rpx 24rpx;
	}

	.load-more {
		padding: 20rpx 0;
	}

	.empty-tip {
		padding: 60rpx 0;
	}

	/* 弹窗内容 */
	.modal-content {
		padding: 32rpx;
	}

	.modal-main {
		font-size: 32rpx;
		font-weight: 600;
		color: #1A1A1A;
		text-align: center;
		margin-bottom: 28rpx;
	}

	.modal-info {
		background: #F5F7F4;
		border-radius: 12rpx;
		padding: 24rpx;
		margin-bottom: 24rpx;
	}

	.info-row {
		font-size: 28rpx;
		color: #666666;
		line-height: 1.8;
	}

	.highlight {
		color: #22A84F;
		font-weight: 600;
		margin: 0 6rpx;
	}

	.modal-tip {
		font-size: 24rpx;
		color: #999999;
		text-align: center;
	}
</style>
