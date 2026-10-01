<!--
 * 优诚配运 - 新增/编辑地址
 * 设计规范：自然清新
-->
<template>
    <pageWrapper>
        <!-- 顶部 -->
        <view class="page-header">
            <view class="header-content">
                <view class="back-btn" @click="goBack">
                    <u-icon name="arrow-left" size="40rpx" color="#FFFFFF"></u-icon>
                </view>
                <text class="page-title">{{ isEdit ? '编辑地址' : '新增地址' }}</text>
            </view>
        </view>

        <!-- 表单 -->
        <view class="form-section">
            <!-- 联系人 -->
            <view class="form-card">
                <view class="form-item">
                    <view class="form-label">
                        <u-icon name="account" size="36rpx" color="#22A84F"></u-icon>
                        <text>联系人</text>
                    </view>
                    <input
                        class="form-input"
                        v-model="formData.userName"
                        placeholder="请输入收货人姓名"
                        placeholder-class="placeholder"
                    />
                </view>
                <view class="form-item">
                    <view class="form-label">
                        <u-icon name="phone" size="36rpx" color="#22A84F"></u-icon>
                        <text>手机号</text>
                    </view>
                    <input
                        class="form-input"
                        v-model="formData.userPhone"
                        type="number"
                        maxlength="11"
                        placeholder="请输入手机号"
                        placeholder-class="placeholder"
                    />
                </view>
            </view>

            <!-- 地址 -->
            <view class="form-card">
                <view class="form-item address-item">
                    <view class="form-label">
                        <u-icon name="map" size="36rpx" color="#22A84F"></u-icon>
                        <text>收货地址</text>
                    </view>
                    <view class="address-select" @click="chooseLocation">
                        <text v-if="formData.detailAddress">{{ formData.detailAddress }}</text>
                        <text v-else class="placeholder">点击选择收货地址</text>
                        <u-icon name="arrow-right" size="32rpx" color="#CCCCCC"></u-icon>
                    </view>
                </view>
                <view class="form-item">
                    <view class="form-label">
                        <u-icon name="file-text" size="36rpx" color="#22A84F"></u-icon>
                        <text>门牌号</text>
                    </view>
                    <input
                        class="form-input"
                        v-model="formData.lableName"
                        placeholder="如：东区3号楼101室"
                        placeholder-class="placeholder"
                    />
                </view>
            </view>

            <!-- 地址标签 -->
            <view class="form-card">
                <view class="form-item">
                    <view class="form-label">
                        <u-icon name="tags" size="36rpx" color="#22A84F"></u-icon>
                        <text>标签</text>
                    </view>
                    <view class="tag-list">
                        <view
                            class="tag-item"
                            :class="{ active: formData.lableName === tag }"
                            v-for="tag in tagOptions"
                            :key="tag"
                            @click="selectTag(tag)"
                        >
                            {{ tag }}
                        </view>
                    </view>
                </view>
            </view>
        </view>

        <!-- 保存按钮 -->
        <view class="save-btn" :class="{ loading }" @click="saveAddress">
            <text v-if="!loading">保存地址</text>
            <text v-else>保存中...</text>
        </view>

        <!-- 删除按钮 -->
        <view class="delete-btn" v-if="isEdit" @click="deleteAddress">
            <text>删除该地址</text>
        </view>

        <!-- 地图选择组件 -->
        <liu-chooseAddress ref="chooseAddress" @submit="submitAddress" @detele="deteleAddress"></liu-chooseAddress>

        <u-toast ref="toast" style="z-index: 9999"></u-toast>
    </pageWrapper>
</template>

