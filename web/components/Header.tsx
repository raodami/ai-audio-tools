import Link from 'next/link';

export default function Header() {
  return (
    <header className="header">
      <div className="container">
        <Link href="/" className="logo">🎵 AudioAI</Link>
        <nav className="nav">
          <Link href="/#features">Features</Link>
          <Link href="/#pricing">Pricing</Link>
          <Link href="/dashboard" className="btn btn-outline" style={{marginLeft: 24}}>Dashboard</Link>
          <Link href="/login" className="btn btn-primary">Get Started</Link>
        </nav>
      </div>
    </header>
  );
}
