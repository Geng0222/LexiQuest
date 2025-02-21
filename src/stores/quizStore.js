//quizStore.js
import { defineStore } from "pinia";
import { loadWordlist } from "../utils/loadWordlist.js";
import { useApiStatus } from "../utils/useApiStatus";

export const useQuizStore = defineStore("quiz", {
  state: () => ({
    words: [], // 存放題庫
    currentIndex: 0, // 當前題目索引
    correctAnswers: 0, // 正確回答數
  }),

  actions: {
    // ✅ 加載完整題庫（靜態模式 & API）
    async loadQuiz(category, filename) {
      const response = await loadWordlist(category, filename);

      if (!response.success) {
        console.error("❌ 題庫載入失敗:", response.error);
        return;
      }

      this.words = response.data;
      this.currentIndex = 0;
      this.correctAnswers = 0;
    },

    // ✅ 加載隨機題庫（API 模式 & 靜態模式）
    async loadRandomQuiz(category, filename, limit = 10) {
      const { isUsingAPI, checkAPIStatus } = useApiStatus();
      await checkAPIStatus(); // ✅ 先確認 API 狀態

      let wordlist = [];

      // ✅ 如果 API 可用，請求 `/api/wordlist/random/:category/:filename/:limit`
      if (isUsingAPI.value) {
        try {
          console.log(`📡 從 API 獲取隨機題目: /api/wordlist/random/${category}/${filename}/${limit}`);
          const response = await fetch(`/api/wordlist/random/${category}/${filename}/${limit}`);

          if (!response.ok) {
            throw new Error(`❌ API 讀取失敗: HTTP ${response.status}`);
          }

          wordlist = await response.json();
          console.log(`✅ API 隨機題目載入成功，共 ${wordlist.length} 條單字`);
        } catch (error) {
          console.error("❌ API 隨機題目載入失敗，回退本地模式:", error.message);
        }
      }

      // ❌ 如果 API 無法使用，則回退到本地模式
      if (wordlist.length === 0) {
        console.log(`📂 讀取完整題庫，並從本地模式隨機抽取 ${limit} 題`);
        const response = await loadWordlist(category, filename);
        if (!response.success) {
          console.error("❌ 無法加載題庫（本地模式）:", response.error);
          return;
        }
        wordlist = response.data.sort(() => Math.random() - 0.5).slice(0, limit);
      }

      this.words = wordlist;
      this.currentIndex = 0;
      this.correctAnswers = 0;
    },

    // ✅ 切換到下一個單字
    nextQuestion() {
      if (this.currentIndex < this.words.length - 1) {
        this.currentIndex++;
      }
    }
  },
});

