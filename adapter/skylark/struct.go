package skylark

import (
	factory "chat/adapter/common"
	"chat/globals"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
)

type ChatInstance struct {
	Instance *arkruntime.Client
}

func NewChatInstance(endpoint, apiKey string) *ChatInstance {
	//https://ark.cn-beijing.volces.com/api/v3
	instance := arkruntime.NewClientWithApiKey(apiKey, arkruntime.WithBaseUrl(endpoint))
	return &ChatInstance{
		Instance: instance,
	}
}

func NewChatInstanceFromConfig(conf globals.ChannelConfig) factory.Factory {
	params := conf.SplitRandomSecret(1)

	return NewChatInstance(
		conf.GetEndpoint(),
		params[0],
	)
}
