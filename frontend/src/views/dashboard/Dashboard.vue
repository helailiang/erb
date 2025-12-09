<template>
  <div class="dashboard">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>欢迎使用日常工作统计系统</span>
        </div>
      </template>
      <div class="welcome-content">
        <p>您好，{{ userStore.user?.real_name }}！</p>
        <p>系统正在开发中，更多功能即将上线...</p>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
</script>

<style scoped lang="scss">
.dashboard {
  .card-header {
    font-weight: bold;
    font-size: 16px;
  }

  .welcome-content {
    p {
      margin: 10px 0;
      font-size: 14px;
      color: #606266;
    }
  }
}
</style>

