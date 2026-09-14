'use client';
import Link from 'next/link';
import { useState, useEffect } from 'react';

export default function Navbar() {
  const [user, setUser] = useState<any>(null);
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (token) {
      fetch('/api/auth/me', { headers: { Authorization: `Bearer ${token}` } })
        .then(r => r.json())
        .then(data => setUser(data))
        .catch(() => {});
    }

    const handleScroll = () => setScrolled(window.scrollY > 20);
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  return (
    <nav style={{
      position: 'fixed',
      top: 0,
      left: 0,
      right: 0,
      zIndex: 100,
      padding: scrolled ? '12px 24px' : '20px 24px',
      background: scrolled ? 'rgba(6,27,49,0.95)' : 'transparent',
      backdropFilter: scrolled ? 'blur(10px)' : 'none',
      borderBottom: scrolled ? '1px solid rgba(255,255,255,0.1)' : 'none',
      transition: 'all 0.3s ease',
    }}>
      <div style={{ maxWidth: 1200, margin: '0 auto', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Link href="/" style={{ fontSize: 24, fontWeight: 700, color: '#533afd', textDecoration: 'none' }}>
          🎵 AudioAI
        </Link>
        
        <div style={{ display: 'flex', alignItems: 'center', gap: 32 }}>
          <Link href="#features" style={{ color: '#8899a6', textDecoration: 'none', fontSize: 15 }}>Features</Link>
          <Link href="#pricing" style={{ color: '#8899a6', textDecoration: 'none', fontSize: 15 }}>Pricing</Link>
          <Link href="#about" style={{ color: '#8899a6', textDecoration: 'none', fontSize: 15 }}>About</Link>
          
          {user ? (
            <>
              <Link href="/dashboard" style={{ color: '#fff', textDecoration: 'none', fontSize: 15 }}>Dashboard</Link>
              <span style={{ color: '#4ade80', fontSize: 12, fontWeight: 600 }}>{user.is_pro ? 'PRO' : 'FREE'}</span>
            </>
          ) : (
            <>
              <Link href="/login" style={{ color: '#fff', textDecoration: 'none', fontSize: 15, padding: '8px 16px', borderRadius: 8, border: '1px solid rgba(255,255,255,0.2)' }}>
                Sign In
              </Link>
              <Link href="/login?mode=register" style={{ background: 'linear-gradient(135deg, #533afd, #7c5cfc)', color: '#fff', textDecoration: 'none', fontSize: 15, padding: '10px 20px', borderRadius: 8, fontWeight: 600 }}>
                Get Started
              </Link>
            </>
          )}
        </div>
      </div>
    </nav>
  );
}
