<!--
 * 优诚配运 - 页面包装器
 * 设计规范：自然清新
-->
<template>
	<view :style="{ minHeight: windowHeight + 'px', backgroundColor: bgColor }">
		<slot></slot>
	</view>
</template>

<script>
	export default {
		name: 'pageWrapper',
		data() {
			return {
				windowHeight: 0
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
					this.windowHeight = res.windowHeight - this.subHeight
				}
			})
		}
	}
</script>

<style lang="scss"></style>
