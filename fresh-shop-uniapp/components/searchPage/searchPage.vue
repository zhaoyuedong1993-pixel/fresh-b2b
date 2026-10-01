<!--
 * 优诚配运 - 搜索组件
 * 设计规范：自然清新
-->
<template>
	<view class="search-panel">
		<view class="search-box">
			<u-icon name="search" size="32rpx" color="#999999"></u-icon>
			<input
				class="search-input"
				v-model="keyword"
				type="text"
				:placeholder="placeholder"
				confirm-type="search"
				@confirm="handleSearch"
				@input="handleInput"
			/>
			<u-icon v-if="keyword" name="close-circle-fill" size="32rpx" color="#CCCCCC" @click="clearKeyword"></u-icon>
		</view>

		<!-- 搜索面板 -->
		<view class="search-panel-content" v-if="showPanel">
			<!-- 搜索历史 -->
			<view class="search-section" v-if="historyList.length > 0">
				<view class="section-header">
					<view class="section-title">
						<u-icon name="clock" size="28rpx" color="#22A84F"></u-icon>
						<text>搜索历史</text>
					</view>
					<view class="clear-btn" @click="clearHistory">
						<u-icon name="trash" size="28rpx" color="#999999"></u-icon>
					</view>
				</view>
				<view class="keyword-list">
					<view
						class="keyword-item"
						v-for="(item, index) in historyList"
						:key="index"
						@click="fillKeyword(item)"
					>
						{{ item }}
					</view>
				</view>
			</view>

			<!-- 热门搜索 -->
			<view class="search-section" v-if="hotKeywords.length > 0">
				<view class="section-header">
					<view class="section-title">
						<u-icon name="fire" size="28rpx" color="#F97316"></u-icon>
						<text>热门搜索</text>
					</view>
				</view>
				<view class="keyword-list">
					<view
						class="keyword-item"
						v-for="(item, index) in hotKeywords"
						:key="index"
						@click="fillKeyword(item)"
					>
						{{ item }}
					</view>
				</view>
			</view>
		</view>
	</view>
</template>

<script>
	export default {
		name: 'searchPage',
		data() {
			return {
				keyword: '',
				historyList: [],
				showPanel: true,
				hotKeywords: [
					'新鲜蔬菜', '时令水果', '肉类海鲜', '粮油调味',
					'冷冻食品', '蛋奶豆制品', '饮料酒水', '休闲零食'
				]
			}
		},
		props: {
			placeholder: {
				type: String,
				default: '搜索商品'
			},
			historyKey: {
				type: String,
				default: 'searchHistoryList'
			}
		},
		created() {
			this.loadHistory()
		},
		methods: {
			loadHistory() {
				try {
					const history = uni.getStorageSync(this.historyKey)
					if (history) {
						this.historyList = JSON.parse(history)
					}
				} catch (e) {
					this.historyList = []
				}
			},
			saveKeyword() {
				if (!this.keyword) return
				try {
					let list = this.historyList.filter(item => item !== this.keyword)
					list.unshift(this.keyword)
					list = list.slice(0, 10)
					uni.setStorageSync(this.historyKey, JSON.stringify(list))
					this.historyList = list
				} catch (e) { }
			},
			clearHistory() {
				uni.showModal({
					title: '清空历史',
					content: '确定清空搜索历史吗？',
					success: (res) => {
						if (res.confirm) {
							uni.removeStorageSync(this.historyKey)
							this.historyList = []
						}
					}
				})
			},
			clearKeyword() {
				this.keyword = ''
				this.$emit('changeKeyword', '')
			},
			handleInput(val) {
				this.$emit('changeKeyword', val)
			},
			handleSearch() {
				this.saveKeyword()
				this.$emit('search', this.keyword)
			},
			fillKeyword(keyword) {
				this.keyword = keyword
				this.saveKeyword()
				this.$emit('fillKeyword', keyword)
			}
		}
	}
</script>

<style lang="scss" scoped>
	.search-panel {
		padding: 20rpx 24rpx;
		background: #FFFFFF;
	}

	.search-box {
		height: 72rpx;
		background: #F5F7F4;
		border-radius: 36rpx;
		display: flex;
		align-items: center;
		padding: 0 24rpx;
		gap: 12rpx;
	}

	.search-input {
		flex: 1;
		height: 100%;
		font-size: 28rpx;
		color: #1A1A1A;
	}

	.search-panel-content {
		margin-top: 24rpx;
	}

	.search-section {
		margin-bottom: 24rpx;
	}

	.section-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 16rpx;
	}

	.section-title {
		display: flex;
		align-items: center;
		gap: 8rpx;

		text {
			font-size: 26rpx;
			font-weight: 600;
			color: #1A1A1A;
		}
	}

	.clear-btn {
		padding: 8rpx;
	}

	.keyword-list {
		display: flex;
		flex-wrap: wrap;
		gap: 16rpx;
	}

	.keyword-item {
		padding: 12rpx 24rpx;
		background: #F5F7F4;
		border-radius: 24rpx;
		font-size: 24rpx;
		color: #666666;

		&:active {
			background: #E8F8EC;
			color: #22A84F;
		}
	}
</style>
