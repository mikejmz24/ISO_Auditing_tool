package middleware

import (
	"ISO_Auditing_Tool/pkg/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				if customErr, ok := err.Err.(*errors.CustomError); ok {
					if customErr.Context != nil {
						c.JSON(customErr.StatusCode, gin.H{"error": customErr.Message})
						return
					}
					c.JSON(customErr.StatusCode, gin.H{"error": customErr.Message})
					return
				}
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
	}
}
