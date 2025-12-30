<template>
  <div class="report-detail" v-loading="loading">
    <el-card v-if="report">
      <template #header>
        <div class="card-header">
          <span>日报详情</span>
          <div>
            <el-button @click="handleBack">返回</el-button>
            <el-button
              v-if="report.status === 'submitted' || report.status === 'rejected'"
              type="primary"
              @click="handleEdit"
            >
              编辑
            </el-button>
            <el-button
              v-if="canApprove"
              type="success"
              @click="showApproveDialog = true"
            >
              审批
            </el-button>
          </div>
        </div>
      </template>

      <div class="report-info">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="日期">
            {{ formatDate(report.report_date) }}
            <el-tag v-if="report.is_late" type="warning" style="margin-left: 10px">补交</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="提交人">
            {{ report.user?.real_name || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(report.status)">
              {{ getStatusText(report.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="提交时间">
            {{ report.submitted_at ? formatDateTime(report.submitted_at) : '-' }}
          </el-descriptions-item>
          <el-descriptions-item v-if="report.late_reason" label="补交原因" :span="2">
            {{ report.late_reason }}
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <div class="report-content">
        <h3 v-if="report.title">{{ report.title }}</h3>
        <div class="content-text" v-html="formatContent(report.content)"></div>
      </div>

      <!-- 审批记录 -->
      <div v-if="report.approvals && report.approvals.length > 0" class="approvals-section">
        <h3>审批记录</h3>
        <el-timeline>
          <el-timeline-item
            v-for="approval in report.approvals"
            :key="approval.id"
            :timestamp="approval.approved_at ? formatDateTime(approval.approved_at) : '待审批'"
            :type="getApprovalType(approval.status)"
          >
            <div>
              <strong>{{ approval.approver?.real_name || '-' }}</strong>
              <el-tag :type="getApprovalType(approval.status)" style="margin-left: 10px">
                {{ getApprovalStatusText(approval.status) }}
              </el-tag>
              <p v-if="approval.comment" style="margin-top: 5px">{{ approval.comment }}</p>
            </div>
          </el-timeline-item>
        </el-timeline>
      </div>

      <!-- 评论区域 -->
      <div class="comments-section">
        <h3>评论</h3>
        <div v-if="comments.length === 0" class="no-comments">暂无评论</div>
        <div v-for="comment in comments" :key="comment.id" class="comment-item">
          <div class="comment-header">
            <strong>{{ comment.user?.real_name || '-' }}</strong>
            <span class="comment-time">{{ formatDateTime(comment.created_at) }}</span>
            <el-tag v-if="comment.is_important" type="danger" size="small">重要</el-tag>
          </div>
          <div class="comment-content">{{ comment.content }}</div>
          <div v-if="comment.replies && comment.replies.length > 0" class="replies">
            <div v-for="reply in comment.replies" :key="reply.id" class="reply-item">
              <strong>{{ reply.user?.real_name || '-' }}</strong>:
              {{ reply.content }}
            </div>
          </div>
        </div>

        <!-- 添加评论 -->
        <el-form :model="commentForm" style="margin-top: 20px">
          <el-form-item>
            <el-input
              v-model="commentForm.content"
              type="textarea"
              :rows="3"
              placeholder="添加评论..."
            />
          </el-form-item>
          <el-form-item>
            <el-checkbox v-model="commentForm.is_important">标记为重要</el-checkbox>
            <el-button type="primary" @click="handleAddComment">发表评论</el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-card>

    <!-- 审批对话框 -->
    <el-dialog v-model="showApproveDialog" title="审批日报" width="500px">
      <el-form :model="approveForm" label-width="80px">
        <el-form-item label="审批结果" required>
          <el-radio-group v-model="approveForm.status">
            <el-radio value="approved">通过</el-radio>
            <el-radio value="rejected">驳回</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="审批意见">
          <el-input
            v-model="approveForm.comment"
            type="textarea"
            :rows="4"
            placeholder="请输入审批意见（可选）"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showApproveDialog = false">取消</el-button>
        <el-button type="primary" :loading="approving" @click="handleApprove">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import { getReportDetail, approveReport, addComment, getComments } from '@/api/report'
import { useUserStore } from '@/stores/user'
import type { DailyReport, Comment, ApproveReportRequest, AddCommentRequest } from '@/api/report'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const approving = ref(false)
const report = ref<DailyReport | null>(null)
const comments = ref<Comment[]>([])
const showApproveDialog = ref(false)

const approveForm = reactive<ApproveReportRequest>({
  status: 'approved',
  comment: '',
})

const commentForm = reactive<AddCommentRequest>({
  content: '',
  is_important: false,
})

// 是否可以审批（当前用户是审批人且状态为pending）
const canApprove = computed(() => {
  if (!report.value || !report.value.approvals) return false
  const pendingApproval = report.value.approvals.find(
    (a) => a.status === 'pending' && a.approver_id === userStore.user?.id
  )
  return !!pendingApproval
})

// 加载日报详情
const loadReportDetail = async () => {
  const id = Number(route.params.id)
  if (!id) {
    ElMessage.error('无效的日报ID')
    router.back()
    return
  }

  loading.value = true
  try {
    const res = await getReportDetail(id)
    report.value = res.data
    await loadComments(id)
  } catch (error: any) {
    ElMessage.error(error.message || '加载日报详情失败')
    router.back()
  } finally {
    loading.value = false
  }
}

// 加载评论
const loadComments = async (reportId: number) => {
  try {
    const res = await getComments(reportId)
    comments.value = res.data
  } catch (error: any) {
    console.error('加载评论失败:', error)
  }
}

// 审批日报
const handleApprove = async () => {
  if (!report.value) return

  approving.value = true
  try {
    await approveReport(report.value.id, approveForm)
    ElMessage.success('审批成功')
    showApproveDialog.value = false
    await loadReportDetail()
  } catch (error: any) {
    ElMessage.error(error.message || '审批失败')
  } finally {
    approving.value = false
  }
}

// 添加评论
const handleAddComment = async () => {
  if (!commentForm.content.trim()) {
    ElMessage.warning('请输入评论内容')
    return
  }

  if (!report.value) return

  try {
    await addComment(report.value.id, commentForm)
    ElMessage.success('评论发表成功')
    commentForm.content = ''
    commentForm.is_important = false
    await loadComments(report.value.id)
  } catch (error: any) {
    ElMessage.error(error.message || '发表评论失败')
  }
}

// 编辑日报
const handleEdit = () => {
  if (report.value) {
    router.push(`/reports/${report.value.id}/edit`)
  }
}

// 返回
const handleBack = () => {
  router.back()
}

// 格式化
const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD')
}

const formatDateTime = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const formatContent = (content: string) => {
  // 简单的换行处理
  return content.replace(/\n/g, '<br>')
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    draft: 'info',
    submitted: 'warning',
    approved: 'success',
    rejected: 'danger',
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    draft: '草稿',
    submitted: '已提交',
    approved: '已审批',
    rejected: '已驳回',
  }
  return map[status] || status
}

const getApprovalType = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    approved: 'success',
    rejected: 'danger',
  }
  return map[status] || 'info'
}

const getApprovalStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待审批',
    approved: '已通过',
    rejected: '已驳回',
  }
  return map[status] || status
}

