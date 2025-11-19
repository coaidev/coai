# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Chat Nio is a Next Generation AIGC One-Stop Business Solution that combines a powerful API distribution system with a rich user interface. It serves as both a C-end (client) chat interface and B-end (business) API proxy/distribution platform, supporting multiple AI model providers (OpenAI, Anthropic, Gemini, Midjourney, and 15+ others).

**Tech Stack:**
- Frontend: React 18 + Redux Toolkit + Radix UI + Tailwind CSS + Vite
- Backend: Go 1.20 + Gin + MySQL + Redis
- Application: PWA + WebSocket for real-time communication

## Development Commands

### Frontend (app/)

```bash
cd app
pnpm install              # Install dependencies (use pnpm, not npm)
pnpm dev                  # Start dev server (Vite)
pnpm build                # Production build (includes TypeScript compilation)
pnpm fast-build           # Build without TypeScript check
pnpm lint                 # ESLint check
pnpm prettier             # Format code
```

### Backend (root)

```bash
# Development
go build -o chatnio       # Build binary
./chatnio                 # Run application (default port: 8094)

# With custom config
MYSQL_HOST=localhost REDIS_HOST=localhost ./chatnio

# Docker development
docker-compose up -d      # Start all services (app + MySQL + Redis)
docker-compose down       # Stop all services
docker-compose pull       # Update images
```

**Note:** The backend requires MySQL and Redis to be running. Configuration is read from `config/config.yaml` or environment variables (e.g., `MYSQL_HOST` overrides `mysql.host`).

## Architecture

### Backend Structure

The backend follows a modular, layered architecture:

**Core Layers:**
- `main.go` - Entry point, initializes managers and registers routes
- `adapter/` - **Adapter pattern** for AI provider integrations
- `channel/` - **Channel management system** with load balancing
- `manager/` - Business logic layer (chat, images, videos, usage tracking)
- `auth/` - Authentication & authorization
- `admin/` - Admin panel APIs
- `middleware/` - HTTP middleware (CORS, auth, rate limiting)
- `connection/` - Database connection management (MySQL/SQLite)
- `utils/` - Utilities including WebSocket handling, logging, etc.

**Key Architectural Patterns:**

1. **Adapter Pattern (adapter/):** Each AI provider has a dedicated adapter implementing the `FactoryCreator` interface. Adapters are registered in `adapter/adapter.go` channelFactories map:
   - `adapter/openai/` - OpenAI and OpenAI-compatible APIs
   - `adapter/claude/` - Anthropic Claude
   - `adapter/midjourney/` - Midjourney image generation
   - New adapters must implement `CreateStreamChatRequest()` and optionally `CreateVideoRequest()`

2. **Channel Management (channel/):** Sophisticated load balancing system supporting:
   - **Priority-based routing** - Channels are tried in priority order
   - **Weight-based distribution** - Load balancing probability for same-priority channels
   - **User grouping** - Different channel sets for different user tiers
   - **Auto-retry on failure** - Automatic failover to next available channel
   - **Model redirection** - Map requested models to available channel models
   - Channel state is managed by `Manager` in `channel/manager.go`

3. **Manager Layer (manager/):** Core business logic for:
   - `chat.go` - Chat session management
   - `chat_completions.go` - OpenAI-compatible completions endpoint
   - `images.go` - Image generation (DALL-E, Midjourney, etc.)
   - `videos.go` - Video generation (Sora, etc.)
   - `conversation/` - Conversation sync and sharing
   - `relay.go` - Request relaying to appropriate adapters

4. **WebSocket & Real-time Communication:**
   - `manager/connection.go` - WebSocket connection management for chat sessions
   - `manager/broadcast/` - Broadcasting system for real-time updates across sessions
   - `utils/websocket.go` - Low-level WebSocket utilities and connection monitoring

5. **Database Layer (connection/):** Database connection management
   - Supports both MySQL (production) and SQLite (fallback for development)
   - Auto-switches to SQLite if MySQL host is not configured
   - Connection pooling and automatic reconnection handled by `worker.go`

### Frontend Structure

React-based SPA with Redux state management:

```
app/src/
├── components/       # Reusable UI components (Radix UI + Tailwind)
├── routes/          # Page components
├── admin/           # Admin dashboard pages
├── store/           # Redux slices (auth, chat, settings)
├── api/             # API client functions
├── types/           # TypeScript type definitions
├── utils/           # Helper functions
├── dialogs/         # Modal dialogs
├── assets/          # Static assets (CSS, images, i18n)
├── router.tsx       # React Router configuration
└── App.tsx          # Root component
```

**State Management:** Redux Toolkit with slices for auth, chat state, user settings, and admin data.

**Styling:** Tailwind CSS + custom LESS (see `assets/globals.less` for theme customization).

## Configuration

**Backend Config (config/config.yaml):**
- MySQL connection (host, port, user, password, database)
- Redis connection
- JWT secret (`secret`)
- Server port (default: 8094)
- `serve_static` - Whether to serve frontend static files (set to `false` for separate frontend deployment)

**Frontend Config (app/src/conf/):**
- API endpoints
- Feature flags
- Site customization

Environment variables override config file values. Example: `MYSQL_HOST` overrides `mysql.host` in config.yaml.

## API Routes

Routes are registered in `main.go` under `/api` prefix (when `serve_static=true`):
- `/api/auth/*` - Authentication endpoints (auth package)
- `/api/admin/*` - Admin panel APIs (admin package)
- `/api/v1/*` - OpenAI-compatible API (adapter package)
- `/api/chat/*` - Chat management (manager package)
- `/api/conversation/*` - Conversation sync (conversation package)
- `/api/addition/*` - Additional features (addition package)

## Database Migrations

Database schema changes are handled by `migration/` package. Migrations run automatically on startup in `admin/InitInstance()`.

## Adding a New AI Provider

1. Create new package in `adapter/your_provider/`
2. Implement `FactoryCreator` interface with `CreateStreamChatRequest()`
3. Add to `channelFactories` map in `adapter/adapter.go`
4. Define channel type constant in `globals/channel.go`
5. Update channel configuration options in admin panel

## Deployment

Default credentials after deployment: username `root`, password `chatnio123456` (change immediately in admin panel).

Docker deployment includes automatic database initialization and volume mounts for:
- `~/db` - MySQL data
- `~/redis` - Redis data
- `~/config` - Configuration files
- `~/logs` - Application logs
- `~/storage` - File uploads and generated content
