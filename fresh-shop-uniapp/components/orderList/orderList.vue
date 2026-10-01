<!--
 * 优诚配运 - 订单列表组件
 * 设计规范：自然清新
-->
<template>
	<view class="order-list-component">
		<scroll-view scroll-y="true" :style="{ height: scrollViewHeight + 'px' }" refresher-enabled="true"
			:refresher-threshold="70" :refresher-triggered="triggered" @refresherrefresh="onRefresh"
			@scrolltolower="scrollTolower" :scroll-anchoring="true">
			<view class="order-items" v-if="list.length > 0">
				<view class="order-card" v-for="(order, index) in list" :key="index" @click="toOrderDetail(order.ID)">
					<!-- 订单头部 -->
					<view class="order-header">
						<view class="header-left">
							<text class="order-sn">订单号：{{ order.orderSn || order.order_sn }}</text>
							<text class="order-time">{{ formatDate(order.createdAt || order.created_at) }}</text>
						</view>
						<view class="order-status" :class="'status-' + order.status">
							{{ getStatusText(order) }}
						</view>
					</view>

					<!-- 商品展示 -->
					<view class="order-goods">
						<!-- 多商品展示 -->
						<view class="goods-grid" v-if="order.details && order.details.length > 1">
							<view class="goods-thumbs">
								<image v-for="(item, idx) in order.details.slice(0, 4)" :key="idx"
									:src="item.goodsImage || '/static/nopicture.jpg'" mode="aspectFill"
									class="thumb-img"></image>
								<view class="thumb-more" v-if="order.details.length > 4">
									+{{ order.details.length - 4 }}
								</view>
							</view>
						</view>

						<!-- 单商品展示 -->
						<view class="goods-single" v-else-if="order.details && order.details.length === 1">
							<image class="single-img" :src="order.details[0].goodsImage || '/static/nopicture.jpg'"
								mode="aspectFill"></image>
							<view class="single-info">
								<text class="goods-name">{{ order.details[0].goodsName }}</text>
								<view class="goods-spec">
									<text>单价：
										<text class="price">¥{{ order.details[0].price }}</text>
									</text>
									<text class="num">x{{ order.details[0].num }}</text>
								</view>
							</view>
						</view>
					</view>

					<!-- 订单底部 -->
					<view class="order-footer">
						<view class="footer-left">
							<text class="goods-count">共 {{ order.num }} 件商品</text>
							<view class="order-amount">
								<text class="amount-label" v-if="order.status !== 0">实付</text>
								<text class="amount-value">
									<text class="symbol" v-if="order.goodsArea === 0">¥</text>
									<text class="symbol" v-else>积分 </text>
									{{ formatAmount(order.status !== 0 ? order.finish : order.total) }}
								</text>
							</view>
						</view>
						<view class="footer-right">
							<view class="settlement-tag" :class="'settle-' + order.settlementType" v-if="order.status !== 0">
								{{ getSettlementText(order) }}
							</view>
						</view>
					</view>

					<!-- 操作按钮 -->
					<view class="order-actions" v-if="order.statusCancel === 0">
						<view class="action-btn outline" v-if="order.status === 0" @click.stop="cancelOrder(order.ID)">
							取消订单
						</view>
						<view class="action-btn outline" v-if="order.status === 2" @click.stop="confirmOrder(order.ID)">
							确认收货
						</view>
						<view class="action-btn danger" v-if="order.status === 3 || order.status === 4"
							@click.stop="deleteOrder(order.ID)">
							删除订单
						</view>
						<view class="action-btn primary" v-if="order.status === 0" @click.stop="payOrder(order.ID)">
							立即付款
						</view>
					</view>
				</view>

				<view class="load-more">
					<u-loadmore :status="loadMore" nomore-text="没有更多订单了" />
				</view>
			</view>

			<!-- 空状态 -->
			<view class="empty-order" v-else>
				<view class="empty-icon">
					<u-icon name="order" size="100" color="#CCCCCC"></u-icon>
				</view>
				<view class="empty-text">暂无相关订单</view>
			</view>
		</scroll-view>

		<!-- 操作弹窗 -->
		<u-modal :show="showModal" :title="modalTitle" :showCancelButton="true" @confirm="confirmAction" @cancel="closeModal"
			@close="closeModal">
			<view class="modal-content">
				<text>{{ modalContent }}</text>
			</view>
		</u-modal>

		<u-toast ref="toast" style="z-index: 9999"></u-toast>
	</view>
</template>

