import request from '@/utils/request'

export interface DailyReport {
  id: number
  user_id: number
  report_date: string
  template_id?: number
  title?: string
  content: string
  status: string
  is_late: boolean
  late_reason?: string
  submitted_at?: string
  created_at: string
  updated_at: string
  user?: any
  approvals?: Approval[]
  comments?: Comment[]
}

export interface Approval {
  id: number
  report_id: number
  approver_id: number
  status: string
  comment?: string
  approved_at?: string
  approver?: any
}

export interface Comment {
  id: number
  report_id: number
  user_id: number
  content: string
  parent_id?: number
  is_important: boolean
  created_at: string
  user?: any
  replies?: Comment[]
}

export interface ReportDraft {
  id: number
  user_id: number
  report_date: string
  template_id?: number
  title?: string
  content: string
  version: number
}

export interface CreateReportRequest {
  report_date: string
  template_id?: number
  title?: string
  content: string
  is_late?: boolean
  late_reason?: string
}

export interface UpdateReportRequest {
  title?: string
  content: string
}

export interface ListReportParams {
  page?: number
  page_size?: number
  user_id?: number
  status?: string
  start_date?: string
  end_date?: string
  department?: string
}

export interface ReportListResponse {
  list: DailyReport[]
  total: number
  page: number
  page_size: number
}

export interface SaveDraftRequest {
  report_date: string
  template_id?: number
  content: string
  title?: string
}

export interface ApproveReportRequest {
  status: 'approved' | 'rejected'
  comment?: string
}

export interface AddCommentRequest {
  content: string
  parent_id?: number
  is_important?: boolean
}

/**
 * 获取日报列表
 */
export function getReportList(params: ListReportParams) {
  return request.get<ReportListResponse>('/v1/reports', { params })
}

/**
 * 获取日报详情
 */
export function getReportDetail(id: number) {
  return request.get<DailyReport>(`/v1/reports/${id}`)
}

/**
 * 创建日报
 */
export function createReport(data: CreateReportRequest) {
  return request.post<DailyReport>('/v1/reports', data)
}

/**
 * 更新日报
 */
export function updateReport(id: number, data: UpdateReportRequest) {
  return request.put<DailyReport>(`/v1/reports/${id}`, data)
}

/**
 * 删除日报
 */
export function deleteReport(id: number) {
  return request.delete(`/v1/reports/${id}`)
}

/**
 * 保存草稿
 */
export function saveDraft(data: SaveDraftRequest) {
  return request.post('/v1/reports/draft', data)
}

/**
 * 获取草稿
 */
export function getDraft(date: string) {
  return request.get<ReportDraft>('/v1/reports/draft', { params: { date } })
}

/**
 * 审批日报
 */
export function approveReport(id: number, data: ApproveReportRequest) {
  return request.post(`/v1/reports/${id}/approve`, data)
}

/**
 * 添加评论
 */
export function addComment(id: number, data: AddCommentRequest) {
  return request.post<Comment>(`/v1/reports/${id}/comments`, data)
}

/**
 * 获取评论列表
 */
export function getComments(id: number) {
  return request.get<Comment[]>(`/v1/reports/${id}/comments`)
}

