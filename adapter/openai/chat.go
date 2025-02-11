package openai

import (
	adaptercommon "chat/adapter/common"
	"chat/globals"
	"chat/utils"
	"errors"
	"fmt"
	"regexp"
)

func (c *ChatInstance) GetChatEndpoint(props *adaptercommon.ChatProps) string {
	if props.Model == globals.GPT3TurboInstruct {
		return fmt.Sprintf("%s/v1/completions", c.GetEndpoint())
	}
	return fmt.Sprintf("%s/v1/chat/completions", c.GetEndpoint())
}

func (c *ChatInstance) GetCompletionPrompt(messages []globals.Message) string {
	result := ""
	for _, message := range messages {
		result += fmt.Sprintf("%s: %s\n", message.Role, message.Content)
	}
	return result
}

func (c *ChatInstance) GetLatestPrompt(props *adaptercommon.ChatProps) string {
	if len(props.Message) == 0 {
		return ""
	}

	return props.Message[len(props.Message)-1].Content
}

func (c *ChatInstance) GetChatBody(props *adaptercommon.ChatProps, stream bool) interface{} {
	if props.Model == globals.GPT3TurboInstruct {
		// for completions
		return CompletionRequest{
			Model:    props.Model,
			Prompt:   c.GetCompletionPrompt(props.Message),
			MaxToken: props.MaxTokens,
			Stream:   stream,
		}
	}

	messages := formatMessages(props)

	// o1, o3 compatibility
	isNewModel := len(props.Model) >= 2 && (props.Model[:2] == "o1" || props.Model[:2] == "o3")

	var temperature *float32
	if isNewModel {
		temp := float32(1.0)
		temperature = &temp
	} else {
		temperature = props.Temperature
	}

	request := ChatRequest{
		Model:            props.Model,
		Messages:         messages,
		Stream:           stream,
		PresencePenalty:  props.PresencePenalty,
		FrequencyPenalty: props.FrequencyPenalty,
		Temperature:      temperature,
		TopP:             props.TopP,
		Tools:            props.Tools,
		ToolChoice:       props.ToolChoice,
	}

	if isNewModel {
		request.MaxCompletionTokens = props.MaxTokens
	} else {
		request.MaxToken = props.MaxTokens
	}
	return request
}

// CreateChatRequest is the native http request body for openai
func (c *ChatInstance) CreateChatRequest(props *adaptercommon.ChatProps) (string, error) {
	if globals.IsOpenAIDalleModel(props.Model) {
		return c.CreateImage(props)
	}

	res, err := utils.Post(
		c.GetChatEndpoint(props),
		c.GetHeader(),
		c.GetChatBody(props, false),
		props.Proxy,
	)

	if err != nil || res == nil {
		return "", fmt.Errorf("openai error: %s", err.Error())
	}

	data := utils.MapToStruct[ChatResponse](res)
	if data == nil {
		return "", fmt.Errorf("openai error: cannot parse response")
	} else if data.Error.Message != "" {
		return "", fmt.Errorf("openai error: %s", data.Error.Message)
	}
	return data.Choices[0].Message.Content, nil
}

func hideRequestId(message string) string {
	// xxx (request id: 2024020311120561344953f0xfh0TX)

	exp := regexp.MustCompile(`\(request id: [a-zA-Z0-9]+\)`)
	return exp.ReplaceAllString(message, "")
}

// CreateStreamChatRequest is the stream response body for openai
func (c *ChatInstance) CreateStreamChatRequest(props *adaptercommon.ChatProps, callback globals.Hook) error {
	if globals.IsOpenAIDalleModel(props.Model) {
		if url, err := c.CreateImage(props); err != nil {
			return err
		} else {
			return callback(&globals.Chunk{
				Content: url,
			})
		}
	}

	isCompletionType := props.Model == globals.GPT3TurboInstruct

	ticks := 0
	err := utils.EventScanner(&utils.EventScannerProps{
		Method:  "POST",
		Uri:     c.GetChatEndpoint(props),
		Headers: c.GetHeader(),
		Body:    c.GetChatBody(props, true),
		Callback: func(data string) error {
			ticks += 1

			partial, err := c.ProcessLine(data, isCompletionType)
			if err != nil {
				return err
			}
			return callback(partial)
		},
	}, props.Proxy)

	if err != nil {
		if form := processChatErrorResponse(err.Body); form != nil {
			if form.Error.Type == "" && form.Error.Message == "" {
				return errors.New(utils.ToMarkdownCode("json", err.Body))
			}

			msg := fmt.Sprintf("%s (type: %s)", form.Error.Message, form.Error.Type)
			return errors.New(hideRequestId(msg))
		}
		return err.Error
	}

	if ticks == 0 {
		return errors.New("no response")
	}

	return nil
}