<script>
    import { getToken } from '@/store/storage.js'
    import { getAddressInfo, createAddress, updateAddress, deleteAddress } from '@/api/address.js'
    import ChooseAddress from '@/uni_modules/liu-chooseAddress/components/liu-chooseAddress/liu-chooseAddress.vue'

    export default {
        components: {
            ChooseAddress
        },
        data() {
            return {
                id: null,
                formData: {
                    userName: '',
                    userPhone: '',
                    detailAddress: '',
                    lableName: '',
                    address: '',
                    latitude: '',
                    longitude: ''
                },
                loading: false,
                tagOptions: ['家', '公司', '学校', '其他']
            }
        },
        computed: {
            isEdit() {
                return !!this.id
            }
        },
        onLoad(options) {
            const token = getToken()
            if (!token) {
                uni.showToast({ title: '请先登录', icon: 'none' })
                setTimeout(() => {
                    uni.navigateBack()
                }, 1500)
                return
            }

            if (options.id) {
                this.id = parseInt(options.id)
                this.loadAddressInfo()
            }
        },
        methods: {
            goBack() {
                uni.navigateBack()
            },
            async loadAddressInfo() {
                uni.showLoading({ title: '加载中...' })
                try {
                    const res = await getAddressInfo({ ID: this.id })
                    if (res.code === 0 && res.data?.reuserAddress) {
                        const addr = res.data.reuserAddress
                        this.formData = {
                            userName: addr.userName || '',
                            userPhone: addr.userPhone || '',
                            detailAddress: addr.detailAddress || '',
                            lableName: addr.lableName || '',
                            address: addr.address || '',
                            latitude: addr.latitude || '',
                            longitude: addr.longitude || ''
                        }
                        this.$refs.chooseAddress?.setData(addr)
                    }
                } catch (e) {
                    uni.showToast({ title: '加载失败', icon: 'none' })
                } finally {
                    uni.hideLoading()
                }
            },
            chooseLocation() {
                // 调用地图选择组件
                this.$refs.chooseAddress?.chooseLocation()
            },
            selectTag(tag) {
                this.formData.lableName = tag
            },
            async submitAddress(data) {
                if (!data.latitude || !data.longitude) {
                    this.$message(this.$refs.toast).error('请选择收货地址')
                    return
                }
                if (!data.userName) {
                    this.$message(this.$refs.toast).error('请输入收货人姓名')
                    return
                }
                if (!data.userPhone) {
                    this.$message(this.$refs.toast).error('请输入手机号')
                    return
                }
                if (!/^1[3-9]\d{9}$/.test(data.userPhone)) {
                    this.$message(this.$refs.toast).error('请输入合法手机号')
                    return
                }

                this.loading = true
                try {
                    const submitData = {
                        ...data,
                        detailAddress: data.address + (data.lableName ? ' ' + data.lableName : '')
                    }

                    let res
                    if (this.isEdit) {
                        submitData.ID = this.id
                        res = await updateAddress(submitData)
                    } else {
                        res = await createAddress(submitData)
                    }

                    if (res.code === 0) {
                        this.$message(this.$refs.toast).success('保存成功')
                        setTimeout(() => {
                            uni.navigateBack()
                        }, 1500)
                    } else {
                        this.$message(this.$refs.toast).error(res.msg || '保存失败')
                    }
                } catch (e) {
                    this.$message(this.$refs.toast).error('保存失败')
                } finally {
                    this.loading = false
                }
            },
            deleteAddress() {
                uni.showModal({
                    title: '删除地址',
                    content: '确定要删除该地址吗？',
                    success: async (res) => {
                        if (res.confirm) {
                            try {
                                const result = await deleteAddress({ ID: this.id })
                                if (result.code === 0) {
                                    this.$message(this.$refs.toast).success('删除成功')
                                    setTimeout(() => {
                                        uni.navigateBack()
                                    }, 1500)
                                }
                            } catch (e) {
                                this.$message(this.$refs.toast).error('删除失败')
                            }
                        }
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .page-header {
        background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
        padding: 48rpx 32rpx 32rpx;
    }

    .header-content {
        display: flex;
        align-items: center;
        gap: 24rpx;
    }

    .back-btn {
        width: 64rpx;
        height: 64rpx;
        background: rgba(255, 255, 255, 0.2);
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .page-title {
        font-size: 36rpx;
        font-weight: 600;
        color: #FFFFFF;
    }

    .form-section {
        padding: 24rpx;
    }

    .form-card {
        background: #FFFFFF;
        border-radius: 24rpx;
        margin-bottom: 20rpx;
        overflow: hidden;
        box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.06);
    }

    .form-item {
        padding: 28rpx 28rpx;
        border-bottom: 1rpx solid #F0F0F0;

        &:last-child {
            border-bottom: none;
        }

        &.address-item {
            flex-direction: column;
            align-items: flex-start;
        }
    }

    .form-label {
        display: flex;
        align-items: center;
        gap: 12rpx;
        margin-bottom: 16rpx;

        text {
            font-size: 28rpx;
            color: #666666;
        }
    }

    .form-input {
        height: 72rpx;
        background: #F5F7F4;
        border-radius: 16rpx;
        padding: 0 24rpx;
        font-size: 30rpx;
        color: #1A1A1A;
    }

    .placeholder {
        color: #CCCCCC;
        font-size: 28rpx;
    }

    .address-select {
        width: 100%;
        min-height: 72rpx;
        background: #F5F7F4;
        border-radius: 16rpx;
        padding: 20rpx 24rpx;
        display: flex;
        align-items: center;
        justify-content: space-between;
        font-size: 28rpx;
        color: #1A1A1A;
        box-sizing: border-box;
    }

    .tag-list {
        display: flex;
        flex-wrap: wrap;
        gap: 16rpx;
    }

    .tag-item {
        padding: 12rpx 28rpx;
        background: #F5F7F4;
        border-radius: 32rpx;
        font-size: 26rpx;
        color: #666666;
        border: 2rpx solid transparent;

        &.active {
            background: #E8F8EC;
            color: #22A84F;
            border-color: #22A84F;
        }
    }

    .save-btn {
        margin: 40rpx 32rpx 24rpx;
        height: 96rpx;
        background: linear-gradient(135deg, #22A84F 0%, #1A9A45 100%);
        border-radius: 48rpx;
        display: flex;
        align-items: center;
        justify-content: center;
        box-shadow: 0 4rpx 24rpx rgba(34, 168, 79, 0.3);

        text {
            font-size: 32rpx;
            font-weight: 600;
            color: #FFFFFF;
        }

        &.loading {
            opacity: 0.7;
        }
    }

    .delete-btn {
        margin: 0 32rpx;
        height: 88rpx;
        background: #FEE2E2;
        border-radius: 44rpx;
        display: flex;
        align-items: center;
        justify-content: center;

        text {
            font-size: 30rpx;
            color: #EF4444;
        }
    }
</style>
