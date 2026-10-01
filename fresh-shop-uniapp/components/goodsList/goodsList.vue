<!--
 * 优诚配运 - 商品卡片组件
 * 设计规范：自然清新
-->
<template>
	<view>
		<!-- 网格模式（首页） -->
		<view class="goods-grid" v-if="!vertical">
			<view
				class="goods-card"
				v-for="(item, index) in lists"
				:key="index"
				@click="goodsClick(item)"
			>
				<!-- 图片 -->
				<view class="card-image-wrap">
					<image
						v-if="item.images && item.images[0]"
						class="card-image"
						:src="item.images[0].url"
						mode="aspectFill"
					></image>
					<image v-else class="card-image" src="/static/nopicture.jpg" mode="aspectFill"></image>
					<!-- 缺货遮罩 -->
					<view class="stock-mask" v-if="item.store <= 0">
						<text class="stock-text">补货中</text>
					</view>
					<!-- 已售标签 -->
					<view class="sale-tag" v-if="item.sale > 0">
						<text>已售 {{ item.sale }}</text>
					</view>
				</view>

				<!-- 内容 -->
				<view class="card-body">
					<text class="goods-name">{{ item.name }}</text>

					<view class="goods-meta" v-if="item.unit">
						<text class="spec">规格：{{ item.weight ? item.weight + 'g/' : '' }}{{ item.unit }}</text>
					</view>

					<view class="card-footer" v-if="isAudit">
						<view class="price-wrap">
							<text class="price-symbol">¥</text>
							<text class="price">{{ item.price > 0 && item.price < item.costPrice ? item.price : item.costPrice || '0' }}</text>
							<text class="price-original" v-if="item.price > 0 && item.price < item.costPrice">
								¥{{ item.costPrice }}
							</text>
						</view>
					</view>
				</view>
			</view>
		</view>

		<!-- 列表模式（分类页） -->
		<view class="goods-list" v-else>
			<view
				class="goods-row"
				v-for="(item, index) in lists"
				:key="index"
				@click="goodsClick(item)"
			>
				<!-- 图片 -->
				<view class="row-image-wrap">
					<image
						v-if="item.images && item.images[0]"
						class="row-image"
						:src="item.images[0].url"
						mode="aspectFill"
					></image>
					<image v-else class="row-image" src="/static/nopicture.jpg" mode="aspectFill"></image>
					<view class="stock-badge" v-if="item.store <= 0">
						<text>补货中</text>
					</view>
				</view>

				<!-- 内容 -->
				<view class="row-body">
					<text class="row-name">{{ item.name }}</text>

					<view class="row-meta">
						<text class="spec">规格：{{ item.weight ? item.weight + 'g/' : '' }}{{ item.unit }}</text>
						<text class="stock">库存：{{ item.store >= item.minCount ? item.store : 0 }} {{ item.unit }}</text>
					</view>

					<view class="row-footer" v-if="isAudit">
						<view class="price-wrap">
							<text class="price-symbol">¥</text>
							<text class="price">{{ item.price > 0 && item.price < item.costPrice ? item.price : item.costPrice || '0' }}</text>
							<text class="price-unit">/{{ item.unit }}</text>
							<text class="price-original" v-if="item.price > 0 && item.price < item.costPrice">
								¥{{ item.costPrice }}
							</text>
						</view>
						<view class="sale-tag" v-if="item.sale > 0">
							<text>已售 {{ item.sale }}</text>
						</view>
					</view>
				</view>
			</view>
		</view>
	</view>
</template>

<script>
export default {
	props: {
		lists: { type: Array, required: true, default: () => [] },
		priceType: { type: String, default: '¥' },
		vertical: { type: Boolean, default: false },
		disableJump: { type: Boolean, default: false },
		isPoint: { type: Boolean, default: false },
		isAudit: { type: Boolean, default: false },
		isAddCart: { type: Boolean, default: false },
		imgWidth: { type: String, default: '' },
		imgHeight: { type: String, default: '' },
		showPayCount: { type: Boolean, default: false }
	},
	methods: {
		goodsClick(goods) {
			if (this.disableJump) {
				this.$emit('onGoods', goods)
			} else {
				uni.navigateTo({ url: `/pages/goods/detail?id=${goods.ID}` })
			}
		},
		goodsLongClick(goods) {
			this.$emit('onGoodsLongClick', goods)
		},
		updateCart(index, cardId, num) {
			this.$emit('updateCart', index, cardId, num)
		}
	}
}
</script>

