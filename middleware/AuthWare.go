package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		sessionID := session.Get("sessionID")

		if sessionID == nil {
			ctx.Redirect(http.StatusFound, "/auth/web")
			ctx.Abort()
			return
		}

		// Xác thực thành công
		ctx.Next()
	}
}
