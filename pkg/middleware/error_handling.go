package middleware

import (
	"ISO_Auditing_Tool/pkg/custom_errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				if customErr, ok := err.Err.(*custom_errors.CustomError); ok {
					respondWithError(c, customErr)
					return
				}
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
	}
}

func respondWithError(c *gin.Context, customErr *custom_errors.CustomError) {
	response := gin.H{"error": customErr.Error()}
	// if customErr.Context != nil {
	// 	response["Context"] = customErr.Context
	// }
	c.JSON(customErr.StatusCode, response)
}
