package website_config

import (
	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/service"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func initAIKeyCipher() *service.AIKeyCipher {
	return service.NewAIKeyCipher(viper.GetString("ai_assistant.encryption_key"))
}

func initAIProviderTester() *service.AIProviderTester {
	return service.NewAIProviderTester(&http.Client{Timeout: 30 * time.Second})
}
