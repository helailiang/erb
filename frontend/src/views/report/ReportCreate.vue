<template>
  <div class="report-create">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>写日报</span>
          <div>
            <el-button @click="handleCancel">取消</el-button>
            <el-button @click="handleSaveDraft">保存草稿</el-button>
            <el-button type="primary" :loading="submitting" @click="handleSubmit">提交</el-button>
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
            :disabled-date="disabledDate"
          />
          <el-checkbox v-model="form.is_late" style="margin-left: 20px">补交</el-checkbox>
        </el-form-item>

        <el-form-item v-if="form.is_late" label="补交原因">
          <el-input
            v-model="form.late_reason"
            type="textarea"
            :rows="2"
            placeholder="请输入补交原因"
          />
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
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import { createReport, saveDraft, getDraft } from '@/api/report'
import type { CreateReportRequest, SaveDraftRequest } from '@/api/report'

const router = useRouter()
const editorRef = ref<HTMLElement>()
const submitting = ref(false)
let editor: any = null

const form = reactive<CreateReportRequest & { late_reason?: string }>({
  report_date: dayjs().format('YYYY-MM-DD'),
  title: '',
  content: '',
  is_late: false,
  late_reason: '',
})

// 禁用未来日期
const disabledDate = (time: Date) => {
  return time.getTime() > Date.now()
}

// 初始化编辑器（使用简单的textarea，后续可以集成富文本编辑器）
onMounted(() => {
  if (editorRef.value) {
    // 创建简单的文本编辑器
    const textarea = document.createElement('textarea')
    textarea.style.width = '100%'
    textarea.style.minHeight = '400px'
    textarea.style.padding = '10px'
    textarea.style.border = '1px solid #dcdfe6'
    textarea.style.borderRadius = '4px'
    textarea.placeholder = '请输入今日工作内容...'
    editorRef.value.appendChild(textarea)
    editor = textarea

    // 尝试加载草稿
    loadDraft()

    // 自动保存草稿（每30秒）
    const autoSaveTimer = setInterval(() => {
      if (form.content) {
        handleAutoSave()
      }
    }, 30000)

    onBeforeUnmount(() => {
      clearInterval(autoSaveTimer)
    })
  }
})

// 加载草稿
const loadDraft = async () => {
  try {
    const res = await getDraft(form.report_date)
    if (res.data) {
      form.title = res.data.title || ''
      form.content = res.data.content
      if (editor) {
        editor.value = res.data.content
      }
      ElMessage.info('已加载草稿')
    }
  } catch (error) {
    // 没有草稿，忽略错误
  }
}

// 自动保存草稿
const handleAutoSave = async () => {
  if (!editor) return
  form.content = editor.value

  try {
    const draftData: SaveDraftRequest = {
      report_date: form.report_date,
      content: form.content,
      title: form.title,
    }
    await saveDraft(draftData)
  } catch (error) {
    // 静默失败
  }
}

// 保存草稿
const handleSaveDraft = async () => {
  if (!editor) return
  form.content = editor.value

  if (!form.content.trim()) {
    ElMessage.warning('请输入日报内容')
    return
  }

  try {
    const draftData: SaveDraftRequest = {
      report_date: form.report_date,
      content: form.content,
      title: form.title,
    }
    await saveDraft(draftData)
    ElMessage.success('草稿保存成功')
  } catch (error: any) {
    ElMessage.error(error.message || '保存草稿失败')
  }
}

// 提交日报
const handleSubmit = async () => {
  if (!editor) return
  form.content = editor.value

  if (!form.content.trim()) {
    ElMessage.warning('请输入日报内容')
    return
  }

  submitting.value = true
  try {
    const reportData: CreateReportRequest = {
      report_date: form.report_date,
      title: form.title || undefined,
      content: form.content,
      is_late: form.is_late,
      late_reason: form.late_reason || undefined,
    }
    console.log(reportData)
    await createReport(reportData)
    ElMessage.success('日报提交成功')
    router.push('/reports')
  } catch (error: any) {
    ElMessage.error(error.message || '提交失败')
  } finally {
    submitting.value = false
  }
}

// 取消
const handleCancel = () => {
  router.back()
}
</script>

<style scoped lang="scss">
.report-create {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: bold;
    font-size: 16px;
  }
}
</style>

