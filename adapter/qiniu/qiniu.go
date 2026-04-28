package qiniu

import (
	adaptercommon "chat/adapter/common"
	"chat/adapter/openai"
	"chat/globals"
	"strings"
)

// DefaultEndpoint is the host base for OpenAI-compatible paths (/v1/chat/completions).
const DefaultEndpoint = "https://api.qnaigc.com"

func normalizeEndpoint(endpoint string) string {
	e := strings.TrimSpace(endpoint)
	if e == "" {
		return DefaultEndpoint
	}
	e = strings.TrimSuffix(e, "/")
	e = strings.TrimSuffix(e, "/v1")
	return e
}

// NewChatInstanceFromConfig wires the Qiniu gateway using the shared OpenAI-compatible adapter.
func NewChatInstanceFromConfig(conf globals.ChannelConfig) adaptercommon.Factory {
	return openai.NewChatInstance(
		normalizeEndpoint(conf.GetEndpoint()),
		conf.GetRandomSecret(),
	)
}
