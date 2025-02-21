//QuizCard.vue
<template>
  <el-card class="quiz-card">
    <h1>📖 測驗模式</h1>

    <el-alert v-if="loading" title="題庫載入中..." type="info" show-icon />
    <el-alert v-else-if="errorMessage" :title="errorMessage" type="error" show-icon />

    <div v-else-if="currentWord" class="quiz-content">
      <el-divider content-position="left">單字</el-divider>
      <p class="word">
        📌 <strong>
          <!-- 利用 v-for 逐字渲染 -->
          <span
            v-for="(item, index) in wordLetters"
            :key="index"
            :class="item.className"
          >
            {{ item.char }}
          </span>
        </strong>
      </p>
      <p class="translation">🔹 翻譯：{{ currentWord.translation }}</p>
      <p class="type">🔸 詞性：{{ currentWord.type }}</p>

      <el-divider content-position="left">你的答案</el-divider>

      <el-input
        v-model="localUserInput"
        placeholder="請輸入答案"
        clearable
        @keyup.enter="$emit('check-answer')"
      />

      <div class="feedback">
        <el-tag v-if="feedback" :type="feedback.includes('✅') ? 'success' : 'danger'">
          {{ feedback }}
        </el-tag>
      </div>

      <p class="time">⏳ 時間：{{ elapsedTime }} 秒</p>

      <el-divider />

      <div class="button-group">
        <el-button type="info" icon="el-icon-back" @click="$emit('go-back')">返回</el-button>
        <el-button type="primary" icon="el-icon-check" @click="$emit('check-answer')">提交</el-button>
        <el-button type="danger" icon="el-icon-warning" @click="$emit('end-quiz')">結算</el-button>
      </div>
    </div>
  </el-card>
</template>

<script>
export default {
  props: {
    userInput: String,
    feedback: String,
    currentWord: Object,
    elapsedTime: Number,
    loading: Boolean,
    errorMessage: String
  },
  emits: ["check-answer", "go-back", "end-quiz", "update:userInput"],
  computed: {
    // v-model 的本地映射
    localUserInput: {
      get() {
        return this.userInput;
      },
      set(value) {
        this.$emit("update:userInput", value);
      },
    },
    // 依據 currentWord.word 與 userInput 產生逐字對應的陣列
    wordLetters() {
      if (!this.currentWord || !this.currentWord.word) return [];
      const letters = this.currentWord.word.split('');
      const input = this.userInput || '';
      return letters.map((char, index) => {
        let className = '';
        if (input[index] === undefined) {
          className = 'pending'; // 尚未輸入
        } else if (input[index] === char) {
          className = 'correct'; // 正確
        } else {
          className = 'wrong';   // 錯誤
        }
        return { char, className };
      });
    }
  },
};
</script>

<style scoped>
.quiz-card {
  width: 600px;
  max-width: 90%;
  padding: 30px;
  border-radius: 10px;
  box-shadow: 0px 6px 12px rgba(0, 0, 0, 0.1);
}

.word {
  font-size: 32px;
  font-weight: bold;
  color: #333;
}

/* 正確字母：綠色 */
.correct {
  color: #4CAF50;
  font-weight: bold;
  transition: color 0.3s ease-in-out;
}

/* 錯誤字母：紅色 */
.wrong {
  color: #FF5252;
  font-weight: bold;
  transition: color 0.3s ease-in-out;
}

/* 尚未輸入：灰色 */
.pending {
  color: #aaa;
}

.translation, .type {
  font-size: 16px;
  color: #666;
}

.feedback {
  margin-top: 10px;
}

.time {
  font-size: 14px;
  color: #888;
}

.button-group {
  display: flex;
  justify-content: space-between;
  margin-top: 20px;
}
</style>
