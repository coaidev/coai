package web

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const parallelEndpoint = "https://search.parallel.ai/mcp"
const parallelUserAgent = "CoAI/2.6.0 (+https://github.com/coaidev/coai)"
const parallelTimeout = 30 * time.Second
const maxMCPResponse = 1024 * 1024

// Only the HTTP connection pool is shared. MCP transport sessions are scoped to
// one augmentation request, and are not conversation identities.
var parallelHTTPClient = &http.Client{Timeout: parallelTimeout}

type mcpSession struct {
	id      int
	session string
	version string
}

type mcpResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (s *mcpSession) request(ctx context.Context, method string, params interface{}, result interface{}) error {
	message := map[string]interface{}{"jsonrpc": "2.0", "method": method, "params": params}
	if result != nil {
		s.id++
		message["id"] = s.id
	}
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, parallelEndpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	req.Header.Set("User-Agent", parallelUserAgent)
	if s.session != "" {
		req.Header.Set("Mcp-Session-Id", s.session)
	}
	if s.version != "" {
		req.Header.Set("MCP-Protocol-Version", s.version)
	}
	resp, err := parallelHTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Parallel MCP returned HTTP %d", resp.StatusCode)
	}
	if method == "initialize" {
		s.session = resp.Header.Get("Mcp-Session-Id")
	}
	if result == nil {
		return nil
	}

	contentType, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	var response mcpResponse
	reader := io.LimitReader(resp.Body, maxMCPResponse+1)
	switch contentType {
	case "application/json":
		data, err := io.ReadAll(reader)
		if err != nil {
			return err
		}
		if len(data) > maxMCPResponse {
			return errors.New("Parallel MCP response exceeds size limit")
		}
		if err := json.Unmarshal(data, &response); err != nil {
			return fmt.Errorf("invalid Parallel MCP response: %w", err)
		}
	case "text/event-stream":
		// Streamable HTTP permits either JSON or SSE for each POST response.
		scanner := bufio.NewScanner(reader)
		scanner.Buffer(make([]byte, 4096), maxMCPResponse)
		var data []string
		found := false
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "data:") {
				data = append(data, strings.TrimPrefix(strings.TrimPrefix(line, "data:"), " "))
			}
			if line == "" && len(data) > 0 {
				response = mcpResponse{}
				if err := json.Unmarshal([]byte(strings.Join(data, "\n")), &response); err != nil {
					return fmt.Errorf("invalid Parallel MCP event: %w", err)
				}
				data = nil
				if response.ID == s.id {
					found = true
					break
				}
			}
		}
		if !found {
			if err := scanner.Err(); err != nil {
				return err
			}
			return errors.New("Parallel MCP stream ended without a response")
		}
	default:
		return errors.New("unsupported Parallel MCP response content type")
	}
	if response.JSONRPC != "2.0" || response.ID != s.id {
		return errors.New("invalid Parallel MCP response identity")
	}
	if response.Error != nil {
		return fmt.Errorf("Parallel MCP error %d: %s", response.Error.Code, response.Error.Message)
	}
	if len(response.Result) == 0 || string(response.Result) == "null" {
		return errors.New("Parallel MCP response has no result")
	}
	return json.Unmarshal(response.Result, result)
}

func parallelSearch(ctx context.Context, query string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, parallelTimeout)
	defer cancel()
	session := &mcpSession{}
	var initialized struct {
		ProtocolVersion string `json:"protocolVersion"`
	}
	if err := session.request(ctx, "initialize", map[string]interface{}{
		"protocolVersion": "2025-03-26",
		"capabilities":    map[string]interface{}{},
		"clientInfo":      map[string]string{"name": "coai", "version": "2.6.0"},
	}, &initialized); err != nil {
		return "", err
	}
	if initialized.ProtocolVersion != "2025-03-26" {
		return "", errors.New("unsupported Parallel MCP protocol version")
	}
	session.version = initialized.ProtocolVersion
	if err := session.request(ctx, "notifications/initialized", map[string]interface{}{}, nil); err != nil {
		return "", err
	}

	tool := "web_search"
	arguments := map[string]interface{}{"objective": query, "search_queries": []string{query}}
	if uri, err := url.ParseRequestURI(query); err == nil && uri.Host != "" && (uri.Scheme == "http" || uri.Scheme == "https") {
		tool = "web_fetch"
		arguments = map[string]interface{}{"urls": []string{query}}
	}
	// session_id is optional. COAI does not persist an opaque search conversation
	// ID, so do not substitute a transport session or derive one from user data.
	var tools struct {
		Tools []struct {
			Name string `json:"name"`
		} `json:"tools"`
	}
	if err := session.request(ctx, "tools/list", map[string]interface{}{}, &tools); err != nil {
		return "", err
	}
	found := false
	for _, available := range tools.Tools {
		if available.Name == tool {
			found = true
		}
	}
	if !found {
		return "", fmt.Errorf("Parallel MCP does not offer %s", tool)
	}
	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		IsError bool `json:"isError"`
	}
	if err := session.request(ctx, "tools/call", map[string]interface{}{"name": tool, "arguments": arguments}, &result); err != nil {
		return "", err
	}
	if result.IsError {
		return "", errors.New("Parallel MCP tool failed; try again later or check the service limits")
	}
	var text []string
	for _, item := range result.Content {
		if item.Type == "text" && strings.TrimSpace(item.Text) != "" {
			text = append(text, item.Text)
		}
	}
	if len(text) == 0 {
		return "", errors.New("Parallel MCP returned no text content")
	}
	return strings.Join(text, "\n"), nil
}
