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

package jwtutil

import (
	"crypto/rand"
	"time"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/aesutil"

	"github.com/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtKey []byte
)

func init() {
	jwtKey = make([]byte, 32) // 生成32字节（256位）的密钥
	if _, err := rand.Read(jwtKey); err != nil {
		panic(err) // 生成密钥时发生错误
	}
}

// GenerateJwt 生成 JWT
func GenerateJwt() (string, int64, error) {
	now := time.Now().Local()
	exp := now.Add(time.Hour * 12)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "https://github.com/chenmingyong0423/fnote",
		Subject:   "",
		Audience:  nil,
		ExpiresAt: jwt.NewNumericDate(exp),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
	})
	signedString, err := t.SignedString(jwtKey)
	if err != nil {
		return "", 0, errors.Wrap(err, "generate jwt failed")
	}

	// aes 加密
	encrypt, err := aesutil.AesEncrypt([]byte(signedString))
	if err != nil {
		return "", 0, err
	}
	return encrypt, exp.Unix(), nil
}

func ParseJwt(jwtStr string) (jwt.Claims, error) {
	claims := &jwt.RegisteredClaims{}
	decrypt, err := aesutil.AesDecrypt(jwtStr)
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(decrypt, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
