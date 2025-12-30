import request from '@/utils/request'

export interface PersonalStatistics {
  total_reports: number
  submitted_reports: number
  approved_reports: number
  rejected_reports: number
  late_reports: number
  submit_rate: number
  average_submit_time: string
  status_distribution: Record<string, number>
  daily_submit_count: DailySubmitCount[]
}

export interface TeamStatistics {
  total_members: number
  active_members: number
  total_reports: number
  submit_rate: number
  member_statistics: MemberStatistics[]
  status_distribution: Record<string, number>
  daily_submit_trend: DailySubmitCount[]
}

export interface DailySubmitCount {
  date: string
  count: number
}

export interface MemberStatistics {
  user_id: number
  user_name: string
  real_name: string
  total_reports: number
  submit_rate: number
  late_count: number
}

export interface StatisticsParams {
  user_id?: number
  department?: string
  start_date?: string
  end_date?: string
}

/**
 * 获取个人统计
 */
export function getPersonalStatistics(params?: StatisticsParams) {
  return request.get<PersonalStatistics>('/v1/statistics/personal', { params })
}

/**
 * 获取团队统计
 */
export function getTeamStatistics(params?: StatisticsParams) {
  return request.get<TeamStatistics>('/v1/statistics/team', { params })
}

