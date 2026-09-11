export default function Hero() {
  return (
    <section className="hero">
      <div className="container">
        <h1>Transform Audio<br />with AI Power</h1>
        <p>Transcribe, summarize, and synthesize audio files instantly. Built for creators, developers, and teams who need reliable audio processing.</p>
        <div style={{ display: 'flex', gap: 16, justifyContent: 'center' }}>
          <Link href="/dashboard" className="btn btn-primary" style={{ padding: '14px 32px', fontSize: '16px' }}>
            Start Free Trial
          </Link>
          <Link href="#features" className="btn btn-outline" style={{ padding: '14px 32px', fontSize: '16px' }}>
            Learn More
          </Link>
        </div>
      </div>
    </section>
  );
}
