'use client';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import AudioWaveform from '@/components/AudioWaveform';

export default function Dashboard() {
  const router = useRouter();
  const [user, setUser] = useState<any>(null);
  const [usage, setUsage] = useState<any>(null);
  const [analytics, setAnalytics] = useState<any>(null);
  const [jobs, setJobs] = useState<any[]>([]);
  const [uploading, setUploading] = useState(false);
  const [selectedFiles, setSelectedFiles] = useState<File[]>([]);
  const [uploadProgress, setUploadProgress] = useState<{id: string, name: string, status: string}[]>([]);
  const [languages, setLanguages] = useState<any[]>([]);
  const [models, setModels] = useState<any[]>([]);
  const [effects, setEffects] = useState<any[]>([]);
  const [selectedLanguage, setSelectedLanguage] = useState('auto');
  const [selectedModel, setSelectedModel] = useState('nova-2');
  const [showOptions, setShowOptions] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) { router.push('/login'); return; }
    
    Promise.all([
      fetch('/api/auth/me', { headers: { Authorization: `Bearer ${token}` } }).then(r => r.json()),
      fetch('/api/user/usage', { headers: { Authorization: `Bearer ${token}` } }).then(r => r.json()),
      fetch('/api/user/analytics', { headers: { Authorization: `Bearer ${token}` } }).then(r => r.json().catch(() => null)),
      fetch('/api/audio/languages', { headers: { Authorization: `Bearer ${token}` } }).then(r => r.json().catch(() => ({languages: []}))),
      fetch('/api/audio/models', { headers: { Authorization: `Bearer ${token}` } }).then(r => r.json().catch(() => ({models: []}))),
      fetch('/api/audio/effects', { headers: { Authorization: `Bearer ${token}` } }).then(r => r.json().catch(() => ({effects: []}))),
    ]).then(([u, u2, a, langs, models, eff]) => { 
      setUser(u); 
      setUsage(u2); 
      setAnalytics(a);
      if (a?.recent_jobs) setJobs(a.recent_jobs);
      setLanguages(langs.languages || []);
      setModels(models.models || []);
      setEffects(eff.effects || []);
    });
  }, [router]);

  async function handleUpload(e: React.ChangeEvent<HTMLInputElement>) {
    const files = e.target.files;
    if (!files || files.length === 0) return;
    
    const fileArray = Array.from(files);
    setSelectedFiles(fileArray);
    setUploading(true);
    setUploadProgress(fileArray.map(f => ({ id: '', name: f.name, status: 'uploading' })));
    
    // Submit with options
    const formData = new FormData();
    formData.append('audio', fileArray[0]);
    formData.append('language', selectedLanguage);
    formData.append('model', selectedModel);
    
    try {
      const res = await fetch('/api/audio/submit', {
        method: 'POST', 
        headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }, 
        body: formData,
      });
      const data = await res.json();
      
      if (data.job_id) {
        const newJob = {
          id: data.job_id,
          name: data.file_name,
          status: 'processing',
          created_at: Math.floor(Date.now()/1000),
          language: selectedLanguage,
          model: selectedModel
        };
        
        setJobs(prev => [newJob, ...prev]);
        setUploadProgress([{ id: data.job_id, name: data.file_name, status: 'processing' }]);
        
        pollJob(data.job_id);
      }
    } finally {
      setUploading(false);
    }
  }

  async function pollJob(jobId: string) {
    const token = localStorage.getItem('token');
    for (let i = 0; i < 15; i++) {
      await new Promise(r => setTimeout(r, 2000));
      const res = await fetch(`/api/audio/${jobId}`, { headers: { Authorization: `Bearer ${token}` } });
      const job = await res.json();
      setJobs(prev => prev.map(j => j.id === jobId ? { ...j, ...job } : j));
      setUploadProgress(prev => prev.map(p => p.id === jobId ? { ...p, status: job.status } : p));
      if (job.status === 'completed' || job.status === 'failed') break;
    }
  }

  async function downloadSubtitle(jobId: string, format: string) {
    const token = localStorage.getItem('token');
    const res = await fetch(`/api/audio/${jobId}/subtitle?format=${format}`, {
      headers: { Authorization: `Bearer ${token}` }
    });
    const data = await res.json();
    
    if (data.content) {
      const blob = new Blob([data.content], { type: 'text/plain' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `subtitle.${format}`;
      a.click();
      URL.revokeObjectURL(url);
    }
  }

  if (!user) return <div style={{ padding: 40, color: '#8899a6' }}>Loading...</div>;

  return (
    <div style={{ minHeight: '100vh', background: '#061b31' }}>
      {/* Header */}
      <header style={{ borderBottom: '1px solid rgba(255,255,255,0.1)', padding: '16px 24px', display: 'flex', justifyContent: 'space-between', alignItems: 'center', backdropFilter: 'blur(10px)', position: 'sticky', top: 0, background: 'rgba(6,27,49,0.9)', zIndex: 100 }}>
        <span style={{ fontSize: 20, fontWeight: 700, color: '#533afd' }}>🎵 AudioAI</span>
        <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
          <span style={{ color: '#8899a6', fontSize: 14 }}>{user.email}</span>
          {user.is_pro && <span style={{ background: 'linear-gradient(135deg, #533afd, #7c5cfc)', padding: '4px 12px', borderRadius: 20, fontSize: 12, fontWeight: 600 }}>PRO</span>}
          <button onClick={() => { localStorage.removeItem('token'); router.push('/login'); }} style={{ background: 'none', border: '1px solid rgba(255,255,255,0.2)', color: '#8899a6', padding: '8px 16px', borderRadius: 8, cursor: 'pointer' }}>Logout</button>
        </div>
      </header>

      <main style={{ maxWidth: 900, margin: '0 auto', padding: '40px 24px' }}>
        {/* Analytics Overview */}
        {analytics && (
          <div style={{ marginBottom: 24 }}>
            <h2 style={{ marginBottom: 16, fontSize: 18, color: '#f8fafc' }}>Processing Analytics</h2>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))', gap: 16 }}>
              <div style={{ background: 'rgba(255,255,255,0.05)', border: '1px solid rgba(255,255,255,0.1)', borderRadius: 12, padding: 20 }}>
                <div style={{ color: '#8899a6', fontSize: 13 }}>Total Jobs</div>
                <div style={{ fontSize: 32, fontWeight: 700, color: '#533afd' }}>{analytics.total_jobs}</div>
              </div>
              <div style={{ background: 'rgba(255,255,255,0.05)', border: '1px solid rgba(255,255,255,0.1)', borderRadius: 12, padding: 20 }}>
                <div style={{ color: '#8899a6', fontSize: 13 }}>Completed</div>
                <div style={{ fontSize: 32, fontWeight: 700, color: '#4ade80' }}>{analytics.completed_jobs}</div>
              </div>
              <div style={{ background: 'rgba(255,255,255,0.05)', border: '1px solid rgba(255,255,255,0.1)', borderRadius: 12, padding: 20 }}>
                <div style={{ color: '#8899a6', fontSize: 13 }}>Failed</div>
                <div style={{ fontSize: 32, fontWeight: 700, color: '#f87171' }}>{analytics.failed_jobs}</div>
              </div>
              <div style={{ background: 'rgba(255,255,255,0.05)', border: '1px solid rgba(255,255,255,0.1)', borderRadius: 12, padding: 20 }}>
                <div style={{ color: '#8899a6', fontSize: 13 }}>Avg Duration</div>
                <div style={{ fontSize: 32, fontWeight: 700, color: '#f59e0b' }}>{Math.round(analytics.avg_duration_sec)}s</div>
              </div>
            </div>
          </div>
        )}

        {/* Usage Card */}
        <div style={{ background: 'rgba(255,255,255,0.05)', border: '1px solid rgba(255,255,255,0.1)', borderRadius: 16, padding: 24, marginBottom: 24 }}>
          <h2 style={{ marginBottom: 16, fontSize: 18 }}>Usage</h2>
          <div style={{ display: 'flex', gap: 32 }}>
            <div>
              <div style={{ color: '#8899a6', fontSize: 13 }}>Free Quota</div>
              <div style={{ fontSize: 32, fontWeight: 700, color: usage?.remaining > 0 ? '#4ade80' : '#f87171' }}>{usage?.remaining_min ?? 30}</div>
              <div style={{ color: '#8899a6', fontSize: 12 }}>min remaining</div>
            </div>
            <div>
              <div style={{ color: '#8899a6', fontSize: 13 }}>Used</div>
              <div style={{ fontSize: 32, fontWeight: 700 }}>{usage?.usage_minutes ?? 0}</div>
              <div style={{ color: '#8899a6', fontSize: 12 }}>min used</div>
            </div>
            {user.is_pro && (
              <div>
                <div style={{ color: '#4ade80', fontSize: 13 }}>Pro Plan</div>
                <div style={{ fontSize: 24, fontWeight: 700 }}>∞</div>
                <div style={{ color: '#8899a6', fontSize: 12 }}>unlimited</div>
              </div>
            )}
          </div>
          {!user.is_pro && (
            <button style={{ marginTop: 16, padding: '10px 24px', background: 'linear-gradient(135deg, #533afd, #7c5cfc)', border: 'none', borderRadius: 8, color: '#fff', fontWeight: 600, cursor: 'pointer' }}>
              Upgrade to Pro — $9.9/mo
            </button>
          )}
        </div>

        {/* Upload Card */}
        <div style={{ background: 'rgba(255,255,255,0.05)', border: '1px solid rgba(255,255,255,0.1)', borderRadius: 16, padding: 24, marginBottom: 24 }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
            <h2 style={{ marginBottom: 0, fontSize: 18 }}>Upload Audio</h2>
            <button 
              onClick={() => setShowOptions(!showOptions)}
              style={{ background: 'none', border: '1px solid rgba(255,255,255,0.2)', color: '#8899a6', padding: '8px 16px', borderRadius: 8, cursor: 'pointer', fontSize: 13 }}
            >
              {showOptions ? 'Hide Options' : 'Advanced Options'}
            </button>
          </div>
          
          {showOptions && (
            <div style={{ marginBottom: 16, padding: 16, background: 'rgba(0,0,0,0.2)', borderRadius: 8 }}>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 16 }}>
                <div>
                  <label style={{ display: 'block', color: '#8899a6', fontSize: 12, marginBottom: 4 }}>Language</label>
                  <select 
                    value={selectedLanguage}
                    onChange={e => setSelectedLanguage(e.target.value)}
                    style={{ width: '100%', padding: '8px 12px', background: 'rgba(255,255,255,0.1)', border: '1px solid rgba(255,255,255,0.2)', borderRadius: 6, color: '#fff' }}
                  >
                    {languages.map((lang: any) => (
                      <option key={lang.code} value={lang.code}>{lang.native} ({lang.name})</option>
                    ))}
                  </select>
                </div>
                <div>
                  <label style={{ display: 'block', color: '#8899a6', fontSize: 12, marginBottom: 4 }}>Model</label>
                  <select 
                    value={selectedModel}
                    onChange={e => setSelectedModel(e.target.value)}
                    style={{ width: '100%', padding: '8px 12px', background: 'rgba(255,255,255,0.1)', border: '1px solid rgba(255,255,255,0.2)', borderRadius: 6, color: '#fff' }}
                  >
                    {models.map((m: any) => (
                      <option key={m.code} value={m.code}>{m.name}</option>
                    ))}
                  </select>
                </div>
              </div>
              <div style={{ marginTop: 12 }}>
                <label style={{ display: 'block', color: '#8899a6', fontSize: 12, marginBottom: 4 }}>Available Effects</label>
                <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8 }}>
                  {effects.map((eff: any) => (
                    <span key={eff.type} style={{ padding: '4px 12px', background: 'rgba(83,58,253,0.2)', borderRadius: 20, fontSize: 12, color: '#a5b4fc' }}>
                      {eff.name}
                    </span>
                  ))}
                </div>
              </div>
            </div>
          )}
          
          <div style={{ border: '2px dashed rgba(255,255,255,0.2)', borderRadius: 12, padding: 40, textAlign: 'center', cursor: 'pointer', transition: 'all 0.2s' }}
            onClick={() => document.getElementById('file-input')?.click()}
            onMouseEnter={e => (e.currentTarget.style.borderColor = '#533afd')}
            onMouseLeave={e => (e.currentTarget.style.borderColor = 'rgba(255,255,255,0.2)')}>
            <div style={{ fontSize: 40, marginBottom: 12 }}>🎙️</div>
            <div style={{ color: '#fff', fontSize: 16, fontWeight: 600 }}>
              {selectedFiles.length > 0 ? `${selectedFiles.length} files selected` : 'Click to upload audio'}
            </div>
            <div style={{ color: '#8899a6', fontSize: 13, marginTop: 4 }}>MP3, WAV, M4A, FLAC, OGG, AAC (Multiple files supported)</div>
          </div>
          <input id="file-input" type="file" accept="audio/*" multiple style={{ display: 'none' }} onChange={handleUpload} />
          {uploading && <div style={{ marginTop: 12, color: '#533afd' }}>⏳ Processing...</div>}
        </div>

        {/* Upload Progress */}
        {uploadProgress.length > 0 && (
          <div style={{ background: 'rgba(255,255,255,0.05)', border: '1px solid rgba(255,255,255,0.1)', borderRadius: 12, padding: 16, marginBottom: 24 }}>
            <h3 style={{ marginBottom: 12, fontSize: 14, color: '#f8fafc' }}>Upload Progress</h3>
            {uploadProgress.map((p, i) => (
              <div key={i} style={{ display: 'flex', justifyContent: 'space-between', padding: '8px 0', borderBottom: i < uploadProgress.length - 1 ? '1px solid rgba(255,255,255,0.05)' : 'none' }}>
                <span style={{ color: '#fff', fontSize: 13 }}>{p.name}</span>
                <span style={{ 
                  color: p.status === 'completed' ? '#4ade80' : p.status === 'processing' ? '#533afd' : p.status === 'failed' ? '#f87171' : '#8899a6',
                  fontSize: 12,
                  fontWeight: 600
                }}>
                  {p.status}
                </span>
              </div>
            ))}
          </div>
        )}

        {/* Jobs List */}
        <div>
          <h2 style={{ marginBottom: 16, fontSize: 18 }}>Recent Jobs</h2>
          {jobs.length === 0 ? (
            <div style={{ color: '#8899a6', textAlign: 'center', padding: 40 }}>No jobs yet. Upload an audio file to get started.</div>
          ) : (
            jobs.map(job => (
              <div key={job.id} style={{ background: 'rgba(255,255,255,0.05)', border: '1px solid rgba(255,255,255,0.1)', borderRadius: 12, padding: 16, marginBottom: 8 }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
                  <div>
                    <div style={{ fontWeight: 600 }}>{job.file_name || job.name}</div>
                    <div style={{ color: '#8899a6', fontSize: 12 }}>
                      {new Date((job.created_at || job.created_at) * 1000).toLocaleString()}
                      {job.language && <span style={{ marginLeft: 12 }}>🌐 {job.language}</span>}
                      {job.model && <span style={{ marginLeft: 8 }}>🤖 {job.model}</span>}
                    </div>
                  </div>
                  <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                    {job.status === 'completed' && (
                      <>
                        <button 
                          onClick={() => downloadSubtitle(job.id, 'srt')}
                          style={{ padding: '4px 12px', background: 'rgba(83,58,253,0.2)', border: '1px solid rgba(83,58,253,0.4)', borderRadius: 6, color: '#a5b4fc', fontSize: 12, cursor: 'pointer' }}
                        >
                          SRT
                        </button>
                        <button 
                          onClick={() => downloadSubtitle(job.id, 'vtt')}
                          style={{ padding: '4px 12px', background: 'rgba(83,58,253,0.2)', border: '1px solid rgba(83,58,253,0.4)', borderRadius: 6, color: '#a5b4fc', fontSize: 12, cursor: 'pointer' }}
                        >
                          VTT
                        </button>
                      </>
                    )}
                    <span style={{ padding: '4px 12px', borderRadius: 20, fontSize: 12, fontWeight: 600, background: job.status === 'completed' ? 'rgba(74,222,128,0.2)' : job.status === 'processing' ? 'rgba(83,58,253,0.2)' : 'rgba(239,68,68,0.2)', color: job.status === 'completed' ? '#4ade80' : job.status === 'processing' ? '#533afd' : '#f87171' }}>
                      {job.status}
                    </span>
                  </div>
                </div>
                {job.status === 'completed' && job.result && (
                  <div style={{ marginTop: 12 }}>
                    <AudioWaveform audioUrl={`/api/audio/${job.id}/download`} height={60} />
                    <div style={{ marginTop: 8, color: '#e2e8f0', fontSize: 13, whiteSpace: 'pre-wrap' }}>
                      {typeof job.result === 'string' ? job.result.substring(0, 200) + '...' : ''}
                    </div>
                  </div>
                )}
              </div>
            ))
          )}
        </div>
      </main>
    </div>
  );
}
