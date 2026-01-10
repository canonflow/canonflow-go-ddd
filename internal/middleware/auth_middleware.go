package middleware

import (
	"net/http"
	"time"

	"github.com/canonflow/canonflow-go-ddd/internal/domain/user/model"
	"github.com/canonflow/canonflow-go-ddd/pkg/jwt"
	"github.com/canonflow/canonflow-go-ddd/pkg/response"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := jwt.GetToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, response.BaseErrorResponse{
				Code:   http.StatusUnauthorized,
				Status: "Unauthorized",
				Error:  "Missing Authorization Token",
			})
			ctx.Abort()
			return
		}

		// TODO: Parse Token
		claims, err := jwt.ParseToken(tokenString, secret)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.BaseErrorResponse{
				Code:   http.StatusInternalServerError,
				Status: "Internal Server Error",
				Error:  err.Error(),
			})
			ctx.Abort()
			return
		}

		// TODO: Check the expiry time
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			jwt.DeleteToken(ctx)
			ctx.JSON(http.StatusUnauthorized, response.BaseErrorResponse{
				Code:   http.StatusUnauthorized,
				Status: "Unauthorized",
				Error:  "Token Expired",
			})
			ctx.Abort()
			return
		}

		// TODO: Set the current User
		ctx.Set(jwt.USER_KEY, model.User{
			ID:       int64(claims["sub"].(float64)),
			Username: claims["username"].(string),
		})

		// TODO: Attach back the token
		ctx.Set(jwt.TOKEN_KEY, tokenString)

		// TODO: Continue
		ctx.Next()
	}
}
