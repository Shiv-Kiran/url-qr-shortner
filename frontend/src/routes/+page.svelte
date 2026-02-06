<script>
  import { onMount } from 'svelte';
  import { getQRCode, shortenURL } from '../lib/api';

  const qrLevels = [
    { value: 'L', label: 'L (7%)' },
    { value: 'M', label: 'M (15%)' },
    { value: 'Q', label: 'Q (25%)' },
    { value: 'H', label: 'H (30%)' }
  ];

  const kpis = [
    { label: 'Total Clicks', value: '128.4K', delta: '+12.6%' },
    { label: 'Unique Visitors', value: '48.1K', delta: '+7.3%' },
    { label: 'Active Links', value: '312', delta: '+3.2%' },
    { label: 'QR Scans', value: '8.9K', delta: '+18.9%' }
  ];

  const recentLinks = [
    { name: 'campaign-alpha', clicks: '12.4K', status: 'Active' },
    { name: 'qr-launch', clicks: '6.7K', status: 'Active' },
    { name: 'bundle-summit', clicks: '3.9K', status: 'Paused' },
    { name: 'a-b-test-x', clicks: '1.2K', status: 'Running' }
  ];

  let originalUrl = '';
  let selectedQRLevel = 'M';
  let qrSize = 256;
  let loading = false;
  let refreshingQR = false;
  let error = '';
  let copied = false;
  let result = null;
  let isDark = true;

  onMount(() => {
    isDark = document.documentElement.classList.contains('dark');
  });

  async function handleShorten() {
    if (!originalUrl.trim()) {
      error = 'Enter a destination URL first.';
      return;
    }

    loading = true;
    error = '';

    try {
      result = await shortenURL({
        original_url: originalUrl.trim(),
        qr_error_correction: selectedQRLevel,
        qr_size: Number(qrSize)
      });
      selectedQRLevel = result.qr_error_correction;
      qrSize = result.qr_size;
      originalUrl = '';
      copied = false;
    } catch (err) {
      error = err.response?.data?.error || 'Failed to shorten URL';
    } finally {
      loading = false;
    }
  }

  async function refreshQRCode() {
    if (!result) return;

    refreshingQR = true;
    error = '';

    try {
      const qrResult = await getQRCode(result.short_code, selectedQRLevel, Number(qrSize));
      result = {
        ...result,
        ...qrResult
      };
      selectedQRLevel = qrResult.qr_error_correction;
      qrSize = qrResult.qr_size;
    } catch (err) {
      error = err.response?.data?.error || 'Failed to refresh QR code';
    } finally {
      refreshingQR = false;
    }
  }

  async function copyShortURL() {
    if (!result?.short_url) return;
    await navigator.clipboard.writeText(result.short_url);
    copied = true;
    setTimeout(() => {
      copied = false;
    }, 1800);
  }

  function downloadQRCode() {
    if (!result?.qr_data_url) return;

    const anchor = document.createElement('a');
    anchor.href = result.qr_data_url;
    anchor.download = `${result.short_code}-qr.png`;
    anchor.click();
  }

  function handleKeyDown(event) {
    if (event.key === 'Enter' && !loading) {
      handleShorten();
    }
  }

  function toggleTheme() {
    isDark = !isDark;
    if (isDark) {
      document.documentElement.classList.add('dark');
      localStorage.setItem('theme', 'dark');
    } else {
      document.documentElement.classList.remove('dark');
      localStorage.setItem('theme', 'light');
    }
  }
</script>

