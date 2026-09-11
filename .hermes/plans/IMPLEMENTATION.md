# AI Audio Tools - Implementation Plan

## Project Overview
- **Path**: `D:/ai-audio-tools`
- **Tech Stack**: Go (Gin) + SQLite + DeepSeek/Deepgram/ElevenLabs + Next.js 14
- **MVP Duration**: 2 weeks
- **Architecture**: See `references/product-architecture-template.md`

---

## Task 1: Project Setup (TDD RED)
**Goal**: Initialize Go module, create directory structure, write failing test
**Files to create**:
- `cmd/server/main.go` — server entry point
- `internal/store/store.go` — SQLite initialization
- `internal/llm/llm.go` — DeepSeek client interface
- `internal/llm/deepgram.go` — Deepgram client interface
- `internal/llm/elevenlabs.go` — ElevenLabs client interface
- `internal/api/router.go` — Gin router setup
- `go.mod`, `go.sum`

**Test file**:
- `internal/store/store_test.go` — fails (no implementation)

**Commands**:
```bash
cd D:/ai-audio-tools
go mod init ai-audio-tools
mkdir -p cmd/server internal/{llm,store,api,webhook}
```

**Acceptance Criteria**:
- [ ] `go test ./...` shows RED (expected failures)
- [ ] Directory structure matches template

---

## Task 2: SQLite Store + Migration (TDD GREEN)
**Goal**: Implement user table, free quota tracking via SQLite
**Implementation**:
```go
// internal/store/store.go
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE,
    created_at INTEGER NOT NULL,
    usage_minutes INTEGER DEFAULT 0,
    is_pro INTEGER DEFAULT 0
);

CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT,
    ip_address TEXT,
    cookie_hash TEXT,
    usage_minutes INTEGER DEFAULT 0,
    created_at INTEGER NOT NULL
);
```

**Tests**:
```go
// internal/store/store_test.go
func TestCreateUser(t *testing.T) { /* ... */ }
func TestIncrementUsage(t *testing.T) { /* ... */ }
func TestGetAnonymousUser(t *testing.T) { /* ... */ }
```

**Acceptance Criteria**:
- [ ] All tests pass (GREEN)
- [ ] In-memory test for migration edge cases
- [ ] PRAGMA table_info checks for safe ALTER

---

## Task 3: Deepgram ASR Integration (TDD GREEN)
**Goal**: Transcribe audio file → text using Deepgram API
**Implementation**:
```go
// internal/llm/deepgram.go
type TranscriptionResponse struct {
    Text string `json:"text"`
    Duration float64 `json:"duration"`
}

func (c *DeepgramClient) Transcribe(audio []byte) (*TranscriptionResponse, error) {
    // Upload to Deepgram, poll result
}
```

**Tests**:
```go
func TestTranscribe_Success(t *testing.T) { /* mock HTTP server */ }
func TestTranscribe_APIError(t *testing.T) { /* returns error */ }
```

**Acceptance Criteria**:
- [ ] Tests mock Deepgram API with httptest
- [ ] Handles 429 rate limit with retry

---

## Task 4: DeepSeek LLM Integration (TDD GREEN)
**Goal**: Summarize transcription, extract key points
**Implementation**:
```go
// internal/llm/llm.go
func (c *DeepSeekClient) Summarize(text string) (string, error) {
    // Call DeepSeek API with system prompt
}
```

**Tests**:
```go
func TestSummarize_Success(t *testing.T) { /* mock DeepSeek */ }
```

**Acceptance Criteria**:
- [ ] Mock HTTP server for DeepSeek API
- [ ] Context window handling (>4k tokens chunking)

---

## Task 5: ElevenLabs TTS Integration (TDD GREEN)
**Goal**: Text-to-speech synthesis
**Implementation**:
```go
// internal/llm/elevenlabs.go
func (c *ElevenLabsClient) Synthesize(text string, voiceID string) ([]byte, error) {
    // Call ElevenLabs API, return MP3 bytes
}
```

**Tests**:
```go
func TestSynthesize_Success(t *testing.T) { /* mock API */ }
```

**Acceptance Criteria**:
- [ ] Returns MP3 bytes for frontend playback
- [ ] Voice ID configurable via env

---

## Task 6: API Routes - User Management (TDD GREEN)
**Goal**: Register, login, quota check endpoints
**Routes**:
- `POST /api/auth/register` — email + password → JWT
- `POST /api/auth/login` — email + password → JWT
- `GET /api/user/usage` — current usage + quota
- `GET /api/auth/me` — current user info