<script>
	import config from '@/config/config.js'
	import { getOrderList, confirmOrder, cancelOrder, deleteOrder as deleteOrderApi, orderPay } from '@/api/order.js'

	export default {
		name: 'orderList',
		data() {
			return {
				list: [],
				scrollViewHeight: 600,
				page: { page: 1, pageSize: 10, total: 0, isMore: true },
				loadMore: 'loadmore',
				showModal: false,
				modalTitle: '',
				modalContent: '',
				currentAction: '',
				currentOrderId: 0
			}
		},
		props: {
			status: { type: [String, Number], default: null }
		},
		computed: {
			triggered() { return false }
		},
		mounted() {
			uni.getSystemInfo({
				success: (res) => {
					this.scrollViewHeight = res.windowHeight - 90
				}
			})
			this.getOrderListData(0)
		},
		methods: {
			async getOrderListData(type = 0) {
				if (type === 0) {
					this.page.page = 1
					this.page.isMore = true
					this.loadMore = 'loadmore'
				}

				const data = { page: this.page.page, pageSize: this.page.pageSize }
				if (this.status !== null && this.status !== 'null') {
					data.status = parseInt(this.status)
				}

				const res = await getOrderList(data, type === 0)
				if (res.code !== 0) return false

				const list = res.data?.list || []
				list.forEach(item => {
					if (item.details) {
						item.details.forEach((d, i) => {
							if (d.goodsImage && d.goodsImage.slice(0, 4) !== 'http') {
								item.details[i].goodsImage = config.baseUrl + '/' + d.goodsImage
							}
						})
					}
				})

				this.page.total = res.data?.total || 0
				this.page.isMore = this.page.page * this.page.pageSize < this.page.total
				this.loadMore = this.page.isMore ? 'loadmore' : 'nomore'

				if (type === 0) {
					this.list = list
				} else {
					this.list = [...this.list, ...list]
				}

				this.page.page++
				return true
			},
			async onRefresh() {
				const b = await this.getOrderListData(0)
				this.$message(this.$refs.toast).success(b ? '刷新成功' : '刷新失败')
			},
			async scrollTolower() {
				if (this.loadMore === 'loading' || !this.page.isMore) return
				this.loadMore = 'loading'
				await this.getOrderListData(1)
				this.loadMore = this.page.isMore ? 'loadmore' : 'nomore'
			},
			getStatusText(order) {
				if (order.statusCancel > 0) {
					const map = { 1: '已取消', 2: '后台取消', 3: '超时取消' }
					return map[order.statusCancel] || '已取消'
				}
				const map = { 0: '待付款', 1: '备货中', 2: '配送中', 3: '已完成' }
				return map[order.status] || '未知'
			},
			getSettlementText(order) {
				const map = { 0: '实付', 1: '月结未结', 2: '已结清' }
				return map[order.settlementType] || ''
			},
			formatAmount(amount) {
				return (Number(amount) || 0).toFixed(2)
			},
			formatDate(dateStr) {
				if (!dateStr) return '-'
				const d = new Date(dateStr)
				return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
			},
			toOrderDetail(id) {
				uni.navigateTo({ url: `/pages/order/detail?id=${id}` })
			},
			cancelOrder(id) {
				this.modalTitle = '取消订单'
				this.modalContent = '确定要取消该订单吗？'
				this.currentAction = 'cancel'
				this.currentOrderId = id
				this.showModal = true
			},
			confirmOrder(id) {
				this.modalTitle = '确认收货'
				this.modalContent = '请确认货物是否完整！'
				this.currentAction = 'confirm'
				this.currentOrderId = id
				this.showModal = true
			},
			deleteOrder(id) {
				this.modalTitle = '删除订单'
				this.modalContent = '确定要删除该订单吗？'
				this.currentAction = 'delete'
				this.currentOrderId = id
				this.showModal = true
			},
			async confirmAction() {
				try {
					let res
					switch (this.currentAction) {
						case 'cancel':
							res = await cancelOrder({ ID: this.currentOrderId })
							break
						case 'confirm':
							res = await confirmOrder({ ID: this.currentOrderId })
							break
						case 'delete':
							res = await deleteOrderApi({ ID: this.currentOrderId })
							break
					}
					if (res.code === 0) {
						this.$message(this.$refs.toast).success('操作成功')
						this.getOrderListData(0)
					} else {
						this.$message(this.$refs.toast).error(res.msg || '操作失败')
					}
				} catch (e) {
					this.$message(this.$refs.toast).error('操作失败')
				}
				this.closeModal()
			},
			closeModal() {
				this.showModal = false
				this.currentAction = ''
				this.currentOrderId = 0
			},
			async payOrder(id) {
				const res = await orderPay({ ID: id })
				if (res.code === 0) {
					this.$message(this.$refs.toast).success('订单已提交')
					this.getOrderListData(0)
				}
			}
		}
	}
	</script>

