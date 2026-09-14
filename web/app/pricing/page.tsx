import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';

export default function Pricing() {
  const router = useRouter();
  const [loading, setLoading] = useState<string | null>(null);

  const plans = [
    {
      name: 'Free',
      price: '$0',
      period: '/month',
      features: ['30 min transcription', 'Basic summaries', '1 TTS voice', 'Email support'],
      cta: 'Get Started',
      action: () => router.push('/login'),
    },
    {
      name: 'Pro',
      price: '$9.9',
      period: '/month',
      features: ['500 min transcription', 'AI summaries', 'All TTS voices', 'Priority support', 'API access'],
      cta: 'Upgrade to Pro',
      action: () => handleUpgrade('pro'),
      popular: true,
    },
    {
      name: 'Team',
      price: '$29.9',
      period: '/month',
      features: ['Unlimited transcription', 'Advanced analytics', 'Team management', 'SSO', 'Dedicated support'],
      cta: 'Contact Sales',
      action: () => router.push('/contact'),
    },
  ];

  async function handleUpgrade(plan: string) {
    const token = localStorage.getItem('token');
    if (!token) {
      router.push('/login');
      return;
    }
    setLoading(plan);
    try {
      const res = await fetch('/api/user/subscribe', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ plan }),
      });
      const data = await res.json();
      if (data.checkout_url) {
        window.location.href = data.checkout_url;
      }
    } finally {
      setLoading(null);
    }
  }

  return (
    <div style={{ minHeight: '100vh', background: '#061b31' }}>
      <main style={{ maxWidth: 1000, margin: '0 auto', padding: '80px 24px' }}>
        <h1 style={{ textAlign: 'center', fontSize: 48, fontWeight: 700, marginBottom: 16, background: 'linear-gradient(135deg, #fff 0%, #a5b4fc 100%)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>
          Simple, Transparent Pricing
        </h1>
        <p style={{ textAlign: 'center', color: '#8899a6', fontSize: 18, marginBottom: 64 }}>
          Start free. Upgrade when you need more.
        </p>

        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 24 }}>
          {plans.map((plan, i) => (
            <div
              key={i}
              style={{
                background: 'rgba(255,255,255,0.05)',
                border: plan.popular ? '1px solid #533afd' : '1px solid rgba(255,255,255,0.1)',
                borderRadius: 16,
                padding: 32,
                position: 'relative',
                display: 'flex',
                flexDirection: 'column',
              }}
            >
              {plan.popular && (
                <div style={{ position: 'absolute', top: -12, left: '50%', transform: 'translateX(-50%)', background: 'linear-gradient(135deg, #533afd, #7c5cfc)', padding: '4px 16px', borderRadius: 20, fontSize: 12, fontWeight: 600 }}>
                  Most Popular
                </div>
              )}
              <h3 style={{ fontSize: 20, marginBottom: 8 }}>{plan.name}</h3>
              <div style={{ fontSize: 48, fontWeight: 700, marginBottom: 8 }}>
                {plan.price}<span style={{ fontSize: 16, color: '#8899a6', fontWeight: 400 }}>{plan.period}</span>
              </div>
              <ul style={{ listStyle: 'none', padding: 0, margin: '24px 0', flex: 1 }}>
                {plan.features.map((f, j) => (
                  <li key={j} style={{ padding: '8px 0', color: '#8899a6', fontSize: 14, display: 'flex', alignItems: 'center', gap: 8 }}>
                    <span style={{ color: '#533afd' }}>✓</span> {f}
                  </li>
                ))}
              </ul>
              <button
                onClick={plan.action}
                disabled={loading !== null}
                style={{
                  width: '100%',
                  padding: '14px 24px',
                  background: plan.popular ? 'linear-gradient(135deg, #533afd, #7c5cfc)' : 'rgba(255,255,255,0.1)',
                  border: 'none',
                  borderRadius: 8,
                  color: '#fff',
                  fontSize: 15,
                  fontWeight: 600,
                  cursor: 'pointer',
                  opacity: loading !== null ? 0.7 : 1,
                }}
              >
                {loading === plan.name ? 'Processing...' : plan.cta}
              </button>
            </div>
          ))}
        </div>
      </main>
    </div>
  );
}
