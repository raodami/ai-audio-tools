export default function Pricing() {
  const plans = [
    { name: 'Free', price: '$0', period: '/month', features: ['30 min transcription', 'Basic summaries', '1 TTS voice'], popular: false },
    { name: 'Pro', price: '$9.9', period: '/month', features: ['500 min transcription', 'AI summaries', 'All TTS voices', 'Priority support'], popular: true },
    { name: 'Team', price: '$29.9', period: '/month', features: ['Unlimited transcription', 'Advanced analytics', 'API access', 'Team management'], popular: false },
  ];

  return (
    <section className="pricing" id="pricing">
      <div className="container">
        <h2>Simple Pricing</h2>
        <div className="pricing-grid">
          {plans.map((plan, i) => (
            <div className={`pricing-card ${plan.popular ? 'featured' : ''}`} key={i}>
              <h3>{plan.name}</h3>
              <div className="price">{plan.price}<span>{plan.period}</span></div>
              <ul>
                {plan.features.map((f, j) => <li key={j}>{f}</li>)}
              </ul>
              <button className="btn btn-primary" style={{ width: '100%' }}>
                {plan.popular ? 'Get Started' : 'Choose Plan'}
              </button>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
