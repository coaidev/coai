package web

import (
	"chat/globals"
	"chat/manager/conversation"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type recordedCall struct {
	Method string `json:"method"`
	ID     int    `json:"id"`
	Params struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	} `json:"params"`
}

func searchConfig(t *testing.T, provider, endpoint string) {
	t.Helper()
	oldProvider, oldEndpoint, oldCrop, oldLength := globals.SearchProvider, globals.SearchEndpoint, globals.SearchCrop, globals.SearchCropLength
	oldClient := parallelHTTPClient
	globals.SearchProvider, globals.SearchEndpoint, globals.SearchCrop = provider, endpoint, false
	t.Cleanup(func() {
		globals.SearchProvider, globals.SearchEndpoint, globals.SearchCrop, globals.SearchCropLength = oldProvider, oldEndpoint, oldCrop, oldLength
		parallelHTTPClient = oldClient
	})
	for _, key := range []string{"PARALLEL_API_KEY", "BRAVE_API_KEY", "TAVILY_API_KEY"} {
		t.Setenv(key, "")
	}
}

func jsonResponse(id int, result interface{}, sse bool) *http.Response {
	data, _ := json.Marshal(map[string]interface{}{"jsonrpc": "2.0", "id": id, "result": result})
	body, contentType := string(data), "application/json"
	if sse {
		body = ": heartbeat\n\ndata: {\"jsonrpc\":\"2.0\",\"method\":\"notifications/progress\"}\n\ndata: " + body + "\n\n"
		contentType = "text/event-stream; charset=utf-8"
	}
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{contentType}}, Body: io.NopCloser(strings.NewReader(body))}
}

// Exercise the production HTTP client and dispatcher while observing every
// outgoing destination, header and tool argument. No model inference is mocked
// or claimed: these are deterministic search-augmentation caller tests.
func mcpFixture(t *testing.T, sse bool, alter func(recordedCall, *http.Response) *http.Response) *[]recordedCall {
	t.Helper()
	calls := &[]recordedCall{}
	var mu sync.Mutex
	parallelHTTPClient = &http.Client{Timeout: parallelTimeout, Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost || r.URL.String() != parallelEndpoint {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL)
		}
		if r.Header.Get("User-Agent") != parallelUserAgent || r.Header.Get("Authorization") != "" {
			t.Errorf("incorrect attribution or ambient credentials: %v", r.Header)
		}
		if r.Header.Get("Content-Type") != "application/json" || r.Header.Get("Accept") != "application/json, text/event-stream" {
			t.Errorf("incorrect MCP content negotiation: %v", r.Header)
		}
		if deadline, ok := r.Context().Deadline(); !ok || time.Until(deadline) > parallelTimeout {
			t.Error("request has no bounded deadline")
		}
		var call recordedCall
		if err := json.NewDecoder(r.Body).Decode(&call); err != nil {
			t.Error(err)
		}
		mu.Lock()
		*calls = append(*calls, call)
		mu.Unlock()
		if call.Method != "initialize" && (r.Header.Get("Mcp-Session-Id") != "transport-session" || r.Header.Get("MCP-Protocol-Version") != "2025-03-26") {
			t.Error("MCP handshake state was not propagated")
		}
		var resp *http.Response
		switch call.Method {
		case "initialize":
			if call.ID != 1 || r.Header.Get("Mcp-Session-Id") != "" {
				t.Error("new augmentation reused transport state")
			}
			resp = jsonResponse(call.ID, map[string]interface{}{"protocolVersion": "2025-03-26"}, sse)
			resp.Header.Set("Mcp-Session-Id", "transport-session")
		case "notifications/initialized":
			if call.ID != 0 {
				t.Error("notification must not have an ID")
			}
			resp = &http.Response{StatusCode: 202, Header: http.Header{}, Body: http.NoBody}
		case "tools/list":
			resp = jsonResponse(call.ID, map[string]interface{}{"tools": []map[string]string{{"name": "web_search"}, {"name": "web_fetch"}}}, sse)
		case "tools/call":
			if _, ok := call.Params.Arguments["session_id"]; ok {
				t.Error("transport state must not be used as a conversation identity")
			}
			resp = jsonResponse(call.ID, map[string]interface{}{"content": []map[string]string{{"type": "text", "text": "COAI source (https://coai.dev): current information"}}}, sse)
		default:
			t.Errorf("unexpected MCP method %q", call.Method)
			return nil, errors.New("unexpected method")
		}
		if alter != nil {
			resp = alter(call, resp)
		}
		return resp, nil
	})}
	return calls
}

