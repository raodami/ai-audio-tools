export default function Hero() {
  return (
    <section style={{ paddingTop: 160, paddingBottom: 80, textAlign: 'center', position: 'relative' }}>
      {/* Background glow */}
      <div style={{ position: 'absolute', top: '50%', left: '50%', transform: 'translate(-50%, -50%)', width: 600, height: 600, background: 'radial-gradient(circle, rgba(83,58,253,0.15) 0%, transparent 70%)', pointerEvents: 'none' }} />
      
      <div style={{ position: 'relative', zIndex: 1 }}>
        <div style={{ display: 'inline-block', padding: '8px 16px', background: 'rgba(83,58,253,0.2)', border: '1px solid rgba(83,58,253,0.3)', borderRadius: 20, fontSize: 13, color: '#a5b4fc', marginBottom: 24 }}>
          🚀 New: AI-Powered Audio Processing
        </div>
        <h1 style={{ fontSize: 64, fontWeight: 700, lineHeight: 1.1, marginBottom: 24, background: 'linear-gradient(135deg, #fff 0%, #a5b4fc 100%)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>
          Transform Audio<br />with AI
        </h1>
        <p style={{ fontSize: 20, color: '#8899a6', maxWidth: 600, margin: '0 auto 40px', lineHeight: 1.6 }}>
          Transcribe, summarize, and synthesize audio files instantly. Built for creators, developers, and teams.
        </p>
        <div style={{ display: 'flex', gap: 16, justifyContent: 'center' }}>
          <a href="/login?mode=register" style={{ background: 'linear-gradient(135deg, #533afd, #7c5cfc)', color: '#fff', padding: '16px 32px', borderRadius: 12, fontSize: 16, fontWeight: 600, textDecoration: 'none', display: 'inline-flex', alignItems: 'center', gap: 8 }}>
            Start Free Trial <span>→</span>
          </a>
          <a href="#features" style={{ background: 'rgba(255,255,255,0.05)', color: '#fff', padding: '16px 32px', borderRadius: 12, fontSize: 16, fontWeight: 600, textDecoration: 'none', border: '1px solid rgba(255,255,255,0.1)', display: 'inline-flex', alignItems: 'center', gap: 8 }}>
            Learn More
          </a>
        </div>
        
        {/* Stats */}
        <div style={{ display: 'flex', justifyContent: 'center', gap: 64, marginTop: 64 }}>
          {[
            { label: 'Audio Processed', value: '10M+' },
            { label: 'Active Users', value: '50K+' },
            { label: 'Uptime', value: '99.9%' },
          ].map((stat, i) => (
            <div key={i}>
              <div style={{ fontSize: 32, fontWeight: 700, color: '#fff' }}>{stat.value}</div>
              <div style={{ fontSize: 14, color: '#8899a6' }}>{stat.label}</div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
