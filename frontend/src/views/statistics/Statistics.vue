<template>
  <div class="statistics">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>统计分析</span>
          <el-radio-group v-model="statisticsType" @change="handleTypeChange">
            <el-radio-button label="personal">个人统计</el-radio-button>
            <el-radio-button label="team">团队统计</el-radio-button>
          </el-radio-group>
        </div>
      </template>

      <!-- 筛选条件 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item v-if="statisticsType === 'team'" label="部门">
          <el-select
            v-model="searchForm.department"
            placeholder="请选择部门"
            clearable
            style="width: 200px"
          >
            <el-option
              v-for="dept in departmentList"
              :key="dept"
              :label="dept"
              :value="dept"
            />
          </el-select>
        </el-form-item>
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
        <el-form-item>
          <el-button type="primary" @click="loadStatistics">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 个人统计 -->
      <div v-if="statisticsType === 'personal'" v-loading="loading">
        <div v-if="personalStats" class="stats-content">
          <!-- 概览卡片 -->
          <el-row :gutter="20" class="overview-cards">
            <el-col :span="6">
              <el-card class="stat-card">
                <div class="stat-value">{{ personalStats.total_reports }}</div>
                <div class="stat-label">总日报数</div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stat-card">
                <div class="stat-value">{{ personalStats.submitted_reports }}</div>
                <div class="stat-label">已提交</div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stat-card">
                <div class="stat-value">{{ personalStats.approved_reports }}</div>
                <div class="stat-label">已审批</div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stat-card">
                <div class="stat-value">{{ personalStats.submit_rate.toFixed(1) }}%</div>
                <div class="stat-label">提交率</div>
              </el-card>
            </el-col>
          </el-row>

          <!-- 图表区域 -->
          <el-row :gutter="20" style="margin-top: 20px">
            <el-col :span="12">
              <el-card>
                <template #header>状态分布</template>
                <div ref="statusChartRef" style="height: 300px"></div>
              </el-card>
            </el-col>
            <el-col :span="12">
              <el-card>
                <template #header>每日提交趋势</template>
                <div ref="dailyChartRef" style="height: 300px"></div>
              </el-card>
            </el-col>
          </el-row>
        </div>
      </div>

      <!-- 团队统计 -->
      <div v-if="statisticsType === 'team'" v-loading="loading">
        <div v-if="teamStats" class="stats-content">
          <!-- 概览卡片 -->
          <el-row :gutter="20" class="overview-cards">
            <el-col :span="6">
              <el-card class="stat-card">
                <div class="stat-value">{{ teamStats.total_members }}</div>
                <div class="stat-label">总成员数</div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stat-card">
                <div class="stat-value">{{ teamStats.active_members }}</div>
                <div class="stat-label">活跃成员</div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stat-card">
                <div class="stat-value">{{ teamStats.total_reports }}</div>
                <div class="stat-label">总日报数</div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stat-card">
                <div class="stat-value">{{ teamStats.submit_rate.toFixed(1) }}%</div>
                <div class="stat-label">团队提交率</div>
              </el-card>
            </el-col>
          </el-row>

          <!-- 成员统计表格 -->
          <el-card style="margin-top: 20px">
            <template #header>成员统计</template>
            <el-table :data="teamStats.member_statistics" border>
              <el-table-column prop="real_name" label="姓名" width="120" />
              <el-table-column prop="total_reports" label="总日报数" width="120" />
              <el-table-column prop="submit_rate" label="提交率" width="120">
                <template #default="{ row }">
                  {{ row.submit_rate.toFixed(1) }}%
                </template>
              </el-table-column>
              <el-table-column prop="late_count" label="补交次数" width="120" />
            </el-table>
          </el-card>

          <!-- 图表区域 -->
          <el-row :gutter="20" style="margin-top: 20px">
            <el-col :span="12">
              <el-card>
                <template #header>状态分布</template>
                <div ref="teamStatusChartRef" style="height: 300px"></div>
              </el-card>
            </el-col>
            <el-col :span="12">
              <el-card>
                <template #header>每日提交趋势</template>
                <div ref="teamDailyChartRef" style="height: 300px"></div>
              </el-card>
            </el-col>
          </el-row>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import * as echarts from 'echarts'
import { getPersonalStatistics, getTeamStatistics } from '@/api/statistics'
import type { PersonalStatistics, TeamStatistics, StatisticsParams } from '@/api/statistics'
import { getDepartments } from '@/api/common'

const loading = ref(false)
const statisticsType = ref<'personal' | 'team'>('personal')
const dateRange = ref<[string, string] | null>(null)
const personalStats = ref<PersonalStatistics | null>(null)
const teamStats = ref<TeamStatistics | null>(null)

