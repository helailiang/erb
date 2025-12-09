/**
 * 用户相关类型定义
 */

export interface User {
  id: number
  username: string
  real_name: string
  email?: string
  phone?: string
  department?: string
  position?: string
  employee_no?: string
  avatar?: string
  status: string
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
  expire_at: string
}

