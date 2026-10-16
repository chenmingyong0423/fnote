package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

var ErrAIKeyEncryption = errors.New("AI 模型密钥无法加解密，请检查启动 YAML 配置中的 ai_assistant.encryption_key（Base64 编码的 32 字节密钥）")

// AIKeyCipher uses a deployment secret so saved credentials survive restarts.
// Validation is deferred until use: an unconfigured AI feature must not stop the blog.
type AIKeyCipher struct {
	encodedKey string
}

func NewAIKeyCipher(encodedKey string) *AIKeyCipher {
	return &AIKeyCipher{encodedKey: encodedKey}
}

func (c *AIKeyCipher) aead() (cipher.AEAD, error) {
	key, err := base64.StdEncoding.DecodeString(c.encodedKey)
	if err != nil || len(key) != 32 {
		return nil, ErrAIKeyEncryption
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrAIKeyEncryption
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, ErrAIKeyEncryption
	}
	return aead, nil
}

func (c *AIKeyCipher) Encrypt(providerId, plainText string) (string, error) {
	aead, err := c.aead()
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", ErrAIKeyEncryption
	}
	sealed := aead.Seal(nonce, nonce, []byte(plainText), []byte("ai_assistant:"+providerId))
	return base64.StdEncoding.EncodeToString(sealed), nil
}

func (c *AIKeyCipher) Decrypt(providerId, encrypted string) (string, error) {
	aead, err := c.aead()
	if err != nil {
		return "", err
	}
	sealed, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil || len(sealed) < aead.NonceSize()+aead.Overhead() {
		return "", ErrAIKeyEncryption
	}
	plainText, err := aead.Open(nil, sealed[:aead.NonceSize()], sealed[aead.NonceSize():], []byte("ai_assistant:"+providerId))
	if err != nil {
		return "", ErrAIKeyEncryption
	}
	return string(plainText), nil
}
