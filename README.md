# AI Audio Tools

AI-powered audio processing platform with transcription, summarization, and text-to-speech capabilities.

## Tech Stack
- **Backend**: Go + Gin + SQLite (modernc)
- **Frontend**: Next.js 14 + TypeScript + Stripe UI
- **AI Services**: Deepgram (ASR), DeepSeek (LLM), ElevenLabs (TTS)

## Quick Start

### 1. Clone and Setup
```bash
cd D:/ai-audio-tools
cp .env.example .env
# Add your API keys to .env
```

### 2. Backend
```bash
go mod tidy
go test ./...
go build -o bin/server ./cmd/server
./bin/server
```

### 3. Frontend
```bash
cd web
npm install
npm run dev
```

Open http://localhost:3000

### 4. Docker
```bash
docker-compose up -d
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/auth/register | Register new user |
| POST | /api/auth/login | Login |
| GET | /api/auth/me | Get current user |
| GET | /api/user/usage | Check usage quota |
| POST | /api/audio/upload | Upload audio file |
| GET | /api/audio/{id} | Get job status/result |

## Pricing
- **Free**: 30 min/month
- **Pro**: $9.9/mo (500 min)
- **Team**: $29.9/mo (unlimited)

## Environment Variables
- `DEEPGRAM_API_KEY` — Speech transcription
- `DEEPSEEK_API_KEY` — AI summarization
- `ELEVENLABS_API_KEY` — Text-to-speech
- `STRIPE_SECRET_KEY` — Payment processing
- `ADMIN_KEY` — Admin authentication

## Project Structure
```
ai-audio-tools/
├── cmd/server/           # Main entry point
├── internal/
│   ├── api/             # HTTP handlers & routes
│   ├── llm/             # DeepSeek client
│   ├── store/           # SQLite database
│   └── webhook/         # Stripe webhooks
├── web/                 # Next.js frontend
├── docs/PRD.md          # Product requirements
├── Dockerfile
└── render.yaml          # Render deployment config
```
