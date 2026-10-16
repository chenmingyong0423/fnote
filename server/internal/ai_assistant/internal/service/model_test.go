package service

import (
	"strings"
	"testing"
)

func TestStreamRedactorHandlesSplitKeysAndUnicode(t *testing.T) {
	r := streamRedactor{secret: "provider-private-key"}
	var output strings.Builder
	for _, part := range []string{"你好，pro", "vider-priv", "ate-key! 继续。"} {
		output.WriteString(r.Write(part, false))
	}
	output.WriteString(r.Write("", true))
	if output.String() != "你好，[密钥已隐藏]! 继续。" {
		t.Fatal("a key split across arbitrary frames must be redacted", output.String())
	}
	r = streamRedactor{secret: "abc"}
	if got := r.Write("a", false) + r.Write("bcabc", false) + r.Write("", true); got != "[密钥已隐藏][密钥已隐藏]" {
		t.Fatal("adjacent keys must be hidden", got)
	}
}