onMounted(() => {
  loadReportDetail()
})
</script>


<style scoped lang="scss">
.report-detail {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: bold;
    font-size: 16px;
  }

  .report-info {
    margin-bottom: 20px;
  }

  .report-content {
    margin: 20px 0;
    padding: 20px;
    background: #f5f7fa;
    border-radius: 4px;

    h3 {
      margin-bottom: 15px;
      color: #303133;
    }

    .content-text {
      line-height: 1.8;
      white-space: pre-wrap;
      word-break: break-word;
    }
  }

  .approvals-section,
  .comments-section {
    margin-top: 30px;
    padding-top: 20px;
    border-top: 1px solid #e4e7ed;

    h3 {
      margin-bottom: 15px;
      color: #303133;
    }
  }

  .no-comments {
    color: #909399;
    text-align: center;
    padding: 20px;
  }

  .comment-item {
    margin-bottom: 15px;
    padding: 10px;
    background: #f5f7fa;
    border-radius: 4px;

    .comment-header {
      display: flex;
      align-items: center;
      gap: 10px;
      margin-bottom: 5px;

      .comment-time {
        color: #909399;
        font-size: 12px;
      }
    }

    .comment-content {
      margin-top: 5px;
      line-height: 1.6;
    }

    .replies {
      margin-top: 10px;
      padding-left: 20px;
      border-left: 2px solid #dcdfe6;

      .reply-item {
        margin-bottom: 5px;
        font-size: 14px;
        color: #606266;
      }
    }
  }
}
</style>
