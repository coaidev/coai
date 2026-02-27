package globals

import (
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const ChatMaxThread = 5
const AnonymousMaxThread = 1

var HttpMaxTimeout = 30 * time.Minute

var AllowedOrigins []string

var DebugMode bool
var NotifyUrl = ""
var ArticlePermissionGroup []string
var GenerationPermissionGroup []string
var CacheAcceptedModels []string
var CacheAcceptedExpire int64
var CacheAcceptedSize int64
var AcceptImageStore bool
var AcceptPromptStore bool
var CloseRegistration bool
var CloseRelay bool

var SearchEndpoint string
var SearchCrop bool
var SearchCropLength int
var SearchEngines string    // e.g. "google,bing"
var SearchImageProxy string // e.g. "True", "False"
var SearchSafeSearch int    // e.g. 0: None, 1: Moderation, 2: Strict

func OriginIsAllowed(uri string) bool {
	if len(AllowedOrigins) == 0 {
		// if allowed origins is empty, allow all origins
		return true
	}

	instance, _ := url.Parse(uri)
	if instance == nil {
		return false
	}

	if instance.Hostname() == "localhost" || instance.Scheme == "file" {
		return true
	}

	if strings.HasPrefix(instance.Host, "www.") {
		instance.Host = instance.Host[4:]
	}

	return in(instance.Host, AllowedOrigins)
}

func OriginIsOpen(c *gin.Context) bool {
	return strings.HasPrefix(c.Request.URL.Path, "/v1") || strings.HasPrefix(c.Request.URL.Path, "/dashboard") || strings.HasPrefix(c.Request.URL.Path, "/mj")
}

const (
	GPT3Turbo                      = "gpt-3.5-turbo"
	GPT3TurboInstruct              = "gpt-3.5-turbo-instruct"
	GPT3Turbo0613                  = "gpt-3.5-turbo-0613"
	GPT3Turbo0301                  = "gpt-3.5-turbo-0301"
	GPT3Turbo1106                  = "gpt-3.5-turbo-1106"
	GPT3Turbo0125                  = "gpt-3.5-turbo-0125"
	GPT3Turbo16k                   = "gpt-3.5-turbo-16k"
	GPT3Turbo16k0613               = "gpt-3.5-turbo-16k-0613"
	GPT3Turbo16k0301               = "gpt-3.5-turbo-16k-0301"
	GPT4                           = "gpt-4"
	GPT4All                        = "gpt-4-all"
	GPT4Vision                     = "gpt-4-v"
	GPT4Dalle                      = "gpt-4-dalle"
	GPT40314                       = "gpt-4-0314"
	GPT40613                       = "gpt-4-0613"
	GPT41106Preview                = "gpt-4-1106-preview"
	GPT40125Preview                = "gpt-4-0125-preview"
	GPT4TurboPreview               = "gpt-4-turbo-preview"
	GPT4VisionPreview              = "gpt-4-vision-preview"
	GPT4Turbo                      = "gpt-4-turbo"
	GPT4Turbo20240409              = "gpt-4-turbo-2024-04-09"
	GPT41106VisionPreview          = "gpt-4-1106-vision-preview"
	GPT432k                        = "gpt-4-32k"
	GPT432k0314                    = "gpt-4-32k-0314"
	GPT432k0613                    = "gpt-4-32k-0613"
	GPT4O                          = "gpt-4o"
	GPT4O0513                      = "gpt-4o-2024-05-13"
	GPT4OMini                      = "gpt-4o-mini"
	GPT4OMiniSearchPreview         = "gpt-4o-mini-search-preview"
	GPT4O1120                      = "gpt-4o-2024-11-20"
	GPT4O0806                      = "gpt-4o-2024-08-06"
	GPT4OLatest                    = "gpt-4o-latest"
	GPTImage1                      = "gpt-image-1"
	Sora2                          = "sora-2"
	GPT41Mini                      = "gpt-4.1-mini"
	O11217                         = "o1-2024-12-17"
	O1Preview                      = "o1-preview"
	O1Mini                         = "o1-mini"
	GPT5                           = "gpt-5"
	GPT5Mini                       = "gpt-5-mini"
	GPT5Nano                       = "gpt-5-nano"
	Dalle                          = "dalle"
	Dalle2                         = "dall-e-2"
	Dalle3                         = "dall-e-3"
	Claude1                        = "claude-1"
	Claude1100k                    = "claude-1.3"
	Claude2                        = "claude-1-100k"
	Claude2100k                    = "claude-2"
	Claude2200k                    = "claude-2.1"
	Claude3                        = "claude-3"
	ClaudeSlack                    = "claude-slack"
	Claude3Opus                    = "claude-3-opus-20240229"
	Claude3Sonnet                  = "claude-3-sonnet-20240229"
	Claude3Haiku                   = "claude-3-haiku-20240307"
	Claude35Sonnet0620             = "claude-3-5-sonnet-20240620"
	Claude35Sonnet                 = "claude-3-5-sonnet-20241022"
	Claude35Haiku                  = "claude-3-5-haiku-20241022"
	Claude37Sonnet                 = "claude-3-7-sonnet-20250219"
	Claude4Sonnet                  = "claude-sonnet-4-20250514"
	Claude4Opus                    = "claude-opus-4-20250514"
	SparkDesk                      = "spark-desk-v1.5"
	SparkDeskV2                    = "spark-desk-v2"
	SparkDeskV3                    = "spark-desk-v3"
	SparkDeskV35                   = "spark-desk-v3.5"
	ChatBison001                   = "chat-bison-001"
	GeminiPro                      = "gemini-pro"
	GeminiProVision                = "gemini-pro-vision"
	Gemini15ProLatest              = "gemini-1.5-pro-latest"
	Gemini15FlashLatest            = "gemini-1.5-flash-latest"
	Gemini20ProExp                 = "gemini-2.0-pro-exp-02-05"
	Gemini20Flash                  = "gemini-2.0-flash"
	Gemini20FlashExp               = "gemini-2.0-flash-exp"
	Gemini20Flash001               = "gemini-2.0-flash-001"
	Gemini20FlashThinkingExp       = "gemini-2.0-flash-thinking-exp-01-21"
	Gemini20FlashLitePreview       = "gemini-2.0-flash-lite-preview-02-05"
	Gemini20FlashThinkingExp1219   = "gemini-2.0-flash-thinking-exp-1219"
	Gemini25ProExp0325             = "gemini-2.5-pro-exp-03-25"
	GeminiExp1206                  = "gemini-exp-1206"
	Gemini25FlashImage             = "gemini-2.5-flash-image"
	Gemini25FlashImagePreview      = "gemini-2.5-flash-image-preview"
	GoogleImagen30Generate002      = "imagen-3.0-generate-002"
	GoogleImagen40Generate001      = "imagen-4.0-generate-001"
	GoogleImagen40UltraGenerate001 = "imagen-4.0-ultra-generate-001"
	GoogleImagen40FastGenerate001  = "imagen-4.0-fast-generate-001"
	BingCreative                   = "bing-creative"
	BingBalanced                   = "bing-balanced"
	BingPrecise                    = "bing-precise"
	ZhiPuChatGLM4                  = "glm-4"
	ZhiPuChatGLM4Vision            = "glm-4v"
	ZhiPuChatGLM3Turbo             = "glm-3-turbo"
	ZhiPuChatGLMTurbo              = "zhipu-chatglm-turbo"
	ZhiPuChatGLMPro                = "zhipu-chatglm-pro"
	ZhiPuChatGLMStd                = "zhipu-chatglm-std"
	ZhiPuChatGLMLite               = "zhipu-chatglm-lite"
	QwenTurbo                      = "qwen-turbo"
	QwenPlus                       = "qwen-plus"
	QwenTurboNet                   = "qwen-turbo-net"
	QwenPlusNet                    = "qwen-plus-net"
	Midjourney                     = "midjourney"
	MidjourneyFast                 = "midjourney-fast"
	MidjourneyTurbo                = "midjourney-turbo"
	Hunyuan                        = "hunyuan"
	GPT360V9                       = "360-gpt-v9"
	Baichuan53B                    = "baichuan-53b"
	SkylarkLite                    = "skylark-lite-public"
	SkylarkPlus                    = "skylark-plus-public"
	SkylarkPro                     = "skylark-pro-public"
	SkylarkChat                    = "skylark-chat"
	DeepseekV3                     = "deepseek-chat"
	DeepseekR1                     = "deepseek-reasoner"
)

var OpenAIDalleModels = []string{
	Dalle, Dalle2, Dalle3, GPTImage1,
}

var GoogleImagenModels = []string{
	GoogleImagen30Generate002, GoogleImagen40Generate001, GoogleImagen40UltraGenerate001, GoogleImagen40FastGenerate001,
}

var VisionModels = []string{
	GPT4VisionPreview, GPT41106VisionPreview, GPT4Turbo, GPT4Turbo20240409, GPT4O, GPT4O0513, GPT4O0806, GPT4O1120, GPT4OMini, GPT4OMiniSearchPreview, GPT4OLatest, GPT41Mini, O11217, GPT5, GPT5Mini, GPT5Nano, // openai
	GeminiProVision, Gemini15ProLatest, Gemini15FlashLatest, Gemini20ProExp, Gemini20Flash, Gemini20Flash001, Gemini20FlashLitePreview, GeminiExp1206, Gemini20FlashThinkingExp, Gemini20FlashThinkingExp1219, Gemini25ProExp0325, Gemini25FlashImage, Gemini25FlashImagePreview, // gemini
	Claude3, Claude3Opus, Claude3Sonnet, Claude3Haiku, Claude35Sonnet0620, Claude35Sonnet, Claude4Sonnet, Claude4Opus, // anthropic
	ZhiPuChatGLM4Vision, // chatglm
}

var VisionSkipModels = []string{
	GPT4TurboPreview,
}

var VideoModels = []string{
	Sora2,
}

func in(value string, slice []string) bool {
	for _, item := range slice {
		if item == value || strings.Contains(value, item) {
			return true
		}
	}
	return false
}

func IsOpenAIDalleModel(model string) bool {
	// using image generation api if model is in dalle models
	return in(model, OpenAIDalleModels) && !strings.Contains(model, "gpt-4-dalle")
}

func IsGoogleImagenModel(model string) bool {
	// using image generation api if model is in imagen models
	return in(model, GoogleImagenModels)
}

func IsVisionModel(model string) bool {
	return in(model, VisionModels) && !in(model, VisionSkipModels)
}
