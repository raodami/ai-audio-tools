const features = [
  {
    icon: '🎙️',
    title: 'Speech-to-Text',
    desc: 'Accurate transcription with Deepgram AI. Support for 30+ languages and all major audio formats.',
  },
  {
    icon: '📝',
    title: 'AI Summarization',
    desc: 'Extract key points from your audio with DeepSeek intelligence. Get concise summaries instantly.',
  },
  {
    icon: '🔊',
    title: 'Text-to-Speech',
    desc: 'Natural-sounding voice synthesis with ElevenLabs. Multiple voices and languages available.',
  },
  {
    icon: '📊',
    title: 'Analytics',
    desc: 'Track your audio processing metrics. Understand usage patterns and optimize your workflow.',
  },
  {
    icon: '🔒',
    title: 'Secure & Private',
    desc: 'Enterprise-grade security. Your data is encrypted at rest and in transit. GDPR compliant.',
  },
  {
    icon: '⚡',
    title: 'Fast Processing',
    desc: 'Process audio files in real-time. Parallel processing for batch operations.',
  },
];

export default function Features() {
  return (
    <section id="features" style={{ padding: '80px 24px', maxWidth: 1200, margin: '0 auto' }}>
      <h2 style={{ textAlign: 'center', fontSize: 40, fontWeight: 700, marginBottom: 64 }}>
        Powerful Features
      </h2>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 24 }}>
        {features.map((f, i) => (
          <div
            key={i}
            style={{
              background: 'rgba(255,255,255,0.03)',
              border: '1px solid rgba(255,255,255,0.08)',
              borderRadius: 16,
              padding: 32,
              transition: 'all 0.3s ease',
            }}
          >
            <div style={{ fontSize: 40, marginBottom: 20 }}>{f.icon}</div>
            <h3 style={{ fontSize: 18, fontWeight: 600, marginBottom: 12 }}>{f.title}</h3>
            <p style={{ color: '#8899a6', fontSize: 14, lineHeight: 1.7 }}>{f.desc}</p>
          </div>
        ))}
      </div>
    </section>
  );
}
