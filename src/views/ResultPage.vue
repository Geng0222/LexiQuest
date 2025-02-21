// ResultPage.vue
<template>
  <div class="result-container">
    <el-card class="result-card">
      <h1 class="title">🎉 測驗結束</h1>

      <div v-if="loading" class="loading">
        <el-skeleton animated :rows="3" />
      </div>

      <div v-else>
        <div class="result-info">
          <p class="score">
            ✅ <strong>正確答案數</strong>：
            <el-tag :type="scoreTagType">{{ score }} / {{ totalQuestions }}</el-tag>
          </p>
          <p class="time">
            ⏳ <strong>花費時間</strong>：{{ time }} 秒
          </p>
          <p class="message">{{ resultMessage }}</p>
        </div>

        <el-divider />

        <!-- ✅ 顯示當前模式 -->
        <p class="mode-info">
          🔍 模式：<el-tag :type="isUsingAPI ? 'success' : 'info'">{{ isUsingAPI ? 'API 模式' : '靜態模式' }}</el-tag>
        </p>

        <div class="button-group">
          <el-button type="primary" icon="el-icon-refresh" @click="restartQuiz">
            🔄 重新開始測驗
          </el-button>
          <el-button type="info" icon="el-icon-arrow-left" @click="goHome">
            ⬅️ 回到首頁
          </el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script>
import { useRouter, useRoute } from "vue-router";
import { useQuizStore } from "../stores/quizStore";
import { computed, ref } from "vue";
import { useApiStatus } from "../utils/useApiStatus";

export default {
  setup() {
    const router = useRouter();
    const route = useRoute();
    const store = useQuizStore();
    const { isUsingAPI } = useApiStatus();

    const loading = ref(false);

    // ✅ 獲取測驗的 `category` 和 `filename`
    const category = ref(route.params.category);
    const filename = ref(route.params.filename);

    // ✅ 獲取 `totalQuestions`，優先從 `store.words.length` 取得
    const totalQuestions = computed(() => store.words.length || route.query.totalQuestions || 1);

    // ✅ `score` 來自 `store.correctAnswers`
    const score = computed(() => store.correctAnswers);

    // ✅ `time` 只在靜態模式下顯示
    const time = ref(isUsingAPI.value ? 0 : (route.query.time || 0));

    function restartQuiz() {
      store.currentIndex = 0;
      store.correctAnswers = 0;
      router.push(`/quiz/${category.value}/${filename.value}`);
    }

    function goHome() {
      router.push("/");
    }

    // ✅ 顯示成績訊息
    const resultMessage = computed(() => {
      if (score.value >= totalQuestions.value * 0.8) return "🎉 太棒了！你答對了大部分題目！";
      if (score.value >= totalQuestions.value * 0.4) return "💡 不錯喔！還可以更進步！";
      return "😅 加油！下次一定能更好！";
    });

    // ✅ 依照成績變更顏色標籤
    const scoreTagType = computed(() => {
      if (score.value >= totalQuestions.value * 0.8) return "success"; // 綠色
      if (score.value >= totalQuestions.value * 0.4) return "warning"; // 黃色
      return "danger"; // 紅色
    });

    // ✅ 調試輸出，確認變數是否正確
    console.log("Score:", score.value);
    console.log("Total Questions:", totalQuestions.value);

    return {
      score,
      totalQuestions,
      time,
      restartQuiz,
      goHome,
      resultMessage,
      scoreTagType,
      loading,
      isUsingAPI,
      category,
      filename,
    };
  }
};
</script>

<style scoped>
/* 🌟 結果頁面 */
.result-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f5f5f5;
}

/* 📌 卡片樣式 */
.result-card {
  width: 400px;
  padding: 20px;
  border-radius: 12px;
  box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.1);
  background-color: white;
  text-align: center;
  animation: fadeIn 0.5s ease-in-out;
}

/* 🏆 標題 */
.title {
  font-size: 24px;
  font-weight: bold;
  color: #2c3e50;
}

/* 🔢 分數 & 時間 */
.result-info {
  font-size: 18px;
  color: #333;
  margin: 15px 0;
}

.score, .time {
  font-size: 18px;
}

.message {
  font-size: 16px;
  font-weight: bold;
  color: #555;
  margin-top: 10px;
}

/* 🔍 模式顯示 */
.mode-info {
  font-size: 14px;
  font-weight: bold;
  color: #666;
  margin-bottom: 10px;
}

/* 🎮 按鈕區 */
.button-group {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
</style>
