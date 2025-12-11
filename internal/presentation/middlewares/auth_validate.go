package middlewares

import (
	"Api/internal"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthConfiguration defines the required roles and an optional user ID check for the endpoint.
type AuthConfiguration struct {
	Roles []string
	Id    string
}

// AuthValidate returns a Gin Gonic middleware (gin.HandlerFunc) that validates a JWT
// and performs authorization checks based on the provided AuthConfiguration.
func AuthValidate(config AuthConfiguration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get authorization header from request and check valid format
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Missing authorization header",
			})
			return
		}

		headerStrings := strings.Split(authHeader, " ")
		if len(headerStrings) != 2 || headerStrings[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid authorization header format",
			})
			return
		}

		// Parse and validate token
		tokenString := headerStrings[1]
		token, err := internal.ValidateJWT(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token: " + err.Error(),
			})
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims or token is not valid",
			})
			return
		}

		// Check for eventual validation passed as parameter
		// Role-based authorization
		if len(config.Roles) > 0 {
			roleAllowed := false
			// TODO: make it valid for Keycloak actual roles
			userRole, exists := claims["email"].(string)
			if exists {
				for _, role := range config.Roles {
					if role == userRole {
						roleAllowed = true
						break
					}
				}
			}

			if !roleAllowed {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "Insufficient permissions (Role check failed)",
				})
				return
			}
		}

		// ID-based authorization (e.g., resource ownership check)
		if config.Id != "" {
			sub, exists := claims["sub"].(string)
			if !exists || sub != config.Id {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "Cannot access other user's resource",
				})
				return
			}
		}

		ctx.Set("user_id", claims["sub"])
		ctx.Set("user_role", claims["email"])

		ctx.Next()
	}
}