func TestParallelChatAugmentation(t *testing.T) {
	for _, sse := range []bool{false, true} {
		t.Run(fmt.Sprintf("sse=%v", sse), func(t *testing.T) {
			searchConfig(t, "", "") // Fresh unconfigured installation.
			calls := mcpFixture(t, sse, nil)
			instance := conversation.NewAnonymousConversation()
			instance.SetEnableWeb(true)
			instance.AddMessageFromUser("latest COAI information")
			messages := ToChatSearchedContext(context.Background(), instance, false)
			if len(messages) != 2 || messages[0].Role != globals.System || !strings.Contains(messages[0].Content, "https://coai.dev") || messages[1].Content != "latest COAI information" {
				t.Fatalf("search result did not reach the model prompt: %+v", messages)
			}
			if instance.GetMessageLength() != 1 {
				t.Fatal("search context must not be stored in conversation history")
			}
			if len(*calls) != 4 || (*calls)[3].Params.Name != "web_search" || (*calls)[3].Params.Arguments["objective"] != "latest COAI information" || !reflect.DeepEqual((*calls)[3].Params.Arguments["search_queries"], []interface{}{"latest COAI information"}) {
				t.Fatalf("incorrect search dispatch: %+v", *calls)
			}
			// Follow-up augmentation through the same HTTP pool gets its own
			// transport handshake, without inventing a per-turn session_id.
			instance.AddMessageFromAssistant("previous reply")
			instance.AddMessageFromUser("COAI deployment updates")
			ToChatSearchedContext(context.Background(), instance, false)
			if len(*calls) != 8 || (*calls)[7].Params.Arguments["objective"] != "COAI deployment updates" {
				t.Fatal("follow-up did not use its own query")
			}
		})
	}
}

func TestParallelURLFetchAndCrop(t *testing.T) {
	searchConfig(t, "parallel", "https://saved-searx.example")
	calls := mcpFixture(t, true, nil)
	message := []globals.Message{{Role: globals.User, Content: "https://coai.dev/docs?a=1&b=2"}}
	result := ToSearchedContext(context.Background(), true, message)
	if len(result) != 2 || !strings.Contains(result[0].Content, "current information") {
		t.Fatalf("fetch content was not consumed: %+v", result)
	}
	call := (*calls)[3]
	if call.Params.Name != "web_fetch" || !reflect.DeepEqual(call.Params.Arguments["urls"], []interface{}{message[0].Content}) {
		t.Fatalf("URL was altered or not fetched: %+v", call)
	}
	globals.SearchCrop, globals.SearchCropLength = true, 10
	content, err := GenerateSearchResult("COAI updates")
	if err != nil || content != "COAI sourc..." {
		t.Fatalf("crop option not preserved: %q %v", content, err)
	}
}

func TestParallelConcurrentConversations(t *testing.T) {
	searchConfig(t, "parallel", "")
	mcpFixture(t, false, func(call recordedCall, resp *http.Response) *http.Response {
		if call.Method == "tools/call" {
			return jsonResponse(call.ID, map[string]interface{}{"content": []map[string]string{{"type": "text", "text": call.Params.Arguments["objective"].(string)}}}, false)
		}
		return resp
	})
	var wg sync.WaitGroup
	for _, query := range []string{"first conversation", "second conversation"} {
		wg.Add(1)
		go func(query string) {
			defer wg.Done()
			instance := conversation.NewAnonymousConversation()
			instance.SetEnableWeb(true)
			instance.AddMessageFromUser(query)
			messages := ToChatSearchedContext(context.Background(), instance, false)
			if len(messages) != 2 || !strings.Contains(messages[0].Content, "search results: "+query) {
				t.Errorf("conversation received another search result: %+v", messages)
			}
		}(query)
	}
	wg.Wait()
}

func TestSearchChoicesAndDisabled(t *testing.T) {
	searchConfig(t, "", "")
	calls := mcpFixture(t, false, nil)
	message := []globals.Message{{Role: globals.User, Content: "COAI information"}}
	if got := ToSearched(false, message); !reflect.DeepEqual(got, message) || len(*calls) != 0 {
		t.Fatal("disabled search made a request")
	}
	instance := conversation.NewAnonymousConversation()
	instance.AddMessageFromUser("COAI information")
	ToChatSearched(instance, false)
	if len(*calls) != 0 || len(ToSearched(true, nil)) != 0 {
		t.Fatal("disabled chat or empty messages made a request")
	}
	if _, err := GenerateSearchResult(" "); err == nil || len(*calls) != 0 {
		t.Fatal("empty query made a request")
	}
	for _, provider := range []string{"searxng", "invalid"} {
		globals.SearchProvider = provider
		if _, err := GenerateSearchResult("COAI"); err == nil || len(*calls) != 0 {
			t.Fatalf("explicit %s unexpectedly fell back", provider)
		}
	}
	var query string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.Query().Get("q")
		fmt.Fprint(w, `{"results":[{"title":"Saved source","url":"https://example.org","content":"SearXNG content"}]}`)
	}))
	defer server.Close()
	for _, provider := range []string{"", "auto", "searxng"} {
		globals.SearchProvider, globals.SearchEndpoint = provider, server.URL
		content, err := GenerateSearchResult("spaces & unicode 世界")
		if err != nil || query != "spaces & unicode 世界" || !strings.Contains(content, "SearXNG content") || len(*calls) != 0 {
			t.Fatalf("saved endpoint was not preserved: %q %q %v", query, content, err)
		}
	}
	globals.SearchEndpoint = ":invalid-url"
	if _, err := GenerateSearchResult("COAI"); err == nil || len(*calls) != 0 {
		t.Fatal("failed explicit SearXNG endpoint silently fell back")
	}
}

