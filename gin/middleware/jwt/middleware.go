package jwt

import (
	"github.com/gin-gonic/gin"
)

const TokenKey = "Authorization"

type GinJwtMiddleware struct {
	Signer       JwtSigner
	Unauthorized interface{}
}

func (g *GinJwtMiddleware) Jwt(c *gin.Context) {
	token := c.Request.Header.Get(TokenKey)
	claims, isUpdate, err := g.Signer.Verification(token)
	if err != nil {
		c.JSON(200, g.Unauthorized)
		c.Abort()
		return
	}

	// 如果需要更新token
	if isUpdate {
		newToken, err := g.Signer.Refresh(claims)
		if err != nil {
			c.JSON(200, g.Unauthorized)
			c.Abort()
			return
		}
		c.Writer.Header().Set(TokenKey, newToken)
	}

	for k, v := range claims {
		c.Set(k, v)
	}
}
