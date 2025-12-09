/**
 * 认证相关工具函数
 */

// 存储token到localStorage
export function setToken(token: string): void {
  localStorage.setItem('token', token)
}

// 从localStorage获取token
export function getToken(): string | null {
  return localStorage.getItem('token')
}

// 移除token
export function removeToken(): void {
  localStorage.removeItem('token')
}

// 检查是否已登录
export function isLoggedIn(): boolean {
  return !!getToken()
}

