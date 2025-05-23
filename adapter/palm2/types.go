package palm2

const (
	GeminiUserType  = "user"
	GeminiModelType = "model"
)

type PalmMessage struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

// PalmChatBody is the native http request body for palm2
type PalmChatBody struct {
	Prompt PalmPrompt `json:"prompt"`
}

type PalmPrompt struct {
	Messages []PalmMessage `json:"messages"`
}

// PalmChatResponse is the native http response body for palm2
type PalmChatResponse struct {
	Candidates []PalmMessage `json:"candidates"`
}

// GeminiChatBody is the native http request body for gemini
type GeminiChatBody struct {
	Contents         []GeminiContent `json:"contents"`
	GenerationConfig GeminiConfig    `json:"generationConfig"`
}

type GeminiConfig struct {
	Temperature     *float32 `json:"temperature,omitempty"`
	MaxOutputTokens *int     `json:"maxOutputTokens,omitempty"`
	TopP            *float32 `json:"topP,omitempty"`
	TopK            *int     `json:"topK,omitempty"`
}

type GeminiContent struct {
	Role  string           `json:"role"`
	Parts []GeminiChatPart `json:"parts"`
}

type GeminiChatPart struct {
	Text       *string           `json:"text,omitempty"`
	InlineData *GeminiInlineData `json:"inline_data,omitempty"`
}

type GeminiInlineData struct {
	MimeType string `json:"mime_type"`
	Data     string `json:"data"`
}

type GeminiChatResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
			Role string `json:"role"`
		} `json:"content"`
	} `json:"candidates"`
}

type GeminiChatErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

type GeminiStreamResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
			Role string `json:"role"`
		} `json:"content"`
	} `json:"candidates"`
}

// ImageRequest is the native http request body for imagen
type ImageRequest struct {
	Instances  []ImageInstance `json:"instances"`
	Parameters ImageParameters `json:"parameters"`
}

type ImageInstance struct {
	Prompt string `json:"prompt"`
}

type ImageParameters struct {
	SampleCount      int    `json:"sampleCount,omitempty"`
	AspectRatio      string `json:"aspectRatio,omitempty"`
	PersonGeneration string `json:"personGeneration,omitempty"`
}

// ImageResponse is the native http response body for imagen
type ImageResponse struct {
	Predictions []ImagePrediction `json:"predictions"`
}

type ImagePrediction struct {
	MimeType           string `json:"mimeType"`
	BytesBase64Encoded string `json:"bytesBase64Encoded"`
}
