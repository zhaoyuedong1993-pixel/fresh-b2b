<!--
 * 优诚配运 - 页面包装器
 * 设计规范：自然清新
-->
<template>
	<view class="page-wrapper" :style="{ backgroundColor: bgColor }">
		<slot></slot>
	</view>
</template>

<script>
	export default {
		name: 'pageWrapper',
		data() {
			return {
				screenHeight: 0
			}
		},
		props: {
			subHeight: {
				type: Number,
				default: 0
			},
			bgColor: {
				type: String,
				default: '#F5F7F4'
			}
		},
		mounted() {
			uni.getSystemInfo({
				success: (res) => {
					this.screenHeight = res.windowHeight
				}
			})
		}
	}
</script>

<style lang="scss">
.page-wrapper {
	min-height: 100vh;
	box-sizing: border-box;
}
</style>
