package middleware

import (
	"isp-engine/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

//处理跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		// 处理请求
		c.Next()
	}
}

func JSONAppErrorReporter() gin.HandlerFunc {
	return jsonAppErrorReporterT(gin.ErrorTypeAny)
}

func jsonAppErrorReporterT(errType gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		detectedErrors := c.Errors.ByType(errType)
		if len(detectedErrors) > 0 {
			err := detectedErrors[0].Err
			var parsedError *utils.ApiError
			switch err.(type) {
			//如果产生的error是自定义的结构体,转换error,返回自定义的code和msg
			case *utils.ApiError:
				parsedError = err.(*utils.ApiError)
			default:
				parsedError = &utils.ApiError{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				}
			}
			c.IndentedJSON(parsedError.Code, parsedError)
			return
		}

	}
}
