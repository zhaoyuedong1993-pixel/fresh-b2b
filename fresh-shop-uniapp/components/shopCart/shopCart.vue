<!--
 * 优诚配运 - 购物车组件
 * 设计规范：自然清新
-->
<template>
	<view class="cart-component">
		<!-- 空购物车 -->
		<view class="empty-cart" v-if="list.length === 0">
			<view class="empty-icon">
				<u-icon name="shopping-cart" size="100" color="#CCCCCC"></u-icon>
			</view>
			<view class="empty-text">购物车为空</view>
			<view class="empty-btn" @click="toHome">去逛逛</view>
		</view>

		<!-- 购物车列表 -->
		<view class="cart-list" v-else>
			<scroll-view scroll-y="true" :style="{ height: height + 'px' }" refresher-enabled="true"
				:refresher-triggered="triggered" @refresherrefresh="onRefresh" :scroll-anchoring="true">
				<view class="cart-items">
					<view class="cart-item" v-for="(cart, index) in list" :key="index" @longpress="showDeleteDialog(index)">
						<!-- 选择框 -->
						<view class="item-check" :class="{ disabled: cart.goods.store <= 0 || cart.goods.store < cart.num }"
							@click.stop="checkedGoods(cart.ID, cart.checked, index)">
							<view class="check-circle" :class="{ active: cart.checked === 1 }">
								<u-icon v-if="cart.checked === 1" name="checkmark" size="20rpx" color="#FFFFFF"></u-icon>
							</view>
						</view>

						<!-- 商品内容 -->
						<view class="item-content" @click.stop="toGoodsDetail(cart.goodsId)">
							<view class="item-image">
								<image v-if="cart.goods.images && cart.goods.images.length > 0"
									:src="cart.goods.images[0].url" mode="aspectFill"></image>
								<image v-else src="/static/nopicture.jpg" mode="aspectFill"></image>
								<view class="stock-badge" v-if="cart.goods.store <= 0">补货中</view>
							</view>

							<view class="item-info">
								<view class="info-header">
									<text class="goods-name">{{ cart.goods.name }}</text>
									<view class="goods-tags">
										<text class="tag unit">{{ cart.goods.unit }}</text>
										<text class="tag stock" :class="{ low: cart.goods.store < 10 }">库存 {{ cart.goods.store }}</text>
									</view>
								</view>

								<view class="info-footer">
									<view class="price-box">
										<text class="price-symbol">¥</text>
										<text class="price-value">{{ formatPrice(getPrice(cart)) }}</text>
									</view>

									<!-- 数量控制 -->
									<view class="quantity-control">
										<view class="qty-btn minus" :class="{ disabled: cart.num <= 1 }"
											@click.stop="updateNum(cart, index, cart.num - 1)">
											<u-icon name="minus" size="20rpx" color="#22A84F"></u-icon>
										</view>
										<text class="qty-num">{{ cart.num }}</text>
										<view class="qty-btn plus" :class="{ disabled: cart.num >= cart.goods.store }"
											@click.stop="updateNum(cart, index, cart.num + 1)">
											<u-icon name="plus" size="20rpx" color="#FFFFFF"></u-icon>
										</view>
									</view>
								</view>
							</view>
						</view>
					</view>
				</view>
			</scroll-view>

			<!-- 底部结算栏 -->
			<view class="checkout-bar">
				<view class="bar-left">
					<view class="select-all" @click="allCheck">
						<view class="check-circle" :class="{ active: isCheckAll }">
							<u-icon v-if="isCheckAll" name="checkmark" size="20rpx" color="#FFFFFF"></u-icon>
						</view>
						<text class="select-text">全选</text>
					</view>
				</view>

				<view class="bar-right">
					<view class="price-info">
						<text class="price-label">合计：</text>
						<text class="total-price">
							<text class="price-symbol">¥</text>
							<text class="price-value">{{ total }}</text>
						</text>
					</view>
					<view class="checkout-btn" :class="{ disabled: selectedCount === 0 }" @click="accounts">
						<text>结算</text>
						<text class="btn-count" v-if="selectedCount > 0">({{ selectedCount }})</text>
					</view>
				</view>
			</view>
		</view>

		<!-- 删除确认弹窗 -->
		<u-modal :show="showDelete" :showCancelButton="true" title="删除商品" @confirm="deleteCart" @cancel="showDelete = false"
			@close="showDelete = false">
			<view class="delete-content">
				<u-icon name="warning-fill" size="48" color="#F97316"></u-icon>
				<text class="delete-text">确定要删除这个商品吗？</text>
			</view>
		</u-modal>

		<u-toast ref="toast" style="z-index: 9999"></u-toast>
	</view>
