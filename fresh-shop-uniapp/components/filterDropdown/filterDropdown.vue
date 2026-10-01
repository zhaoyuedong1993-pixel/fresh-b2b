<!--
 * 优诚配运 - 筛选下拉组件
 * 设计规范：自然清新
-->
<template>
	<view class="filter-dropdown">
		<view class="filter-bar" :style="{ height: height + 'rpx', backgroundColor: bgcolor }">
			<view
				class="filter-item"
				v-for="(item, index) in listArr"
				:key="index"
				@click="itemClick(index)"
			>
				<view class="item-text" :class="{ active: currentIndex === index && show }">
					<text>{{ getShowText(index) }}</text>
					<u-icon
						name="arrow-down"
						size="20rpx"
						:color="currentIndex === index && show ? '#22A84F' : '#666666'"
						:class="{ 'rotate': currentIndex === index && show }"
					></u-icon>
				</view>
			</view>
		</view>

		<!-- 下拉面板 -->
		<view class="dropdown-panel" v-if="show" :style="{ top: height + 'rpx' }">
			<view class="dropdown-mask" @click="maskClose"></view>
			<view class="dropdown-content">
				<scroll-view scroll-y="true" class="dropdown-list">
					<view
						class="dropdown-item"
						v-for="(opt, optIndex) in currentOptions"
						:key="optIndex"
						@click="subItemClick(optIndex)"
					>
						<text :class="{ active: isOptionSelected(optIndex) }">{{ opt[showTag] }}</text>
						<u-icon v-if="isOptionSelected(optIndex)" name="checkmark" size="28rpx" color="#22A84F"></u-icon>
					</view>
				</scroll-view>

				<view class="dropdown-actions">
					<view class="action-btn reset" @click="clearClick">
						<text>重置</text>
					</view>
					<view class="action-btn confirm" @click="confirmClick">
						<text>确定</text>
					</view>
				</view>
			</view>
		</view>
	</view>
</template>

<script>
	export default {
		name: 'filterDropdown',
		data() {
			return {
				currentIndex: 0,
				show: false,
				selectedIndex: {},
				tempSelected: {},
				updateArr: []
			}
		},
		props: {
			height: {
				type: Number,
				default: 88
			},
			bgcolor: {
				type: String,
				default: '#FFFFFF'
			},
			listArr: {
				type: Array,
				default: () => ['智能排序', '全部分类', '全部品牌']
			},
			itemArr: {
				type: Array,
				default: () => [
					[{ text: '智能排序', value: 0 }, { text: '最新', value: 1 }, { text: '销量', value: 2 }],
					[{ text: '全部分类', value: 0 }],
					[{ text: '全部品牌', value: 0 }]
				]
			},
			showTag: {
				type: String,
				default: 'text'
			}
		},
		computed: {
			currentOptions() {
				return this.itemArr[this.currentIndex] || []
			}
		},
		watch: {
			listArr: {
				immediate: true,
				handler(val) {
					this.updateArr = [...val]
				}
			}
		},
		mounted() {
			this.updateArr = [...this.listArr]
			this.initSelected()
		},
		methods: {
			initSelected() {
				this.listArr.forEach((_, index) => {
					this.selectedIndex[index] = 0
					this.tempSelected[index] = 0
				})
			},
			itemClick(index) {
				if (this.currentIndex === index) {
					this.show = !this.show
				} else {
					this.currentIndex = index
					this.tempSelected = { ...this.selectedIndex }
					this.show = true
				}
			},
			getShowText(index) {
				const opts = this.itemArr[index]
				if (!opts || opts.length === 0) return this.listArr[index]
				const selectedIdx = this.selectedIndex[index] || 0
				return opts[selectedIdx]?.[this.showTag] || this.listArr[index]
			},
			subItemClick(optIndex) {
				this.tempSelected[this.currentIndex] = optIndex
			},
			isOptionSelected(optIndex) {
				return this.tempSelected[this.currentIndex] === optIndex
			},
			confirmClick() {
				this.selectedIndex = { ...this.tempSelected }

				const opts = this.itemArr[this.currentIndex]
				const selectedIdx = this.selectedIndex[this.currentIndex]
				if (opts && opts[selectedIdx]) {
					this.$emit('finish', {
						'$index': this.currentIndex,
						...opts[selectedIdx]
					})
				}

				this.show = false
			},
			clearClick() {
				this.tempSelected[this.currentIndex] = 0
			},
			maskClose() {
				this.show = false
				this.$emit('clear')
			},
			updateTitle(arrIndex, title) {
				this.$set(this.updateArr, arrIndex, title)
			}
		}
	}
</script>

<style lang="scss" scoped>
	.filter-dropdown {
		position: relative;
	}

	.filter-bar {
		display: flex;
		align-items: center;
		position: relative;
		z-index: 100;

		&::after {
			content: '';
			position: absolute;
			bottom: 0;
			left: 0;
			right: 0;
			height: 1rpx;
			background: #EEEEEE;
		}
	}

	.filter-item {
		flex: 1;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.item-text {
		display: flex;
		align-items: center;
		gap: 6rpx;
		font-size: 26rpx;
		color: #666666;

		&.active {
			color: #22A84F;
			font-weight: 600;
		}

		/deep/.rotate {
			transform: rotate(180deg);
		}
	}

	.dropdown-panel {
		position: fixed;
		left: 0;
		right: 0;
		bottom: 0;
		z-index: 99;
	}

	.dropdown-mask {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.5);
	}

	.dropdown-content {
		position: fixed;
		left: 0;
		right: 0;
		bottom: 0;
		background: #FFFFFF;
		border-radius: 24rpx 24rpx 0 0;
		max-height: 60vh;
		display: flex;
		flex-direction: column;
	}

	.dropdown-list {
		flex: 1;
		padding: 16rpx 0;
	}

	.dropdown-item {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 28rpx 40rpx;
		font-size: 28rpx;
		color: #1A1A1A;

		&:active {
			background: #F5F7F4;
		}

		text.active {
			color: #22A84F;
			font-weight: 600;
		}
	}

	.dropdown-actions {
		display: flex;
		gap: 20rpx;
		padding: 24rpx 40rpx;
		padding-bottom: calc(env(safe-area-inset-bottom) + 24rpx);
		border-top: 1rpx solid #EEEEEE;
	}

	.action-btn {
		flex: 1;
		height: 80rpx;
		border-radius: 40rpx;
		display: flex;
		align-items: center;
		justify-content: center;

		text {
			font-size: 28rpx;
			font-weight: 600;
		}

		&.reset {
			background: #F5F5F5;

			text {
				color: #666666;
			}
		}

		&.confirm {
			background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);

			text {
				color: #FFFFFF;
			}
		}
	}
</style>
