// WordlistPage.vue
<template>
    <div class="wordlist-container">
      <div class="wordlist-header">
        <h1>📖 {{ formattedTitle }}</h1>
        <el-tag :type="isUsingAPI ? 'success' : 'danger'" size="large">
          {{ isUsingAPI ? '🔗 API模式' : '🌐 靜態模式' }}
        </el-tag>
      </div>
  
      <div v-if="loading" class="loading">
        <el-skeleton animated :rows="5" />
      </div>
  
      <div v-else-if="errorMessage" class="error">
        <p>{{ errorMessage }}</p>
      </div>
  
      <div v-else>
        <el-table :data="wordlist" border>
          <el-table-column prop="word" label="單字" width="250" />
          <el-table-column prop="translation" label="翻譯" width="200" />
          <el-table-column prop="type" label="詞性" width="200" />
        </el-table>
  
        <!-- ✅ 新增「開始測驗」按鈕 -->
        <el-button type="primary" @click="startQuiz" class="start-quiz-btn">
          🎯 開始測驗
        </el-button>
      </div>
    </div>
  </template>
  
  <script>
  import { computed, ref, onMounted } from "vue";
  import { useRoute, useRouter } from "vue-router";
  import { loadWordlist } from "../utils/loadWordlist";
  import { useApiStatus } from "../utils/useApiStatus"; // ✅ 取得 API 狀態
  
  export default {
    setup() {
      const route = useRoute();
      const router = useRouter();
      const { isUsingAPI, checkAPIStatus } = useApiStatus();
  
      const wordlist = ref([]);
      const loading = ref(true);
      const errorMessage = ref("");
  
      const category = ref(route.params.category);
      const filename = ref(route.params.filename);
  
      async function fetchWordlist() {
        await checkAPIStatus(); // ✅ 先檢查 API 狀態
        loading.value = true;
  
        try {
          const response = await loadWordlist(category.value, filename.value);
  
          if (!response.success) {
            throw new Error(response.error || "無法加載題庫");
          }
  
          wordlist.value = response.data;
        } catch (error) {
          errorMessage.value = `❌ 無法讀取題庫: ${error.message}`;
        } finally {
          loading.value = false;
        }
      }
  
      onMounted(fetchWordlist);
  
      // ✅ 產生題庫標題
      const formattedTitle = computed(() => {
        return `單字列表 - ${category.value} (${filename.value})`;
      });
  
      // ✅ 開始測驗按鈕功能
      function startQuiz() {
        router.push(`/quiz/${category.value}/${filename.value}`);
      }
  
      return { wordlist, loading, errorMessage, formattedTitle, isUsingAPI, startQuiz };
    }
  };
  </script>
  
  <style scoped>
  .wordlist-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 20px;
  }
  
  .wordlist-header {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 20px;
  }
  
  .start-quiz-btn {
    margin-top: 20px;
    font-size: 18px;
    padding: 10px 20px;
  }
  </style>
  