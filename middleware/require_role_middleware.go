package middleware

import (
	"net/http"

	"github.com/ardiannm/go/models"
	"github.com/ardiannm/go/utils"
	"github.com/gin-gonic/gin"
)

func RequireRole(role models.UserRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole, err := utils.GetUserRoleFromContext(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User role not found"})
			ctx.Abort()
			return
		}
		if userRole != role {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
