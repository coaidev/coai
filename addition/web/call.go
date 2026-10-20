package web

import (
	"chat/globals"
	"chat/manager/conversation"
	"chat/utils"
	"context"
	"fmt"
	"time"
)

type Hook func(message []globals.Message, token int) (string, error)

func toWebSearchingMessage(ctx context.Context, message []globals.Message) []globals.Message {
	if len(message) == 0 {
		return message
	}
	data, err := GenerateSearchResultContext(ctx, message[len(message)-1].Content)
	if err != nil {
		return message
	}

	return utils.Insert(message, 0, globals.Message{
		Role: globals.System,
		Content: fmt.Sprintf("You will play the role of an AI Q&A assistant, where your knowledge base is not offline, but can be networked in real time, and you can provide real-time networked information with links to networked search sources."+
			"Current time: %s, Real-time internet search results: %s",
			time.Now().Format("2006-01-02 15:04:05"), data,
		),
	})
}

func ToChatSearched(instance *conversation.Conversation, restart bool) []globals.Message {
	return ToChatSearchedContext(context.Background(), instance, restart)
}

func ToChatSearchedContext(ctx context.Context, instance *conversation.Conversation, restart bool) []globals.Message {
	segment := conversation.CopyMessage(instance.GetChatMessage(restart))

	if instance.IsEnableWeb() {
		segment = toWebSearchingMessage(ctx, segment)
	}

	return segment
}

func ToSearched(enable bool, message []globals.Message) []globals.Message {
	return ToSearchedContext(context.Background(), enable, message)
}

func ToSearchedContext(ctx context.Context, enable bool, message []globals.Message) []globals.Message {
	if enable {
		return toWebSearchingMessage(ctx, message)
	}

	return message
}
