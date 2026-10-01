<!--
 * 优诚配运 - 收货地址列表
 * 设计规范：自然清新
-->
<template>
    <pageWrapper>
        <!-- 顶部 -->
        <view class="page-header">
            <text class="page-title">收货地址</text>
            <text class="page-desc">管理您的收货地址</text>
        </view>

        <!-- 地址列表 -->
        <view class="address-list" v-if="addressList.length > 0">
            <view
                class="address-card"
                v-for="item in addressList"
                :key="item.ID"
                @click="selectAddress(item)"
            >
                <view class="address-main">
                    <view class="address-user">
                        <text class="name">{{ item.userName }}</text>
                        <text class="phone">{{ item.userPhone }}</text>
                    </view>
                    <view class="address-detail">
                        <text class="tag" v-if="item.lableName">{{ item.lableName }}</text>
                        <text>{{ item.detailAddress || item.address }}</text>
                    </view>
                </view>
                <view class="address-actions">
                    <view class="action-btn" @click.stop="editAddress(item)">
                        <u-icon name="edit-pen" size="32rpx" color="#666666"></u-icon>
                        <text>编辑</text>
                    </view>
                    <view class="action-btn delete" @click.stop="deleteAddress(item)">
                        <u-icon name="trash" size="32rpx" color="#EF4444"></u-icon>
                        <text>删除</text>
                    </view>
                </view>
                <u-icon
                    v-if="item.isDefault === 1 || item.agreeState === 1"
                    name="checkmark-circle-fill"
                    size="40rpx"
                    color="#22A84F"
                    class="default-icon"
                ></u-icon>
            </view>
        </view>

        <!-- 空状态 -->
        <view class="empty-wrap" v-else-if="!loading">
            <view class="empty-icon">
                <u-icon name="map" size="120rpx" color="#CCCCCC"></u-icon>
            </view>
            <view class="empty-text">暂无收货地址</view>
            <view class="empty-tip">点击下方按钮添加新地址</view>
        </view>

        <!-- 加载中 -->
        <view class="loading-wrap" v-if="loading">
            <text>加载中...</text>
        </view>

        <!-- 新增地址按钮 -->
        <view class="add-btn" @click="addAddress">
            <u-icon name="plus" size="36rpx" color="#FFFFFF"></u-icon>
            <text>新增地址</text>
        </view>

        <!-- 底部占位 -->
        <view class="bottom-placeholder"></view>

        <!-- 删除确认弹窗 -->
        <u-modal
            :show="showDeleteModal"
            :showCancelButton="true"
            title="删除地址"
            content="确定要删除该地址吗？"
            @confirm="confirmDelete"
            @cancel="showDeleteModal = false"
            @close="showDeleteModal = false"
        ></u-modal>

        <u-toast ref="toast" style="z-index: 9999"></u-toast>
    </pageWrapper>
</template>

