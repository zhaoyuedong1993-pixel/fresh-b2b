import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

const store = new Vuex.Store({
	state: {
		userInfo: null,
		isLogin: false
	},
	mutations: {
		login(state, userInfo) {
			state.userInfo = userInfo
			state.isLogin = true
		},
		logout(state) {
			state.userInfo = null
			state.isLogin = false
		}
	},
	actions: {
		initLoginState({ commit }) {
			const token = uni.getStorageSync('token')
			const userInfo = uni.getStorageSync('userInfo')
			if (token && userInfo) {
				commit('login', userInfo)
			}
		}
	}
})

export default store
