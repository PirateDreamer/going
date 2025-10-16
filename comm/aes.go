package comm

import (
	"crypto/aes"
	"crypto/cipher"
)

// aesDecrypt AES解密函数
func AesDecrypt(ciphertext []byte, key []byte, iv []byte) ([]byte, error) {
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
