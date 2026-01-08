package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger ログミドルウェア
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// リクエスト開始時刻
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// リクエストの詳細をログ出力
		log.Printf("[%s] %s %s - Request started", method, path, c.Request.URL.RawQuery)

		// リクエストボディの読み取り（可能な場合）
		if c.Request.Body != nil {
			// ボディは後続のハンドラーで使用されるため、ここではログのみ
			log.Printf("[%s] %s - Request body present", method, path)
		}

		// 次のハンドラーを実行
		c.Next()

		// レスポンスの詳細をログ出力
		latency := time.Since(start)
		status := c.Writer.Status()

		log.Printf("[%s] %s - Status: %d, Latency: %v", method, path, status, latency)

		// エラーが発生した場合
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Printf("[%s] %s - Error: %v", method, path, err.Error())
			}
		}
	}
}
