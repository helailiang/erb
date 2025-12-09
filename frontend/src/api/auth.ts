import request from '@/utils/request'
import type { LoginRequest, LoginResponse, User } from '@/types/user'

/**
 * 用户登录
 */
export function login(data: LoginRequest) {
  return request.post<LoginResponse>('/v1/auth/login', data)
}

/**
 * 获取当前用户信息
 */
export function getUserInfo() {
  return request.get<User>('/v1/auth/me')
}

