package globals

type Hook func(data *Chunk) error

type Message struct {
	Role             string        `json:"role"`
	Content          string        `json:"content"`
	Name             *string       `json:"name,omitempty"`
	FunctionCall     *FunctionCall `json:"function_call,omitempty"`     // only `function` role
	ToolCallId       *string       `json:"tool_call_id,omitempty"`      // only `tool` role
	ToolCalls        *ToolCalls    `json:"tool_calls,omitempty"`        // only `assistant` role
	ReasoningContent *string       `json:"reasoning_content,omitempty"` // only for deepseek reasoner models
	MultiModal       *MultiModal   `json:"multi_mod_content,omitempty"` // gemini inline data (e.g., image base64)
}

type Chunk struct {
	Content      string        `json:"content"`
	ToolCall     *ToolCalls    `json:"tool_call,omitempty"`
	FunctionCall *FunctionCall `json:"function_call,omitempty"`
}

type InlineData struct {
	Data     string `json:"data"`
	MimeType string `json:"mime_type"`
}

type MultiModalContent struct {
	InlineData *InlineData `json:"inline_data,omitempty"`
}

type MultiModal []MultiModalContent

type ChatSegmentResponse struct {
	Conversation int64   `json:"conversation"`
	Quota        float32 `json:"quota"`
	Keyword      string  `json:"keyword"`
	Message      string  `json:"message"`
	End          bool    `json:"end"`
	Plan         bool    `json:"plan"`
}

type GenerationSegmentResponse struct {
	Quota   float32 `json:"quota"`
	Message string  `json:"message"`
	Hash    string  `json:"hash"`
	End     bool    `json:"end"`
	Error   string  `json:"error"`
}

type ListModels struct {
	Object string           `json:"object"`
	Data   []ListModelsItem `json:"data"`
}

type ListModelsItem struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

type ProxyConfig struct {
	ProxyType int    `json:"proxy_type" mapstructure:"proxytype"`
	Proxy     string `json:"proxy" mapstructure:"proxy"`
	Username  string `json:"username" mapstructure:"username"`
	Password  string `json:"password" mapstructure:"password"`
}
