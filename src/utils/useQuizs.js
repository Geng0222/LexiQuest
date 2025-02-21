// utils/useResults.js
export async function submitQuizResults(category, filename, results) {
  try {
    console.log("📡 [API請求] 提交測驗結果...");
    console.log("📌 category:", category);
    console.log("📌 filename:", filename);
    console.log("📌 results (送出前):", JSON.stringify(results, null, 2));

    const response = await fetch("/api/quiz/submit", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        category,
        filename,
        results // ✅ 確保 `results` 只包含一個題目
      })
    });

    console.log("📡 [API回應] HTTP 狀態碼:", response.status);

    if (!response.ok) {
      throw new Error(`HTTP 錯誤 ${response.status}`);
    }

    const resultData = await response.json();
    console.log("✅ [API成功] 測驗提交成功:", resultData);
    return { success: true, data: resultData };
  } catch (error) {
    console.error("❌ [API錯誤] 測驗提交失敗:", error);
    return { success: false, error: "測驗結果提交失敗，請稍後再試" };
  }
}