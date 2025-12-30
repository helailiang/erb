<template>
  <el-dialog
    v-model="dialogVisible"
    title="权限管理"
    width="700px"
    @close="handleClose"
  >
    <div v-loading="loading">
      <el-tree
        ref="treeRef"
        :data="permissionTree"
        :props="{ children: 'children', label: 'name' }"
        show-checkbox
        node-key="id"
        :default-checked-keys="checkedKeys"
        :default-expand-all="true"
      />
    </div>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getRoleDetail, updateRole, getAllPermissions } from '@/api/role'
import type { Role, Permission } from '@/api/role'

const props = defineProps<{
  modelValue: boolean
  roleId: number | null
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  success: []
}>()

const dialogVisible = ref(props.modelValue)
const loading = ref(false)
const saving = ref(false)
const treeRef = ref()
const permissionTree = ref<any[]>([])
const checkedKeys = ref<number[]>([])

// 加载权限数据
const loadPermissions = async () => {
  if (!props.roleId) return

  loading.value = true
  try {
    // 获取所有权限
    const permissionsRes = await getAllPermissions()
    const allPermissions = permissionsRes.data

    // 获取角色详情（包含已分配的权限）
    const roleRes = await getRoleDetail(props.roleId)
    const role = roleRes.data

    // 构建权限树（按模块分组）
    const moduleMap = new Map<string, Permission[]>()
    allPermissions.forEach(perm => {
      const module = perm.module || '其他'
      if (!moduleMap.has(module)) {
        moduleMap.set(module, [])
      }
      moduleMap.get(module)!.push(perm)
    })

    const tree: any[] = []
    moduleMap.forEach((perms, module) => {
      tree.push({
        id: `module_${module}`,
        name: module,
        children: perms.map(perm => ({
          id: perm.id,
          name: `${perm.name} (${perm.code})`,
        })),
      })
    })
    permissionTree.value = tree

    // 设置已选中的权限
    if (role.permissions) {
      checkedKeys.value = role.permissions.map((p: Permission) => p.id)
    }
  } catch (error: any) {
    ElMessage.error(error.message || '加载权限数据失败')
  } finally {
    loading.value = false
  }
}

// 保存权限
const handleSave = async () => {
  if (!props.roleId) return

  saving.value = true
  try {
    const checkedNodes = treeRef.value?.getCheckedNodes(false, true) || []
    const permissionIds = checkedNodes
      .filter((node: any) => typeof node.id === 'number')
      .map((node: any) => node.id)

    await updateRole(props.roleId, {
      permission_ids: permissionIds,
    })
    ElMessage.success('权限保存成功')
    emit('success')
    handleClose()
  } catch (error: any) {
    ElMessage.error(error.message || '保存权限失败')
  } finally {
    saving.value = false
  }
}

// 关闭对话框
const handleClose = () => {
  dialogVisible.value = false
  emit('update:modelValue', false)
  permissionTree.value = []
  checkedKeys.value = []
}

// 监听对话框显示状态
watch(() => props.modelValue, (val) => {
  dialogVisible.value = val
  if (val && props.roleId) {
    loadPermissions()
  }
})

watch(dialogVisible, (val) => {
  emit('update:modelValue', val)
})
</script>

<style scoped lang="scss">
:deep(.el-tree-node__label) {
  font-size: 14px;
}
</style>

