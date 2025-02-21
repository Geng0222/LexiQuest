package GoApiFunc

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const userFilePath = "./data/userdata/user.json"

// UserData 定義了單一用戶的完整資料，包括 username 與學習進度
type UserData struct {
	Username string                                   `json:"username"`
	Progress map[string]map[string]map[string]float64 `json:"progress"`
}

// EnsureCategoryAndFilename 確保進度資料中存在指定的分類與題庫
func EnsureCategoryAndFilename(userProgress map[string]map[string]map[string]float64, category, filename string) {
	if _, exists := userProgress[category]; !exists {
		userProgress[category] = make(map[string]map[string]float64)
	}
	if _, exists := userProgress[category][filename]; !exists {
		userProgress[category][filename] = make(map[string]float64)
	}
}

// EnsureUserFile 確保 user.json 存在，若不存在則建立包含預設資料的檔案
func EnsureUserFile() {
	dir := filepath.Dir(userFilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	if _, err := os.Stat(userFilePath); os.IsNotExist(err) {
		defaultUser := UserData{
			Username: "",
			Progress: make(map[string]map[string]map[string]float64),
		}
		file, _ := os.Create(userFilePath)
		json.NewEncoder(file).Encode(defaultUser)
		file.Close()
	}
}

// GetUserData 讀取 user.json 並返回完整的 UserData 結構
func GetUserData() (*UserData, error) {
	EnsureUserFile()

	file, err := os.Open(userFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data UserData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

// SaveUserData 將完整的 UserData 寫入 user.json
func SaveUserData(data *UserData) error {
	EnsureUserFile()

	file, err := os.Create(userFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(data)
}

// GetUserProgress 返回 UserData 中的 Progress 部分
func GetUserProgress() (map[string]map[string]map[string]float64, error) {
	data, err := GetUserData()
	if err != nil {
		return nil, err
	}
	return data.Progress, nil
}

// SaveUserProgress 僅更新 UserData 中的 Progress 部分
func SaveUserProgress(progress map[string]map[string]map[string]float64) error {
	data, err := GetUserData()
	if err != nil {
		// 若讀取失敗，則初始化一個新的 UserData
		data = &UserData{
			Username: "",
			Progress: make(map[string]map[string]map[string]float64),
		}
	}
	data.Progress = progress
	return SaveUserData(data)
}