<script>
    import { getAddressList, deleteAddress as deleteAddressApi } from '@/api/address.js'
    import { getToken } from '@/store/storage.js'

    export default {
        data() {
            return {
                addressList: [],
                loading: false,
                showDeleteModal: false,
                currentDeleteId: null,
                selectMode: false
            }
        },
        onLoad(options) {
            // 判断是否为选择模式（从订单提交页进入）
            if (options.select) {
                this.selectMode = true
            }
            this.loadAddressList()
        },
        onShow() {
            // 返回时刷新列表
            this.loadAddressList()
        },
        methods: {
            async loadAddressList() {
                const token = getToken()
                if (!token) {
                    uni.showToast({ title: '请先登录', icon: 'none' })
                    return
                }

                this.loading = true
                try {
                    const res = await getAddressList()
                    if (res.code === 0) {
                        this.addressList = res.data || []
                    }
                } catch (e) {
                    uni.showToast({ title: '加载失败', icon: 'none' })
                } finally {
                    this.loading = false
                }
            },
            addAddress() {
                uni.navigateTo({ url: '/pages/address/addressForm' })
            },
            editAddress(item) {
                uni.navigateTo({ url: `/pages/address/addressForm?id=${item.ID}` })
            },
            selectAddress(item) {
                if (this.selectMode) {
                    // 选择地址并返回上一页
                    const pages = getCurrentPages()
                    const prevPage = pages[pages.length - 2]
                    if (prevPage) {
                        prevPage.address = item
                        prevPage.addressId = item.ID
                    }
                    uni.navigateBack()
                }
            },
            deleteAddress(item) {
                this.currentDeleteId = item.ID
                this.showDeleteModal = true
            },
            async confirmDelete() {
                if (!this.currentDeleteId) return

                try {
                    const res = await deleteAddressApi({ ID: this.currentDeleteId })
                    if (res.code === 0) {
                        this.$message(this.$refs.toast).success('删除成功')
                        this.loadAddressList()
                    } else {
                        this.$message(this.$refs.toast).error(res.msg || '删除失败')
                    }
                } catch (e) {
                    this.$message(this.$refs.toast).error('删除失败')
                } finally {
                    this.showDeleteModal = false
                    this.currentDeleteId = null
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .page-header {
        background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
        padding: 48rpx 32rpx 32rpx;
    }

    .page-title {
        display: block;
        font-size: 48rpx;
        font-weight: 700;
        color: #FFFFFF;
    }

    .page-desc {
        display: block;
        font-size: 26rpx;
        color: rgba(255, 255, 255, 0.7);
        margin-top: 8rpx;
    }

    .address-list {
        padding: 24rpx;
    }

    .address-card {
        background: #FFFFFF;
        border-radius: 24rpx;
        padding: 28rpx;
        margin-bottom: 20rpx;
        box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.06);
        position: relative;

        &:active {
            background: #FAFAFA;
        }
    }

    .address-main {
        margin-bottom: 20rpx;
    }

    .address-user {
        display: flex;
        align-items: center;
        gap: 16rpx;
        margin-bottom: 12rpx;

        .name {
            font-size: 32rpx;
            font-weight: 600;
            color: #1A1A1A;
        }

        .phone {
            font-size: 28rpx;
            color: #666666;
        }
    }

    .address-detail {
        font-size: 26rpx;
        color: #666666;
        line-height: 1.5;
        display: flex;
        align-items: flex-start;
        gap: 12rpx;

        .tag {
            background: #E8F8EC;
            color: #22A84F;
            font-size: 22rpx;
            padding: 4rpx 12rpx;
            border-radius: 8rpx;
            flex-shrink: 0;
        }
    }

    .address-actions {
        display: flex;
        align-items: center;
        gap: 32rpx;
        padding-top: 20rpx;
        border-top: 1rpx solid #F0F0F0;
    }

    .action-btn {
        display: flex;
        align-items: center;
        gap: 8rpx;
        font-size: 26rpx;
        color: #666666;

        &.delete {
            color: #EF4444;
        }
    }

    .default-icon {
        position: absolute;
        top: 28rpx;
        right: 28rpx;
    }

    .empty-wrap {
        display: flex;
        flex-direction: column;
        align-items: center;
        padding-top: 160rpx;
        gap: 16rpx;
    }

    .empty-text {
        font-size: 32rpx;
        color: #666666;
        font-weight: 500;
    }

    .empty-tip {
        font-size: 26rpx;
        color: #999999;
        margin-top: 8rpx;
    }

    .loading-wrap {
        text-align: center;
        padding: 48rpx;
        color: #999999;
        font-size: 26rpx;
    }

    .add-btn {
        position: fixed;
        bottom: 40rpx;
        left: 32rpx;
        right: 32rpx;
        height: 96rpx;
        background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
        border-radius: 48rpx;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 12rpx;
        box-shadow: 0 4rpx 24rpx rgba(34, 168, 79, 0.3);

        text {
            font-size: 32rpx;
            font-weight: 600;
            color: #FFFFFF;
        }
    }

    .bottom-placeholder {
        height: 180rpx;
    }
</style>
