const features = [
  { icon: '🎙️', title: 'Speech-to-Text', desc: 'Accurate transcription with Deepgram AI. Support for 30+ languages and audio formats.' },
  { icon: '📝', title: 'AI Summarization', desc: 'Extract key points from your audio with DeepSeek intelligence. Get concise summaries instantly.' },
  { icon: '🔊', title: 'Text-to-Speech', desc: 'Natural-sounding voice synthesis with ElevenLabs. Multiple voices and languages available.' },
];

export default function Features() {
  return (
    <section className="features" id="features">
      <div className="container">
        <h2>Powerful Features</h2>
        <div className="features-grid">
          {features.map((f, i) => (
            <div className="feature-card" key={i}>
              <div className="feature-icon">{f.icon}</div>
              <h3>{f.title}</h3>
              <p>{f.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
