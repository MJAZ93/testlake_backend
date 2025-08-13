package middleware

import (
	"net/http"

	"testlake/inout"
	"testlake/utils"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		err := utils.ValidateJWT(context)
		if err != nil {
			response := inout.BaseResponse{
				ErrorCode:        401,
				ErrorDescription: err.Error(),
			}
			context.JSON(http.StatusUnauthorized, response)
			context.Abort()
			return
		}
		context.Next()
	}
}