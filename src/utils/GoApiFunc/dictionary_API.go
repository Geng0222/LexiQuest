// dictionary_API.go
package GoApiFunc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// 定義 API 資料結構
type Phonetic struct {
	Text  string `json:"text"`
	Audio string `json:"audio,omitempty"`
}

type Definition struct {
	Definition string `json:"definition"`
	Example    string `json:"example,omitempty"`
}

type Meaning struct {
	PartOfSpeech string       `json:"partOfSpeech"`
	Definitions  []Definition `json:"definitions"`
}

type WordData struct {
	Word      string     `json:"word"`
	Phonetic  string     `json:"phonetic"`
	Phonetics []Phonetic `json:"phonetics"`
	Meanings  []Meaning  `json:"meanings"`
}

// 過濾後的回傳格式
type FilteredWordData struct {
	Word     string            `json:"word"`
	Phonetic string            `json:"phonetic"`
	Audio    string            `json:"audio,omitempty"`
	Meanings []FilteredMeaning `json:"meanings"`
}

type FilteredMeaning struct {
	PartOfSpeech string `json:"partOfSpeech"`
	Definition   string `json:"definition"`
	Example      string `json:"example,omitempty"`
}

// ✅ 使用 sync.Map 優化快取
var wordCache sync.Map

// ✅ `FetchWordData` 直接檢查網路錯誤，移除 `isInternetAvailable`
func FetchWordData(word string) ([]WordData, error) {
	client := http.Client{Timeout: 5 * time.Second}

	url := fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", word)
	resp, err := client.Get(url)
	if err != nil {
		if _, ok := err.(net.Error); ok {
			return nil, fmt.Errorf("network unreachable") // 直接判斷網路錯誤
		}
		return nil, fmt.Errorf("failed to fetch word data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch word data: status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var wordData []WordData
	if err = json.Unmarshal(body, &wordData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	return wordData, nil
}

// ✅ 優化 `FilterWordData`，確保 `example` 只在有內容時加入
func FilterWordData(data []WordData) FilteredWordData {
	if len(data) == 0 {
		return FilteredWordData{}
	}

	word := data[0]

	// ✅ 取得主要音標（如果 `word.Phonetic` 為空，則從 `Phonetics` 內抓取）
	phonetic := word.Phonetic
	if phonetic == "" && len(word.Phonetics) > 0 {
		for _, p := range word.Phonetics {
			if p.Text != "" {
				phonetic = p.Text
				break
			}
		}
	}

	// ✅ 取得音訊連結（優先選擇有值的）
	audio := ""
	for _, p := range word.Phonetics {
		if p.Audio != "" {
			audio = p.Audio
			break
		}
	}

	// ✅ 處理 `meanings`
	var filteredMeanings []FilteredMeaning
	for _, m := range word.Meanings {
		if len(m.Definitions) > 0 {
			def := m.Definitions[0] // 取第一個定義
			filteredMeaning := FilteredMeaning{
				PartOfSpeech: m.PartOfSpeech,
				Definition:   def.Definition,
			}

			// **只有當 `example` 存在且非空時才加入**
			if def.Example != "" {
				filteredMeaning.Example = def.Example
			}

			filteredMeanings = append(filteredMeanings, filteredMeaning)
		}
	}

	return FilteredWordData{
		Word:     word.Word,
		Phonetic: phonetic,
		Audio:    audio,
		Meanings: filteredMeanings,
	}
}

// ✅ 精簡 `GetWordHandler`，移除不必要的鎖定邏輯
func GetWordHandler(c *gin.Context) {
	word := c.Param("word")

	// ✅ 先從快取讀取
	if cachedData, found := wordCache.Load(word); found {
		fmt.Println("🟢 使用快取資料:", word)
		c.JSON(http.StatusOK, cachedData.(FilteredWordData))
		return
	}

	// ❌ 如果快取沒有，則請求 API
	wordData, err := FetchWordData(word)
	if err != nil {
		status := http.StatusInternalServerError
		msg := "Failed to fetch word data"
		if err.Error() == "network unreachable" {
			status = http.StatusServiceUnavailable
			msg = "Network Unreachable, please check your internet connection."
		}
		c.JSON(status, gin.H{"error": msg})
		return
	}

	filteredData := FilterWordData(wordData)

	// ✅ 更新快取
	wordCache.Store(word, filteredData)

	fmt.Println("🔵 請求 API 取得新資料:", word)
	c.JSON(http.StatusOK, filteredData)
}
