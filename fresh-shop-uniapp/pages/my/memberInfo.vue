<template>
    <pageWrapper>
        <view class="page-content">
            <view class="info-card">
                <view class="card-header">
                    <text class="card-title">个人信息</text>
                </view>
                <view class="info-item">
                    <text class="label">联系人</text>
                    <text class="value">{{ userInfo && userInfo.originContactName || '-' }}</text>
                </view>
                <view class="info-item">
                    <text class="label">手机号</text>
                    <text class="value">{{ userInfo && userInfo.phone || '-' }}</text>
                </view>
                <view class="info-item">
                    <text class="label">客户类型</text>
                    <text class="value">{{ getCustomerType() }}</text>
                </view>
                <view class="info-item last">
                    <text class="label">审核状态</text>
                    <view class="value">
                        <text class="status-tag pass" v-if="auditStatus === 1">已通过</text>
                        <text class="status-tag pending" v-else-if="auditStatus === 2 || auditStatus === 3">审核中</text>
                        <text class="status-tag reject" v-else-if="auditStatus === 4">未通过</text>
                        <text class="status-tag" v-else>-</text>
                    </view>
                </view>
            </view>
            <view class="info-card" v-if="companyInfo">
                <view class="card-header">
                    <text class="card-title">公司信息</text>
                </view>
                <view class="info-item">
                    <text class="label">公司名称</text>
                    <text class="value">{{ companyInfo.name || '-' }}</text>
                </view>
                <view class="info-item">
                    <text class="label">联系人</text>
                    <text class="value">{{ companyInfo.contact || '-' }}</text>
                </view>
                <view class="info-item">
                    <text class="label">联系电话</text>
                    <text class="value">{{ companyInfo.phone || '-' }}</text>
                </view>
                <view class="info-item">
                    <text class="label">公司地址</text>
                    <text class="value">{{ companyInfo.address || '-' }}</text>
                </view>
                <view class="info-item last">
                    <text class="label">客户类型</text>
                    <text class="value">{{ companyInfo.companyType === 'monthly' ? '月度结算' : '零售' }}</text>
                </view>
            </view>
        </view>
        <u-toast ref="toast" style="z-index:9998;"></u-toast>
    </pageWrapper>
</template>

<script>
import { getUser, getRole } from "@/store/storage"
import { getUserAuditStatus } from "@/api/user"
import { getCompanyById } from "@/api/login"

export default {
    data() {
        return {
            userInfo: null,
            role: {},
            companyInfo: null,
            auditStatus: 0
        }
    },
    computed: {
        customerTypeText() {
            if (!this.role || !this.role.authorityName) return '-'
            return this.role.authorityName.replace('客户', '')
        }
    },
    onLoad() {
        const token = getUser()
        if (!token) {
            uni.redirectTo({ url: '/pages/my/my' })
            return
        }
        this.userInfo = getUser()
        this.role = getRole()
        this.loadData()
    },
    methods: {
        getCustomerType() {
            if (!this.role || !this.role.authorityName) return '-'
            return this.role.authorityName.replace('客户', '')
        },
        async loadData() {
            try {
                const res = await getUserAuditStatus()
                if (res.code === 0) {
                    this.auditStatus = res.data.auditStatus
                }
            } catch (e) { }
            if (this.userInfo && this.userInfo.companyId) {
                try {
                    const companyRes = await getCompanyById(this.userInfo.companyId)
                    if (companyRes.code === 0) {
                        this.companyInfo = companyRes.data
                    }
                } catch (e) { }
            }
        }
    }
}
</script>

<style lang="scss" scoped>
.page-content {
    padding: 24rpx;
    background: #F5F7F4;
    min-height: 100vh;
}
.info-card {
    background: #FFFFFF;
    border-radius: 24rpx;
    margin-bottom: 24rpx;
}
.card-header {
    padding: 28rpx 32rpx 20rpx;
    border-bottom: 1rpx solid #EEEEEE;
}
.card-title {
    font-size: 32rpx;
    font-weight: 600;
    color: #1A1A1A;
}
.info-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 28rpx 32rpx;
    border-bottom: 1rpx solid #EEEEEE;
}
.info-item.last { border-bottom: none; }
.label { font-size: 28rpx; color: #999999; flex-shrink: 0; }
.value { font-size: 28rpx; color: #1A1A1A; text-align: right; max-width: 70%; }
.status-tag {
    font-size: 24rpx;
    padding: 6rpx 20rpx;
    border-radius: 20rpx;
}
.status-tag.pass { background: #E8F8EC; color: #22A84F; }
.status-tag.pending { background: #FEF3E2; color: #F97316; }
.status-tag.reject { background: #FEE2E2; color: #EF4444; }
</style>
