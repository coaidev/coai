package globals

const (
	System    = "system"
	User      = "user"
	Assistant = "assistant"
	Tool      = "tool"
	Function  = "function"
)

const (
	OpenAIChannelType      = "openai"
	AzureOpenAIChannelType = "azure"
	ClaudeChannelType      = "claude"
	SlackChannelType       = "slack"
	SparkdeskChannelType   = "sparkdesk"
	ChatGLMChannelType     = "chatglm"
	HunyuanChannelType     = "hunyuan"
	QwenChannelType        = "qwen"
	ZhinaoChannelType      = "zhinao"
	BaichuanChannelType    = "baichuan"
	SkylarkChannelType     = "skylark"
	BingChannelType        = "bing"
	PalmChannelType        = "palm"
	MidjourneyChannelType  = "midjourney"
	MoonshotChannelType    = "moonshot"
	GroqChannelType        = "groq"
	DeepseekChannelType    = "deepseek"
	DifyChannelType        = "dify"
	CozeChannelType        = "coze"
	QiniuChannelType       = "qiniu"
)

const (
	NonBilling   = "non-billing"
	TimesBilling = "times-billing"
	TokenBilling = "token-billing"
)

const (
	AnonymousType = "anonymous"
	NormalType    = "normal"
	BasicType     = "basic"    // basic subscription
	StandardType  = "standard" // standard subscription
	ProType       = "pro"      // pro subscription
	AdminType     = "admin"
)

const (
	NoneProxyType = iota
	HttpProxyType
	HttpsProxyType
	Socks5ProxyType
)

const (
	WebTokenType = "web"
	ApiTokenType = "api"
	SystemToken  = "system"
)
