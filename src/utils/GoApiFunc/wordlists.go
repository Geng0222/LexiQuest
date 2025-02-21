package GoApiFunc

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 定義題庫的根目錄路徑
const wordlistPath = "./data/wordlists"

// ListAllWordlists 讀取並回傳題庫目錄下的所有分類與檔案
func ListAllWordlists(c *gin.Context) {
	categories := make(map[string][]string)

	// 檢查題庫目錄是否存在
	if _, err := os.Stat(wordlistPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "題庫目錄不存在"})
		return
	}

	// 讀取題庫根目錄中的所有項目
	files, err := os.ReadDir(wordlistPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法讀取題庫目錄"})
		return
	}

	// 逐個處理目錄，並收集裡面的 .txt 檔案名稱（不含副檔名）
	for _, file := range files {
		if file.IsDir() {
			category := file.Name()
			subFiles, err := os.ReadDir(filepath.Join(wordlistPath, category))
			if err != nil {
				continue // 忽略讀取錯誤的目錄
			}

			var filenames []string
			for _, subFile := range subFiles {
				if strings.HasSuffix(subFile.Name(), ".txt") {
					filename := strings.TrimSuffix(subFile.Name(), ".txt")
					filenames = append(filenames, filename)
				}
			}
			categories[category] = filenames
		}
	}

	c.JSON(http.StatusOK, categories)
}

// LoadWordlistHandler 讀取並解析指定的題庫文件
func LoadWordlistHandler(c *gin.Context) {
	category := c.Param("category")
	filename := c.Param("filename")

	// 檢查檔案是否存在
	filePath := filepath.Join(wordlistPath, category, filename+".txt")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("題庫文件不存在: %s", filePath)})
		return
	}

	// 讀取並解析 .txt 檔案
	words, err := parseWordlistFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法解析題庫"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"category": category, "filename": filename, "data": words})
}

// parseWordlistFile 解析 .txt 文件並回傳單字資料
func parseWordlistFile(filePath string) ([]map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []map[string]string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 跳過空行或以 # 開頭的註解行
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) < 3 {
			fmt.Printf("⚠️ 格式錯誤: %s\n", line)
			continue // 跳過格式錯誤的行
		}

		words = append(words, map[string]string{
			"word":        strings.TrimSpace(parts[0]),
			"translation": strings.TrimSpace(parts[1]),
			"type":        strings.TrimSpace(parts[2]),
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

// GetRandomWordlist 採用加權隨機選擇出題，並透過 Query 參數 "last" 排除上一次出現的單字
func GetRandomWordlist(c *gin.Context) {
	category := c.Param("category")
	filename := c.Param("filename")
	lastWord := c.Query("last") // ✅ 仍然支援 Query 參數來排除上一次的單字

	// ✅ 修正這行，從 Path 參數讀取 `limit`
	limitStr := c.Param("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // 預設返回 10 題
	}

	// 讀取題庫內容
	filePath := filepath.Join(wordlistPath, category, filename+".txt")
	words, err := parseWordlistFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法解析題庫"})
		return
	}

	// **✅ 確保 `words` 不是空的**
	if len(words) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "題庫為空，無可用單字"})
		return
	}

	// 讀取用戶學習進度
	userProgress, err := GetUserProgress()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "無法讀取用戶學習進度"})
		return
	}

	// 確保進度數據存在，避免 nil 錯誤
	if _, exists := userProgress[category]; !exists {
		userProgress[category] = make(map[string]map[string]float64)
	}
	if _, exists := userProgress[category][filename]; !exists {
		userProgress[category][filename] = make(map[string]float64)
	}

	// 加權隨機選擇多個單字，排除 lastWord
	selectedWords := make([]map[string]string, 0, limit)
	usedWords := make(map[string]bool) // 紀錄已選過的單字，避免重複

	// **✅ 先確保第一個選擇不會重複 `lastWord`**
	firstWord := weightedRandomSelection(words, userProgress[category][filename], lastWord)
	if firstWord != nil {
		usedWords[firstWord["word"]] = true
		selectedWords = append(selectedWords, firstWord)
	}

	// **✅ 修正：確保選不超過 `words` 的總數量**
	for len(selectedWords) < limit && len(selectedWords) < len(words) {
		word := weightedRandomSelection(words, userProgress[category][filename], "")
		if word == nil {
			break // **避免無可用單字時進入無限迴圈**
		}
		// 避免重複
		if _, exists := usedWords[word["word"]]; exists {
			continue
		}
		usedWords[word["word"]] = true
		selectedWords = append(selectedWords, word)
	}

	// **✅ 確保 API 一定回傳至少 1 個單字**
	if len(selectedWords) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "無可用單字"})
		return
	}

	c.JSON(http.StatusOK, selectedWords)
}

// weightedRandomSelection 根據各單字的權重進行隨機選擇；若 exclude 不為空且候選不只有一項，則先排除該單字
func weightedRandomSelection(words []map[string]string, progress map[string]float64, exclude string) map[string]string {
	filteredWords := words
	if exclude != "" && len(words) > 1 {
		var temp []map[string]string
		for _, word := range words {
			if word["word"] != exclude {
				temp = append(temp, word)
			}
		}
		if len(temp) > 0 {
			filteredWords = temp
		}
	}

	var totalWeight float64
	weightedWords := []struct {
		word   map[string]string
		weight float64
	}{}

	// 產生加權單字列表
	for _, word := range filteredWords {
		wordText := word["word"]
		weight, exists := progress[wordText]
		if !exists {
			weight = 10.0 // 預設權重
		}
		if weight < 1.0 {
			weight = 1.0 // 最低權重
		}
		totalWeight += weight
		weightedWords = append(weightedWords, struct {
			word   map[string]string
			weight float64
		}{word, weight})
	}

	if totalWeight == 0 {
		return nil
	}

	// 根據總權重產生隨機閾值，並依序累計選出單字
	rand.Seed(time.Now().UnixNano())
	threshold := rand.Float64() * totalWeight
	var cumulativeWeight float64
	for _, item := range weightedWords {
		cumulativeWeight += item.weight
		if cumulativeWeight >= threshold {
			return item.word
		}
	}

	return nil
}