func TestParallelFailuresLeaveMessagesUnchanged(t *testing.T) {
	cases := []string{"http", "rpc", "tool", "malformed", "identity", "missing-tool", "no-text", "oversize", "version", "content-type", "truncated-stream"}
	for _, failure := range cases {
		t.Run(failure, func(t *testing.T) {
			searchConfig(t, "auto", "")
			mcpFixture(t, false, func(call recordedCall, resp *http.Response) *http.Response {
				if call.Method == "initialize" && failure == "version" {
					return jsonResponse(call.ID, map[string]string{"protocolVersion": "invalid"}, false)
				}
				if call.Method == "tools/list" && failure == "missing-tool" {
					return jsonResponse(call.ID, map[string]interface{}{"tools": []interface{}{}}, false)
				}
				if call.Method != "tools/call" {
					return resp
				}
				switch failure {
				case "http":
					resp.StatusCode = 429
				case "rpc":
					resp.Body = io.NopCloser(strings.NewReader(fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"error":{"code":-32603,"message":"unavailable"}}`, call.ID)))
				case "tool":
					return jsonResponse(call.ID, map[string]interface{}{"isError": true}, false)
				case "malformed":
					resp.Body = io.NopCloser(strings.NewReader("invalid JSON"))
				case "identity":
					return jsonResponse(42, map[string]string{}, false)
				case "no-text":
					return jsonResponse(call.ID, map[string]interface{}{"content": []interface{}{}}, false)
				case "oversize":
					resp.Body = io.NopCloser(strings.NewReader(strings.Repeat("x", maxMCPResponse+1)))
				case "content-type":
					resp.Header.Set("Content-Type", "text/html")
				case "truncated-stream":
					resp.Header.Set("Content-Type", "text/event-stream")
					resp.Body = io.NopCloser(strings.NewReader("data: {}\n\n"))
				}
				return resp
			})
			message := []globals.Message{{Role: globals.User, Content: "COAI information"}}
			if got := ToSearched(true, message); !reflect.DeepEqual(got, message) {
				t.Fatalf("failed search injected fabricated evidence: %+v", got)
			}
		})
	}
}

func TestParallelCancellationAndDeadline(t *testing.T) {
	for _, deadline := range []bool{false, true} {
		t.Run(fmt.Sprint(deadline), func(t *testing.T) {
			searchConfig(t, "parallel", "")
			parallelHTTPClient = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				<-r.Context().Done()
				return nil, r.Context().Err()
			})}
			ctx, cancel := context.WithCancel(context.Background())
			if deadline {
				ctx, cancel = context.WithTimeout(context.Background(), 5*time.Millisecond)
			} else {
				cancel()
			}
			defer cancel()
			message := []globals.Message{{Role: globals.User, Content: "COAI"}}
			if got := ToSearchedContext(ctx, true, message); !reflect.DeepEqual(got, message) {
				t.Fatal("canceled search changed the prompt")
			}
		})
	}
}

func TestAdminSearchUsesResolver(t *testing.T) {
	searchConfig(t, "auto", "")
	calls := mcpFixture(t, false, nil)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest("GET", "/search?query=COAI", nil)
	TestSearch(c)
	if !strings.Contains(recorder.Body.String(), `"status":true`) || !strings.Contains(recorder.Body.String(), "https://coai.dev") || len(*calls) != 4 {
		t.Fatalf("admin search failed: %s", recorder.Body.String())
	}
}

func TestLiveParallelCaller(t *testing.T) {
	if os.Getenv("COAI_LIVE_SEARCH") != "1" {
		t.Skip("set COAI_LIVE_SEARCH=1 to exercise the anonymous endpoint")
	}
	searchConfig(t, "", "")
	parallelHTTPClient = &http.Client{Timeout: parallelTimeout, Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		t.Logf("outbound: %s %s User-Agent=%q Authorization=%q", r.Method, r.URL, r.UserAgent(), r.Header.Get("Authorization"))
		resp, err := http.DefaultTransport.RoundTrip(r)
		if err == nil && resp.StatusCode >= 400 {
			data, _ := io.ReadAll(io.LimitReader(resp.Body, maxMCPResponse))
			resp.Body.Close()
			resp.Body = io.NopCloser(strings.NewReader(string(data)))
			t.Logf("live response: status=%d headers=%v body=%s", resp.StatusCode, resp.Header, data)
		}
		return resp, err
	})}
	message := []globals.Message{{Role: globals.User, Content: "What is COAI.Dev?"}}
	result := ToSearchedContext(context.Background(), true, message)
	if len(result) != 2 || !strings.Contains(result[0].Content, "http") {
		t.Fatal("live search did not return sources to the chat prompt")
	}
	t.Logf("live prompt excerpt: %.1000s", result[0].Content)
}
