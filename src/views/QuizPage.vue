//QuizPage.vue
<!-- QuizPage.vue -->
<template>
  <div class="quiz-container">
    <div class="quiz-header">
      <h1>📖 {{ formattedTitle }}</h1>
      <el-tag :type="isUsingAPI ? 'success' : 'danger'" size="large">
        {{ isUsingAPI ? '🔗 API模式' : '🌐 靜態模式' }}
      </el-tag>
    </div>

    <!-- 題目進度 -->
    <div class="progress">
      <p>📌 題目進度：{{ currentQuestionNumber }} / {{ totalQuestions }}</p>
    </div>

    <div v-if="loading" class="loading">
      <el-skeleton animated :rows="3" />
    </div>

    <div v-else-if="errorMessage" class="error">
      <p>{{ errorMessage }}</p>
    </div>

    <div v-else>
      <!-- 渲染題目卡片 -->
      <QuizCard
        :currentWord="currentWord"
        v-model:userInput="userInput"
        :feedback="feedback"
        :elapsedTime="elapsedTime"
        :loading="loading"
        :errorMessage="errorMessage"
        @check-answer="checkAnswer"
        @go-back="goBack"
        @next-question="goToNextQuestion"
        @end-quiz="endQuiz"
      />
    </div>

    <!-- 切換字典視窗的按鈕 -->
    <el-button
      class="toggle-dict-btn"
      type="text"
      @click="toggleDictionary"
      v-if="currentWord"
    >
      {{ showDictionary ? "隱藏字典資訊" : "顯示字典資訊" }}
    </el-button>

    <!-- 漂浮字典視窗 -->
    <DictionaryFloat
      :show="showDictionary"
      :currentData="currentWord"
      @close="showDictionary = false"
    />
  </div>
</template>

<script>
import { useQuizStore } from "../stores/quizStore";
import { computed, ref, onMounted, onUnmounted, nextTick, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useApiStatus } from "../utils/useApiStatus";
import { submitQuizResults } from "../utils/useQuizs";
import QuizCard from "../components/QuizCard.vue";
import DictionaryFloat from "../components/DictionaryFloat.vue";

export default {
  components: { QuizCard, DictionaryFloat },
  setup() {
    const store = useQuizStore();
    const route = useRoute();
    const router = useRouter();
    const { isUsingAPI, checkAPIStatus } = useApiStatus();

    const userInput = ref("");
    const feedback = ref("");
    const loading = ref(true);
    const errorMessage = ref("");
    const startTime = ref(null);
    const elapsedTime = ref(0);
    let timer = null;
    // 🔹 使用單一 Audio 實例以避免重複建立新的對象
    const audioInstance = new Audio();

    const category = ref(route.params.category || "daily");
    const filename = ref(route.params.filename || "basic");

    // 🔹 預設開啟字典視窗
    const showDictionary = ref(true);
    const toggleDictionary = () => {
      showDictionary.value = !showDictionary.value;
    };

    // 🔹 取得當前單字
    const currentWord = computed(() => store.words[store.currentIndex] || null);
    const totalQuestions = computed(() => store.words.length);
    const currentQuestionNumber = computed(() => store.currentIndex + 1);

    // 🔹 播放語音（確保同時只有一個音檔播放）
    function playAudio() {
      if (!currentWord.value?.audio) return;

      // 若音檔正在播放，先暫停並重置播放位置
      if (!audioInstance.paused) {
        audioInstance.pause();
        audioInstance.currentTime = 0;
      }
      // 更新音檔來源，並播放
      audioInstance.src = currentWord.value.audio;
      audioInstance.play().catch(err => console.warn("🎵 播放音檔失敗:", err));
    }

    // 🔹 當題目索引變更時播放音檔
    watch(() => store.currentIndex, () => {
      playAudio();
    });

    onMounted(async () => {
      startTime.value = Date.now();
      timer = setInterval(() => {
        elapsedTime.value = Math.floor((Date.now() - startTime.value) / 1000);
      }, 1000);

      try {
        await checkAPIStatus();
        await store.loadRandomQuiz(category.value, filename.value, 10);
        if (store.words.length === 0) {
          errorMessage.value = "❌ 題庫為空或無法讀取";
        } else {
          playAudio(); // 🔹 進入第一題時播放音檔
        }
      } catch (error) {
        errorMessage.value = "❌ 載入題庫時發生錯誤";
        console.error(error);
      } finally {
        loading.value = false;
      }
    });

    onUnmounted(() => {
      clearInterval(timer);
      // 結束前暫停播放
      if (!audioInstance.paused) {
        audioInstance.pause();
      }
    });

    async function checkAnswer() {
      if (!currentWord.value) return;

      if (
        userInput.value.trim().toLowerCase() ===
        currentWord.value.word.toLowerCase()
      ) {
        feedback.value = "✅ 正確!";
        store.correctAnswers++;

        if (isUsingAPI.value) {
          await submitCurrentResult();
        }

        await nextTick();
        goToNextQuestion(); // 答對後跳至下一題
      } else {
        feedback.value = "❌ 錯誤，請再試一次！";
      }
    }

    function goToNextQuestion() {
      if (store.currentIndex < store.words.length - 1) {
        store.currentIndex++;
        userInput.value = "";
      } else {
        endQuiz();
      }
    }

    async function submitCurrentResult() {
      const result = await submitQuizResults(category.value, filename.value, [currentWord.value.word]);
      if (!result.success) {
        console.error("❌ [API錯誤] 提交失敗:", result.error);
        errorMessage.value = "⚠️ 提交失敗，請稍後再試";
      }
    }

    function goBack() {
      router.go(-1);
    }

    function endQuiz() {
      clearInterval(timer);
      if (!audioInstance.paused) {
        audioInstance.pause();
      }
      elapsedTime.value = Math.floor((Date.now() - startTime.value) / 1000);
      router.push({
        path: `/result/${category.value}/${filename.value}`,
        query: {
          time: elapsedTime.value,
          score: store.correctAnswers,
        },
      });
    }

    const formattedTitle = computed(() => `測驗 - ${category.value} (${filename.value})`);

    return {
      userInput,
      feedback,
      currentWord,
      checkAnswer,
      goToNextQuestion,
      loading,
      errorMessage,
      elapsedTime,
      goBack,
      endQuiz,
      isUsingAPI,
      formattedTitle,
      totalQuestions,
      currentQuestionNumber,
      showDictionary,
      toggleDictionary,
    };
  },
};
</script>

<style scoped>
.quiz-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 20px;
}
.quiz-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 20px;
}
.progress {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 15px;
  color: #333;
}
.toggle-dict-btn {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 999;
}
</style>
