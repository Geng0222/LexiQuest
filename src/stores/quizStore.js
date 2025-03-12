//quizStore.js
import { defineStore } from "pinia";
import { loadWordlist } from "../utils/loadWordlist.js";
import { useApiStatus } from "../utils/useApiStatus";

// 新增一個函式，根據單字呼叫字典 API
async function fetchDictionaryEntry(word) {
  try {
    const response = await fetch(`http://localhost:28080/dictionary/${word}`);
    if (!response.ok) {
      throw new Error(`無法取得 ${word} 的字典資料，HTTP 狀態: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error("獲取字典資料失敗:", error);
    return null;
  }
}

// 整合題庫中的每個單字資料
async function enrichWordsWithDictionaryData(words) {
  const enrichedWords = await Promise.all(
    words.map(async (wordObj) => {
      const dictData = await fetchDictionaryEntry(wordObj.word);
      if (dictData) {
        // 將字典 API 回傳的資料整合進原本的單字物件中
        return {
          ...wordObj,
          phonetic: dictData.phonetic,
          audio: dictData.audio,
          meanings: dictData.meanings,
        };
      }
      // 若無法取得資料則返回原單字物件
      return wordObj;
    })
  );
  return enrichedWords;
}

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
    
      // 使用上面的函式整合字典 API 資料
      this.words = await enrichWordsWithDictionaryData(response.data);
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

      this.words = await enrichWordsWithDictionaryData(wordlist);
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

