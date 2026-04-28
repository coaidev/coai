import { isUrl } from "@/utils/base.ts";
import { Model } from "@/api/types.tsx";

import {
  Claude,
  Gemini,
  Gemma,
  OpenAI,
  Spark,
  Qwen,
  Baichuan,
  ByteDance,
  Meta,
  Copilot,
  Hunyuan,
  Midjourney,
  Stability,
  Moonshot,
  LLaVA,
  DeepSeek,
  Grok,
  Minimax,
  Mistral,
  Dalle,
  Rwkv,
  Cloudflare,
  Cohere,
  Fireworks,
  Groq,
  OpenRouter,
  Perplexity,
  GithubCopilot,
  Suno,
  Qingyan,
  IconAvatarProps,
  Azure,
  Coze,
  Dify,
  Qiniu
} from "@lobehub/icons";
import React from "react";
import { cn } from "@/components/ui/lib/utils.ts";

type ModelAvatarProps = {
  model: Model | { id: string; name: string; avatar?: string };
  className?: string;
  size?: number;
};

const builtinAvatars: Record<string, React.ExoticComponent<IconAvatarProps>> = {
  openai: OpenAI.Avatar,
  "gpt-3.5": OpenAI.Avatar,
  "gpt-4": OpenAI.Avatar,
  dalle: Dalle.Avatar,
  "dall-e": Dalle.Avatar,

  azure: Azure.Avatar,

  claude: Claude.Avatar,
  anthropic: Claude.Avatar,

  gemini: Gemini.Avatar,
  palm: Gemma.Avatar,
  gemma: Gemma.Avatar,
  "chat-bison": Gemma.Avatar, // "chat-bision" is a typo, but we need to keep it for compatibility
  google: Gemini.Avatar,

  glm: Qingyan.Avatar,
  zhipu: Qingyan.Avatar,

  spark: Spark.Avatar,

  tongyi: Qwen.Avatar,
  qwen: Qwen.Avatar,

  baichuan: Baichuan.Avatar,

  byte: ByteDance.Avatar,
  bytedance: ByteDance.Avatar,
  skylark: ByteDance.Avatar,

  meta: Meta.Avatar,
  llama: Meta.Avatar,

  bing: Copilot.Avatar,

  hunyuan: Hunyuan.Avatar,

  midjourney: Midjourney.Avatar,

  stability: Stability.Avatar,
  "stable-diffusion": Stability.Avatar,
  stablediffusion: Stability.Avatar,
  sd: Stability.Avatar,

  moonshot: Moonshot.Avatar,
  kimi: Moonshot.Avatar,

  llava: LLaVA.Avatar,

  deepseek: DeepSeek.Avatar,
  "deep-seek": DeepSeek.Avatar,

  coze: Coze.Avatar,

  dify: Dify.Avatar,

  grok: Grok.Avatar,
  minimax: Minimax.Avatar,
  abab: Minimax.Avatar,
  mistral: Mistral.Avatar,

  rwkv: Rwkv.Avatar,

  cf: Cloudflare.Combine,
  cloudflare: Cloudflare.Combine,

  command: Cohere.Avatar,
  cohere: Cohere.Avatar,

  firework: Fireworks.Avatar,

  groq: Groq.Avatar,

  qiniu: Qiniu.Avatar,

  router: OpenRouter.Avatar,

  perplexity: Perplexity.Avatar,

  copilot: GithubCopilot.Avatar,

  suno: Suno.Avatar,
};

function getAvatarType(id: string): string | undefined {
  if (id.includes("gpt-3.5")) return "gpt3";
  if (id.includes("gpt-4") || id.includes("o1")) return "gpt4";
}

function ModelAvatar({ model, className, size }: ModelAvatarProps) {
  const avatarSize = size ?? 42;
  
  if (isUrl(model.avatar ?? "")) {
    return (
      <div 
        style={{
          width: avatarSize,
          height: avatarSize,
          minWidth: avatarSize,
          minHeight: avatarSize
        }}
        className={cn(
          "relative flex items-center justify-center overflow-hidden",
          // using scale to make the avatar smaller
          className?.includes("h-4") && "scale-[0.85]",
          className
        )}
      >
        <img
          src={model.avatar}
          alt={model.name}
          className="rounded-full object-cover w-full h-full"
          style={{
            transform: className?.includes("h-4") ? "scale(1.15)" : "none"
          }}
        />
      </div>
    );
  }

  // if key is include, return value (reactelement)
  const id = model.id.toLowerCase();
  const key = Object.keys(builtinAvatars).find((key) => id.includes(key));
  const Avatar = key ? builtinAvatars[key] : OpenAI.Avatar;

  return (
    <Avatar
      size={avatarSize}
      className={className}
      // @ts-ignore
      type={getAvatarType(id)}
    />
  );
}

export default ModelAvatar;

export type ChannelTypeAvatarProps = {
  type: string;
  size?: number;
  className?: string;
};

export function ChannelTypeAvatar({
  type,
  size,
  className,
}: ChannelTypeAvatarProps) {
  const key = Object.keys(builtinAvatars).find((key) => type.includes(key));
  const Avatar = key ? builtinAvatars[key] : OpenAI.Avatar;

  return <Avatar size={size ?? 42} className={className} />;
}
