package service

import (
	"bytes"
	"encoding/base64"
	"errors"
	"testing"
)

func testAIKeyCipher() *AIKeyCipher {
	return NewAIKeyCipher(base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{42}, 32)))
}

func TestAIKeyCipherRoundTripAndRestart(t *testing.T) {
	first := testAIKeyCipher()
	encrypted, err := first.Encrypt("provider-1", "test-provider-secret")
	if err != nil {
		t.Fatal(err)
	}
	second, err := first.Encrypt("provider-1", "test-provider-secret")
	if err != nil || second == encrypted {
		t.Fatal("encryption must use a new nonce")
	}
	plain, err := testAIKeyCipher().Decrypt("provider-1", encrypted)
	if err != nil || plain != "test-provider-secret" {
		t.Fatal("a new instance with the same deployment secret must decrypt saved keys")
	}
	if _, err = first.Decrypt("provider-2", encrypted); !errors.Is(err, ErrAIKeyEncryption) {
		t.Fatal("ciphertext must be bound to its provider")
	}
	sealed, _ := base64.StdEncoding.DecodeString(encrypted)
	sealed[len(sealed)-1] ^= 1
	if _, err = first.Decrypt("provider-1", base64.StdEncoding.EncodeToString(sealed)); !errors.Is(err, ErrAIKeyEncryption) {
		t.Fatal("tampered ciphertext must fail authentication")
	}
	other := NewAIKeyCipher(base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{43}, 32)))
	if _, err = other.Decrypt("provider-1", encrypted); !errors.Is(err, ErrAIKeyEncryption) {
		t.Fatal("a different deployment secret must not decrypt the key")
	}
}

func TestAIKeyCipherRejectsInvalidInput(t *testing.T) {
	for _, encoded := range []string{"", "not-base64", base64.StdEncoding.EncodeToString([]byte("too-short"))} {
		if _, err := NewAIKeyCipher(encoded).Encrypt("p", "secret"); !errors.Is(err, ErrAIKeyEncryption) {
			t.Fatal("missing or invalid deployment secret must be rejected")
		}
	}
	for _, encrypted := range []string{"", "invalid-base64", base64.StdEncoding.EncodeToString([]byte{1, 2, 3})} {
		if _, err := testAIKeyCipher().Decrypt("p", encrypted); !errors.Is(err, ErrAIKeyEncryption) {
			t.Fatal("malformed ciphertext must return a safe error")
		}
	}
}
