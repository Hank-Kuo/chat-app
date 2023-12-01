package http

import (
	"strings"

	"chat-app/pkg/customError"
	httpResponse "chat-app/pkg/response/http_response"
	"chat-app/pkg/utils"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) AuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		tokenStr := strings.Split(token, "Bearer ")
		if len(tokenStr) > 1 {
			claims, err := utils.ValidJwt(m.cfg, tokenStr[1])
			if err != nil {
				httpResponse.Fail(err, m.logger).ToJSON(c)
				return
			}
			c.Set("userID", claims.UserID)
			c.Next()
			return
		}
		httpResponse.Fail(customError.ErrNotGetJWTToken, m.logger).ToJSON(c)
	}
}
