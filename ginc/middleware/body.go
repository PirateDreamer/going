package middleware

import (
	"bytes"
	"encoding/base64"
	"io"
	"time"

	"github.com/PirateDreamer/going/comm"
	"github.com/PirateDreamer/going/ginc"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

type BodyAesDecryptParam struct {
	Key          string // 长度必须为16位
	Iv           string // 长度必须为16位
	Milliseconds int64  // 接口数据时间戳与当前时间戳的差异,大于此值则返回错误
}

// BodyAesDecrypt 请求body数据解密，aes使用CBC解密器，去除PKCS7填充
func BodyAesDecrypt(param BodyAesDecryptParam) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 读取原始请求body数据
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			ginc.ResFail(c.Request.Context(), c, err)
			return
		}

		// base64解密
		if body, err = base64.StdEncoding.DecodeString(string(body)); err != nil {
			ginc.ResFail(c.Request.Context(), c, err)
			return
		}

		// 执行AES解密
		decryptedBody, err := comm.AesDecrypt(body, []byte(param.Key), []byte(param.Iv))
		if err != nil {
			ginc.ResFail(c.Request.Context(), c, err)
			return
		}

		// 判断接口里的时间和当前时间的差异,大于秒数则返回错误
		if (time.Now().UnixMilli() - gjson.Get(string(decryptedBody), "timestamp").Int()) > param.Milliseconds {
			ginc.ResFail(c.Request.Context(), c, err)
			return
		}

		// 将解密后的数据重新放入请求body中
		c.Request.Body = io.NopCloser(bytes.NewBuffer(decryptedBody))

		c.Next()
	}
}