</template>

<script>
	import { addCart, updateCart, selectAllCart, clearSelectAllCart, deleteCartByIds } from '@/api/cart.js'

	export default {
		name: 'shopCart',
		data() {
			return {
				isCheckAll: false,
				total: '0.00',
				selectedCount: 0,
				showDelete: false,
				currentDeleteIndex: 0
			}
		},
		props: {
			list: {
				type: Array,
				default: () => []
			},
			height: {
				type: Number,
				default: 0
			},
			triggered: {
				type: Boolean,
				default: false
			}
		},
		watch: {
			list: {
				handler() {
					this.updateStats()
				},
				deep: true,
				immediate: true
			}
		},
		methods: {
			getPrice(cart) {
				const goods = cart.goods
				if (goods.price > 0 && goods.price < goods.costPrice) {
					return goods.price
				}
				return goods.costPrice || 0
			},
			formatPrice(price) {
				return (Number(price) || 0).toFixed(2)
			},
			updateStats() {
				let total = 0
				let count = 0
				let allChecked = true

				this.list.forEach(item => {
					if (item.checked === 1) {
						total += this.getPrice(item) * item.num
						count += item.num
					} else {
						allChecked = false
					}
				})

				this.isCheckAll = allChecked && this.list.length > 0
				this.total = total.toFixed(2)
				this.selectedCount = count
			},
			async checkedGoods(id, checked, index) {
				const newChecked = checked === 1 ? 0 : 1
				this.list[index].checked = newChecked
				this.updateStats()

				const res = await updateCart({ ID: id, checked: newChecked })
				if (res.code !== 0) {
					this.list[index].checked = checked
					this.updateStats()
				}
			},
			async allCheck() {
				if (this.isCheckAll) {
					// 取消全选
					this.list.forEach(item => { item.checked = 0 })
					this.updateStats()
					await clearSelectAllCart()
				} else {
					// 全选
					this.list.forEach(item => {
						if (item.goods.store > 0 && item.goods.store >= item.num) {
							item.checked = 1
						}
					})
					this.updateStats()
					await selectAllCart()
				}
			},
			async updateNum(cart, index, num) {
				if (num < cart.goods.minCount) {
					uni.showToast({ title: `商品最低购买${cart.goods.minCount}件`, icon: 'none' })
					return
				}
				if (num < 1) num = 0

				const originalNum = cart.num
				cart.num = num
				this.updateStats()

				const res = await addCart({
					goodsId: cart.goods.ID,
					specType: 0,
					num: num
				})

				if (res.code !== 0) {
					cart.num = originalNum
					this.updateStats()
				}
			},
			showDeleteDialog(index) {
				this.currentDeleteIndex = index
				this.showDelete = true
			},
			async deleteCart() {
				const item = this.list[this.currentDeleteIndex]
				this.list.splice(this.currentDeleteIndex, 1)
				this.updateStats()

				const res = await deleteCartByIds({ ids: [item.ID] })
				if (res.code !== 0) {
					this.list.splice(this.currentDeleteIndex, 0, item)
					this.updateStats()
					uni.showToast({ title: '删除失败', icon: 'none' })
				}
			},
			accounts() {
				if (this.selectedCount === 0) return
				uni.navigateTo({ url: '/pages/order/submit' })
			},
			onRefresh() {
				this.$emit('onRefresh')
			},
			toHome() {
				uni.switchTab({ url: '/pages/index/index' })
			},
			toGoodsDetail(id) {
				uni.navigateTo({ url: `/pages/goods/detail?id=${id}` })
			}
		}
	}
</script>

