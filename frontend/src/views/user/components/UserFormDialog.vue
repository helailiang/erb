<template>
  <el-dialog
    v-model="dialogVisible"
    :title="userId ? '编辑用户' : '新增用户'"
    width="600px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item label="用户名" prop="username">
        <el-input
          v-model="form.username"
          placeholder="请输入用户名"
          :disabled="!!userId"
        />
      </el-form-item>

      <el-form-item v-if="!userId" label="密码" prop="password">
        <el-input
          v-model="form.password"
          type="password"
          placeholder="请输入密码（至少8位）"
          show-password
        />
      </el-form-item>

      <el-form-item label="真实姓名" prop="real_name">
        <el-input v-model="form.real_name" placeholder="请输入真实姓名" />
      </el-form-item>

      <el-form-item label="邮箱" prop="email">
        <el-input v-model="form.email" placeholder="请输入邮箱" />
      </el-form-item>

      <el-form-item label="手机号" prop="phone">
        <el-input v-model="form.phone" placeholder="请输入手机号" />
      </el-form-item>

      <el-form-item label="部门" prop="department">
        <el-select v-model="form.department" placeholder="请选择部门" style="width: 100%">
          <el-option
            v-for="dept in departmentList"
            :key="dept"
            :label="dept"
            :value="dept"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="职位" prop="position">
        <el-input v-model="form.position" placeholder="请输入职位" />
      </el-form-item>

      <el-form-item label="员工工号" prop="employee_no">
        <el-input v-model="form.employee_no" placeholder="请输入员工工号" />
      </el-form-item>

      <el-form-item v-if="userId" label="状态" prop="status">
        <el-select v-model="form.status" placeholder="请选择状态">
          <el-option label="正常" value="normal" />
          <el-option label="已禁用" value="disabled" />
        </el-select>
      </el-form-item>

      <el-form-item label="角色" prop="role_ids">
        <el-select
          v-model="form.role_ids"
          multiple
          placeholder="请选择角色"
          style="width: 100%"
        >
          <el-option
            v-for="role in roleList"
            :key="role.id"
            :label="role.name"
            :value="role.id"
          />
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
import { ref, reactive, watch, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { createUser, updateUser, getUserDetail } from '@/api/user'
import { getRoleList } from '@/api/role'
import { getDepartments } from '@/api/common'
import type { CreateUserRequest, UpdateUserRequest } from '@/api/user'
import type { Role } from '@/api/role'

const props = defineProps<{
  modelValue: boolean
  userId?: number | null
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  success: []
}>()

const dialogVisible = ref(props.modelValue)
const formRef = ref<FormInstance>()
const loading = ref(false)
const roleList = ref<Role[]>([])
const departmentList = ref<string[]>([])

const form = reactive<CreateUserRequest & UpdateUserRequest>({
  username: '',
  password: '',
  real_name: '',
  email: '',
  phone: '',
  department: '',
  position: '',
  employee_no: '',
  status: 'normal',
  role_ids: [],
})

const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 8, message: '密码长度至少8位', trigger: 'blur' },
  ],
  real_name: [
    { required: true, message: '请输入真实姓名', trigger: 'blur' },
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' },
  ],
}

// 加载角色列表
const loadRoleList = async () => {
  try {
    const res = await getRoleList({ page: 1, page_size: 100 })
    roleList.value = res.data.list
  } catch (error: any) {
    ElMessage.error(error.message || '加载角色列表失败')
  }
}

// 加载部门列表
const loadDepartmentList = async () => {
  try {
    const res = await getDepartments()
    departmentList.value = res.data
  } catch (error: any) {
    ElMessage.error(error.message || '加载部门列表失败')
  }
}

// 加载用户详情
const loadUserDetail = async () => {
  if (!props.userId) return

  try {
    const res = await getUserDetail(props.userId)
    Object.assign(form, {
      username: res.data.username,
      real_name: res.data.real_name,
      email: res.data.email,
      phone: res.data.phone,
      department: res.data.department,
      position: res.data.position,
      employee_no: res.data.employee_no,
      status: res.data.status,
      role_ids: (res.data as any).roles?.map((r: Role) => r.id) || [],
    })
  } catch (error: any) {
    ElMessage.error(error.message || '加载用户详情失败')
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      if (props.userId) {
        const { password, username, ...updateData } = form
        await updateUser(props.userId, updateData)
        ElMessage.success('更新成功')
      } else {
        await createUser(form as CreateUserRequest)
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
    username: '',
    password: '',
    real_name: '',
    email: '',
    phone: '',
    department: '',
    position: '',
    employee_no: '',
    status: 'normal',
    role_ids: [],
  })
}

// 监听对话框显示状态
watch(() => props.modelValue, (val) => {
  dialogVisible.value = val
  if (val) {
    loadRoleList()
    loadDepartmentList()
    if (props.userId) {
      loadUserDetail()
    }
  }
})

onMounted(() => {
  loadDepartmentList()
})

watch(dialogVisible, (val) => {
  emit('update:modelValue', val)
})
</script>

