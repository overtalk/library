package jwt_test

import (
	"testing"
	"time"

	. "web-layout/utils/gin/middleware/jwt"
	. "web-layout/utils/parse"
)

func TestJWTAuth(t *testing.T) {
	j := JwtSigner{
		Key: []byte("sdfasdfasdfasd"),
	}

	if err := j.Init(); err != nil {
		t.Error(err)
		return
	}

	token, err := j.Gen("1")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("--- before refresh ---")
	t.Log(token)
	// show details
	data, _, err := j.Verification(token)
	if err != nil {
		t.Error(err)
		return
	}
	// detail
	for k, v := range data {
		if k == "expire" || k == "sign_ts" {
			t.Log(k, time.Unix(Int(v), 0).String()[:19])
		}
	}

	time.Sleep(3 * time.Second)

	newToken, err := j.Refresh(data)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("--- after refresh ---")
	t.Log(newToken)
	// show details
	data, _, err = j.Verification(newToken)
	if err != nil {
		t.Error(err)
		return
	}
	// detail
	for k, v := range data {
		if k == "expire" || k == "sign_ts" {
			t.Log(k, time.Unix(Int(v), 0).String()[:19])
		}
	}
}

//func TestServer(t *testing.T) {
//	j := JwtSigner{
//		Key: []byte("sdfasdfasdfasd"),
//	}
//
//	if err := j.Init(); err != nil {
//		t.Error(err)
//		return
//	}
//
//	jwt := GinJwtMiddleware{
//		Signer: j,
//		Unauthorized: map[string]interface{}{
//			"errcode": 111,
//			"message": "token error",
//		},
//	}
//
//	r := gin.New()
//
//	r.POST("/login", func(c *gin.Context) {
//		token, err := j.Gen("1")
//		if err != nil {
//			c.JSON(200, gin.H{
//				"errocde": 112,
//				"message": "gen token error",
//			})
//			return
//		}
//
//		c.Writer.Header().Set(TokenKey, token)
//		c.JSON(200, gin.H{
//			"errocde": 112,
//			"message": "生成 token 成功",
//		})
//	}) //login不验JWT
//
//	auth := r.Group("/token", jwt.Jwt) //这个组全认证
//	auth.POST("/test", func(c *gin.Context) {
//		c.JSON(200, gin.H{
//			"errocde": 0,
//			"message": "token 解析成功",
//		})
//	})
//
//	r.Run(":9999")
//}
