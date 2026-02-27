export const modelColorMapper: Record<string, string> = {
  // OpenAI & Azure OpenAI
  "gpt-3.5-turbo": "green-500",
  "gpt-3.5-turbo-16k": "green-600",
  "gpt-4": "purple-600",

  dalle: "green-600",
  "dall-e-2": "green-600",
  "dall-e-3": "purple-700",

  whisper: "gray-300",
  tts: "gray-300",
  openai: "gray-300",
  azure: "gray-300",

  // Anthropic Claude
  "claude-3": "orange-500",
  claude: "orange-400",
  anthropic: "orange-400",

  // Spark Desk
  "spark-desk": "blue-400",
  sparkdesk: "blue-400",

  // Moonshot
  moonshot: "black-500",
  kimi: "black-500",

  // Midjourney
  midjourney: "indigo-600",
  "mid-journey": "indigo-600",
  niji: "indigo-600",

  // Stable Diffusion
  "stable-diffusion": "gray-400",
  stablediffusion: "gray-400",
  stability: "gray-400",

  // Groq Cloud
  "llama2-70b-4096": "red-500",
  "mixtral-8x7b-32768": "red-500",
  "gemma-7b-it": "red-500",

  // Google Gemini & Gemma
  "chat-bison-001": "red-500",
  palm: "red-500",
  gemini: "red-500",
  gemma: "red-500",

  // DeepSeek
  deepseek: "blue-700",

  // New Bing
  bing: "blue-700",

  // ChatGLM
  zhipu: "lime-500",
  glm: "lime-500",

  // Tongyi Qwen
  qwen: "indigo-600",
  tongyi: "indigo-600",

  // Meta LLaMA
  llama: "sky-400",

  // Tencent Hunyuan
  hunyuan: "blue-500",

  // 360 GPT
  "360": "stone-500",

  // Baichuan AI
  baichuan: "orange-700",

  // ByteDance Skylark / Doubao / Coze
  skylark: "sky-300",
  doubao: "sky-300",
  coze: "sky-300",

  // Dify
  dify: "gray-300",

  // OpenRouter
  openrouter: "purple-600",

  // Novita AI
  novita: "violet-500",
};

const unknownColors = [
  "gray-700",
  "indigo-600",
  "green-500",
  "green-600",
  "purple-600",
  "purple-700",
  "orange-400",
  "blue-400",
  "red-500",
  "blue-700",
  "lime-500",
  "sky-400",
  "stone-500",
  "orange-700",
  "sky-300",
];

export function getUnknownModelColor(model: string): string {
  const char = model.length > 0 ? model[0] : "A";
  const code = char.charCodeAt(0);

  return unknownColors[code % unknownColors.length];
}

export function getModelColor(model: string): string {
  for (const key in modelColorMapper) {
    if (model.toLowerCase().includes(key)) {
      return modelColorMapper[key];
    }
  }

  return getUnknownModelColor(model);
}
