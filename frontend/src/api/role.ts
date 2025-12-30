import request from '@/utils/request'

export interface Role {
  id: number
  name: string
  description?: string
  status: string
  is_system: boolean
  created_at: string
  updated_at: string
  permissions?: Permission[]
}

export interface Permission {
  id: number
  code: string
  name: string
  description?: string
  type: string
  module?: string
  status: string
}

export interface CreateRoleRequest {
  name: string
  description?: string
  status?: string
  permission_ids?: number[]
}

export interface UpdateRoleRequest {
  name?: string
  description?: string
  status?: string
  permission_ids?: number[]
}

export interface ListRoleParams {
  page?: number
  page_size?: number
  name?: string
  status?: string
}

export interface RoleListResponse {
  list: Role[]
  total: number
  page: number
  page_size: number
}

/**
 * 获取角色列表
 */
export function getRoleList(params: ListRoleParams) {
  return request.get<RoleListResponse>('/v1/roles', { params })
}

/**
 * 获取角色详情
 */
export function getRoleDetail(id: number) {
  return request.get<Role>(`/v1/roles/${id}`)
}

/**
 * 创建角色
 */
export function createRole(data: CreateRoleRequest) {
  return request.post<Role>('/v1/roles', data)
}

/**
 * 更新角色
 */
export function updateRole(id: number, data: UpdateRoleRequest) {
  return request.put<Role>(`/v1/roles/${id}`, data)
}

/**
 * 删除角色
 */
export function deleteRole(id: number) {
  return request.delete(`/v1/roles/${id}`)
}

/**
 * 获取所有权限
 */
export function getAllPermissions() {
  return request.get<Permission[]>('/v1/permissions')
}



