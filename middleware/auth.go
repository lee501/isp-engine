package middleware

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var (
	jwtSecret string
	identityKey = "id"
	globalRoutes = make([]*Router, 0)
)

func init() {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 20)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	jwtSecret = string(b)
}

type Router struct {
	HttpMethod   string
	RelativePath string
	Handlers     []gin.HandlerFunc
}

func registerGlobal(method, path string, handlers ...gin.HandlerFunc) {
	for _, r := range globalRoutes {
		if r.HttpMethod == method && r.RelativePath == path {
			return
		}
	}

	r := &Router{
		HttpMethod:   method,
		RelativePath: path,
		Handlers:     handlers,
	}
	globalRoutes = append(globalRoutes, r)
}

func Auth() gin.HandlerFunc {

	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:           "jin jwt",
		Key:             []byte(jwtSecret),
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour * 24,
		Authenticator:   authenticator,
		Authorizator:    authorizator,
		IdentityHandler: identityHandler,
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		PayloadFunc: loginPayload,
		SigningAlgorithm: "HS256",
		TokenLookup: "header:Authorization",
		TimeFunc: time.Now,
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",
		TokenHeadName: "Bearer",
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, gin.H{
				"code":   http.StatusOK,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		},
		RefreshResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, gin.H{
				"code":   http.StatusOK,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		},
	}
	authMiddleware.MiddlewareInit()
	registerGlobal("POST", "/api/login", authMiddleware.LoginHandler)
	registerGlobal("POST", "/api/logout", func(*gin.Context) {})
	return authMiddleware.MiddlewareFunc()
}

type AuthBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func authenticator(c *gin.Context) (interface{}, error) {
	var body AuthBody
	if err := c.ShouldBind(&body); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	userID := body.Username
	password := body.Password

	if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
		return &AuthBody{
			Username:  userID,
			Password:  password,
		}, nil
	}

	return nil, jwt.ErrFailedAuthentication
}

func loginPayload(data interface{}) jwt.MapClaims {
	if v, ok := data.(*AuthBody); ok {
		return jwt.MapClaims{
			identityKey: v.Username,
		}
	}
	return jwt.MapClaims{}
}

func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &AuthBody{
		Username: claims[identityKey].(string),
	}
}

func authorizator(data interface{}, c *gin.Context) bool {
	if data == nil {
		return false
	}
	return true
}

func GetglobalRoutes() []*Router {
	return globalRoutes
}