export default function Footer() {
  return (
    <footer style={{ borderTop: '1px solid rgba(255,255,255,0.1)', padding: '40px 24px', textAlign: 'center', color: '#8899a6' }}>
      <div style={{ maxWidth: 1200, margin: '0 auto' }}>
        <p style={{ marginBottom: 16 }}>© 2026 AudioAI. All rights reserved.</p>
        <div style={{ display: 'flex', justifyContent: 'center', gap: 24, fontSize: 14 }}>
          <a href="/privacy" style={{ color: '#8899a6', textDecoration: 'none' }}>Privacy Policy</a>
          <a href="/terms" style={{ color: '#8899a6', textDecoration: 'none' }}>Terms of Service</a>
          <a href="/contact" style={{ color: '#8899a6', textDecoration: 'none' }}>Contact</a>
        </div>
      </div>
    </footer>
  );
}
