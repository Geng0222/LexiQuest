package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"lexiquest/src/utils/GoApiFunc"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 取得區域網 IP
func getLocalIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("❌ 無法取得網路介面:", err)
		return "127.0.0.1"
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				ip := ipNet.IP.To4()
				if ip != nil {
					// 限制只取得區域網 IP
					if (ip[0] == 192 && ip[1] == 168) || // 192.168.x.x
						(ip[0] == 10) || // 10.x.x.x
						(ip[0] == 172 && ip[1] >= 16 && ip[1] <= 31) { // 172.16.x.x ~ 172.31.x.x
						return ip.String()
					}
				}
			}
		}
	}

	return "127.0.0.1"
}

func main() {
	r := gin.Default()

	// ✅ 允許區域網內所有設備存取
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		AllowCredentials: true,
	}))

	// ✅ API 路由
	api := r.Group("/api")
	{
		api.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "API is running"})
		})
		api.GET("/wordlists/all", GoApiFunc.ListAllWordlists)
		api.GET("/wordlist/:category/:filename", GoApiFunc.LoadWordlistHandler)
		api.GET("/wordlist/random/:category/:filename/:limit", GoApiFunc.GetRandomWordlist)

		api.POST("/quiz/submit", GoApiFunc.SubmitQuiz)
		api.DELETE("/quiz/progress/:category/:filename", GoApiFunc.DeleteWordlistProgress)

		api.GET("/user", GoApiFunc.GetUser)
		api.POST("/user", GoApiFunc.CreateUser)
		api.PUT("/user", GoApiFunc.UpdateUser)
		api.GET("/user/progress", GoApiFunc.GetUserProgressHandler)
	}

	// ✅ 用戶資料與題庫目錄
	userDataPath := "./data"
	if _, err := os.Stat(userDataPath); os.IsNotExist(err) {
		fmt.Println("📂 `data/` 資料夾不存在，正在建立...")
		err := os.Mkdir(userDataPath, os.ModePerm)
		if err != nil {
			fmt.Println("❌ 無法建立 `data/`:", err)
			return
		}
		fmt.Println("✅ `data/` 資料夾已建立")
	}

	// ✅ API 讓前端存取 `data/`
	api.GET("/userdata/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		filePath := filepath.Join(userDataPath, filename)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
			return
		}

		c.File(filePath)
	})

	api.POST("/userdata/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		filePath := filepath.Join(userDataPath, filename)

		file, err := os.Create(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法建立文件"})
			return
		}
		defer file.Close()

		_, err = file.WriteString("測試內容")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法寫入文件"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "文件已建立"})
	})

	// ✅ **提供 Vue 靜態文件 (`/lexiquest/`)**
	r.Static("/lexiquest", "./dist") // **確保對應 Vite 設定**

	r.NoRoute(func(c *gin.Context) {
		c.File("./dist/index.html")
	})

	// ✅ 取得區域網 IP
	localIP := getLocalIP()
	port := "28080"
	addr := "0.0.0.0:" + port

	server := &http.Server{Addr: addr, Handler: r}

	// 🚀 啟動伺服器
	go func() {
		fmt.Println("🚀 伺服器啟動成功！")
		fmt.Println("🌍 區域網存取: http://" + localIP + ":" + port + "/lexiquest/")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("❌ 伺服器錯誤:", err)
		}
	}()

	// 🛑 等待中斷訊號
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	fmt.Println("🛑 伺服器關閉")
	server.Shutdown(context.Background())
}
