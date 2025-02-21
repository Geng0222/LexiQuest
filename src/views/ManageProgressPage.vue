<template>
  <div class="manage-progress">
    <ProgressTable :progressList="displayedProgressList" @delete="confirmDelete" />

    <!-- 🔹 刪除確認對話框 -->
    <el-dialog v-model="showDeleteDialog" title="確認刪除" width="400px">
      <p class="text-lg">確定要刪除 <strong>{{ selectedItem?.category }}/{{ selectedItem?.filename }}</strong> 的學習進度嗎？</p>
      <template #footer>
        <el-button @click="showDeleteDialog = false">取消</el-button>
        <el-button type="danger" @click="handleDelete">確定刪除</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, onMounted, computed } from "vue";
import { useApiStatus } from "../utils/useApiStatus"; // ✅ 引入 API 檢查
import { fetchProgressList, deleteProgress } from "../utils/useProgress";
import ProgressTable from "../components/ProgressTable.vue";

export default {
  components: { ProgressTable },
  setup() {
    const { isUsingAPI, checkAPIStatus } = useApiStatus();
    const progressList = ref([]);
    const staticProgressList = ref([
      { category: "daily", filename: "advanced", words: 2 },
      { category: "daily", filename: "basic", words: 3 },
    ]);

    const showDeleteDialog = ref(false);
    const selectedItem = ref(null);

    async function loadProgressList() {
      if (!isUsingAPI.value) return;
      try {
        progressList.value = await fetchProgressList();
      } catch (error) {
        console.error("❌ 獲取學習進度失敗:", error);
      }
    }

    function confirmDelete(item) {
      selectedItem.value = item;
      showDeleteDialog.value = true;
    }

    async function handleDelete() {
      if (!selectedItem.value || !isUsingAPI.value) return;
      try {
        await deleteProgress(selectedItem.value.category, selectedItem.value.filename);
        showDeleteDialog.value = false;
        loadProgressList();
      } catch (error) {
        console.error("❌ 刪除失敗:", error);
      }
    }

    onMounted(async () => {
      await checkAPIStatus();
      await loadProgressList();
    });

    // ✅ 根據 `isUsingAPI` 決定顯示 API 或靜態資料
    const displayedProgressList = computed(() =>
      isUsingAPI.value ? progressList.value : staticProgressList.value
    );

    return {
      isUsingAPI,
      displayedProgressList,
      showDeleteDialog,
      selectedItem,
      confirmDelete,
      handleDelete,
    };
  },
};
</script>

<style scoped>
.manage-progress {
  padding: 20px;
  display: flex;
  justify-content: center;
}
</style>
