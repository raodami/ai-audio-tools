export default function Dashboard() {
  return (
    <div style={{ padding: '40px', maxWidth: '800px', margin: '0 auto' }}>
      <h1 style={{ marginBottom: 24 }}>Dashboard</h1>
      <div style={{ background: 'rgba(255,255,255,0.05)', borderRadius: 12, padding: 24, marginBottom: 24 }}>
        <h3 style={{ marginBottom: 16 }}>Upload Audio</h3>
        <input type="file" accept="audio/*" style={{ marginBottom: 16, display: 'block' }} />
        <button style={{ background: '#533afd', color: 'white', border: 'none', padding: '10px 20px', borderRadius: 8, cursor: 'pointer' }}>
          Process
        </button>
      </div>
      <div style={{ background: 'rgba(255,255,255,0.05)', borderRadius: 12, padding: 24 }}>
        <h3 style={{ marginBottom: 16 }}>Usage</h3>
        <p style={{ color: '#8899a6' }}>Free: 30 min/month</p>
        <div style={{ background: 'rgba(255,255,255,0.1)', borderRadius: 8, height: 8, marginTop: 8 }}>
          <div style={{ background: '#533afd', width: '0%', height: '100%', borderRadius: 8 }} />
        </div>
      </div>
    </div>
  );
}
