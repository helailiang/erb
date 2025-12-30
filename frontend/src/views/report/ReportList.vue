<template>
  <div class="report-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>日报管理</span>
          <el-button type="primary" @click="handleCreate">写日报</el-button>
        </div>
      </template>

      <!-- 搜索表单 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="日期范围">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            @change="handleDateChange"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
            <el-option label="草稿" value="draft" />
            <el-option label="已提交" value="submitted" />
            <el-option label="已审批" value="approved" />
            <el-option label="已驳回" value="rejected" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 日报表格 -->
      <el-table :data="reportList" v-loading="loading" border>
        <el-table-column prop="report_date" label="日期" width="120">
          <template #default="{ row }">
            {{ formatDate(row.report_date) }}
          </template>
        </el-table-column>
        <el-table-column prop="user.real_name" label="提交人" width="120" />
        <el-table-column prop="title" label="标题" min-width="150" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_late" label="补交" width="80">
          <template #default="{ row }">
            <el-tag v-if="row.is_late" type="warning">补交</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="submitted_at" label="提交时间" width="180">
          <template #default="{ row }">
            {{ row.submitted_at ? formatDateTime(row.submitted_at) : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleView(row)">查看</el-button>
            <el-button
              v-if="row.status === 'submitted' || row.status === 'rejected'"
              link
              type="primary"
              @click="handleEdit(row)"
            >
              编辑
            </el-button>
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import { getReportList, deleteReport } from '@/api/report'
import type { DailyReport, ListReportParams } from '@/api/report'

const router = useRouter()
const loading = ref(false)
const reportList = ref<DailyReport[]>([])
const dateRange = ref<[string, string] | null>(null)

const searchForm = reactive<ListReportParams>({
  status: '',
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0,
})

// 加载日报列表
const loadReportList = async () => {
  loading.value = true
  try {
    const params: ListReportParams = {
      page: pagination.page,
      page_size: pagination.pageSize,
      ...searchForm,
    }
    const res = await getReportList(params)
    reportList.value = res.data.list
    pagination.total = res.data.total
  } catch (error: any) {
    ElMessage.error(error.message || '加载日报列表失败')
  } finally {
    loading.value = false
  }
}

// 日期范围变化
const handleDateChange = (dates: [string, string] | null) => {
  if (dates) {
    searchForm.start_date = dates[0]
    searchForm.end_date = dates[1]
  } else {
    searchForm.start_date = undefined
    searchForm.end_date = undefined
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadReportList()
}

// 重置
const handleReset = () => {
  searchForm.status = ''
  dateRange.value = null
  searchForm.start_date = undefined
  searchForm.end_date = undefined
  pagination.page = 1
  loadReportList()
}

// 新增日报
const handleCreate = () => {
  router.push('/reports/create')
}

// 查看日报
const handleView = (row: DailyReport) => {
  router.push(`/reports/${row.id}`)
}

// 编辑日报
const handleEdit = (row: DailyReport) => {
  router.push(`/reports/${row.id}/edit`)
}

// 删除日报
const handleDelete = async (row: DailyReport) => {
  try {
    await ElMessageBox.confirm(`确定要删除 ${formatDate(row.report_date)} 的日报吗？`, '提示', {
      type: 'warning',
    })
    await deleteReport(row.id)
    ElMessage.success('删除成功')
    loadReportList()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

// 分页变化
const handleSizeChange = () => {
  loadReportList()
}

const handlePageChange = () => {
  loadReportList()
}

// 状态相关
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

// 格式化
const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD')
}

const formatDateTime = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

onMounted(() => {
  loadReportList()
})
</script>

<style scoped lang="scss">
.report-list {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: bold;
    font-size: 16px;
  }

  .search-form {
    margin-bottom: 20px;
  }

  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>



