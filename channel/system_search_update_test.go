package channel_test

import (
	"bytes"
	"chat/addition/web"
	"chat/channel"
	"chat/globals"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type searchTransport func(*http.Request) (*http.Response, error)

func (f searchTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestSavedProviderReachesAdminSearch(t *testing.T) {
	oldSettings, oldSystem, oldTransport := globals.GetSearchSettings(), channel.SystemInstance, http.DefaultTransport
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		globals.SetSearchSettings(oldSettings)
		channel.SystemInstance, http.DefaultTransport = oldSystem, oldTransport
		viper.Reset()
		if err := os.Chdir(workingDir); err != nil {
			t.Error(err)
		}
	})
	t.Setenv("SYSTEM_SEARCH_PROVIDER", "")
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir("config", 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("config/config.yaml", []byte("system:\n  search:\n    provider: searxng\n    endpoint: https://saved-searx.example\n"), 0600); err != nil {
		t.Fatal(err)
	}
	viper.Reset()
	viper.SetConfigFile("config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		t.Fatal(err)
	}
	channel.SystemInstance = channel.NewSystemConfig()
	if globals.GetSearchSettings().Provider != "searxng" {
		t.Fatal("legacy provider did not load")
	}
	update := func(provider string) *httptest.ResponseRecorder {
		t.Helper()
		body, _ := json.Marshal(map[string]interface{}{"search": map[string]interface{}{"provider": provider, "endpoint": "https://saved-searx.example"}})
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request = httptest.NewRequest("POST", "/system", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		channel.UpdateConfig(c)
		return recorder
	}
	if response := update("parallel"); !strings.Contains(response.Body.String(), `"status":true`) {
		t.Fatalf("save failed: %s", response.Body.String())
	}
	settings := globals.GetSearchSettings()
	if settings.Provider != "parallel" || settings.Endpoint != "https://saved-searx.example" {
		t.Fatalf("save did not activate selected provider: %+v", settings)
	}
	calls := 0
	http.DefaultTransport = searchTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.URL.String() != "https://search.parallel.ai/mcp" || r.UserAgent() != "CoAI/2.6.0 (+https://github.com/coaidev/coai)" {
			return nil, fmt.Errorf("unexpected destination or attribution")
		}
		return &http.Response{StatusCode: 503, Header: http.Header{}, Body: io.NopCloser(strings.NewReader("unavailable"))}, nil
	})
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest("GET", "/search?query=COAI", nil)
	web.TestSearch(c)
	if calls != 1 || !strings.Contains(recorder.Body.String(), "Parallel MCP returned HTTP 503") {
		t.Fatalf("admin test used stale provider or fell back: %s, calls=%d", recorder.Body.String(), calls)
	}
	// A deployment override remains authoritative even after an admin save.
	t.Setenv("SYSTEM_SEARCH_PROVIDER", "parallel")
	if response := update("searxng"); !strings.Contains(response.Body.String(), `"status":true`) || globals.GetSearchSettings().Provider != "parallel" {
		t.Fatal("save bypassed the environment override")
	}
	t.Setenv("SYSTEM_SEARCH_PROVIDER", "")
	// A failed persistence operation must not activate the pending selection.
	if err := os.Mkdir("config/config.tmp.yaml", 0700); err != nil {
		t.Fatal(err)
	}
	if response := update("searxng"); !strings.Contains(response.Body.String(), `"status":false`) || globals.GetSearchSettings().Provider != "parallel" {
		t.Fatal("failed save changed the active provider")
	}
}
