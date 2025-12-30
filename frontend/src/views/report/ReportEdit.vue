<template>
  <div class="report-edit" v-loading="loading">
    <el-card v-if="report">
      <template #header>
        <div class="card-header">
          <span>编辑日报</span>
          <div>
            <el-button @click="handleCancel">取消</el-button>
            <el-button type="primary" :loading="submitting" @click="handleSubmit">保存</el-button>
          </div>
        </div>
      </template>

      <el-form :model="form" label-width="100px">
        <el-form-item label="日期">
          <el-date-picker
            v-model="form.report_date"
            type="date"
            placeholder="选择日期"
            value-format="YYYY-MM-DD"
            disabled
          />
          <el-tag v-if="report.is_late" type="warning" style="margin-left: 20px">补交</el-tag>
        </el-form-item>

        <el-form-item label="标题">
          <el-input v-model="form.title" placeholder="请输入日报标题（可选）" />
        </el-form-item>

        <el-form-item label="内容" required>
          <div ref="editorRef" style="min-height: 400px"></div>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import { getReportDetail, updateReport } from '@/api/report'
import type { DailyReport, UpdateReportRequest } from '@/api/report'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const submitting = ref(false)
const report = ref<DailyReport | null>(null)
const editorRef = ref<HTMLElement>()
let editor: any = null

const form = reactive({
  report_date: '',
  title: '',
  content: '',
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
    form.report_date = res.data.report_date
    form.title = res.data.title || ''
    form.content = res.data.content

    // 初始化编辑器
    if (editorRef.value && !editor) {
      const textarea = document.createElement('textarea')
      textarea.style.width = '100%'
      textarea.style.minHeight = '400px'
      textarea.style.padding = '10px'
      textarea.style.border = '1px solid #dcdfe6'
      textarea.style.borderRadius = '4px'
      textarea.value = form.content
      editorRef.value.appendChild(textarea)
      editor = textarea
    } else if (editor) {
      editor.value = form.content
    }
  } catch (error: any) {
    ElMessage.error(error.message || '加载日报详情失败')
    router.back()
  } finally {
    loading.value = false
  }
}

// 提交修改
const handleSubmit = async () => {
  if (!editor || !report.value) return

  form.content = editor.value

  if (!form.content.trim()) {
    ElMessage.warning('请输入日报内容')
    return
  }

  submitting.value = true
  try {
    const updateData: UpdateReportRequest = {
      title: form.title,
      content: form.content,
    }
    await updateReport(report.value.id, updateData)
    ElMessage.success('修改成功')
    router.push(`/reports/${report.value.id}`)
  } catch (error: any) {
    ElMessage.error(error.message || '修改失败')
  } finally {
    submitting.value = false
  }
}

// 取消
const handleCancel = () => {
  router.back()
}

onMounted(() => {
  loadReportDetail()
})
</script>

<style scoped lang="scss">
.report-edit {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: bold;
    font-size: 16px;
  }
}
</style>
