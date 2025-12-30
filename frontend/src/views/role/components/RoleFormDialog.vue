<template>
  <el-dialog
    v-model="dialogVisible"
    :title="roleId ? '编辑角色' : '新增角色'"
    width="600px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item label="角色名" prop="name">
        <el-input
          v-model="form.name"
          placeholder="请输入角色名"
          :disabled="isSystemRole"
        />
      </el-form-item>

      <el-form-item label="描述" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="3"
          placeholder="请输入角色描述"
        />
      </el-form-item>

      <el-form-item label="状态" prop="status">
        <el-select v-model="form.status" placeholder="请选择状态">
          <el-option label="启用" value="enabled" />
          <el-option label="禁用" value="disabled" />
        </el-select>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" :loading="loading" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { createRole, updateRole, getRoleDetail } from '@/api/role'
import type { CreateRoleRequest, UpdateRoleRequest, Role } from '@/api/role'

const props = defineProps<{
  modelValue: boolean
  roleId?: number | null
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  success: []
}>()

const dialogVisible = ref(props.modelValue)
const formRef = ref<FormInstance>()
const loading = ref(false)
const roleDetail = ref<Role | null>(null)

const form = reactive<CreateRoleRequest & UpdateRoleRequest>({
  name: '',
  description: '',
  status: 'enabled',
})

const isSystemRole = computed(() => {
  return roleDetail.value?.is_system || false
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入角色名', trigger: 'blur' },
  ],
}

// 加载角色详情
const loadRoleDetail = async () => {
  if (!props.roleId) return

  try {
    const res = await getRoleDetail(props.roleId)
    roleDetail.value = res.data
    Object.assign(form, {
      name: res.data.name,
      description: res.data.description,
      status: res.data.status,
    })
  } catch (error: any) {
    ElMessage.error(error.message || '加载角色详情失败')
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      if (props.roleId) {
        await updateRole(props.roleId, form)
        ElMessage.success('更新成功')
      } else {
        await createRole(form as CreateRoleRequest)
        ElMessage.success('创建成功')
      }
      emit('success')
      handleClose()
    } catch (error: any) {
      ElMessage.error(error.message || '操作失败')
    } finally {
      loading.value = false
    }
  })
}

// 关闭对话框
const handleClose = () => {
  dialogVisible.value = false
  emit('update:modelValue', false)
  formRef.value?.resetFields()
  Object.assign(form, {
    name: '',
    description: '',
    status: 'enabled',
  })
  roleDetail.value = null
}

// 监听对话框显示状态
watch(() => props.modelValue, (val) => {
  dialogVisible.value = val
  if (val && props.roleId) {
    loadRoleDetail()
  }
})

watch(dialogVisible, (val) => {
  emit('update:modelValue', val)
})
</script>

