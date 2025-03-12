// utils/useApiStatus.js
import { ref, onMounted, watchEffect } from "vue";

export function useApiStatus() {
  const isUsingAPI = ref(false);

  async function checkAPIStatus() {
    try {
      const response = await fetch("/api/status"); // ✅ 使用 Vite Proxy 確保請求正確轉發
      if (response.ok) {
        isUsingAPI.value = true;
        console.log("✅ API 存活，切換到 API 模式");
        return true;
      }
    } catch (error) {
      console.warn("🔴 API 不可用，切換到靜態模式", error);
    }
    isUsingAPI.value = false;
    return false;
  }

  onMounted(() => {
    console.log("🔄 檢查 API 狀態...");
    checkAPIStatus();
  });

  watchEffect(() => {
    console.log("📡 API 狀態變更：", isUsingAPI.value ? "🔗 API模式" : "🌐 靜態模式");
  });

  return { isUsingAPI, checkAPIStatus };
}
