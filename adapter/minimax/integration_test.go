//go:build integration

package minimax

import (
	adaptercommon "chat/adapter/common"
	"chat/globals"
	"os"
	"strings"
	"testing"
)

func getTestInstance(t *testing.T) *ChatInstance {
	apiKey := os.Getenv("MINIMAX_API_KEY")
	if apiKey == "" {
		t.Skip("MINIMAX_API_KEY not set, skipping integration test")
	}
	return NewChatInstance("https://api.minimax.io", apiKey)
}

func TestIntegrationCreateChatRequest(t *testing.T) {
	instance := getTestInstance(t)

	maxTokens := 50
	temp := float32(0.7)
	props := &adaptercommon.ChatProps{
		Model: "MiniMax-M2.5",
		Message: []globals.Message{
			{Role: "user", Content: "Say hello in one sentence."},
		},
		MaxTokens:   &maxTokens,
		Temperature: &temp,
	}

	response, err := instance.CreateChatRequest(props)
	if err != nil {
		t.Fatalf("CreateChatRequest failed: %v", err)
	}
	if response == "" {
		t.Error("expected non-empty response")
	}
	t.Logf("Response: %s", response)
}

func TestIntegrationCreateStreamChatRequest(t *testing.T) {
	instance := getTestInstance(t)

	maxTokens := 50
	temp := float32(0.7)
	props := &adaptercommon.ChatProps{
		Model: "MiniMax-M2.5",
		Message: []globals.Message{
			{Role: "user", Content: "Say hello in one sentence."},
		},
		MaxTokens:   &maxTokens,
		Temperature: &temp,
	}

	var chunks []string
	callback := func(chunk *globals.Chunk) error {
		chunks = append(chunks, chunk.Content)
		return nil
	}

	err := instance.CreateStreamChatRequest(props, callback)
	if err != nil {
		t.Fatalf("CreateStreamChatRequest failed: %v", err)
	}

	fullResponse := strings.Join(chunks, "")
	if fullResponse == "" {
		t.Error("expected non-empty streaming response")
	}
	t.Logf("Streaming response: %s", fullResponse)
}

func TestIntegrationTemperatureClamping(t *testing.T) {
	instance := getTestInstance(t)

	maxTokens := 30
	temp := float32(2.0) // above max, should be clamped to 1.0
	props := &adaptercommon.ChatProps{
		Model: "MiniMax-M2.5",
		Message: []globals.Message{
			{Role: "user", Content: "Say hi."},
		},
		MaxTokens:   &maxTokens,
		Temperature: &temp,
	}

	response, err := instance.CreateChatRequest(props)
	if err != nil {
		t.Fatalf("CreateChatRequest with clamped temperature failed: %v", err)
	}
	if response == "" {
		t.Error("expected non-empty response with clamped temperature")
	}
}
