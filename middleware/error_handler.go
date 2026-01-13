package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler エラーハンドリングミドルウェア
func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// パニックをキャッチして統一されたエラーレスポンスを返す
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		c.Abort()
	})
}
