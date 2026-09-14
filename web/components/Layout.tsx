import Hero from './Hero';
import Features from './Features';
import Pricing from './Pricing';
import Footer from './Footer';
import Navbar from './Navbar';

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <div style={{ minHeight: '100vh', background: '#061b31', color: '#fff' }}>
      <Navbar />
      {children}
      <Hero />
      <Features />
      <Pricing />
      <Footer />
    </div>
  );
}
