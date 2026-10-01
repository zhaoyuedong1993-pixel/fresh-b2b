/*
 * @Author: dalefeng
 * @Date: 2023-03-23 18:13:11
 * @LastEditors: dalefeng
 * @LastEditTime: 2023-04-23 14:14:32
 */
// uniapp request 请求
import toast from '@/utils/toast.js'
import config from '@/config/config.js'
import { setToken, setUser } from '@/store/storage.js'

/**
 * 封装 Uniapp 的请求工具类
 * @param {Object} options 请求配置参数
 * @param {string} options.url 请求的地址
 * @param {string} [options.method='GET'] 请求的方法
 * @param {Object} [options.data] 请求的数据
 * @param {Object} [options.header] 请求的头部信息
 * @param {boolean} [options.showLoading=false] 是否显示加载动画
 * @param {boolean} [options.showError=true] 是否显示错误提示
 * @param {boolean} [options.toLogin=false] 是否自动跳转到登录页
 * @returns {Promise} Promise 对象
 * @param toastRefs
 */
// 获取API地址（动态）
function getApiUrl() {
    // #ifdef H5
    return 'http://' + location.hostname + ':48888'
    // #endif
    // #ifndef H5
    return 'http://192.168.1.12:48888'
    // #endif
}

const request = (options, toastRefs) => {
	options.method = options.method || 'GET'
	options.loading = options.loading === undefined ? false : options.loading
	options.toLogin = options.toLogin === undefined ? false : options.toLogin
	toastRefs = toastRefs === undefined ? null : toastRefs
	if (options.loading) {
		toast.message(toastRefs).loading('加载中...')
	}
	options.header = Object.assign({
		'content-type': 'application/json',
		'x-token': uni.getStorageSync('token')
	}, options.header)

	options.url = config.baseUrl + options.url

	// 处理 params 参数，手动添加到 URL
	if (options.params) {
		const params = options.params
		const queryParts = []
		for (const key in params) {
			queryParts.push(key + '=' + encodeURIComponent(params[key]))
		}
		if (queryParts.length > 0) {
			options.url += '?' + queryParts.join('&')
		}
		delete options.params
	}

	// 发起请求
	return new Promise((resolve, reject) => {
		uni.request({
			...options,
			success: (res) => {
				// 隐藏加载动画
				if (options.loading) {
					toast.message(toastRefs).hide()
				}

				// 如果响应头部中返回了新的 token，则重新写入本地存储中
				if (res.header['new-token']) {
					uni.setStorageSync('token', res.header['new-token'])
				}

				// 根据响应状态码处理响应结果
				switch (res.statusCode) {
					case 200:
						if (res.data.code !== 0) {
							// 授权已过期
							if (res.data.code === 401) {
								toast.message(toastRefs).error('授权已过期，请重新登陆').then(() => {
									if (options.toLogin) {
										setToken('')
										setUser(null)
										// 如果需要自动跳转到登录页，则跳转到登录页
										uni.reLaunch({
											url: '/pages/my/my'
										})
									}
								})
								setToken('')
								setUser(null)
								resolve(res.data)
								return
							} else {
								toast.message(toastRefs).error(res.data.msg)
								resolve(res.data)
								return
							}
						}
						resolve(res.data)
						return
					case 401:
						if (options.toLogin) {
							console.log(res.data.msg)
							// 如果需要自动跳转到登录页，则跳转到登录页
							uni.reLaunch({
								url: '/pages/my/my'
							})
						}
						reject(res.data)
						return
					case 500:
						toast.message(toastRefs).error('服务器异常，请稍后重试！')
						reject(res.data)
						return
					default:
						toast.message(toastRefs).error('请求失败，请稍后重试！')
						reject(res.data)
						return
				}
			},
			fail: (err) => {
				// 隐藏加载动画
				if (options.loading) {
					toast.message(toastRefs).hide()
				}
				console.log(`uni.request ${options.url} fail`, err)
				toast.message(toastRefs).error('服务器异常，请稍后重试！')
				reject(err)
			},
		})
	})
}

export default request
