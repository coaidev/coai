package channel

import (
	"chat/globals"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestSearchProviderConfigAndEnvironment(t *testing.T) {
	defer viper.Reset()
	oldProvider, oldEndpoint := globals.SearchProvider, globals.SearchEndpoint
	defer func() { globals.SearchProvider, globals.SearchEndpoint = oldProvider, oldEndpoint }()
	for _, test := range []struct {
		name, config, override, want string
	}{
		{"legacy", "endpoint: https://searx.example", "", ""},
		{"new", "provider: auto\n    endpoint: ''", "", "auto"},
		{"explicit", "provider: searxng\n    endpoint: https://searx.example", "", "searxng"},
		{"legacy-env-override", "endpoint: https://searx.example", "parallel", "parallel"},
	} {
		t.Run(test.name, func(t *testing.T) {
			viper.Reset()
			viper.SetConfigType("yaml")
			viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
			viper.AutomaticEnv()
			t.Setenv("SYSTEM_SEARCH_PROVIDER", test.override)
			if err := viper.ReadConfig(strings.NewReader("system:\n  search:\n    " + test.config)); err != nil {
				t.Fatal(err)
			}
			config := NewSystemConfig()
			if config.Search.Provider != test.want || globals.SearchProvider != test.want {
				t.Fatalf("provider = %q / %q, want %q", config.Search.Provider, globals.SearchProvider, test.want)
			}
			if globals.SearchEndpoint != config.Search.Endpoint {
				t.Fatal("configured endpoint not loaded")
			}
		})
	}
}
