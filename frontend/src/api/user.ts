import request from '@/utils/request'
import type { User } from '@/types/user'

export interface CreateUserRequest {
  username: string
  password: string
  real_name: string
  email?: string
  phone?: string
  department?: string
  position?: string
  employee_no?: string
  join_date?: string
  manager_id?: number
  role_ids?: number[]
}

export interface UpdateUserRequest {
  real_name?: string
  email?: string
  phone?: string
  department?: string
  position?: string
  employee_no?: string
  join_date?: string
  manager_id?: number
  status?: string
  avatar?: string
}

export interface ListUserParams {
  page?: number
  page_size?: number
  username?: string
  real_name?: string
  department?: string
  status?: string
}

export interface UserListResponse {
  list: User[]
  total: number
  page: number
  page_size: number
}

/**
 * 获取用户列表
 */
export function getUserList(params: ListUserParams) {
  return request.get<UserListResponse>('/v1/users', { params })
}

/**
 * 获取用户详情
 */
export function getUserDetail(id: number) {
  return request.get<User>(`/v1/users/${id}`)
}

/**
 * 创建用户
 */
export function createUser(data: CreateUserRequest) {
  return request.post<User>('/v1/users', data)
}

/**
 * 更新用户
 */
export function updateUser(id: number, data: UpdateUserRequest) {
  return request.put<User>(`/v1/users/${id}`, data)
}

/**
 * 删除用户
 */
export function deleteUser(id: number) {
  return request.delete(`/v1/users/${id}`)
}

/**
 * 重置密码
 */
export function resetPassword(id: number, password: string) {
  return request.post(`/v1/users/${id}/reset-password`, { password })
}



