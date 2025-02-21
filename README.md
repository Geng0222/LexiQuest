# 📌 LexiQuest

🔤 **LexiQuest** 是一款單字打字測驗應用程式，支援離線運行，並允許用戶使用自訂題庫，以提升單字記憶能力。

## 📍 專案概述

### 目標
- **開發一款離線運行的單字測驗應用程式**
- **提供 GitHub Pages 版 Demo 體驗**
- **支援本地題庫，允許用戶自訂單字內容**
- **提升用戶的單字記憶與打字能力**

### 版本
- **線上版（GitHub Pages）**
  - 提供 Demo 題庫，讓用戶體驗測驗流程
  - 不儲存學習紀錄，測驗結束後數據即消失
- **本地版**
  - 支援讀取用戶自訂題庫 (`data/wordlists/`)
  - 可儲存學習進度，並根據進度隨機出題

## 📍 核心功能

### 打字測驗
- 顯示單字資訊（原文 + 翻譯 + 詞性）
- 使用者輸入單字，並檢查答案：
  - ✅ 正確 → 進入下一個單字
  - ❌ 錯誤 → 顯示錯誤訊息，要求重新輸入

### 計時功能
- 測驗開始後，自動計算完成測驗所需時間（不保存）

### 學習紀錄
- **線上版**：不儲存學習紀錄，數據不會被記錄
- **本地版（未來規劃）**：存儲於 JSON 檔案 (`data/userdata/`)

### 題庫管理
- **線上版**：讀取內建 Demo 題庫 (`public/wordlists/`)
- **本地版（未來規劃）**：允許用戶自訂 `.txt` 題庫 (`data/wordlists/`)

## 📍 技術架構

| 項目 | 技術選擇 |
|------|----------|
| 框架 | Vue 3 + Vite + Go (Gin) |
| 線上部署 | GitHub Pages |
| 狀態管理 | Pinia |
| UI 框架 | Element Plus |
| 路由管理 | Vue Router |
| 資料存儲 | 本地版：JSON 檔案，線上版：不存儲 |
| API 模擬 | JSON Server（開發模式） |
| 題庫來源 | `public/wordlists/`、`data/wordlists/` |

## 📍 開發安裝與執行

```bash
# 安裝 Node.js 依賴
npm install

# 啟動vite+vue開發環境
npm run dev

# 啟動GO開發環境
go run main.go

```

## 📍 本地使用

### 🚀 快速開始

1. 下載並解壓縮 ZIP 檔案
2. 執行 `LexiQuest.exe`

### 📚 個人題庫設定

將您的題庫檔案放置於 `data/wordlists/` 目錄下

#### 題庫格式範例：
```txt
apple,蘋果,noun
book,書本, noun
run,跑步,verb