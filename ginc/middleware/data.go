package middleware

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

// AesApiDataDecrypt 接口数据aes解密,key和iv长度必须为16位
func AesApiDataDecrypt(key, iv string, milliseconds int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 读取原始请求body数据
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "无法读取请求数据"})
			return
		}

		// base64解密
		if body, err = base64.StdEncoding.DecodeString(string(body)); err != nil {
			fmt.Println("base64 decode ciphertext error:", err)
			return
		}

		// 执行AES解密
		decryptedBody, err := aesDecrypt(body, []byte(key), []byte(iv))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "数据解密失败"})
			return
		}

		// 判断接口里的时间和当前时间的差异,大于秒数则返回错误
		if (time.Now().UnixMilli() - gjson.Get(string(decryptedBody), "timestamp").Int()) > int64(milliseconds) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "接口错误"})
			return
		}

		// 将解密后的数据重新放入请求body中
		c.Request.Body = io.NopCloser(bytes.NewBuffer(decryptedBody))

		c.Next()
	}
}

// aesDecrypt AES解密函数
func aesDecrypt(ciphertext []byte, key []byte, iv []byte) ([]byte, error) {
	// 创建AES密码块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建CBC解密器
	mode := cipher.NewCBCDecrypter(block, iv)

	// 解密数据
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// 去除填充
	plaintext = pkcs7Unpadding(plaintext)

	return plaintext, nil
}

// pkcs7Unpadding 去除PKCS7填充
func pkcs7Unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
