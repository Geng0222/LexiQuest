import { ref, computed } from "vue";
import { useApiStatus } from "./useApiStatus"; // ✅ 依賴 API 狀態檢查

export function useWordlists() {
  const { isUsingAPI, checkAPIStatus } = useApiStatus(); // ✅ 確保 API 狀態可用
  const loading = ref(true);
  const wordlists = ref({});
  
  // 📌 本地靜態題庫（當 API 無法連線時使用）
  const staticWordlists = {
    "daily": ["basic", "advanced"],
    "tech": ["AI", "computer"]
  };

  // 📌 類別標籤
  const categoryLabels = {
    "daily": "📖 日常單字",
    "tech": "🖥️ 科技單字"
  };

  // ✅ 封裝 API 請求函數
  async function fetchWordlistsFromAPI() {
    try {
      console.log("📡 使用 API 讀取題庫...");
      const response = await fetch("/api/wordlists/all");

      if (!response.ok) {
        throw new Error(`HTTP 錯誤: ${response.status}`);
      }

      const data = await response.json();
      console.log("✅ API 題庫加載成功", data);
      return data;
    } catch (error) {
      console.error("❌ API 題庫載入失敗，切換到靜態模式", error);
      return null; // API 失敗時返回 `null`
    }
  }

  // ✅ 加載題庫（支援 API 與離線模式）
  async function loadWordlists() {
    console.log("🔍 檢查 API 狀態...");
    await checkAPIStatus(); // ✅ 先確認 API 狀態

    if (isUsingAPI.value) {  
      const apiData = await fetchWordlistsFromAPI();
      wordlists.value = apiData || staticWordlists; // ✅ API 失敗時回退到離線題庫
    } else {
      console.warn("⚠️ 使用靜態題庫");
      wordlists.value = staticWordlists;
    }
    
    loading.value = false;
  }

  // ✅ 計算屬性: 轉換題庫為標籤格式
  const wordlistsWithLabels = computed(() => {
    return Object.entries(wordlists.value).map(([category, books]) => ({
      category,  // ✅ 儲存原始類別名稱 (如 `daily`)
      label: categoryLabels[category] || category, // ✅ 使用標籤名稱，若無則用原類別名稱
      books  // ✅ 這個類別底下的所有題庫
    }));
  });

  return { wordlistsWithLabels, categoryLabels, loading, loadWordlists, isUsingAPI };
}
