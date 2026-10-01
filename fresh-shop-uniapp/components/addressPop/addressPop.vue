<!--
 * 优诚配运 - 地址选择弹窗
 * 设计规范：自然清新
-->
<template>
	<u-popup :show="show" @close="close" mode="bottom" round="24rpx" :closeable="true">
		<view class="address-popup">
			<view class="popup-header">
				<text class="popup-title">选择收货地址</text>
			</view>

			<view class="popup-content">
				<scroll-view scroll-y="true" class="address-list">
					<view
						class="address-item"
						v-for="item in list"
						:key="item.ID"
						@click="checkedAddress(item)"
					>
						<view class="item-check">
							<view class="check-circle" :class="{ active: item.ID === addressId }">
								<u-icon v-if="item.ID === addressId" name="checkmark" size="24rpx" color="#FFFFFF"></u-icon>
							</view>
						</view>
						<view class="item-info">
							<view class="item-user">
								<text class="user-name">{{ item.userName }}</text>
								<text class="user-phone">{{ item.userPhone }}</text>
								<view class="default-tag" v-if="item.isDefault === 1">默认</view>
							</view>
							<view class="item-address">
								{{ item.address || '' }}{{ item.detailAddress || '' }}{{ item.lableName || '' }}
							</view>
						</view>
					</view>

					<view class="empty-tip" v-if="list.length === 0">
						<u-icon name="map" size="80rpx" color="#CCCCCC"></u-icon>
						<text>暂无收货地址</text>
					</view>
				</scroll-view>
			</view>

			<view class="popup-footer">
				<view class="add-btn" @click="toCreateAddress">
					<u-icon name="plus" size="32rpx" color="#FFFFFF"></u-icon>
					<text>添加新地址</text>
				</view>
			</view>
		</view>

		<u-toast ref="toast" style="z-index: 9999;"></u-toast>
	</u-popup>
</template>

<script>
	import { getAddressList } from '@/api/address.js'

	export default {
		name: 'addressPop',
		data() {
			return {
				list: []
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
			},
			addressId: {
				type: Number,
				default: null
			}
		},
		watch: {
			show(val) {
				if (val) {
					this.getAddressListData()
				}
			}
		},
		methods: {
			async getAddressListData() {
				try {
					const res = await getAddressList()
					if (res.code === 0) {
						this.list = res.data || []
					}
				} catch (e) {
					this.$message(this.$refs.toast).error('加载失败')
				}
			},
			checkedAddress(addressInfo) {
				this.$emit('checked', addressInfo)
			},
			close() {
				this.$emit('close')
			},
			toCreateAddress() {
				this.close()
				uni.navigateTo({
					url: '/pages/address/addressForm?submit=1'
				})
			}
		}
	}
</script>

<style lang="scss" scoped>
	.address-popup {
		display: flex;
		flex-direction: column;
		max-height: 70vh;
	}

	.popup-header {
		padding: 32rpx;
		border-bottom: 1rpx solid #EEEEEE;
		text-align: center;
	}

	.popup-title {
		font-size: 32rpx;
		font-weight: 600;
		color: #1A1A1A;
	}

	.popup-content {
		flex: 1;
		overflow: hidden;
	}

	.address-list {
		height: 50vh;
		padding: 0 24rpx;
	}

	.address-item {
		display: flex;
		align-items: flex-start;
		padding: 28rpx 0;
		border-bottom: 1rpx solid #F0F0F0;

		&:last-child {
			border-bottom: none;
		}
	}

	.item-check {
		padding-top: 4rpx;
		margin-right: 20rpx;
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

	.item-info {
		flex: 1;
	}

	.item-user {
		display: flex;
		align-items: center;
		gap: 12rpx;
		margin-bottom: 12rpx;
	}

	.user-name {
		font-size: 28rpx;
		font-weight: 600;
		color: #1A1A1A;
	}

	.user-phone {
		font-size: 26rpx;
		color: #666666;
	}

	.default-tag {
		background: #E8F8EC;
		color: #22A84F;
		font-size: 20rpx;
		padding: 4rpx 12rpx;
		border-radius: 8rpx;
	}

	.item-address {
		font-size: 26rpx;
		color: #666666;
		line-height: 1.5;
	}

	.empty-tip {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding-top: 100rpx;
		gap: 20rpx;

		text {
			font-size: 28rpx;
			color: #999999;
		}
	}

	.popup-footer {
		padding: 24rpx 32rpx;
		padding-bottom: calc(env(safe-area-inset-bottom) + 24rpx);
		border-top: 1rpx solid #EEEEEE;
	}

	.add-btn {
		height: 88rpx;
		background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
		border-radius: 44rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 12rpx;
		box-shadow: 0 4rpx 20rpx rgba(34, 168, 79, 0.3);

		text {
			font-size: 30rpx;
			font-weight: 600;
			color: #FFFFFF;
		}
	}
</style>
