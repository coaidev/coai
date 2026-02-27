<div align="center">

![chatnio](/app/public/logo.png)

# [🥳 CoAI.Dev](https://coai.dev)

#### 🚀 下一代 AIGC 一站式商业解决方案
#### *“ CoAI.Dev > [Next Web](https://github.com/ChatGPTNextWeb/ChatGPT-Next-Web) + [One API](https://github.com/songquanpeng/one-api) ”*


[English](./README.md) · 简体中文 · [文档](https://coai.dev) · [Discord](https://discord.gg/rpzNSmqaF2) · [部署文档](https://coai.dev/docs/deploy)

[![CoAI.Dev: #1 Repo Of The Day](https://trendshift.io/api/badge/repositories/6369)](https://trendshift.io/repositories/6369)

<img alt="CoAI.Dev Preview" src="./screenshot/coai.png" width="100%" style="border-radius: 8px">

</div>

## 📝 功能
1. 🤖️ **丰富模型支持**: 多模型服务商支持 (OpenAI / Anthropic / Gemini / Midjourney 等十余种格式兼容 & 私有化 LLM 支持)
2. 🤯 **美观 UI 设计**: UI 兼容 PC / Pad / 移动三端，遵循 [Shadcn UI](https://ui.shadcn.com) & [Tremor Charts](https://blocks.tremor.so) 设计规范，丰富美观的界面设计和后台仪表盘
3. 🎃 **完整 Markdown 支持**: 支持 **LaTeX 公式** / **Mermaid 思维导图** / 表格渲染 / 代码高亮 / 图表绘制 / 进度条等进阶 Markdown 语法支持
4. 👀 **多主题支持**: 支持多种主题切换，包含亮色主题的**明亮模式**和暗色主题的**深色模式**。 👉 [自定义配色](https://github.com/coaidev/coai/blob/main/app/src/assets/globals.less)
5. 📚 **国际化支持**: 支持国际化，支持多语言切换 🇨🇳 🇺🇸 🇯🇵 🇷🇺 👉 欢迎贡献翻译 [Pull Request](https://github.com/coaidev/coai/pulls) 
6. 🎨 **文生图支持**: 支持多种文生图模型: **OpenAI DALL-E**✅ & **Midjourney** (支持 **U/V/R** 操作)✅ & Stable Diffusion✅ 等
7. 📡 **强大对话同步**: **用户 0 成本对话跨端同步支持**，支持**对话分享** (支持链接分享 & 保存为图片 & 分享管理), **无需 WebDav / WebRTC 等依赖和复杂学习成本**
8. 🎈 **模型市场 & 预设系统**: 支持后台可自定义的模型市场, 可提供模型介绍、标签等参数, 站长可根据情况自定义模型简介。同时支持预设系统，包含 **自定义预设** 和 **云端同步** 功能。
9. 📖 **丰富文件解析**: **开箱即用**, 支持**所有模型**的文件解析 (PDF / Docx / Pptx / Excel / 图片等格式解析), **支持更多云端图片存储方案** (S3 / R2 / MinIO 等), **支持 OCR 图片识别** 👉 详情参见项目 [CoAI.Dev Blob Service](https://github.com/coaidev/blob-service) (支持 Vercel / Docker 一键部署)
10. 🌏 **全模型联网搜索**: 基于 [SearXNG](https://github.com/searxng/searxng) 开源引擎, 支持 Google / Bing / DuckDuckGo / Yahoo / WikiPedia / Arxiv / Qwant 等丰富搜索引擎搜索, 支持安全搜索模式, 内容截断, 图片代理, 测试搜索可用性等功能。
11. 💕 **渐进式 Web 应用 (PWA)**: 支持 PWA 应用 & 支持桌面端 (桌面端基于 [Tauri](https://github.com/tauri-apps/tauri))
12. 🤩 **齐全后台管理**: 支持美观丰富的仪表盘, 公告 & 通知管理, 用户管理, 订阅管理, 礼品码 & 兑换码管理, 价格设定, 订阅设定, 自定义模型市场, 自定义站点名称 & Logo, SMTP 发件设置等功能
13. 🤑 **多种计费方式**: 支持 💴 **订阅制** 和 💴 **弹性计费** 两种计费方式, 弹性计费支持 次数计费 / Token 计费 / 不计费 / 可匿名调用 和 **最小请求点数** 检测等强大功能
14. 🎉 **创新模型缓存**: 支持开启模型缓存：即同一个请求入参 Hash 下, 如果之前已请求过, 将直接返回缓存结果 (击中缓存将不计费), 减少请求次数。可自行自定义是否缓存的模型、缓存时间、多种缓存结果数等高级缓存设置
15. 🥪 **附加功能** (停止支持): 🍎 **AI 项目生成器功能** / 📂 **批量文章生成功能** / 🥪 **AI 卡片功能** (已废弃)
16. 😎 **优秀渠道管理**: 自写优秀渠道算法, 支持⚡ **多渠道管理**, 支持🥳**优先级**设置渠道的调用顺序, 支持🥳**权重**设置同一优先级下的渠道均衡负载分配概率, 支持🥳**用户分组**, 🥳**失败自动重试**, 🥳**模型重定向**, 🥳**内置上游隐藏**, 🥳**渠道状态管理**等强大**企业级功能**
17. ⭐ **OpenAI API 分发 & 中转系统**: 支持以 **OpenAI API** 标准格式调用各种大模型, 集成强大的渠道管理功能, 仅需部署一个站点即可实现同时发展 B/C 端业务💖
18. 👌 **快速同步上游**: 渠道设置、模型市场、价格设定等设置都可快速同步上游站点，以此基础修改自己的站点配置，快速搭建自己的站点，省时省力，一键同步，快速上线
19. 👋 **SEO 优化**: 支持 SEO 优化，支持自定义站点名称、站点 Logo 等 SEO 优化设设置使搜索引擎更快的爬取，你的站点与众不同👋
20. 🎫 **多种兑换码体系**: 支持多种兑换码体系，支持礼品码和兑换码，支持批量生成，礼品码适合宣传分发，兑换码适合发卡销售，礼品码一个类型的多个码一个用户仅能兑换一个码，在宣传中一定程度上减少一个用户兑换多次的情况😀
21. 🥰 **商用友好协议**: 采用 **Apache-2.0** 开源协议, 商用二开 & 分发友好 (也请遵守 Apache-2.0 协议的规定, 请勿用于违法用途)

> ### ✨ CoAI.Dev 商业版
> ![商业版预览](./screenshot/coai-pro.png)
>
> - ✅ 美观商业级 UI, 漂亮的前端界面与后台管理
> - ✅ 支持 TTS & STT, 插件市场, RAG 知识库等丰富功能与模块
> - ✅ 更多支付供应商, 更多计费模式和高级订单管理
> - ✅ 支持更多鉴权方式，包括短信登录、OAuth 登录等
> - ✅ 支持模型监控，渠道健康检测，故障告警自动渠道切换
> - ✅ 支持多租户 API Key 分发系统, 企业级令牌权限管理与访问者限制
> - ✅ 支持安全审核, 日志记录, 模型限速, API Gateway 等高级功能
> - ✅ 支持推广奖励，专业数据统计，用户画像分析等商业分析能力
> - ✅ 支持Discord/Telegram/飞书等机器人对接集成能力 (扩展模块)
> - ...
>
> [👉 了解更多](https://www.coai.dev/docs/contact)


## 🔨 支持模型
1. OpenAI & Azure OpenAI *(✅ Vision ✅ Function Calling)*
2. Anthropic Claude *(✅ Vision ✅ Function Calling)*
3. Google Gemini & PaLM2 *(✅ Vision)*
4. Midjourney *(✅ Mode Toggling ✅ U/V/R Actions)*
5. 讯飞星火 SparkDesk *(✅ Vision ✅ Function Calling)*
6. 智谱清言 ChatGLM *(✅ Vision)*
7. 通义千问 Tongyi Qwen
8. 腾讯混元 Tencent Hunyuan
9. 百川大模型 Baichuan AI
10. 月之暗面 Moonshot AI (👉 OpenAI)
11. 深度求索 DeepSeek AI (👉 OpenAI)
12. 字节云雀 ByteDance Skylark *(✅ Function Calling)*
13. Groq Cloud AI
14. OpenRouter (👉 OpenAI)
15. 360 GPT
16. LocalAI / Ollama (👉 OpenAI)

## 👻 中转 OpenAI 兼容 API
   - [x] Chat Completions _(/v1/chat/completions)_
   - [x] Image Generation _(/v1/images)_
   - [x] Model List _(/v1/models)_
   - [x] Dashboard Billing _(/v1/billing)_


## 📦 部署方式
> [!TIP]
> **部署成功后, 管理员账号为 `root`, 密码默认为 `chatnio123456`**

### ✨ Zeabur (一键部署)
[![Deploy on Zeabur](https://zeabur.com/button.svg)](https://zeabur.com/templates/M86XJI)

> Zeabur 提供一定的免费额度, 可以使用非付费区域进行一键部署，同时也支持计划订阅和弹性计费等方式弹性扩展。
> 1. 点击 `Deploy` 进行部署, 并输入你希望绑定的域名，等待部署完成。
> 2. 部署完成后, 请访问你的域名, 并使用用户名 `root` 密码 `chatnio123456` 登录后台管理，请按照提示在 chatnio 后台及时修改密码。

### 🐳 宝塔面板 (一键部署)

1. 安装宝塔面板，前往 [宝塔面板官网](https://www.bt.cn/new/download.html) 进行安装，选择正式版脚本安装。
2. 登录面板，点击左侧 **Docker** 进入 Docker 管理。
3. 如提示未安装 Docker / Docker Compose， 可根据上方引导安装。
4. 安装完成后，进入 **应用商城**，搜索 `CoAI` 并点击 **安装**。
5. 配置应用基本信息，如您的域名，端口等配置，并点击 **确认** (可使用默认配置)。
6. 首次安装可能需要等待 1-2 分钟完成数据库初始化。如遇到问题，请查看面板运行日志进行排查。
7. 访问您配置的域名或服务器 `http://[ip]:[port]`，使用用户名 `root` 和密码 `chatnio123456` 登录后台管理。

### 阿里云计算巢 (一键部署)
 [![Deploy on AlibabaCloud ComputeNest](https://service-info-public.oss-cn-hangzhou.aliyuncs.com/computenest.svg)](https://computenest.console.aliyun.com/service/instance/create/default?type=user&ServiceName=CoAI社区版)
1. 访问计算巢CoAI[部署链接](https://computenest.console.aliyun.com/service/instance/create/cn-hangzhou?type=user&ServiceName=CoAI社区版)，按提示填写部署参数
2. 选择付费类型，填写实例参数与网络参数，点击 **确认订单**
3. 确认部署参数并查看预估价格后，点击立即创建，等待服务实例部署完成
4. 点击左侧 **服务实例** 等待服务实例部署完成后，点击实例ID进入到详情界面
5. 点击详情界面**立即使用**中的链接，可进入CoAI社区版界面。默认用户名为`root`，密码`为chatnio123456` 登录后台管理。
6. 更多操作详情与付费信息，参见：[服务详情](https://computenest.console.aliyun.com/service/detail/cn-hangzhou/service-bfbf676bd89d434691fc/1?type=user&isRecommend=true)

### ClawCloud Run (一键部署)
[![Run on ClawCloud](https://raw.githubusercontent.com/ClawCloud/Run-Template/refs/heads/main/Run-on-ClawCloud.svg)](https://template.run.claw.cloud/?referralCode=WZVLKDWQ4CEA&openapp=system-fastdeploy%3FtemplateName%3Dchatnio)

### ⚡ Docker Compose 安装 (推荐)
> [!NOTE]
> 运行成功后, 宿主机映射地址为 `http://localhost:8000`

 ```shell
 git clone --depth=1 --branch=main --single-branch https://github.com/coaidev/coai.git
 cd chatnio
 docker-compose up -d # 运行服务
# 如需使用 stable 版本, 请使用 docker-compose -f docker-compose.stable.yaml up -d 替代
# 如需使用 watchtower 自动更新, 请使用 docker-compose -f docker-compose.watch.yaml up -d 替代
```
   
版本更新（_开启 Watchtower 自动更新的情况下, 无需手动更新_）：
```shell
docker-compose down 
docker-compose pull
docker-compose up -d
```

> - MySQL 数据库挂载目录项目 ~/**db**
> - Redis 数据库挂载目录项目 ~/**redis**
> - 配置文件挂载目录项目 ~/**config**

### ⚡ Docker 安装 (轻量运行时, 常用于外置 _MYSQL/RDS_ 服务)
> [!NOTE]
> 运行成功后, 宿主机地址为 `http://localhost:8094`。
> 
> 如需使用 stable 版本, 请使用 `programzmh/chatnio:stable` 替代 `programzmh/chatnio:latest`  

```shell
docker run -d --name chatnio \
   --network host \
   -v ~/config:/config \
   -v ~/logs:/logs \
   -v ~/storage:/storage \
   -e MYSQL_HOST=localhost \
   -e MYSQL_PORT=3306 \
   -e MYSQL_DB=chatnio \
   -e MYSQL_USER=root \
   -e MYSQL_PASSWORD=chatnio123456 \
   -e REDIS_HOST=localhost \
   -e REDIS_PORT=6379 \
   -e SECRET=secret \
   -e SERVE_STATIC=true \
   programzmh/chatnio:latest
```

> - *--network host* 指使用宿主机网络, 使 Docker 容器使用宿主机的网络, 可自行修改
> - SECRET: JWT 密钥, 自行生成随机字符串修改
> - SERVE_STATIC: 是否启用静态文件服务 (正常情况下不需要更改此项, 详见下方常见问题解答)
> - *-v ~/config:/config* 挂载配置文件, *-v ~/logs:/logs* 挂载日志文件的宿主机目录, *-v ~/storage:/storage* 挂载附加功能的生成文件
> - 需配置 MySQL 和 Redis 服务, 请自行参考上方信息修改环境变量
 
 版本更新 （_开启 Watchtower 后无需手动更新, 执行后按照上述步骤重新运行即可_）：
 ```shell
docker stop chatnio
docker rm chatnio
docker pull programzmh/chatnio:latest
```

### ⚒ 编译安装
> [!NOTE]
> 部署成功后, 默认端口为 **8094**, 访问地址为 `http://localhost:8094`
> 
> Config 配置项 (~/config/**config.yaml**) 可以使用环境变量进行覆盖, 如 `MYSQL_HOST` 环境变量可覆盖 `mysql.host` 配置项

```shell
git clone https://github.com/coaidev/coai.git
cd chatnio

cd app
npm install -g pnpm
pnpm install
pnpm build

cd ..
go build -o chatnio

# e.g. using nohup (you can also use systemd or other service manager)
nohup ./chatnio > output.log & # using nohup to run in background
```

## ❓ 常见问题 Q&A
1. **为什么我部署后的站点可以访问页面, 可以登录注册, 但是无法使用聊天 (一直在转圈)？**
   - 聊天等此类功能通过 websocket 进行通信, 请确保你的服务支持 websocket。 (Tip: 中转通过 Http 实现, 无需 websocket 支持)
   - 如果你使用了 Nginx, Apache 等反向代理, 请确保已配置 websocket 支持。
   - 如果使用了端口映射, 端口转发, CDN, API Gateway 等服务, 请确保你的服务支持并开启 websocket。
2. **我配置的 Midjourney Proxy 格式的渠道一直转圈或报错 `please provide available notify url`**
   - 若为转圈，请确保你的 Midjourney Proxy 服务已正常运行, 并且已配置正确的上游地址。
   - **Midjourney 要填渠道类型要用 Midjourney 而不是 OpenAI (不知道为什么很多人填成了 OpenAI 类型格式然后过来反馈为什么empty response, mj-chat 类除外)**
   - 排查完这些问题后, 请查看你的系统设置中的**后端域名**是否已经配置并配置正确。如果不配置, 将导致 Midjourney Proxy 服务无法正常回调。
3. **此项目有什么外部依赖？**
   - MySQL: 存储用户信息, 对话记录, 管理员信息等持久化数据。
   - Redis: 存储用户快速鉴权信息, IP 速率限制, 订阅配额, 邮箱验证码等数据。
   - 环境未配置好的情况下, 会导致服务无法正常运行, 请确保你的 MySQL 和 Redis 服务已正常运行 (Docker 部署, 编译部署需自行搭建外部服务)。
4. **我的机器为 ARM 架构, 该项目支持 ARM 架构吗？**
   - 支持。CoAI.Dev 项目使用 BuildX 构建多架构镜像, 你可以直接使用 docker-compose 或 docker 运行, 无需额外配置。
   - 如果你使用编译安装, 直接在 ARM 机器上编译即可, 无需欸外配置。如果你使用 x86 机器编译, 请使用 `GOARCH=arm64 go build -o chatnio` 进行交叉编译并上传至 ARM 机器上运行。
5. **如何修改 Root 默认密码？**
   - 请点击右上角头像或侧边栏底部用户框进入后台管理, 点击系统设置下常规设置操作栏的 修改 Root 密码 进行修改。或者选择在 用户管理 中选定 root 用户进行修改密码操作。
6. **系统设置中的后端域名是什么？**
   - 后端域名是指后端 API 服务的地址, 默认为你访问站点后加 `/api` 的地址, 如 `https://example.com/api` 。
   - 如果设置为非 *SERVE_STATIC* 模式, 开启前后端分离部署, 请将后端域名设置为你的后端 API 服务地址, 如 `https://api.example.com`。
   - 后端域名此处用于 Midjourney Proxy 服务的后端回调地址, 如无需使用 Midjourney Proxy 服务, 请忽略此设置。
7. **如何配置支付方式？**
   - CoAI.Dev 开源版支持发卡模式, 设置系统设置中的购买链接为你的发卡地址即可。卡密可通过用户管理中兑换码管理中批量生成。
8. **礼品码和兑换码有什么区别？**
   - 礼品码一种类型只能一个用户只能绑定一次, 而非 aff code, 发福利等方式可使用礼品码, 可在头像下拉菜单中的礼品码中兑换。
   - 兑换码一种类型可以多个用户绑定, 可作为正常购买和发卡使用, 可在用户管理中的兑换码管理中批量生成, 在头像下拉菜单的点数（菜单第一个）内输入兑换码进行兑换。
   - 一个例子：比如我发了一个类型为 *新年快乐* 的福利, 此时推荐使用礼品码, 假设发放 100 个 66 点数, 如果为兑换码, 手快的一个用户就批量把所有兑换码的 6600 点数都用完了, 而礼品码则可以保证每个用户只能使用一次 (获得 66 点数)。
   - 而搭建发卡的时, 如果用礼品码, 因为一个类型只能兑换一次, 购买多个礼品码会导致兑换失败, 而兑换码则可以在此场景下使用。
9. **该项目支持 Vercel 部署吗？**
   - CoAI.Dev 本身并不支持 Vercel 部署, 但是你可以使用前后端分离模式,  Vercel 部署前端部分, 后端部分使用 Docker 部署或编译部署。
10. **前后端分离部署模式是什么？**
    - 正常情况下, 前后端在同一服务内, 后端地址为 `/api`。前后端分离部署指前端和后端分别部署在不同的服务上, 前端服务为静态文件服务, 后端服务为 API 服务。
      - 举个例子, 前端使用 Nginx (或 Vercel 等) 部署, 部署的域名为 `https://www.chatnio.net`。
      - 后端使用 Docker 部署, 部署的域名为 `https://api.chatnio.net`。
    - 此种部署方式需自行打包前端, 配置环境变量 `VITE_BACKEND_ENDPOINT` 为你的后端地址, 如 `https://api.chatnio.net`。
    - 配置后端环境变量的 `SERVE_STATIC=false` 使后端服务不提供静态文件服务。
11. **弹性计费和订阅详解**
    - 弹性计费, 即 `点数`, 其图标类似于**云**, 模型计费通用方式, 为了防止虚假汇率, 写死 10 点数 = 1 元, 汇率可以在计费规则中的 **应用内置模板** 中自定义汇率。
    - 订阅, 即订阅计划, 为固定价格计费方式按次配额, 订阅计费扣取点数 (举例: 如果站点的用户想订阅 32 元的计划, 则需要保证点数大于等于 320 点数)
    - 订阅是 Item 的组合, 每个 Item 都可设置涵盖的模型, 订阅配额 (-1 为无限使用), 名称, ID (用于区分不同的 Item), 图标等。可在后台的订阅管理中进行操作, 是否开启订阅, 订阅价格等, 修改每个订阅等级的 Item, 以及支持直接导入其他订阅等级的 Item。
    - 订阅支持分层并写死为三个等级。 等级分别为: _普通用户 (0)_, _基础版订阅 (1)_, _标准版订阅 (2)_, _专业版订阅 (3)_, 订阅等级即为用户分组, 可在渠道管理中进行高级设置, 选择勾选可使用此模型的用户分组。
    - 订阅配额设置, 可在订阅管理中进行操作, 是否支持中转 API (默认关闭)
12. **可请求最小点数检测 `user quota is not enough` 详解**
    - 为防止站点用户滥用站点模型, 当请求点数低于最小请求点数时将返回点数不足的错误信息, 大于等于最小请求点数时将正常请求。
    - 模型的最小可请求点数规则: 
        - 不计费模型无限制
        - 次数计费模型最小点数为该模型的 1 次请求点数 (e.g. 若一个模型的单次请求点数为 0.1 点数, 则最小请求点数为 0.1 点数)
        - Token 弹性计费模型为 1K 输入 Tokens 价格 + 1K 输出 Tokens 价格 (e.g. 若一个模型的 1K 输入 Tokens 价格为 0.05 点数, 1K 输出 Tokens 价格 0.1 点数, 则最小请求点数为 0.15 点数)
13. **为何我的 GPT-4-All 等逆向模型无法使用上传文件中的图片?**
    - 上传模型图片为 Base64 格式, 如果逆向不支持 Base64 格式, 请使用 URL 直链而非上传文件做法。
14. **如何开始域名严格跨域检测?**
    - 正常情况下，后端对所有域名开放跨域。如果非特殊需求，无需开启严格跨域检测。
    - 如果需要开启严格跨域检测，可以在后端环境变量中 并配置 `ALLOW_ORIGINS`, 如 `ALLOW_ORIGINS=chatnio.net,chatnio.app` （不需要加协议前缀, www 解析无需手动添加, 后端将自动识别并允许跨域）, 这样就会支持严格跨域检测 (如 *http://www.chatnio.app*, *https://chatnio.net* 等将会被允许, 其他域名将会被拒绝)。
    - 即使在开启严格跨域检测的情况下, /v1 接口会被仍然允许所有域的跨域请求, 以保证中转 API 的正常使用。
15. **模型映射功能是如何使用的？**
    - 渠道内的模型映射格式为 `[from]>[to]`, 多个映射之间换行, **from** 为请求的模型, **to** 为真实向上游发送的模型并且需要上游真实支持
    - 如: 我有一个逆向渠道, 填写 `gpt-4-all>gpt-4`, 则我的用户请求 **gpt-4-all** 模型到该渠道时, 后端则会模型映射至 **gpt-4** 向该渠道请求 **gpt-4**, 此时该渠道支持 2 个模型, **gpt-4** 和 **gpt-4-all** (本质上都为 **gpt-4**)
    - 如果我不想让我的这个逆向渠道影响到 **gpt-4** 的渠道组, 可以加前缀 `!gpt-4-all>gpt-4`, 该渠道 **gpt-4** 则会被忽略, 此时该渠道将只支持 1 个模型, **gpt-4-all** (但本质上为 **gpt-4**)

## 📦 技术栈
- 🥗 前端: React + Redux + Radix UI + Tailwind CSS
- 🍎 后端: Golang + Gin + Redis + MySQL
- 🍒 应用技术: PWA + WebSocket

## 🤯 为什么写此项目 & 项目优势
我们发现，市面上的 AIGC 商业站点，大多数都是偏向于前端轻量部署的项目，有精美的 UI 界面设计，
比如 [Next Chat](https://github.com/ChatGPTNextWeb/ChatGPT-Next-Web) 的二开商业版本，
由于其偏向个人私有化的设计，在二开商业化时有一定的局限性，呈现出一些问题，比如：
  - **对话同步难**, 比如需要 WebDav 等服务，用户学习成本高，跨端实时同步困难。
  - **计费不够完善**, 比如只支持弹性计费或只支持订阅制，无法满足不同用户的需求。
  - **文件解析不便捷**, 比如只支持先在图床上传图片，返回站点后再在输入框中输入 URL 直链，无内置文件解析功能。
  - **不支持对话 URL 分享**, 比如只支持对话截图分享，无法支持对话 URL 分享 (或仅支持 ShareGPT 等工具，无法对站点起到推广作用)。
  - **渠道管理不够强大**, 比如后台仅支持 OpenAI 格式渠道，兼容其他格式渠道困难。且只能填入一个渠道，无法支持多渠道管理。
  - **不支持 API 调用**, 比如只支持用户界面调用，无法支持 API 中转和管理。

另一种是偏向于 API 分发的站点，有强大的分发系统，比如基于 [One API](https://github.com/songquanpeng/one-api) 等项目，
这类项目虽然支持强大的 API 中转和管理，但是缺少界面设计，且缺少一些 C 端功能，比如：
  - **用户界面不够丰富**, 比如只支持 API 调用，不内置用户界面聊天。用户界面聊天需要自行复制密钥并前往其他站点才能使用，这对于普通用户来说，学习成本较高。
  - **没有订阅制**, 比如只支持弹性计费，缺少对 C 端用户的计费设计，无法满足用户的不同需求，对于无基础的用户来说，成本感知不够友好。
  - **C 端功能不够丰富**, 比如只支持 API 调用，不支持对话同步，不支持对话分享，不支持文件解析等功能。
  - **均衡负载不够强大**, 开源版不支持**权重**参数, 无法实现同优先级的渠道均衡负载分配概率 ([New API](https://github.com/Calcium-Ion/new-api) 也解决了此痛点, UI 也更美观)。

因此，我们希望能够将这两种项目的优势结合起来，做出一个既有强大的 API 分发系统，又有丰富的用户界面设计的项目，
这样既能满足 C 端用户的需求，又能发展 B 端业务，提高用户体验，降低用户学习成本，提高用户粘性。

于是，**CoAI.Dev** 应运而生，我们希望能够做出一个既有强大的 API 分发系统，又有丰富的用户界面设计的项目，成为下一代开源 AIGC 项目的商业一站式解决方案。


## ❤ 捐助
如果您觉得这个项目对您有所帮助, 您可以点个 Star 支持一下！
