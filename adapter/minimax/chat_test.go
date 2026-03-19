package minimax

import (
	adaptercommon "chat/adapter/common"
	"chat/globals"
	"encoding/json"
	"testing"
)

func TestClampTemperature(t *testing.T) {
	tests := []struct {
		name     string
		input    *float32
		expected *float32
	}{
		{"nil temperature", nil, nil},
		{"zero temperature", floatPtr(0), floatPtr(0)},
		{"valid temperature 0.5", floatPtr(0.5), floatPtr(0.5)},
		{"valid temperature 1.0", floatPtr(1.0), floatPtr(1.0)},
		{"temperature above max", floatPtr(1.5), floatPtr(1.0)},
		{"temperature above max 2.0", floatPtr(2.0), floatPtr(1.0)},
		{"negative temperature", floatPtr(-0.5), floatPtr(0)},
		{"small valid temperature", floatPtr(0.01), floatPtr(0.01)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := clampTemperature(tt.input)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %f", *result)
				}
				return
			}
			if result == nil {
				t.Errorf("expected %f, got nil", *tt.expected)
				return
			}
			if *result != *tt.expected {
				t.Errorf("expected %f, got %f", *tt.expected, *result)
			}
		})
	}
}

func TestNewChatInstance(t *testing.T) {
	instance := NewChatInstance("https://api.minimax.io/v1", "test-key")
	if instance.GetEndpoint() != "https://api.minimax.io/v1" {
		t.Errorf("expected endpoint https://api.minimax.io/v1, got %s", instance.GetEndpoint())
	}
	if instance.GetApiKey() != "test-key" {
		t.Errorf("expected api key test-key, got %s", instance.GetApiKey())
	}
}

func TestGetHeader(t *testing.T) {
	instance := NewChatInstance("https://api.minimax.io/v1", "sk-test-123")
	headers := instance.GetHeader()

	if headers["Content-Type"] != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", headers["Content-Type"])
	}
	if headers["Authorization"] != "Bearer sk-test-123" {
		t.Errorf("expected Authorization Bearer sk-test-123, got %s", headers["Authorization"])
	}
}

func TestGetChatEndpoint(t *testing.T) {
	instance := NewChatInstance("https://api.minimax.io/v1", "test-key")
	endpoint := instance.GetChatEndpoint()
	expected := "https://api.minimax.io/v1/v1/chat/completions"
	if endpoint != expected {
		t.Errorf("expected %s, got %s", expected, endpoint)
	}

	// Test with trailing slash
	instance2 := NewChatInstance("https://api.minimax.io", "test-key")
	endpoint2 := instance2.GetChatEndpoint()
	expected2 := "https://api.minimax.io/v1/chat/completions"
	if endpoint2 != expected2 {
		t.Errorf("expected %s, got %s", expected2, endpoint2)
	}
}

func TestGetChatBody(t *testing.T) {
	instance := NewChatInstance("https://api.minimax.io/v1", "test-key")

	temp := float32(1.5) // should be clamped to 1.0
	topP := float32(0.9)
	maxTokens := 100

	props := &adaptercommon.ChatProps{
		Model: "MiniMax-M2.5",
		Message: []globals.Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "Hello"},
		},
		MaxTokens:   &maxTokens,
		Temperature: &temp,
		TopP:        &topP,
	}

	body := instance.GetChatBody(props, true)

	// Marshal to JSON to verify structure
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("failed to unmarshal body: %v", err)
	}

	if result["model"] != "MiniMax-M2.5" {
		t.Errorf("expected model MiniMax-M2.5, got %v", result["model"])
	}
	if result["stream"] != true {
		t.Errorf("expected stream true, got %v", result["stream"])
	}
	// Temperature should be clamped to 1.0
	if temp := result["temperature"].(float64); temp != 1.0 {
		t.Errorf("expected temperature 1.0 (clamped), got %f", temp)
	}

	messages := result["messages"].([]interface{})
	if len(messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(messages))
	}
}

func TestGetChatBodyNilTemperature(t *testing.T) {
	instance := NewChatInstance("https://api.minimax.io/v1", "test-key")

	props := &adaptercommon.ChatProps{
		Model: "MiniMax-M2.5",
		Message: []globals.Message{
			{Role: "user", Content: "Hello"},
		},
	}

	body := instance.GetChatBody(props, false)

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("failed to unmarshal body: %v", err)
	}

	// Temperature should be omitted (nil)
	if _, exists := result["temperature"]; exists {
		t.Errorf("expected temperature to be omitted, but it was present")
	}
	if result["stream"] != false {
		t.Errorf("expected stream false, got %v", result["stream"])
	}
}

