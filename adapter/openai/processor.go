package openai

import (
	adaptercommon "chat/adapter/common"
	"chat/globals"
	"chat/utils"
	"errors"
	"fmt"
	"regexp"
)

func formatMessages(props *adaptercommon.ChatProps) interface{} {
	if globals.IsVisionModel(props.Model) {
		return utils.Each[globals.Message, Message](props.Message, func(message globals.Message) Message {
			if message.Role == globals.User {
				content, urls := utils.ExtractImages(message.Content, true)
				images := utils.EachNotNil[string, MessageContent](urls, func(url string) *MessageContent {
					obj, err := utils.NewImage(url)
					props.Buffer.AddImage(obj)
					if err != nil {
						globals.Info(fmt.Sprintf("cannot process image: %s (source: %s)", err.Error(), utils.Extract(url, 24, "...")))
					}

					return &MessageContent{
						Type: "image_url",
						ImageUrl: &ImageUrl{
							Url: url,
						},
					}
				})

				return Message{
					Role: message.Role,
					Content: utils.Prepend(images, MessageContent{
						Type: "text",
						Text: &content,
					}),
					Name:         message.Name,
					FunctionCall: message.FunctionCall,
					ToolCalls:    message.ToolCalls,
					ToolCallId:   message.ToolCallId,
				}
			}

			return Message{
				Role: message.Role,
				Content: MessageContents{
					MessageContent{
						Type: "text",
						Text: &message.Content,
					},
				},
				Name:         message.Name,
				FunctionCall: message.FunctionCall,
				ToolCalls:    message.ToolCalls,
				ToolCallId:   message.ToolCallId,
			}
		})
	}

	return props.Message
}

func processChatResponse(data string) *ChatStreamResponse {
	return utils.UnmarshalForm[ChatStreamResponse](data)
}

func processCompletionResponse(data string) *CompletionResponse {
	return utils.UnmarshalForm[CompletionResponse](data)
}

func processChatErrorResponse(data string) *ChatStreamErrorResponse {
	return utils.UnmarshalForm[ChatStreamErrorResponse](data)
}

func getChoices(form *ChatStreamResponse) *globals.Chunk {
	if len(form.Choices) == 0 {
		return &globals.Chunk{Content: ""}
	}

	choice := form.Choices[0].Delta

	content := choice.Content

	// handle gemini multi-modal inline data (base64 image)
	if len(content) == 0 && choice.MultiModal != nil {
		for _, item := range *choice.MultiModal {
			if item.InlineData == nil || len(item.InlineData.Data) == 0 {
				continue
			}

			mime := item.InlineData.MimeType
			if len(mime) == 0 {
				mime = "image/png"
			}
			// use markdown image to let frontend直接渲染而不是显示base64文本
			content = fmt.Sprintf("![image](data:%s;base64,%s)", mime, item.InlineData.Data)
			break
		}
	}

	return &globals.Chunk{
		Content:      content,
		ToolCall:     choice.ToolCalls,
		FunctionCall: choice.FunctionCall,
	}
}

func getCompletionChoices(form *CompletionResponse) string {
	if len(form.Choices) == 0 {
		return ""
	}

	return form.Choices[0].Text
}

func getRobustnessResult(chunk string) string {
	exp := `\"content\":\"(.*?)\"`
	compile, err := regexp.Compile(exp)
	if err != nil {
		return ""
	}

	matches := compile.FindStringSubmatch(chunk)
	if len(matches) > 1 {
		return utils.ProcessRobustnessChar(matches[1])
	} else {
		return ""
	}
}

func (c *ChatInstance) ProcessLine(data string, isCompletionType bool) (*globals.Chunk, error) {
	if isCompletionType {
		// openai legacy support
		if completion := processCompletionResponse(data); completion != nil {
			return &globals.Chunk{
				Content: getCompletionChoices(completion),
			}, nil
		}

		globals.Warn(fmt.Sprintf("openai error: cannot parse completion response: %s", data))
		return &globals.Chunk{Content: ""}, errors.New("parser error: cannot parse completion response")
	}

	if form := processChatResponse(data); form != nil {
		return getChoices(form), nil
	}

	if form := processChatErrorResponse(data); form != nil {
		return &globals.Chunk{Content: ""}, fmt.Errorf("openai error: %s (type: %s)", form.Error.Message, form.Error.Type)
	}

	globals.Warn(fmt.Sprintf("openai error: cannot parse chat completion response: %s", data))
	return &globals.Chunk{Content: ""}, errors.New("parser error: cannot parse chat completion response")
}
