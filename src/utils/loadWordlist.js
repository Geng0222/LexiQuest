// loadWordlist.js
import { useApiStatus } from "./useApiStatus";

export async function loadWordlist(category, filename) {
    const { isUsingAPI, checkAPIStatus } = useApiStatus();
    await checkAPIStatus(); // ✅ 確保 API 狀態可用

    // ✅ 先嘗試使用 API 模式
    if (isUsingAPI.value) {
        try {
            console.log(`📡 嘗試從 API 讀取題庫: /api/wordlist/${category}/${filename}`);
            const response = await fetch(`/api/wordlist/${category}/${filename}`);

            if (!response.ok) {
                throw new Error(`❌ API 讀取失敗: HTTP ${response.status}`);
            }

            const data = await response.json();
            console.log(`✅ API 題庫載入成功，共 ${data.data.length} 條單字`);
            return { success: true, data: data.data };
        } catch (error) {
            console.error("❌ API 題庫載入失敗，回退本地模式:", error.message);
        }
    }

    // ❌ 如果 API 無法使用，則回退本地 TXT 模式
    return await loadLocalWordlist(category, filename);
}

// 📌 本地 TXT 模式
async function loadLocalWordlist(category, filename) {
    try {
        const baseURL = import.meta.env.BASE_URL || "/";
        const filePath = `${baseURL}wordlists/${category}/${filename}.txt`;

        console.log(`📂 嘗試讀取本地題庫: ${filePath}`);

        const response = await fetch(filePath);
        if (!response.ok) {
            throw new Error(`❌ 無法讀取題庫: ${filePath} (HTTP ${response.status})`);
        }

        // ✅ 強制解析為 UTF-8，確保沒有亂碼
        const text = await response.text();
        const utf8Decoder = new TextDecoder("utf-8");
        const utf8Text = utf8Decoder.decode(new TextEncoder().encode(text));

        // ✅ 解析 TXT 轉換為 JSON
        const jsonData = parseWordlist(utf8Text);

        console.log(`✅ 成功讀取本地題庫，共 ${jsonData.length} 條單字`);
        return { success: true, data: jsonData };
    } catch (error) {
        console.error("❌ 讀取本地題庫失敗:", error.message);
        return { success: false, error: error.message };
    }
}

// 📌 解析 TXT 題庫格式（格式：單字,翻譯,詞性）
function parseWordlist(text) {
    return text
        .trim()
        .split("\n")
        .map(line => {
            const parts = line.split(",");
            if (parts.length < 3) {
                console.warn(`⚠️ 格式錯誤: ${line}`);
                return null;
            }
            return { word: parts[0].trim(), translation: parts[1].trim(), type: parts[2].trim() };
        })
        .filter(entry => entry !== null); // 過濾掉格式錯誤的行
}
