<template>
	<pageWrapper>
		<view class="box2">
			<view class="section-title">个人信息</view>
			<view class="info-item">
				<view class="info-label">联系人</view>
				<view class="info-value">{{ user.originContactName || '-' }}</view>
			</view>
			<view class="info-item">
				<view class="info-label">手机号</view>
				<view class="info-value">{{ user.phone || '-' }}</view>
			</view>
			<view class="info-item">
				<view class="info-label">客户类型</view>
				<view class="info-value">{{ customerTypeText }}</view>
			</view>
			<view class="info-item">
				<view class="info-label">审核状态</view>
				<view class="info-value">
					<text v-if="auditStatus === 1" style="color:#67c23a">已通过</text>
					<text v-else-if="[2,3].includes(auditStatus)" style="color:#e6a23c">审核中</text>
					<text v-else-if="auditStatus === 4" style="color:#f56c6c">未通过</text>
					<text v-else>-</text>
				</view>
			</view>
		</view>

		<view class="box2" v-if="company">
			<view class="section-title">公司信息</view>
			<view class="info-item">
				<view class="info-label">公司名称</view>
				<view class="info-value">{{ company.name || '-' }}</view>
			</view>
			<view class="info-item">
				<view class="info-label">联系人</view>
				<view class="info-value">{{ company.contact || '-' }}</view>
			</view>
			<view class="info-item">
				<view class="info-label">联系电话</view>
				<view class="info-value">{{ company.phone || '-' }}</view>
			</view>
			<view class="info-item">
				<view class="info-label">公司地址</view>
				<view class="info-value">{{ company.address || '-' }}</view>
			</view>
			<view class="info-item">
				<view class="info-label">客户类型</view>
				<view class="info-value">{{ company.companyType === 'monthly' ? '月度结算' : '零售' }}</view>
			</view>
		</view>

		<u-toast style="z-index:9998;" ref="toast"></u-toast>
	</pageWrapper>
</template>

<script>
import { getToken, getUser, getRole } from "@/store/storage";
import { getUserAuditStatus } from "@/api/user";
import { getCompanyById } from "@/api/login";

export default {
	data() {
		return {
			user: null,
			role: {},
			company: null,
			auditStatus: 0
		}
	},
	computed: {
		customerTypeText() {
			if (!this.role.authorityName) return '-'
			return this.role.authorityName.replace('客户', '')
		}
	},
	onLoad() {
		const token = getToken()
		if (!token) {
			uni.redirectTo({ url: '/pages/my/my' })
			return
		}
		this.user = getUser()
		this.role = getRole()
		this.loadData()
	},
	methods: {
		async loadData() {
			// 审核状态
			try {
				const res = await getUserAuditStatus()
				if (res.code === 0) {
					this.auditStatus = res.data.auditStatus
					this.user.auditStatus = res.data.auditStatus
				}
			} catch (e) {}

			// 公司信息
			if (this.user && this.user.companyId) {
				try {
					const companyRes = await getCompanyById(this.user.companyId)
					if (companyRes.code === 0) {
						this.company = companyRes.data
					}
				} catch (e) {}
			}
		}
	}
}
</script>

<style scoped lang="scss">
.box1 {
	margin: 10px 10px 12px;
	border-radius: 10px;
	box-shadow: 0px 0px 20px #f4f3f3;
	background: #FFFFFF;
	padding: 20px 10px 20px;
}

.box2 {
	margin: 10px 10px 12px;
	border-radius: 10px;
	box-shadow: 0px 0px 20px #f4f3f3;
	background: #FFFFFF;
	padding: 20px 20px;
}

.section-title {
	font-size: 16px;
	font-weight: bold;
	color: #333;
	margin-bottom: 12px;
	padding-bottom: 8px;
	border-bottom: 1px solid #f0f0f0;
}

.user-info {
	display: flex;
	align-items: center;
	width: 100%;
	height: 100%;
	justify-content: center;

	.face {
		flex-shrink: 0;
		width: 20vw;
		height: 20vw;

		image {
			width: 20vw;
			height: 100%;
			border-radius: 100%
		}
	}
}

.username {
	font-size: 17px;
	width: 100%;
	text-align: center;
	margin: 16px 0 12px 0;
}

.title {
	font-size: 14px;
	width: 100%;
	text-align: center;
	color: #999;
}

.info-item {
	display: flex;
	align-items: flex-start;
	padding: 10px 0;
	border-bottom: 1px solid #f7f7f7;

	&:last-child {
		border-bottom: none;
	}
}

.info-label {
	width: 80px;
	font-size: 14px;
	color: #999;
	flex-shrink: 0;
}

.info-value {
	flex: 1;
	font-size: 14px;
	color: #333;
	word-break: break-all;
}
</style>