**Tests**:
```go
func TestRegister_Success(t *testing.T) { /* 201, returns token */ }
func TestRegister_DuplicateEmail(t *testing.T) { /* 409 */ }
func TestUsageCheck_FreeQuota(t *testing.T) { /* checks cookie/session */ }
```

**Acceptance Criteria**:
- [ ] bcrypt password hashing (not SHA-256!)
- [ ] JWT expires in 30 days
- [ ] Anonymous user tracked via cookie

---

## Task 7: API Routes - Audio Processing (TDD GREEN)
**Goal**: Upload audio, transcribe, summarize, TTS
**Routes**:
- `POST /api/audio/upload` — multipart upload → returns job_id
- `GET /api/audio/{id}/result` — polling for transcription
- `POST /api/audio/{id}/summarize` — trigger summary generation
- `POST /api/audio/{id}/tts` — generate speech

**Tests**:
```go
func TestUpload_Success(t *testing.T) { /* 202, returns job_id */ }
func TestUpload_FileTooLarge(t *testing.T) { /* 413 */ }
func TestResult_NotReady(t *testing.T) { /* 202 pending */ }
```

**Acceptance Criteria**:
- [ ] Job status: pending → processing → completed/failed
- [ ] 50MB file size limit
- [ ] Supported formats: mp3, wav, m4a

---

## Task 8: Stripe Integration (TDD GREEN)
**Goal**: Checkout session, webhook handling
**Implementation**:
```go
// internal/webhook/stripe.go
func HandleCheckoutEvent(session *checkout.Session) error {
    // Activate pro subscription in DB
}

func HandleSubscriptionDeleted(sub *subscription.Subscription) error {
    // Downgrade to free tier
}
```

**Tests**:
```go
func TestHandleCheckout_Completed(t *testing.T) { /* verifies DB update */ }
func TestWebhook_VerifySignature(t *testing.T) { /* 400 on bad sig */ }
```

**Acceptance Criteria**:
- [ ] Stripe signature verification
- [ ] idempotent webhook handling
- [ ] Error logging for failed webhooks

---

## Task 9: Frontend - Landing Page (Next.js 14)
**Goal**: Stripe-style landing page with pricing
**Components**:
- `app/page.js` — Hero + Features + Pricing
- `components/ThemeProvider.js` — dark/light mode
- `components/Toast.js` — notification system

**Acceptance Criteria**:
- [ ] Source Sans 3 font, #533afd primary color
- [ ] Pricing cards: Free / Pro $9.9/mo / Team $29.9/mo
- [ ] CTA buttons → login/register

---

## Task 10: Frontend - Dashboard (Next.js 14)
**Goal**: User dashboard for audio processing
**Pages**:
- `app/dashboard/page.js` — upload form + results
- `app/dashboard/usage/page.js` — quota tracker
- `app/dashboard/history/page.js` — past sessions

**Acceptance Criteria**:
- [ ] File upload with drag-drop
- [ ] Progress indicator for async jobs
- [ ] Usage bar: X min used / 30 min free

---

## Task 11: Auth & Quota Enforcement (Integration)
**Goal**: Protect routes, enforce free limits
**Implementation**:
```go
// middleware/quota.go
func RequireQuota(next http.Handler) http.Handler {
    // Check usage against limit
    // Return 402 if exceeded
}
```

**Tests**:
```go
func TestQuotaExceeded_Returns402(t *testing.T) { /* mock exceeded usage */ }
```

**Acceptance Criteria**:
- [ ] Free users blocked after 30 min
- [ ] Pro users unlimited
- [ ] 402 response triggers Stripe checkout

---

## Task 12: Deployment - Docker + Render
**Goal**: Production-ready deployment
**Files**:
- `Dockerfile` — multi-stage build
- `render.yaml` — Render Blueprint
- `.env.example` — environment variables

**Acceptance Criteria**:
- [ ] `docker build -t ai-audio-tools .` succeeds
- [ ] Render auto-deploys on git push
- [ ] Health check: GET /health → 200

---

## Pitfalls & Solutions
1. **Windows OOM during build**: `GOFLAGS="-p=1"` for serial compilation
2. **SQLite old DB add column**: `PRAGMA table_info` + conditional ALTER
3. **Stripe v81 API**: Package-level functions `checkoutsession.New`
4. **Deepgram async**: Poll endpoint or use streaming
5. **ElevenLabs rate limit**: Queue with retry logic

## Verification Checklist
- [ ] All unit tests green
- [ ] Integration test with mock APIs
- [ ] Stripe webhook test mode works
- [ ] Frontend loads without CORS errors
- [ ] Docker compose up passes health checks
