<template>
  <el-card shadow="hover" class="progress-card">
    <template #header>
      <div class="header-container">
        <!-- 標題 & API 模式標籤 -->
        <div class="title-wrapper">
          <span class="title-text">📊 學習進度管理</span>
          <el-tag :type="isUsingAPI ? 'success' : 'danger'" size="large">
            {{ isUsingAPI ? '🔗 API模式' : '🌐 靜態模式' }}
          </el-tag>
        </div>

        <!-- 搜尋框 -->
        <el-input
          v-model="searchQuery"
          placeholder="🔍 搜尋題庫..."
          prefix-icon="el-icon-search"
          size="small"
          class="search-input"
        />
      </div>
    </template>

    <!-- 靜態模式警告 -->
    <el-alert v-if="!isUsingAPI" title="目前處於靜態模式，無法刪除學習進度" type="warning" show-icon />

    <el-table :data="filteredProgressList" border stripe class="progress-table">
      <el-table-column prop="category" label="分類" width="180" sortable />
      <el-table-column prop="filename" label="題庫名稱" width="220" sortable />
      <el-table-column prop="words" label="單字數量" width="120" align="center" />

      <el-table-column label="刪除學習進度" width="150" align="center">
        <template #default="{ row }">
          <el-tooltip content="刪除學習進度" placement="top">
            <el-button
              :type="isUsingAPI ? 'danger' : 'info'"
              size="small"
              @click="isUsingAPI ? $emit('delete', row) : null"
              class="delete-btn"
              :disabled="!isUsingAPI"
            >
              刪除
            </el-button>
          </el-tooltip>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>

<script>
import { ref, computed, onMounted } from "vue";
import { useApiStatus } from "../utils/useApiStatus";

export default {
  props: {
    progressList: Array,
  },
  emits: ["delete"],
  setup(props) {
    const { isUsingAPI, checkAPIStatus } = useApiStatus();
    const searchQuery = ref("");

    const filteredProgressList = computed(() =>
      searchQuery.value
        ? props.progressList.filter((item) =>
            item.filename.toLowerCase().includes(searchQuery.value.toLowerCase())
          )
        : props.progressList
    );

    onMounted(checkAPIStatus);

    return {
      isUsingAPI,
      searchQuery,
      filteredProgressList,
    };
  },
};
</script>

<style scoped>
.progress-card {
  max-width: 800px;
  margin: 20px auto;
  padding: 20px;
}

/* 🚀 調整標題與搜尋框的排版 */
.header-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px; /* 增加間距 */
}

/* 🚀 標題與標籤的排版 */
.title-wrapper {
  display: flex;
  align-items: center;
  gap: 10px;
}

.title-text {
  font-size: 1.2rem;
  font-weight: bold;
}

/* 🚀 調整搜尋框寬度 */
.search-input {
  flex-grow: 1;
  min-width: 200px;
  max-width: 250px;
}

/* 表格間距 */
.progress-table {
  margin-top: 10px;
}

/* 刪除按鈕 */
.delete-btn {
  font-weight: bold;
  width: 60px;
  padding: 5px 10px;
}

.delete-btn:disabled {
  background-color: #d3d3d3;
  color: #ffffff;
}
</style>