const searchForm = reactive<StatisticsParams>({
  department: '',
})

// 图表引用
const statusChartRef = ref<HTMLElement>()
const dailyChartRef = ref<HTMLElement>()
const teamStatusChartRef = ref<HTMLElement>()
const teamDailyChartRef = ref<HTMLElement>()

const departmentList = ref<string[]>([])

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

// 统计类型变化
const handleTypeChange = () => {
  loadStatistics()
}

// 加载统计数据
const loadStatistics = async () => {
  loading.value = true
  try {
    if (statisticsType.value === 'personal') {
      const res = await getPersonalStatistics(searchForm)
      personalStats.value = res.data
      await nextTick()
      renderPersonalCharts()
    } else {
      const res = await getTeamStatistics(searchForm)
      teamStats.value = res.data
      await nextTick()
      renderTeamCharts()
    }
  } catch (error: any) {
    ElMessage.error(error.message || '加载统计数据失败')
  } finally {
    loading.value = false
  }
}

// 渲染个人统计图表
const renderPersonalCharts = () => {
  if (!personalStats.value) return

  // 状态分布饼图
  if (statusChartRef.value) {
    const chart = echarts.init(statusChartRef.value)
    const option = {
      tooltip: {
        trigger: 'item',
      },
      series: [
        {
          type: 'pie',
          data: Object.entries(personalStats.value.status_distribution).map(([name, value]) => ({
            name: getStatusText(name),
            value,
          })),
        },
      ],
    }
    chart.setOption(option)
  }

  // 每日提交趋势折线图
  if (dailyChartRef.value) {
    const chart = echarts.init(dailyChartRef.value)
    const data = personalStats.value.daily_submit_count.sort((a, b) => 
      dayjs(a.date).valueOf() - dayjs(b.date).valueOf()
    )
    const option = {
      tooltip: {
        trigger: 'axis',
      },
      xAxis: {
        type: 'category',
        data: data.map(item => item.date),
      },
      yAxis: {
        type: 'value',
      },
      series: [
        {
          type: 'line',
          data: data.map(item => item.count),
          smooth: true,
        },
      ],
    }
    chart.setOption(option)
  }
}

// 渲染团队统计图表
const renderTeamCharts = () => {
  if (!teamStats.value) return

  // 状态分布饼图
  if (teamStatusChartRef.value) {
    const chart = echarts.init(teamStatusChartRef.value)
    const option = {
      tooltip: {
        trigger: 'item',
      },
      series: [
        {
          type: 'pie',
          data: Object.entries(teamStats.value.status_distribution).map(([name, value]) => ({
            name: getStatusText(name),
            value,
          })),
        },
      ],
    }
    chart.setOption(option)
  }

  // 每日提交趋势折线图
  if (teamDailyChartRef.value) {
    const chart = echarts.init(teamDailyChartRef.value)
    const data = teamStats.value.daily_submit_trend.sort((a, b) => 
      dayjs(a.date).valueOf() - dayjs(b.date).valueOf()
    )
    const option = {
      tooltip: {
        trigger: 'axis',
      },
      xAxis: {
        type: 'category',
        data: data.map(item => item.date),
      },
      yAxis: {
        type: 'value',
      },
      series: [
        {
          type: 'line',
          data: data.map(item => item.count),
          smooth: true,
        },
      ],
    }
    chart.setOption(option)
  }
}

// 重置
const handleReset = () => {
  searchForm.department = ''
  dateRange.value = null
  searchForm.start_date = undefined
  searchForm.end_date = undefined
  loadStatistics()
}

// 状态文本映射
const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    draft: '草稿',
    submitted: '已提交',
    approved: '已审批',
    rejected: '已驳回',
  }
  return map[status] || status
}

// 加载部门列表
const loadDepartmentList = async () => {
  try {
    const res = await getDepartments()
    if (res.data && res.data.length > 0) {
      departmentList.value = res.data
    } else {
      console.warn('部门列表为空')
    }
  } catch (error: any) {
    console.error('加载部门列表失败:', error)
    ElMessage.error('加载部门列表失败: ' + (error.message || '未知错误'))
  }
}

onMounted(() => {
  loadStatistics()
  loadDepartmentList()
})
</script>

<style scoped lang="scss">
.statistics {
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

  .overview-cards {
    margin-bottom: 20px;

    .stat-card {
      text-align: center;

      .stat-value {
        font-size: 32px;
        font-weight: bold;
        color: #409eff;
        margin-bottom: 10px;
      }

      .stat-label {
        font-size: 14px;
        color: #909399;
      }
    }
  }
}
</style>