<style lang="scss" scoped>
	.cart-component {
		width: 100%;
		height: 100%;
	}

	.empty-cart {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding-top: 120rpx;
		gap: 24rpx;
	}

	.empty-icon {
		width: 200rpx;
		height: 200rpx;
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

	.empty-btn {
		background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
		color: #FFFFFF;
		font-size: 28rpx;
		padding: 20rpx 60rpx;
		border-radius: 40rpx;
		margin-top: 20rpx;
	}

	.cart-list {
		display: flex;
		flex-direction: column;
		height: 100%;
	}

	.cart-items {
		padding: 16rpx 20rpx;
		padding-bottom: 140rpx;
	}

	.cart-item {
		display: flex;
		align-items: flex-start;
		background: #FFFFFF;
		border-radius: 20rpx;
		padding: 20rpx;
		margin-bottom: 16rpx;
		box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.06);
	}

	.item-check {
		padding: 40rpx 16rpx 0 0;

		&.disabled {
			opacity: 0.4;
		}
	}

	.check-circle {
		width: 40rpx;
		height: 40rpx;
		border-radius: 50%;
		border: 2rpx solid #DDDDDD;
		display: flex;
		align-items: center;
		justify-content: center;

		&.active {
			background: #22A84F;
			border-color: #22A84F;
		}
	}

	.item-content {
		flex: 1;
		display: flex;
		gap: 16rpx;
	}

	.item-image {
		position: relative;
		width: 160rpx;
		height: 160rpx;
		border-radius: 16rpx;
		overflow: hidden;
		flex-shrink: 0;

		image {
			width: 100%;
			height: 100%;
		}
	}

	.stock-badge {
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.6);
		display: flex;
		align-items: center;
		justify-content: center;

		font-size: 22rpx;
		color: #FFFFFF;
	}

	.item-info {
		flex: 1;
		display: flex;
		flex-direction: column;
		justify-content: space-between;
	}

	.info-header {
		display: flex;
		flex-direction: column;
		gap: 8rpx;
	}

	.goods-name {
		font-size: 28rpx;
		font-weight: 600;
		color: #1A1A1A;
		line-height: 1.4;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}

	.goods-tags {
		display: flex;
		gap: 8rpx;
	}

	.tag {
		font-size: 20rpx;
		padding: 4rpx 12rpx;
		border-radius: 8rpx;

		&.unit {
			background: #F5F7F4;
			color: #666666;
		}

		&.stock {
			background: #E8F8EC;
			color: #22A84F;

			&.low {
				background: #FEF3E2;
				color: #F97316;
			}
		}
	}

	.info-footer {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.price-box {
		display: flex;
		align-items: baseline;
	}

	.price-symbol {
		font-size: 24rpx;
		color: #F97316;
		font-weight: 600;
	}

	.price-value {
		font-size: 32rpx;
		color: #F97316;
		font-weight: 700;
	}

	.quantity-control {
		display: flex;
		align-items: center;
		gap: 12rpx;
	}

	.qty-btn {
		width: 48rpx;
		height: 48rpx;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;

		&.minus {
			background: #E8F8EC;
			border: 2rpx solid #22A84F;
		}

		&.plus {
			background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
		}

		&.disabled {
			opacity: 0.4;
		}
	}

	.qty-num {
		font-size: 30rpx;
		font-weight: 600;
		color: #1A1A1A;
		min-width: 48rpx;
		text-align: center;
	}

	.checkout-bar {
		position: fixed;
		bottom: 100rpx;
		left: 0;
		right: 0;
		height: 100rpx;
		background: #FFFFFF;
		box-shadow: 0 -2rpx 20rpx rgba(0, 0, 0, 0.06);
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0 24rpx;
		z-index: 100;
	}

	.bar-left {
		display: flex;
		align-items: center;
	}

	.select-all {
		display: flex;
		align-items: center;
		gap: 12rpx;
	}

	.select-text {
		font-size: 26rpx;
		color: #666666;
	}

	.bar-right {
		display: flex;
		align-items: center;
		gap: 20rpx;
	}

	.price-info {
		display: flex;
		align-items: baseline;
		gap: 4rpx;
	}

	.price-label {
		font-size: 26rpx;
		color: #666666;
	}

	.total-price {
		display: flex;
		align-items: baseline;
	}

	.checkout-btn {
		height: 72rpx;
		padding: 0 40rpx;
		background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
		border-radius: 36rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 4rpx;
		box-shadow: 0 4rpx 16rpx rgba(34, 168, 79, 0.3);

		text {
			font-size: 28rpx;
			font-weight: 600;
			color: #FFFFFF;
		}

		&.disabled {
			opacity: 0.5;
		}
	}

	.delete-content {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 20rpx 0;
		gap: 16rpx;
	}

	.delete-text {
		font-size: 28rpx;
		color: #666666;
	}
</style>