<main class="shell">
  <header class="topbar">
    <div>
      <p class="eyebrow">Analytics Studio</p>
      <h1>URL Intelligence Hub</h1>
      <p class="subtitle">Dark-first, glassy dashboards for link intelligence, QR previews, and A/B experiments.</p>
    </div>
    <div class="topbar-actions">
      <button class="ghost" on:click={toggleTheme} aria-label="Toggle theme">
        {#if isDark}
          <span class="icon">
            <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 2a1 1 0 0 1 1 1v2a1 1 0 1 1-2 0V3a1 1 0 0 1 1-1zm6.36 3.64a1 1 0 0 1 1.41 0l1.42 1.42a1 1 0 0 1-1.42 1.41l-1.41-1.41a1 1 0 0 1 0-1.42zM21 11a1 1 0 0 1 1 1v2a1 1 0 1 1-2 0v-2a1 1 0 0 1 1-1zm-9 3a4 4 0 1 0 0-8 4 4 0 0 0 0 8zm-9 1a1 1 0 0 1-1-1v-2a1 1 0 1 1 2 0v2a1 1 0 0 1-1 1zm1.22-9.36a1 1 0 0 1 0 1.41L2.8 8.46a1 1 0 0 1-1.41-1.41l1.42-1.42a1 1 0 0 1 1.41 0zm14.14 14.14a1 1 0 0 1 0 1.41l-1.42 1.42a1 1 0 0 1-1.41-1.41l1.41-1.42a1 1 0 0 1 1.42 0zM12 19a1 1 0 0 1 1 1v2a1 1 0 1 1-2 0v-2a1 1 0 0 1 1-1zM5.64 18.36a1 1 0 0 1 1.41 0l1.41 1.42a1 1 0 0 1-1.41 1.41l-1.42-1.41a1 1 0 0 1 0-1.42z" fill="currentColor"/></svg>
          </span>
          Dark
        {:else}
          <span class="icon">
            <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M21 15.5A8.5 8.5 0 0 1 8.5 3a.75.75 0 0 0-.71 1.03A7 7 0 0 0 12 21a7 7 0 0 0 6.97-5.79.75.75 0 0 0-.97-.71Z" fill="currentColor"/></svg>
          </span>
          Light
        {/if}
      </button>
      <button class="primary">Create New Link</button>
    </div>
  </header>

  <section class="grid">
    <div class="card hero">
      <div class="hero-header">
        <div>
          <p class="eyebrow">Routing - IDs - Persistence</p>
          <h2>Create a short link + QR package</h2>
          <p class="subtitle">Generate a short URL, tune QR error correction, and ship the scan-ready code.</p>
        </div>
        <div class="hero-badge">Live</div>
      </div>

      <div class="form">
        <div>
          <label class="label" for="url-input">Destination URL</label>
          <input
            id="url-input"
            class="field"
            type="url"
            placeholder="https://example.com/launch/spring-2026"
            bind:value={originalUrl}
            on:keydown={handleKeyDown}
            disabled={loading}
          />
        </div>
        <div class="field-row">
          <div>
            <label class="label" for="qr-level">QR Error Correction</label>
            <select id="qr-level" class="field" bind:value={selectedQRLevel} disabled={loading || refreshingQR}>
              {#each qrLevels as level}
                <option value={level.value}>{level.label}</option>
              {/each}
            </select>
          </div>
          <div>
            <label class="label" for="qr-size">QR Size</label>
            <input
              id="qr-size"
              class="field"
              type="number"
              min="128"
              max="1024"
              step="32"
              bind:value={qrSize}
              disabled={loading || refreshingQR}
            />
          </div>
        </div>
        <button class="primary" on:click={handleShorten} disabled={loading}>
          {loading ? 'Generating...' : 'Create Short URL + QR'}
        </button>
        {#if error}
          <p class="error">{error}</p>
        {/if}
      </div>
    </div>

    <div class="card kpi">
      <p class="eyebrow">Realtime KPIs</p>
      <div class="kpi-grid">
        {#each kpis as metric}
          <div class="kpi-tile">
            <p class="kpi-label">{metric.label}</p>
            <h3>{metric.value}</h3>
            <span class="chip">{metric.delta}</span>
          </div>
        {/each}
      </div>
    </div>

    <div class="card chart">
      <div class="card-header">
        <div>
          <p class="eyebrow">Traffic Pulse</p>
          <h3>Hourly Clicks</h3>
        </div>
        <span class="chip">Last 24h</span>
      </div>
      <div class="sparkline">
        <span></span><span></span><span></span><span></span><span></span><span></span><span></span><span></span>
      </div>
      <p class="muted">Peak activity during launch window. Keep an eye on morning spikes.</p>
    </div>

    <div class="card map">
      <div class="card-header">
        <div>
          <p class="eyebrow">Geo Map</p>
          <h3>Top Regions</h3>
        </div>
        <span class="chip">+18%</span>
      </div>
      <div class="globe">
        <div class="glow"></div>
        <div class="orbit"></div>
        <div class="dot dot-1"></div>
        <div class="dot dot-2"></div>
        <div class="dot dot-3"></div>
      </div>
      <div class="region-list">
        <div><span>North America</span><strong>42%</strong></div>
        <div><span>Europe</span><strong>29%</strong></div>
        <div><span>APAC</span><strong>21%</strong></div>
      </div>
    </div>

    <div class="card devices">
      <div class="card-header">
        <div>
          <p class="eyebrow">Devices</p>
          <h3>Device Mix</h3>
        </div>
        <span class="chip">Mobile 62%</span>
      </div>
      <div class="donut">
        <div class="ring"></div>
        <div class="ring-inner"></div>
      </div>
      <div class="device-breakdown">
        <div><span>Mobile</span><strong>62%</strong></div>
        <div><span>Desktop</span><strong>28%</strong></div>
        <div><span>Tablet</span><strong>10%</strong></div>
      </div>
    </div>

    <div class="card qr">
      <div class="card-header">
        <div>
          <p class="eyebrow">QR Preview</p>
          <h3>Live Scan</h3>
        </div>
        <span class="chip">{result ? result.qr_error_correction : 'M'}</span>
      </div>
      <div class="qr-preview">
        {#if result}
          <div class="scanline"></div>
          <img src={result.qr_data_url} alt={`QR code for ${result.short_url}`} />
        {:else}
          <div class="placeholder">Generate a QR code to preview</div>
        {/if}
      </div>
      <div class="qr-actions">
        <button class="ghost" on:click={refreshQRCode} disabled={!result || refreshingQR}>Refresh</button>
        <button class="ghost" on:click={downloadQRCode} disabled={!result}>Download</button>
      </div>
    </div>

    <div class="card links">
      <div class="card-header">
        <div>
          <p class="eyebrow">Recent Links</p>
          <h3>Latest Releases</h3>
        </div>
        <span class="chip">4 Active</span>
      </div>
      <div class="link-list">
        {#each recentLinks as link}
          <div class="link-row">
            <div>
              <strong>{link.name}</strong>
              <span class="muted">{link.clicks} clicks</span>
            </div>
            <span class="status">{link.status}</span>
          </div>
        {/each}
      </div>
      <div class="link-actions">
        <button class="ghost" on:click={copyShortURL} disabled={!result}>{copied ? 'Copied' : 'Copy Latest'}</button>
        <button class="ghost">View All</button>
      </div>
    </div>

    <div class="card experiment">
      <div class="card-header">
        <div>
          <p class="eyebrow">A/B Testing</p>
          <h3>Variant Performance</h3>
        </div>
        <span class="chip">Variant B +9%</span>
      </div>
      <div class="bar">
        <span style="--fill:72%">Variant A</span>
        <span style="--fill:81%">Variant B</span>
      </div>
      <p class="muted">Rotate traffic smoothly while tracking conversion lift.</p>
    </div>

    <div class="card bundles">
      <div class="card-header">
        <div>
          <p class="eyebrow">Link Bundles</p>
          <h3>Multi-link Packs</h3>
        </div>
        <span class="chip">3 Active</span>
      </div>
      <div class="bundle-grid">
        <div class="bundle">Creator Kit</div>
        <div class="bundle">Launch Stack</div>
        <div class="bundle">Partner Links</div>
      </div>
      <p class="muted">Group multiple destinations under a single short link.</p>
    </div>
  </section>
</main>

<style>
  @import url('https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@400;500;600;700&family=Playfair+Display:wght@600;700&display=swap');

  .shell {
    min-height: 100vh;
    padding: clamp(24px, 5vw, 64px);
    background:
      radial-gradient(circle at 10% 20%, rgba(110, 220, 255, 0.15), transparent 35%),
      radial-gradient(circle at 80% 0%, rgba(150, 110, 255, 0.15), transparent 38%),
      radial-gradient(circle at 70% 80%, rgba(110, 255, 210, 0.1), transparent 38%),
      linear-gradient(180deg, rgba(10, 12, 18, 0.92), rgba(6, 8, 14, 0.98)),
      hsl(var(--bg));
    color: hsl(var(--foreground));
  }

  .topbar {
    display: flex;
    justify-content: space-between;
    gap: 24px;
    align-items: center;
    margin-bottom: 32px;
  }

  .topbar h1 {
    font-family: 'Playfair Display', 'Times New Roman', serif;
    font-size: clamp(2.2rem, 5vw, 3.6rem);
    margin: 6px 0 6px;
    letter-spacing: -0.02em;
  }

  .subtitle {
    color: hsl(var(--muted));
    max-width: 520px;
  }

  .eyebrow {
    letter-spacing: 0.2em;
    text-transform: uppercase;
    font-size: 0.68rem;
    margin: 0;
    color: hsl(var(--muted));
  }

  .topbar-actions {
    display: flex;
    gap: 12px;
    align-items: center;
  }

  .grid {
    display: grid;
    grid-template-columns: repeat(12, minmax(0, 1fr));
    gap: 20px;
  }

  .card {
    background: rgba(18, 22, 30, 0.7);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 22px;
    padding: 22px;
    backdrop-filter: blur(18px);
    box-shadow: var(--shadow);
    transition: transform 200ms ease, box-shadow 200ms ease, border 200ms ease;
  }

  .card:hover {
    transform: translateY(-3px);
    border-color: rgba(110, 220, 255, 0.25);
    box-shadow: 0 30px 80px rgba(8, 12, 24, 0.7);
  }

  .hero {
    grid-column: span 7;
  }

  .kpi {
    grid-column: span 5;
  }

  .chart {
    grid-column: span 4;
  }

  .map {
    grid-column: span 4;
  }

  .devices {
    grid-column: span 4;
  }

  .qr {
    grid-column: span 4;
  }

  .links {
    grid-column: span 4;
  }

  .experiment {
    grid-column: span 4;
  }

  .bundles {
    grid-column: span 4;
  }

  .hero-header {
    display: flex;
    justify-content: space-between;
    gap: 18px;
    align-items: flex-start;
    margin-bottom: 20px;
  }

  h2 {
    font-size: clamp(1.4rem, 2.8vw, 2rem);
    margin: 6px 0 0;
  }

  h3 {
    margin: 6px 0 10px;
    font-size: 1.2rem;
  }

  .hero-badge {
    padding: 6px 14px;
    border-radius: 999px;
    background: linear-gradient(135deg, hsl(var(--accent-1)), hsl(var(--accent-2)));
    color: #0b0d12;
    font-weight: 700;
    font-size: 0.72rem;
    letter-spacing: 0.08em;
  }

  .form {
    display: grid;
    gap: 16px;
  }

  .label {
    display: block;
    margin-bottom: 8px;
    font-size: 0.72rem;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    color: hsl(var(--muted));
  }

  .field {
    width: 100%;
    border-radius: 14px;
    padding: 12px 14px;
    border: 1px solid rgba(255, 255, 255, 0.12);
    background: rgba(10, 12, 18, 0.65);
    color: hsl(var(--foreground));
    transition: border 200ms ease, box-shadow 200ms ease;
  }

  .field:focus {
    outline: none;
    border-color: rgba(110, 220, 255, 0.5);
    box-shadow: 0 0 0 3px rgba(110, 220, 255, 0.12);
  }

  .field-row {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 12px;
  }

  .primary {
    border: none;
    border-radius: 14px;
    padding: 12px 18px;
    color: #0a0d12;
    font-weight: 700;
    background: linear-gradient(130deg, hsl(var(--accent-1)), hsl(var(--accent-2)));
    box-shadow: 0 16px 30px rgba(91, 203, 255, 0.35);
    cursor: pointer;
    transition: transform 200ms ease, filter 200ms ease, box-shadow 200ms ease;
  }

  .primary:hover {
    transform: translateY(-1px) scale(1.01);
    filter: brightness(1.08);
  }

  .primary:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    box-shadow: none;
  }

  .ghost {
    border-radius: 12px;
    border: 1px solid rgba(255, 255, 255, 0.12);
    background: rgba(12, 16, 24, 0.6);
    color: hsl(var(--foreground));
    padding: 10px 14px;
    cursor: pointer;
    display: inline-flex;
    align-items: center;
    gap: 8px;
    transition: border 200ms ease, transform 200ms ease;
  }

  .ghost:hover {
    border-color: rgba(110, 220, 255, 0.45);
    transform: translateY(-1px);
  }

  .icon {
    display: inline-flex;
    width: 18px;
    height: 18px;
  }

  .error {
    margin: 0;
    color: #ff9d8a;
    font-weight: 500;
  }

  .kpi-grid {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 12px;
  }

  .kpi-tile {
    padding: 14px;
    border-radius: 16px;
    background: rgba(12, 16, 24, 0.55);
    border: 1px solid rgba(255, 255, 255, 0.08);
  }

  .kpi-label {
    color: hsl(var(--muted));
    margin: 0 0 8px;
    font-size: 0.84rem;
  }

  .chip {
    padding: 4px 10px;
    border-radius: 999px;
    background: rgba(110, 220, 255, 0.15);
    color: hsl(var(--accent-1));
    font-size: 0.72rem;
    font-weight: 600;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
  }

  .sparkline {
    display: grid;
    grid-template-columns: repeat(8, minmax(0, 1fr));
    gap: 8px;
    align-items: end;
    height: 120px;
    margin-bottom: 10px;
  }

  .sparkline span {
    background: linear-gradient(180deg, hsl(var(--accent-1)), rgba(255, 255, 255, 0));
    border-radius: 8px;
    height: 100%;
    transform: scaleY(0.4);
    transform-origin: bottom;
  }

  .sparkline span:nth-child(2) { transform: scaleY(0.55); }
  .sparkline span:nth-child(3) { transform: scaleY(0.72); }
  .sparkline span:nth-child(4) { transform: scaleY(0.35); }
  .sparkline span:nth-child(5) { transform: scaleY(0.85); }
  .sparkline span:nth-child(6) { transform: scaleY(0.62); }
  .sparkline span:nth-child(7) { transform: scaleY(0.78); }
  .sparkline span:nth-child(8) { transform: scaleY(0.5); }

  .muted {
    color: hsl(var(--muted));
  }

  .globe {
    height: 160px;
    border-radius: 20px;
    background: radial-gradient(circle at 40% 40%, rgba(110, 220, 255, 0.35), rgba(6, 10, 18, 0.2));
    position: relative;
    overflow: hidden;
    margin-bottom: 14px;
  }

  .glow {
    position: absolute;
    width: 200px;
    height: 200px;
    border-radius: 50%;
    background: rgba(110, 220, 255, 0.2);
    top: -40px;
    left: -20px;
    filter: blur(20px);
  }

  .orbit {
    position: absolute;
    inset: 12px;
    border-radius: 50%;
    border: 1px dashed rgba(110, 220, 255, 0.4);
    animation: spin 18s linear infinite;
  }

  .dot {
    position: absolute;
    width: 10px;
    height: 10px;
    border-radius: 50%;
    background: hsl(var(--accent-3));
    box-shadow: 0 0 10px rgba(110, 255, 210, 0.5);
  }

  .dot-1 { top: 32%; left: 28%; }
  .dot-2 { top: 45%; left: 55%; }
  .dot-3 { top: 62%; left: 38%; }

  .region-list {
    display: grid;
    gap: 8px;
  }

  .region-list div {
    display: flex;
    justify-content: space-between;
    color: hsl(var(--muted));
  }

  .donut {
    height: 150px;
    position: relative;
    margin-bottom: 12px;
    display: grid;
    place-items: center;
  }

  .ring {
    width: 140px;
    height: 140px;
    border-radius: 50%;
    background: conic-gradient(
      hsl(var(--accent-1)) 0 62%,
      hsl(var(--accent-2)) 62% 90%,
      rgba(255, 255, 255, 0.1) 90% 100%
    );
  }

  .ring-inner {
    position: absolute;
    width: 90px;
    height: 90px;
    border-radius: 50%;
    background: rgba(10, 12, 18, 0.9);
  }

  .device-breakdown {
    display: grid;
    gap: 6px;
  }

  .device-breakdown div {
    display: flex;
    justify-content: space-between;
    color: hsl(var(--muted));
  }

  .qr-preview {
    height: 220px;
    border-radius: 18px;
    background: rgba(12, 16, 24, 0.6);
    border: 1px dashed rgba(255, 255, 255, 0.2);
    display: grid;
    place-items: center;
    position: relative;
    overflow: hidden;
  }

  .qr-preview img {
    width: 160px;
    height: 160px;
    border-radius: 12px;
    background: #fff;
    padding: 10px;
  }

  .scanline {
    position: absolute;
    width: 100%;
    height: 2px;
    background: linear-gradient(90deg, transparent, rgba(110, 220, 255, 0.8), transparent);
    animation: scan 3s ease-in-out infinite;
  }

  .placeholder {
    color: hsl(var(--muted));
  }

  .qr-actions {
    display: flex;
    gap: 10px;
    margin-top: 12px;
  }

  .link-list {
    display: grid;
    gap: 12px;
  }

  .link-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 12px;
    border-radius: 14px;
    background: rgba(12, 16, 24, 0.55);
  }

  .link-row strong {
    display: block;
  }

  .status {
    padding: 4px 10px;
    border-radius: 999px;
    border: 1px solid rgba(110, 220, 255, 0.35);
    color: hsl(var(--accent-1));
    font-size: 0.72rem;
  }

  .link-actions {
    display: flex;
    gap: 10px;
    margin-top: 12px;
  }

  .bar span {
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: rgba(10, 12, 18, 0.6);
    border-radius: 999px;
    padding: 8px 12px;
    margin-bottom: 10px;
    position: relative;
    overflow: hidden;
  }

  .bar span::after {
    content: '';
    position: absolute;
    left: 0;
    top: 0;
    height: 100%;
    width: var(--fill);
    background: linear-gradient(90deg, rgba(110, 220, 255, 0.4), rgba(165, 120, 255, 0.4));
    z-index: 0;
  }

  .bar span {
    color: hsl(var(--foreground));
  }

  .bar span * {
    position: relative;
    z-index: 1;
  }

  .bundle-grid {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 10px;
  }

  .bundle {
    padding: 12px;
    border-radius: 14px;
    background: rgba(12, 16, 24, 0.55);
    text-align: center;
    font-weight: 600;
  }

  @media (max-width: 1100px) {
    .hero { grid-column: span 12; }
    .kpi { grid-column: span 12; }
    .chart,
    .map,
    .devices,
    .qr,
    .links,
    .experiment,
    .bundles { grid-column: span 6; }
  }

  @media (max-width: 760px) {
    .topbar {
      flex-direction: column;
      align-items: flex-start;
    }

    .grid {
      grid-template-columns: repeat(1, minmax(0, 1fr));
    }

    .card { grid-column: span 1; }
    .field-row { grid-template-columns: 1fr; }
    .topbar-actions { width: 100%; justify-content: flex-start; }
  }

  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }

  @keyframes scan {
    0% { transform: translateY(-60px); opacity: 0; }
    40% { opacity: 1; }
    100% { transform: translateY(200px); opacity: 0; }
  }
</style>
