import request from '@/utils/request'

/**
 * 通用API接口
 */

/**
 * 获取部门列表
 */
export function getDepartments() {
  return request.get<string[]>('/v1/common/departments')
}

