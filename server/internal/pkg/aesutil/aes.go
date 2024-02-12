// Copyright 2024 chenmingyong0423

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aesutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

var key []byte

func init() {
	var err error
	key, err = generateRandomBytes(32) // 生成256位密钥
	if err != nil {
		panic(err)
	}
}

// generateRandomBytes 生成指定长度的随机字节
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return nil, err
	}
	return b, nil
}

// AesEncrypt 加密给定的消息
func AesEncrypt(plainText []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 填充原文以满足AES块大小
	padding := block.BlockSize() - len(plainText)%block.BlockSize()
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	paddedText := append(plainText, padText...)

	// 初始化向量IV必须是唯一的，但不需要保密
	iv, err := generateRandomBytes(block.BlockSize())
	if err != nil {
		return "", err
	}

	// 加密
	ciphertext := make([]byte, len(paddedText))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedText)

	// 将IV附加到密文前以便解密时使用
	encrypted := base64.StdEncoding.EncodeToString(append(iv, ciphertext...))
	return encrypted, nil
}

// AesDecrypt 解密给定的消息
func AesDecrypt(encrypted string) (string, error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(encryptedBytes) < block.BlockSize() {
		return "", fmt.Errorf("ciphertext too short")
	}

	// 提取IV
	iv := encryptedBytes[:block.BlockSize()]
	encryptedBytes = encryptedBytes[block.BlockSize():]

	// 解密
	decrypted := make([]byte, len(encryptedBytes))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, encryptedBytes)

	// 移除填充
	padding := decrypted[len(decrypted)-1]
	if int(padding) > len(decrypted) || padding == 0 {
		return "", fmt.Errorf("invalid padding")
	}
	padLen := int(padding)
	for _, val := range decrypted[len(decrypted)-padLen:] {
		if val != padding {
			return "", fmt.Errorf("invalid padding")
		}
	}
	decrypted = decrypted[:len(decrypted)-padLen]

	return string(decrypted), nil
}
