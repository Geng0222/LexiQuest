package GoApiFunc

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 參數設定：
// decayFactor：每次練習後，單字權重乘以此因子（例如 0.9 表示權重下降至原來的 90%）
// defaultWeight：未曾練習時的預設權重
// minWeight：單字權重的最低值，避免降到 0
const (
	decayFactor   = 0.9
	defaultWeight = 10.0
	minWeight     = 1.0
)

// SubmitQuiz 提交測驗結果（僅記錄練習次數，不區分正確與錯誤）
// 請求的 Results 為練習過的單字列表
func SubmitQuiz(c *gin.Context) {
	var request struct {
		Category string   `json:"category"`
		Filename string   `json:"filename"`
		Results  []string `json:"results"` // 例如 ["apple", "banana"]
	}

	// 🚀 Debug: 確保請求格式正確
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("❌ JSON 解析失敗:", err) // ✅ 檢查 JSON 解析錯誤
		c.JSON(http.StatusBadRequest, gin.H{"error": "請提供有效的測驗結果"})
		return
	}

	// 讀取用戶進度（需由 user_progress.go 提供 GetUserProgress）
	userProgress, err := GetUserProgress()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法讀取用戶進度"})
		return
	}

	// 確保指定題庫進度資料存在
	EnsureCategoryAndFilename(userProgress, request.Category, request.Filename)

	// 對每個練習過的單字進行乘法衰減更新
	for _, word := range request.Results {
		oldWeight := userProgress[request.Category][request.Filename][word]
		if oldWeight == 0 {
			oldWeight = defaultWeight
		}
		newWeight := oldWeight * decayFactor
		if newWeight < minWeight {
			newWeight = minWeight
		}
		userProgress[request.Category][request.Filename][word] = newWeight
	}

	// 儲存更新後的進度
	if err := SaveUserProgress(userProgress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法儲存學習進度"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "測驗結果已提交",
		"progress": userProgress[request.Category][request.Filename],
	})
}

// DeleteWordlistProgress 刪除指定題庫的所有學習進度
func DeleteWordlistProgress(c *gin.Context) {
	category := c.Param("category")
	filename := c.Param("filename")

	userProgress, err := GetUserProgress()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法讀取用戶進度"})
		return
	}

	// 確保類別和題庫存在
	if _, exists := userProgress[category]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "該類別不存在"})
		return
	}
	if _, exists := userProgress[category][filename]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "該題庫不存在"})
		return
	}

	// 刪除該題庫的學習進度
	delete(userProgress[category], filename)

	// 如果該類別下已無任何題庫，則刪除整個類別
	if len(userProgress[category]) == 0 {
		delete(userProgress, category)
	}

	// 儲存變更
	if err := SaveUserProgress(userProgress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法刪除學習進度"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("題庫 %s/%s 的學習進度已刪除", category, filename),
	})
}
