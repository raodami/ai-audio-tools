# PRD: AI Audio Tools SaaS

## 产品定位
AI-powered audio processing platform for content creators, podcasters, and musicians.

## 目标用户画像
- **播客主**：需要去噪、转文字、生成摘要
- **音乐制作人**：人声分离、分轨处理
- **视频创作者**：音频转录+字幕生成
- **外语学习者**：TTS 语音合成

## MVP 功能范围

### 核心功能
1. **Vocal Remover** — 人声分离（伴奏/ vocals）
2. **Audio Transcription** — 音频转文字（Deepgram API）
3. **Text Summary** — 转录文本 AI 摘要（DeepSeek API）
4. **Text-to-Speech** — 多语言 TTS（ElevenLabs API）

### 用户系统
- 匿名用户：每日 30 分钟免费额度（Cookie 计数）
- 注册用户：按订阅解锁更多额度
- 付费墙：超出额度触发 Stripe Checkout

### 技术架构
```
Go + SQLite (modernc) + DeepSeek + Deepgram + ElevenLabs
前端：Next.js 14 + Stripe 风格 UI
部署：Docker + Render Blueprint
```

## 定价策略
- **Free**: 30 min/月，基础功能
- **Pro**: $9.9/月 或 $59/年（前端显示 $4.9/月）
- **Team**: $29.9/月（多人协作）

## 成功指标
- MRR: $1,000 以内验证PMF
- 转化率: 免费→付费 > 3%
- 留存率: 30日 > 40%

## 风险
- API 成本波动（Deepgram ~$0.00014/sec, ElevenLabs ~$0.30/1k chars）
- 视频平台集成（YouTube/TikTok）— 二期功能

## 竞品分析
| 竞品 | 优势 | 劣势 |
|------|------|------|
| Clumi AI | +3140% 增长 | 仅 vocal remover |
| Descript | 全功能 | 功能臃肿，$16/月 |
| Lalal.ai | 专注分离 | 无 TTS/ASR |

## 差异化
- **All-in-one**: 分离 + 转录 + TTS 一站式
- **低价**: $9.9/月 vs Descript $16/月
- **API-first**: 支持开发者集成

## 下一步
写实施计划 → TDD 开发
