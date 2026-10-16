package website_config

import (
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestAIKeyCipherFromShippedYAML(t *testing.T) {
	defer viper.Reset()
	for _, name := range []string{"fnote.yaml", "fnote.docker.yaml", "fnote.nginx.yaml", "fnote.prod.yaml"} {
		t.Run(name, func(t *testing.T) {
			viper.Reset()
			viper.SetConfigFile(filepath.Join("..", "..", "config", name))
			if err := viper.ReadInConfig(); err != nil {
				t.Fatal(err)
			}
			encrypted, err := initAIKeyCipher().Encrypt("provider", "test-provider-secret")
			if err != nil {
				t.Fatal("shipped YAML must contain a usable AI encryption key:", err)
			}
			plain, err := initAIKeyCipher().Decrypt("provider", encrypted)
			if err != nil || plain != "test-provider-secret" {
				t.Fatal("a new cipher from YAML must decrypt saved provider keys")
			}
		})
	}
}
