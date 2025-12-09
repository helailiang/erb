import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login, getUserInfo } from '@/api/auth'
import { setToken, removeToken, getToken } from '@/utils/auth'
import type { LoginRequest, User } from '@/types/user'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(getToken() || '')
  const user = ref<User | null>(null)

  // 登录
  const loginAction = async (loginData: LoginRequest) => {
    const res = await login(loginData)
    token.value = res.data.token
    user.value = res.data.user
    setToken(res.data.token)
    return res
  }

  // 获取用户信息
  const getUserInfoAction = async () => {
    const res = await getUserInfo()
    user.value = res.data
    return res
  }

  // 登出
  const logout = () => {
    token.value = ''
    user.value = null
    removeToken()
  }

  return {
    token,
    user,
    loginAction,
    getUserInfoAction,
    logout,
  }
})

