package GoApiFunc

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUser 取得完整用戶資料（包含 username 與進度）
// 直接使用 GetUserData() 讀取存檔資料
func GetUser(c *gin.Context) {
	userData, err := GetUserData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法讀取用戶資料"})
		return
	}
	c.JSON(http.StatusOK, userData)
}

// CreateUserRequest 定義建立用戶的請求資料，僅接收 username
type CreateUserRequest struct {
	Username string `json:"username"`
}

// CreateUser 建立新用戶，僅使用前端傳入的 username，進度初始化為空
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "請提供有效的 username"})
		return
	}

	newUser := UserData{
		Username: req.Username,
		Progress: make(map[string]map[string]map[string]float64),
	}

	if err := SaveUserData(&newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法建立用戶"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用戶建立成功", "user": newUser})
}

// UpdateUserRequest 定義更新用戶的請求資料，僅更新 username
type UpdateUserRequest struct {
	Username string `json:"username"`
}

// UpdateUser 更新用戶的 username，保留原有進度資料
func UpdateUser(c *gin.Context) {
	userData, err := GetUserData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法讀取用戶資料"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "請提供有效的 username"})
		return
	}

	userData.Username = req.Username

	if err := SaveUserData(userData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法更新用戶資料"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用戶資料更新成功", "user": userData})
}

// GetUserProgressHandler 取得僅用戶的學習進度
// 保持原有 GetUserProgress() 接口以向後兼容
func GetUserProgressHandler(c *gin.Context) {
	progress, err := GetUserProgress()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法讀取用戶進度"})
		return
	}
	c.JSON(http.StatusOK, progress)
}