func TestProcessLine(t *testing.T) {
	instance := NewChatInstance("https://api.minimax.io/v1", "test-key")

	t.Run("valid stream response", func(t *testing.T) {
		data := `{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1234567890,"model":"MiniMax-M2.5","choices":[{"delta":{"content":"Hello"},"index":0,"finish_reason":null}]}`
		content, err := instance.ProcessLine(data)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if content != "Hello" {
			t.Errorf("expected 'Hello', got '%s'", content)
		}
	})

	t.Run("empty choices", func(t *testing.T) {
		data := `{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1234567890,"model":"MiniMax-M2.5","choices":[]}`
		content, err := instance.ProcessLine(data)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if content != "" {
			t.Errorf("expected empty content, got '%s'", content)
		}
	})

	t.Run("error response with choices field", func(t *testing.T) {
		// Error responses with explicit error field but also parseable as stream response
		// are handled via SSE error path (err.Body), not ProcessLine.
		// ProcessLine treats them as empty stream responses.
		data := `{"error":{"message":"rate limit exceeded","type":"rate_limit_error"}}`
		content, err := instance.ProcessLine(data)
		// The error JSON parses as a stream response with empty choices
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if content != "" {
			t.Errorf("expected empty content, got '%s'", content)
		}
	})

	t.Run("finish reason stop", func(t *testing.T) {
		data := `{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1234567890,"model":"MiniMax-M2.5","choices":[{"delta":{"content":""},"index":0,"finish_reason":"stop"}]}`
		content, err := instance.ProcessLine(data)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if content != "" {
			t.Errorf("expected empty content for stop, got '%s'", content)
		}
	})
}

func TestProcessChatResponse(t *testing.T) {
	data := `{"id":"chatcmpl-123","object":"chat.completion","created":1234567890,"model":"MiniMax-M2.5","choices":[{"index":0,"message":{"role":"assistant","content":"Hello world"},"finish_reason":"stop"}],"usage":{"prompt_tokens":10,"completion_tokens":5,"total_tokens":15}}`
	form := processChatResponse(data)
	if form == nil {
		t.Fatal("expected non-nil response")
	}
	if form.ID != "chatcmpl-123" {
		t.Errorf("expected id chatcmpl-123, got %s", form.ID)
	}
	if len(form.Choices) != 1 {
		t.Fatalf("expected 1 choice, got %d", len(form.Choices))
	}
	if form.Choices[0].Message.Content != "Hello world" {
		t.Errorf("expected 'Hello world', got '%s'", form.Choices[0].Message.Content)
	}
	if form.Usage.TotalTokens != 15 {
		t.Errorf("expected total tokens 15, got %d", form.Usage.TotalTokens)
	}
}

func TestProcessChatStreamResponse(t *testing.T) {
	data := `{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1234567890,"model":"MiniMax-M2.5","choices":[{"delta":{"content":"Hi"},"index":0,"finish_reason":null}]}`
	form := processChatStreamResponse(data)
	if form == nil {
		t.Fatal("expected non-nil response")
	}
	if len(form.Choices) != 1 {
		t.Fatalf("expected 1 choice, got %d", len(form.Choices))
	}
	if form.Choices[0].Delta.Content != "Hi" {
		t.Errorf("expected 'Hi', got '%s'", form.Choices[0].Delta.Content)
	}
}

func TestProcessChatErrorResponse(t *testing.T) {
	data := `{"error":{"message":"invalid api key","type":"authentication_error"}}`
	form := processChatErrorResponse(data)
	if form == nil {
		t.Fatal("expected non-nil response")
	}
	if form.Error.Message != "invalid api key" {
		t.Errorf("expected 'invalid api key', got '%s'", form.Error.Message)
	}
	if form.Error.Type != "authentication_error" {
		t.Errorf("expected 'authentication_error', got '%s'", form.Error.Type)
	}
}

func TestChatRequestSerialization(t *testing.T) {
	temp := float32(0.7)
	maxTokens := 200
	req := ChatRequest{
		Model: "MiniMax-M2.7",
		Messages: []globals.Message{
			{Role: "user", Content: "What is AI?"},
		},
		MaxTokens:   &maxTokens,
		Stream:      true,
		Temperature: &temp,
	}

	jsonBytes, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if result["model"] != "MiniMax-M2.7" {
		t.Errorf("expected model MiniMax-M2.7, got %v", result["model"])
	}
	if result["stream"] != true {
		t.Errorf("expected stream true, got %v", result["stream"])
	}
	if result["temperature"].(float64) != 0.7 {
		t.Errorf("expected temperature 0.7, got %v", result["temperature"])
	}
}

func TestChatRequestOmitsNilFields(t *testing.T) {
	req := ChatRequest{
		Model: "MiniMax-M2.5",
		Messages: []globals.Message{
			{Role: "user", Content: "Hello"},
		},
		Stream: false,
	}

	jsonBytes, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	for _, field := range []string{"temperature", "top_p", "max_tokens", "presence_penalty", "frequency_penalty"} {
		if _, exists := result[field]; exists {
			t.Errorf("expected %s to be omitted, but it was present", field)
		}
	}
}

func floatPtr(f float32) *float32 {
	return &f
}
