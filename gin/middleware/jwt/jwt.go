package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	IDKey     = "id"      // 用户唯一标识
	ExpireKey = "expire"  // 过期时间
	SignTSKey = "sign_ts" // token签发时间
)

// JwtSigner 签名结构
type JwtSigner struct {
	// signing algorithm - possible values are HS256, HS384, HS512
	// Optional, default is HS256.
	SignAlgorithm string
	Key           []byte
	Timeout       time.Duration
	MaxRefresh    time.Duration // 最大刷新有效期
	PayloadFunc   func(data interface{}) jwt.MapClaims
}

// JwtSigner 初始化
func (j *JwtSigner) Init() error {
	switch j.SignAlgorithm {
	case "HS256", "HS384", "HS512":
	default:
		j.SignAlgorithm = "HS256"
	}

	if j.Timeout == 0 {
		j.Timeout = time.Hour
	}

	if len(j.Key) == 0 {
		return ErrMissingSecretKey
	}

	if j.PayloadFunc == nil {
		j.PayloadFunc = func(data interface{}) jwt.MapClaims {
			return jwt.MapClaims{
				IDKey: data,
			}
		}
	}

	return nil
}

// 生成Token
func (j *JwtSigner) Gen(data interface{}) (string, error) {
	token := jwt.New(jwt.GetSigningMethod(j.SignAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	if j.PayloadFunc != nil {
		for key, value := range j.PayloadFunc(data) {
			claims[key] = value
		}
	}

	now := time.Now()
	expire := now.Add(j.Timeout)
	claims[ExpireKey] = expire.Unix()
	claims[SignTSKey] = now.Unix()
	return token.SignedString(j.Key)
}

// 解析认证Token
// 未过期的Token & 过期可以Refresh的Token 返回claims
// 过期Token返回error
func (j *JwtSigner) Verification(token string) (jwt.MapClaims, bool, error) {
	claims, isUpdate, err := j.checkIfTokenExpire(token)
	if err != nil {
		return nil, false, err
	}

	m := make(jwt.MapClaims)
	for k, v := range claims {
		m[k] = v
	}
	return m, isUpdate, nil
}

// 更新token，重新设置一些过期时间
func (j *JwtSigner) Refresh(claims jwt.MapClaims) (string, error) {
	// Create the token
	newToken := jwt.New(jwt.GetSigningMethod(j.SignAlgorithm))
	newClaims := newToken.Claims.(jwt.MapClaims)

	for key := range claims {
		newClaims[key] = claims[key]
	}

	now := time.Now()
	expire := now.Add(j.Timeout)
	newClaims[ExpireKey] = expire.Unix()
	newClaims[SignTSKey] = now.Unix()

	return newToken.SignedString(j.Key)
}

// token存在三种状态：未过期、过期、过期但可以Refresh
// err ！= nil ：过去
// err ！= nil && true ：过期可以refresh
func (j *JwtSigner) checkIfTokenExpire(t string) (jwt.MapClaims, bool, error) {
	token, err := j.parseToken(t)
	if err != nil {
		return nil, false, err
	}
	claims := token.Claims.(jwt.MapClaims)

	expTS := int64(claims[ExpireKey].(float64))
	if expTS > time.Now().Unix() {
		return claims, false, nil
	}

	signTS := int64(claims[SignTSKey].(float64))
	if signTS > time.Now().Add(-(j.MaxRefresh + j.Timeout)).Unix() {
		return claims, true, nil
	}

	return nil, false, ErrExpiredToken
}

// parseToken : 将 string 解析成 jwt.Token
func (j *JwtSigner) parseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(j.SignAlgorithm) != t.Method {
			return nil, ErrInvalidSigningAlgorithm
		}
		return j.Key, nil
	})
}
