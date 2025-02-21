//QuizPage.vue
<template>
  <div class="quiz-container">
    <div class="quiz-header">
      <h1>📖 {{ formattedTitle }}</h1>
      <el-tag :type="isUsingAPI ? 'success' : 'danger'" size="large">
        {{ isUsingAPI ? '🔗 API模式' : '🌐 靜態模式' }}
      </el-tag>
    </div>

    <!-- ✅ 顯示題目進度 -->
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
  </div>
</template>

<script>
import { useQuizStore } from "../stores/quizStore";
import { computed, ref, onMounted, onUnmounted, nextTick } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useApiStatus } from "../utils/useApiStatus";
import { submitQuizResults } from "../utils/useQuizs";
import QuizCard from "../components/QuizCard.vue";

export default {
  components: { QuizCard },
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

    const category = ref(route.params.category || "daily");
    const filename = ref(route.params.filename || "basic");

    onMounted(async () => {
      startTime.value = Date.now();
      timer = setInterval(() => {
        elapsedTime.value = Math.floor((Date.now() - startTime.value) / 1000);
      }, 1000);

      try {
        await checkAPIStatus();
        await store.loadRandomQuiz(category.value, filename.value, 10); // 預設 10 題

        if (store.words.length === 0) {
          errorMessage.value = "❌ 題庫為空或無法讀取";
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
    });

    const currentWord = computed(() => store.words[store.currentIndex] || null);
    const totalQuestions = computed(() => store.words.length);
    const currentQuestionNumber = computed(() => store.currentIndex + 1);

    async function checkAnswer() {
      if (!currentWord.value) return;

      if (userInput.value.trim().toLowerCase() === currentWord.value.word.toLowerCase()) {
        feedback.value = "✅ 正確!";
        store.correctAnswers++;

        if (isUsingAPI.value) {
          await submitCurrentResult();
        }

        await nextTick();
        goToNextQuestion();
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
      elapsedTime.value = Math.floor((Date.now() - startTime.value) / 1000);
      
      router.push({
        path: `/result/${category.value}/${filename.value}`,
        query: {
          time: elapsedTime.value,
          score: store.correctAnswers
        }
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
      currentQuestionNumber
    };
  }
};
</script>

<style scoped>
.quiz-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100vh;
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
</style>
