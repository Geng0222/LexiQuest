// useProfile.js
// ✅ 取得用戶個人資料
export async function fetchUserProfile() {
  try {
    const response = await fetch(`/api/user`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP 錯誤狀態碼: ${response.status}`);
    }

    const data = await response.json();

    return {
      username: data.username || "未知用戶",
      progress: data.progress || {}, // ✅ 確保 progress 至少是空物件
    };
  } catch (error) {
    console.error("❌ 無法獲取用戶個人資料:", error.message);
    return { username: "錯誤", progress: {} }; // ✅ 返回預設值，避免 UI 崩潰
  }
}

// ✅ 更新用戶名稱
export async function updateUserProfile(newUsername) {
  try {
    const response = await fetch(`/api/user`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username: newUsername }),
    });

    if (!response.ok) {
      throw new Error(`HTTP 錯誤狀態碼: ${response.status}`);
    }

    const data = await response.json();
    console.log("✅ 用戶名稱更新成功:", data);
    return data;
  } catch (error) {
    console.error("❌ 更新用戶名稱失敗:", error.message);
    throw error;
  }
}
