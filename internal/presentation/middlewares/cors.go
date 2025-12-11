package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cors returns a Gin Gonic middleware (gin.HandlerFunc) that sets the necessary
// Access-Control headers for Cross-Origin Resource Sharing.
func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Set the CORS headers
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		ctx.Writer.Header().Set("Access-Control-Max-Age", "3600")

		// If it's an OPTIONS request, respond with OK status and abort the chain
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusOK)
			return
		}

		ctx.Next()
	}
}