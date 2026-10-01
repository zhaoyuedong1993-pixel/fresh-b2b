<!--
 * 优诚配运 - 底部导航栏
 * 设计规范：自然清新
-->
<template>
    <view class="tabbar-wrapper">
        <u-tabbar
            :value="tabsId"
            @change="tabsChange"
            :fixed="true"
            :placeholder="true"
            :safeAreaInsetBottom="true"
            :activeColor="'#22A84F'"
            :inactiveColor="'#999999'"
        >
            <u-tabbar-item
                v-for="(item, index) in tabs"
                :key="index"
                :text="item.text"
                :icon="item.icon"
            ></u-tabbar-item>
        </u-tabbar>
    </view>
</template>

<script>
    export default {
        name: 'tabbar',
        props: {
            tabsId: {
                type: Number,
                default: 0
            }
        },
        data() {
            return {
                tabs: [
                    {
                        text: '首页',
                        icon: 'home',
                        url: '/pages/index/index'
                    },
                    {
                        text: '分类',
                        icon: 'grid',
                        url: '/pages/category/category'
                    },
                    {
                        text: '快速下单',
                        icon: 'star-fill',
                        url: '/pages/goods/quickPay'
                    },
                    {
                        text: '购物车',
                        icon: 'shopping-cart',
                        url: '/pages/cart/cart'
                    },
                    {
                        text: '我的',
                        icon: 'account',
                        url: '/pages/my/my'
                    }
                ]
            }
        },
        methods: {
            tabsChange(index) {
                uni.redirectTo({
                    url: this.tabs[index].url
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .tabbar-wrapper {
        z-index: 999;
    }
</style>
