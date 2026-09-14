import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'AudioAI - AI-Powered Audio Processing Platform',
  description: 'Transcribe, summarize, and synthesize audio files with AI. Built for creators and developers.',
  keywords: ['AI', 'audio', 'transcription', 'summarization', 'text-to-speech'],
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className={inter.className}>
      <body>{children}</body>
    </html>
  );
}