<style lang="scss" scoped>
/* ========================
   网格模式（首页两列）
   ======================== */
.goods-grid {
	display: grid;
	grid-template-columns: repeat(2, 1fr);
	gap: 20rpx;
	padding: 4rpx;
}

.goods-card {
	background: #FFFFFF;
	border-radius: 20rpx;
	overflow: hidden;
	box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.06);
	transition: all 0.2s;

	&:active {
		transform: scale(0.98);
		box-shadow: 0 1rpx 6rpx rgba(0, 0, 0, 0.08);
	}
}

.card-image-wrap {
	position: relative;
	width: 100%;
	aspect-ratio: 1;
}

.card-image {
	width: 100%;
	height: 100%;
	border-radius: 20rpx 20rpx 0 0;
}

.stock-mask {
	position: absolute;
	inset: 0;
	background: rgba(0, 0, 0, 0.45);
	display: flex;
	align-items: center;
	justify-content: center;
	border-radius: 20rpx 20rpx 0 0;
}

.stock-text {
	font-size: 26rpx;
	color: #FFFFFF;
	background: rgba(0, 0, 0, 0.5);
	padding: 8rpx 24rpx;
	border-radius: 24rpx;
}

.sale-tag {
	position: absolute;
	top: 12rpx;
	right: 12rpx;
	background: rgba(239, 68, 68, 0.9);
	padding: 4rpx 12rpx;
	border-radius: 8rpx;

	text {
		font-size: 20rpx;
		color: #FFFFFF;
	}
}

.card-body {
	padding: 16rpx 20rpx 20rpx;
}

.goods-name {
	font-size: 28rpx;
	font-weight: 500;
	color: #1A1A1A;
	line-height: 1.4;
	display: -webkit-box;
	-webkit-box-orient: vertical;
	-webkit-line-clamp: 2;
	overflow: hidden;
}

.goods-meta {
	margin-top: 8rpx;
}

.spec {
	font-size: 22rpx;
	color: #999999;
}

.card-footer {
	margin-top: 12rpx;
}

.price-wrap {
	display: flex;
	align-items: baseline;
	flex-wrap: wrap;
}

.price-symbol {
	font-size: 24rpx;
	color: #F97316;
	font-weight: 600;
}

.price {
	font-size: 36rpx;
	color: #F97316;
	font-weight: 700;
	line-height: 1;
}

.price-unit {
	font-size: 22rpx;
	color: #666666;
}

.price-original {
	font-size: 22rpx;
	color: #CCCCCC;
	text-decoration: line-through;
	margin-left: 8rpx;
}

/* ========================
   列表模式（分类页单列）
   ======================== */
.goods-list {
	display: flex;
	flex-direction: column;
	gap: 16rpx;
	padding: 4rpx;
}

.goods-row {
	display: flex;
	background: #FFFFFF;
	border-radius: 20rpx;
	overflow: hidden;
	box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.06);
	transition: all 0.2s;
	padding: 16rpx;

	&:active {
		transform: scale(0.99);
		background: #FAFAFA;
	}
}

.row-image-wrap {
	position: relative;
	flex-shrink: 0;
	width: 180rpx;
	height: 180rpx;
	border-radius: 16rpx;
	overflow: hidden;
}

.row-image {
	width: 100%;
	height: 100%;
}

.stock-badge {
	position: absolute;
	bottom: 0;
	left: 0;
	right: 0;
	background: rgba(0, 0, 0, 0.5);
	padding: 4rpx 0;
	text-align: center;

	text {
		font-size: 20rpx;
		color: #FFFFFF;
	}
}

.row-body {
	flex: 1;
	margin-left: 20rpx;
	display: flex;
	flex-direction: column;
	justify-content: space-between;
}

.row-name {
	font-size: 30rpx;
	font-weight: 600;
	color: #1A1A1A;
	line-height: 1.4;
	display: -webkit-box;
	-webkit-box-orient: vertical;
	-webkit-line-clamp: 2;
	overflow: hidden;
}

.row-meta {
	display: flex;
	flex-direction: column;
	gap: 6rpx;
	margin-top: 8rpx;
}

.spec {
	font-size: 24rpx;
	color: #999999;
}

.stock {
	font-size: 24rpx;
	color: #666666;
}

.row-footer {
	display: flex;
	align-items: center;
	justify-content: space-between;
	margin-top: auto;
}

.row-footer .price-wrap {
	flex: 1;
}
</style>