<style lang="scss" scoped>
	.order-list-component {
		width: 100%;
		height: 100%;
	}

	.order-items {
		padding: 20rpx;
	}

	.order-card {
		background: #FFFFFF;
		border-radius: 20rpx;
		padding: 24rpx;
		margin-bottom: 20rpx;
		box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.06);
	}

	.order-header {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		margin-bottom: 20rpx;
	}

	.header-left {
		display: flex;
		flex-direction: column;
		gap: 6rpx;
	}

	.order-sn {
		font-size: 26rpx;
		color: #1A1A1A;
	}

	.order-time {
		font-size: 22rpx;
		color: #999999;
	}

	.order-status {
		font-size: 26rpx;
		font-weight: 600;

		&.status-0 { color: #F97316; }
		&.status-1 { color: #22A84F; }
		&.status-2 { color: #3B82F6; }
		&.status-3 { color: #666666; }
	}

	.order-goods {
		margin-bottom: 20rpx;
	}

	.goods-grid {
		overflow: hidden;
	}

	.goods-thumbs {
		display: flex;
		gap: 8rpx;
		flex-wrap: wrap;
	}

	.thumb-img {
		width: 140rpx;
		height: 140rpx;
		border-radius: 12rpx;
		background: #F5F7F4;
	}

	.thumb-more {
		width: 140rpx;
		height: 140rpx;
		border-radius: 12rpx;
		background: #F5F7F4;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 28rpx;
		color: #999999;
	}

	.goods-single {
		display: flex;
		gap: 16rpx;
	}

	.single-img {
		width: 160rpx;
		height: 160rpx;
		border-radius: 16rpx;
		background: #F5F7F4;
		flex-shrink: 0;
	}

	.single-info {
		flex: 1;
		display: flex;
		flex-direction: column;
		justify-content: space-between;
	}

	.goods-name {
		font-size: 28rpx;
		color: #1A1A1A;
		line-height: 1.4;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}

	.goods-spec {
		display: flex;
		align-items: center;
		gap: 16rpx;
		font-size: 24rpx;
		color: #666666;

		.price {
			color: #F97316;
			font-weight: 600;
		}

		.num {
			color: #999999;
		}
	}

	.order-footer {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding-top: 20rpx;
		border-top: 1rpx solid #F0F0F0;
	}

	.footer-left {
		display: flex;
		align-items: center;
		gap: 20rpx;
	}

	.goods-count {
		font-size: 24rpx;
		color: #999999;
	}

	.order-amount {
		display: flex;
		align-items: baseline;
		gap: 4rpx;
	}

	.amount-label {
		font-size: 24rpx;
		color: #666666;
	}

	.amount-value {
		font-size: 28rpx;
		font-weight: 600;
		color: #F97316;
	}

	.symbol {
		font-size: 22rpx;
	}

	.settlement-tag {
		font-size: 22rpx;
		padding: 6rpx 16rpx;
		border-radius: 20rpx;

		&.settle-0 {
			background: #E8F8EC;
			color: #22A84F;
		}

		&.settle-1 {
			background: #FEF3E2;
			color: #F97316;
		}

		&.settle-2 {
			background: #F5F7F4;
			color: #666666;
		}
	}

	.order-actions {
		display: flex;
		align-items: center;
		justify-content: flex-end;
		gap: 16rpx;
		margin-top: 20rpx;
		padding-top: 20rpx;
		border-top: 1rpx solid #F0F0F0;
	}

	.action-btn {
		height: 60rpx;
		padding: 0 28rpx;
		border-radius: 30rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 26rpx;
		font-weight: 500;

		&.outline {
			border: 1rpx solid #DDDDDD;
			color: #666666;
			background: transparent;
		}

		&.primary {
			background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
			color: #FFFFFF;
			box-shadow: 0 4rpx 16rpx rgba(34, 168, 79, 0.3);
		}

		&.danger {
			background: #FEE2E2;
			color: #EF4444;
		}
	}

	.load-more {
		padding: 20rpx 0;
	}

	.empty-order {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding-top: 160rpx;
		gap: 20rpx;
	}

	.empty-icon {
		width: 180rpx;
		height: 180rpx;
		background: #F5F7F4;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.empty-text {
		font-size: 28rpx;
		color: #999999;
	}

	.modal-content {
		text-align: center;
		padding: 20rpx;
		font-size: 28rpx;
		color: #666666;
	}
</style>
