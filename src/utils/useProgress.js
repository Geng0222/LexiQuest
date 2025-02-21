// useProgress.js
/**
 * 取得所有學習進度
 * @returns {Promise<Array>} 返回格式化後的學習進度
 */
export async function fetchProgressList() {
  try {
    const response = await fetch(`/api/user/progress`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP 錯誤狀態碼: ${response.status}`);
    }

    const data = await response.json();
    return transformProgressData(data);
  } catch (error) {
    console.error("❌ 無法獲取學習進度:", error.message);
    return []; // ✅ 返回空陣列，避免 UI 崩潰
  }
}

/**
 * 刪除指定的學習進度
 * @param {string} category - 類別名稱
 * @param {string} filename - 題庫名稱
 * @returns {Promise<void>}
 */
export async function deleteProgress(category, filename) {
  try {
    const response = await fetch(`/api/quiz/progress/${category}/${filename}`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP 錯誤狀態碼: ${response.status}`);
    }

    console.log(`✅ 已刪除學習進度: ${category}/${filename}`);
  } catch (error) {
    console.error("❌ 刪除學習進度失敗:", error.message);
    throw error;
  }
}

/**
 * 轉換 API 回傳的學習進度數據
 * @param {Object} data - 原始 API 回傳數據
 * @returns {Array} 格式化後的學習進度列表
 */
function transformProgressData(data) {
  const result = [];
  for (const category in data) {
    for (const filename in data[category]) {
      const words = Object.keys(data[category][filename]);
      result.push({
        category,
        filename,
        words: words.length, // 顯示有多少單字
      });
    }
  }
  return result;
}
