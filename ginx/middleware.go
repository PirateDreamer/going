package ginx

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/url"
	"time"

	"github.com/PirateDreamer/going/gredis"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.Request.Header.Get("X-Request-ID")
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "reqId", requestId))
		c.Next()
	}
}

// Cache 接口缓存中间件
func Cache(seconds, db int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if viper.GetBool("cache.disable") {
			c.Next()
			return
		}

		key := "ApiCache_" + c.Request.URL.Path
		parsedURL, err := url.Parse(c.Request.URL.String())
		if err != nil {
			ResFail(context.Background(), c, err)
			return
		}
		// Platform
		platform := c.Request.Header.Get("Platform")
		// 压缩长度
		hash := md5.New()
		hash.Write([]byte(parsedURL.RawQuery + platform))
		hashValue := hash.Sum(nil)
		key += hex.EncodeToString(hashValue)
		cacheData, err := gredis.Redis[db].Get(c, key).Result()
		if err == nil && cacheData != "" {
			var resData map[string]any
			err := json.Unmarshal([]byte(cacheData), &resData)
			if err == nil {
				ResSuccess(context.Background(), c, resData["data"])
				c.Abort()
				return
			}
		}

		// 创建一个 ResponseRecorder
		recorder := &responseRecorder{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = recorder

		c.Next()

		// 读取响应体
		responseBody := recorder.body.String()
		if responseBody != "" {
			// 缓存
			gredis.Redis[db].SetNX(c, key, responseBody, time.Duration(seconds)*time.Second)
		}
	}
}

// 自定义 ResponseRecorder，用于记录响应体
type responseRecorder struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// 接口限流中间件
