package adapter

import (
	"chat/adapter/azure"
	"chat/adapter/baichuan"
	"chat/adapter/bing"
	"chat/adapter/claude"
	adaptercommon "chat/adapter/common"
	"chat/adapter/coze"
	"chat/adapter/dashscope"
	"chat/adapter/deepseek"
	"chat/adapter/dify"
	"chat/adapter/hunyuan"
	"chat/adapter/midjourney"
	"chat/adapter/openai"
	"chat/adapter/palm2"
	"chat/adapter/qiniu"
	"chat/adapter/skylark"
	"chat/adapter/slack"
	"chat/adapter/sparkdesk"
	"chat/adapter/zhinao"
	"chat/adapter/zhipuai"
	"chat/globals"
	"fmt"
)

var channelFactories = map[string]adaptercommon.FactoryCreator{
	globals.OpenAIChannelType:      openai.NewChatInstanceFromConfig,
	globals.AzureOpenAIChannelType: azure.NewChatInstanceFromConfig,
	globals.ClaudeChannelType:      claude.NewChatInstanceFromConfig,
	globals.SlackChannelType:       slack.NewChatInstanceFromConfig,
	globals.BingChannelType:        bing.NewChatInstanceFromConfig,
	globals.PalmChannelType:        palm2.NewChatInstanceFromConfig,
	globals.SparkdeskChannelType:   sparkdesk.NewChatInstanceFromConfig,
	globals.ChatGLMChannelType:     zhipuai.NewChatInstanceFromConfig,
	globals.QwenChannelType:        dashscope.NewChatInstanceFromConfig,
	globals.HunyuanChannelType:     hunyuan.NewChatInstanceFromConfig,
	globals.BaichuanChannelType:    baichuan.NewChatInstanceFromConfig,
	globals.SkylarkChannelType:     skylark.NewChatInstanceFromConfig,
	globals.ZhinaoChannelType:      zhinao.NewChatInstanceFromConfig,
	globals.MidjourneyChannelType:  midjourney.NewChatInstanceFromConfig,
	globals.DeepseekChannelType:    deepseek.NewChatInstanceFromConfig,
	globals.DifyChannelType:        dify.NewChatInstanceFromConfig,
	globals.CozeChannelType:        coze.NewChatInstanceFromConfig,
	globals.QiniuChannelType:       qiniu.NewChatInstanceFromConfig,

	globals.MoonshotChannelType: openai.NewChatInstanceFromConfig, // openai format
	globals.GroqChannelType:     openai.NewChatInstanceFromConfig, // openai format
}

func createChatRequest(conf globals.ChannelConfig, props *adaptercommon.ChatProps, hook globals.Hook) error {
	props.Model = conf.GetModelReflect(props.OriginalModel)
	props.Proxy = conf.GetProxy()

	factoryType := conf.GetType()
	if factory, ok := channelFactories[factoryType]; ok {
		return factory(conf).CreateStreamChatRequest(props, hook)
	}

	return fmt.Errorf("unknown channel type %s (channel #%d)", conf.GetType(), conf.GetId())
}

func createVideoRequest(conf globals.ChannelConfig, props *adaptercommon.VideoProps, hook globals.Hook) error {
	props.Model = conf.GetModelReflect(props.OriginalModel)
	props.Proxy = conf.GetProxy()

	factoryType := conf.GetType()
	if creator, ok := channelFactories[factoryType]; ok {
		inst := creator(conf)
		if v, ok := inst.(adaptercommon.VideoFactory); ok {
			return v.CreateVideoRequest(props, hook)
		}
		return fmt.Errorf("video request not supported by channel type %s (channel #%d)", conf.GetType(), conf.GetId())
	}

	return fmt.Errorf("unknown channel type %s (channel #%d)", conf.GetType(), conf.GetId())
}
