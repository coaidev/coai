package minimax

import (
	adaptercommon "chat/adapter/common"
	"chat/globals"
	"chat/utils"
	"errors"
	"fmt"
)

type ChatInstance struct {
	Endpoint string
	ApiKey   string
}

func (c *ChatInstance) GetEndpoint() string {
	return c.Endpoint
}

func (c *ChatInstance) GetApiKey() string {
	return c.ApiKey
}

func (c *ChatInstance) GetHeader() map[string]string {
	return map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", c.GetApiKey()),
	}
}

func NewChatInstance(endpoint, apiKey string) *ChatInstance {
	return &ChatInstance{
		Endpoint: endpoint,
		ApiKey:   apiKey,
	}
}

func NewChatInstanceFromConfig(conf globals.ChannelConfig) adaptercommon.Factory {
	return NewChatInstance(
		conf.GetEndpoint(),
		conf.GetRandomSecret(),
	)
}

func (c *ChatInstance) GetChatEndpoint() string {
	return fmt.Sprintf("%s/v1/chat/completions", c.GetEndpoint())
}

// clampTemperature ensures temperature stays within MiniMax's valid range [0, 1.0]
func clampTemperature(temperature *float32) *float32 {
	if temperature == nil {
		return nil
	}
	t := *temperature
	if t < 0 {
		t = 0
	} else if t > 1.0 {
		t = 1.0
	}
	return &t
}

func (c *ChatInstance) GetChatBody(props *adaptercommon.ChatProps, stream bool) interface{} {
	return ChatRequest{
		Model:            props.Model,
		Messages:         props.Message,
		MaxTokens:        props.MaxTokens,
		Stream:           stream,
		Temperature:      clampTemperature(props.Temperature),
		TopP:             props.TopP,
		PresencePenalty:  props.PresencePenalty,
		FrequencyPenalty: props.FrequencyPenalty,
	}
}

func processChatResponse(data string) *ChatResponse {
	return utils.UnmarshalForm[ChatResponse](data)
}

func processChatStreamResponse(data string) *ChatStreamResponse {
	return utils.UnmarshalForm[ChatStreamResponse](data)
}

func processChatErrorResponse(data string) *ChatStreamErrorResponse {
	return utils.UnmarshalForm[ChatStreamErrorResponse](data)
}

func (c *ChatInstance) ProcessLine(data string) (string, error) {
	if form := processChatStreamResponse(data); form != nil {
		if len(form.Choices) == 0 {
			return "", nil
		}
		return form.Choices[0].Delta.Content, nil
	}

	if form := processChatErrorResponse(data); form != nil {
		if form.Error.Message != "" {
			return "", fmt.Errorf("minimax error: %s (type: %s)", form.Error.Message, form.Error.Type)
		}
	}

	return "", nil
}

func (c *ChatInstance) CreateChatRequest(props *adaptercommon.ChatProps) (string, error) {
	res, err := utils.Post(
		c.GetChatEndpoint(),
		c.GetHeader(),
		c.GetChatBody(props, false),
		props.Proxy,
	)

	if err != nil || res == nil {
		return "", fmt.Errorf("minimax error: %s", err.Error())
	}

	data := utils.MapToStruct[ChatResponse](res)
	if data == nil {
		return "", fmt.Errorf("minimax error: cannot parse response")
	}

	if len(data.Choices) == 0 {
		return "", fmt.Errorf("minimax error: no choices")
	}

	return data.Choices[0].Message.Content, nil
}

func (c *ChatInstance) CreateStreamChatRequest(props *adaptercommon.ChatProps, callback globals.Hook) error {
	err := utils.EventScanner(&utils.EventScannerProps{
		Method:  "POST",
		Uri:     c.GetChatEndpoint(),
		Headers: c.GetHeader(),
		Body:    c.GetChatBody(props, true),
		Callback: func(data string) error {
			partial, err := c.ProcessLine(data)
			if err != nil {
				return err
			}
			return callback(&globals.Chunk{Content: partial})
		},
	}, props.Proxy)

	if err != nil {
		if form := processChatErrorResponse(err.Body); form != nil {
			if form.Error.Type == "" && form.Error.Message == "" {
				return errors.New(utils.ToMarkdownCode("json", err.Body))
			}
			return fmt.Errorf("minimax error: %s (type: %s)", form.Error.Message, form.Error.Type)
		}
		return err.Error
	}

	return nil
}
