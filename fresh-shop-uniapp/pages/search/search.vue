<!--
 * 优诚配运 - 搜索页
 * 设计规范：自然清新
-->
<template>
    <pageWrapper>
        <!-- 顶部搜索栏 -->
        <view class="search-header">
            <view class="back-btn" @click="goBack">
                <u-icon name="arrow-left" size="36rpx" color="#1A1A1A"></u-icon>
            </view>
            <view class="search-box">
                <u-icon name="search" size="32rpx" color="#999999"></u-icon>
                <input
                    class="search-input"
                    v-model="keyword"
                    type="text"
                    confirm-type="search"
                    @confirm="handleSearch"
                    placeholder="搜索商品"
                    placeholder-class="placeholder"
                    focus
                />
                <u-icon v-if="keyword" name="close-circle-fill" size="32rpx" color="#CCCCCC" @click="clearKeyword"></u-icon>
            </view>
            <view class="search-btn" @click="handleSearch">
                <text>搜索</text>
            </view>
        </view>

        <!-- 热门搜索 -->
        <view class="hot-search" v-if="!keyword && hotKeywords.length > 0">
            <view class="section-header">
                <view class="section-title">
                    <u-icon name="fire" size="32rpx" color="#F97316"></u-icon>
                    <text>热门搜索</text>
                </view>
            </view>
            <view class="hot-list">
                <view
                    class="hot-item"
                    v-for="(item, index) in hotKeywords"
                    :key="index"
                    @click="goGoodsList(item)"
                >
                    <text>{{ item }}</text>
                </view>
            </view>
        </view>

        <!-- 搜索历史 -->
        <view class="search-history" v-if="!keyword && historyList.length > 0">
            <view class="section-header">
                <view class="section-title">
                    <u-icon name="clock" size="32rpx" color="#22A84F"></u-icon>
                    <text>搜索历史</text>
                </view>
                <view class="clear-btn" @click="clearHistory">
                    <u-icon name="trash" size="28rpx" color="#999999"></u-icon>
                    <text>清空</text>
                </view>
            </view>
            <view class="history-list">
                <view
                    class="history-item"
                    v-for="(item, index) in historyList"
                    :key="index"
                    @click="goGoodsList(item)"
                >
                    <u-icon name="search" size="28rpx" color="#CCCCCC"></u-icon>
                    <text>{{ item }}</text>
                </view>
            </view>
        </view>

        <!-- 搜索结果提示 -->
        <view class="search-tip" v-if="keyword">
            <text>搜索 "{{ keyword }}"</text>
        </view>
    </pageWrapper>
</template>

<script>
    export default {
        data() {
            return {
                keyword: '',
                hotKeywords: [
                    '新鲜蔬菜', '时令水果', '肉类海鲜', '粮油调味',
                    '冷冻食品', '蛋奶豆制品', '饮料酒水', '休闲零食'
                ],
                historyList: []
            }
        },
        onLoad() {
            this.loadHistory()
        },
        methods: {
            goBack() {
                uni.navigateBack()
            },
            handleSearch() {
                if (!this.keyword.trim()) {
                    uni.showToast({ title: '请输入关键词', icon: 'none' })
                    return
                }
                this.saveHistory(this.keyword)
                this.goGoodsList(this.keyword)
            },
            goGoodsList(keyword) {
                uni.redirectTo({
                    url: `/pages/goods/goods?keyword=${encodeURIComponent(keyword)}`
                })
            },
            clearKeyword() {
                this.keyword = ''
            },
            loadHistory() {
                try {
                    const history = uni.getStorageSync('search_history') || '[]'
                    this.historyList = JSON.parse(history)
                } catch (e) {
                    this.historyList = []
                }
            },
            saveHistory(keyword) {
                try {
                    let list = this.historyList.filter(item => item !== keyword)
                    list.unshift(keyword)
                    list = list.slice(0, 10)
                    uni.setStorageSync('search_history', JSON.stringify(list))
                    this.historyList = list
                } catch (e) { }
            },
            clearHistory() {
                uni.showModal({
                    title: '清空历史',
                    content: '确定清空搜索历史吗？',
                    success: (res) => {
                        if (res.confirm) {
                            uni.removeStorageSync('search_history')
                            this.historyList = []
                        }
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .search-header {
        display: flex;
        align-items: center;
        gap: 16rpx;
        padding: 16rpx 24rpx;
        background: #FFFFFF;
        position: sticky;
        top: 0;
        z-index: 100;
        box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);
    }

    .back-btn {
        width: 56rpx;
        height: 56rpx;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .search-box {
        flex: 1;
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

    .placeholder {
        color: #CCCCCC;
        font-size: 28rpx;
    }

    .search-btn {
        padding: 0 24rpx;
        height: 56rpx;
        display: flex;
        align-items: center;
        justify-content: center;

        text {
            font-size: 28rpx;
            color: #22A84F;
            font-weight: 500;
        }
    }

    .section-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 0 24rpx 20rpx;
    }

    .section-title {
        display: flex;
        align-items: center;
        gap: 8rpx;

        text {
            font-size: 28rpx;
            font-weight: 600;
            color: #1A1A1A;
        }
    }

    .clear-btn {
        display: flex;
        align-items: center;
        gap: 6rpx;

        text {
            font-size: 24rpx;
            color: #999999;
        }
    }

    .hot-search {
        padding: 32rpx 0 24rpx;
        background: #FFFFFF;
        margin-bottom: 16rpx;
    }

    .hot-list {
        display: flex;
        flex-wrap: wrap;
        padding: 0 24rpx;
        gap: 16rpx;
    }

    .hot-item {
        padding: 12rpx 28rpx;
        background: #F5F7F4;
        border-radius: 32rpx;
        font-size: 26rpx;
        color: #666666;

        &:active {
            background: #E8F8EC;
            color: #22A84F;
        }
    }

    .search-history {
        padding: 32rpx 0 24rpx;
        background: #FFFFFF;
    }

    .history-list {
        padding: 0 24rpx;
    }

    .history-item {
        display: flex;
        align-items: center;
        gap: 12rpx;
        padding: 20rpx 0;
        border-bottom: 1rpx solid #F0F0F0;

        &:last-child {
            border-bottom: none;
        }

        text {
            font-size: 28rpx;
            color: #666666;
        }
    }

    .search-tip {
        padding: 32rpx 24rpx;
        text-align: center;

        text {
            font-size: 28rpx;
            color: #999999;
        }
    }
</style>
