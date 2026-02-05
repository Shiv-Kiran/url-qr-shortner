<script>
  import { getQRCode, shortenURL } from '../lib/api';

  const qrLevels = [
    { value: 'L', label: 'L - 7% recovery' },
    { value: 'M', label: 'M - 15% recovery' },
    { value: 'Q', label: 'Q - 25% recovery' },
    { value: 'H', label: 'H - 30% recovery' }
  ];

  let originalUrl = '';
  let selectedQRLevel = 'M';
  let qrSize = 256;
  let loading = false;
  let refreshingQR = false;
  let error = '';
  let copied = false;
  let result = null;

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
</script>

<main class="page">
  <section class="panel creator">
    <p class="eyebrow">Routing - IDs - Persistence</p>
    <h1>URL Shortener + Custom QR Studio</h1>
    <p class="lead">
      Create a short URL, choose QR error correction, and generate a scannable code you can copy or download.
    </p>

    <label class="field-label" for="url-input">Destination URL</label>
    <input
      id="url-input"
      type="url"
      placeholder="https://example.com/very/long/path"
      bind:value={originalUrl}
      on:keydown={handleKeyDown}
      disabled={loading}
    />

    <div class="options">
      <label class="field-label" for="qr-level">QR Error Correction</label>
      <select id="qr-level" bind:value={selectedQRLevel} disabled={loading || refreshingQR}>
        {#each qrLevels as level}
          <option value={level.value}>{level.label}</option>
        {/each}
      </select>

      <label class="field-label" for="qr-size">QR Size (128-1024 px)</label>
      <input id="qr-size" type="number" min="128" max="1024" step="32" bind:value={qrSize} disabled={loading || refreshingQR} />
    </div>

    <button class="primary" on:click={handleShorten} disabled={loading}>
      {loading ? 'Creating short URL...' : 'Create Short URL + QR'}
    </button>

    {#if error}
      <p class="error">{error}</p>
    {/if}
  </section>

  <section class="panel results">
    {#if result}
      <div class="result-head">
        <p class="eyebrow">Generated</p>
        <h2>{result.short_code}</h2>
      </div>

      <label class="field-label" for="short-url">Short URL</label>
      <input id="short-url" type="text" value={result.short_url} readonly />

      <div class="action-row">
        <button class="ghost" on:click={copyShortURL}>{copied ? 'Copied' : 'Copy URL'}</button>
        <button class="ghost" on:click={refreshQRCode} disabled={refreshingQR || loading}>
          {refreshingQR ? 'Refreshing QR...' : 'Refresh QR'}
        </button>
        <button class="ghost" on:click={downloadQRCode}>Download QR</button>
      </div>

      <div class="qr-card">
        <img src={result.qr_data_url} alt={`QR code for ${result.short_url}`} width={result.qr_size} height={result.qr_size} />
      </div>

      <p class="meta">
        Error correction: <strong>{result.qr_error_correction}</strong> - Size: <strong>{result.qr_size}px</strong>
      </p>
    {:else}
      <div class="empty-state">
        <h2>No Short URL Yet</h2>
        <p>Your generated short link and QR code will appear here.</p>
      </div>
    {/if}
  </section>
</main>

<style>
  @import url('https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@400;500;700&family=DM+Serif+Display:ital@0;1&display=swap');

  :global(:root) {
    --bg-ink: #0f1e2e;
    --card-edge: #e6d5bc;
    --accent: #f08a4b;
    --ink: #11243a;
    --muted: #5d6f82;
    --danger-bg: #ffe5d6;
    --danger: #8f2b13;
  }

  .page {
    min-height: 100vh;
    padding: clamp(24px, 5vw, 56px);
    display: grid;
    gap: 22px;
    grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
    background:
      radial-gradient(circle at 20% 10%, #ffd89f 0%, transparent 30%),
      radial-gradient(circle at 80% 85%, #95d8ff 0%, transparent 35%),
      linear-gradient(160deg, var(--bg-ink), #1a3551 45%, #0f293f 100%);
    font-family: 'Space Grotesk', 'Trebuchet MS', sans-serif;
    color: var(--ink);
  }

  .panel {
    background: linear-gradient(170deg, rgba(255, 250, 242, 0.96), rgba(244, 239, 228, 0.98));
    border: 1px solid var(--card-edge);
    border-radius: 22px;
    padding: clamp(20px, 3vw, 34px);
    box-shadow:
      0 26px 50px rgba(7, 19, 31, 0.26),
      inset 0 1px 0 rgba(255, 255, 255, 0.75);
    animation: rise 420ms ease-out;
  }

  .results {
    animation-delay: 90ms;
  }

  .eyebrow {
    margin: 0;
    letter-spacing: 0.09em;
    text-transform: uppercase;
    font-size: 0.74rem;
    color: var(--muted);
    font-weight: 700;
  }

  h1,
  h2 {
    margin: 8px 0 10px;
    font-family: 'DM Serif Display', Georgia, serif;
    color: var(--ink);
    line-height: 1.05;
  }

  h1 {
    font-size: clamp(2rem, 5vw, 3rem);
  }

  h2 {
    font-size: clamp(1.6rem, 3vw, 2.2rem);
  }

  .lead {
    margin: 0 0 20px;
    color: var(--muted);
    line-height: 1.5;
  }

  .field-label {
    display: block;
    margin-bottom: 6px;
    margin-top: 12px;
    font-size: 0.82rem;
    letter-spacing: 0.06em;
    text-transform: uppercase;
    color: #4f6072;
    font-weight: 700;
  }

  input,
  select,
  button {
    font: inherit;
  }

  input,
  select {
    width: 100%;
    padding: 12px 14px;
    border-radius: 11px;
    border: 1px solid #c5d0dc;
    background: #fff;
    color: var(--ink);
  }

  input:focus,
  select:focus {
    outline: 2px solid #7ec8ff;
    outline-offset: 1px;
    border-color: #7ec8ff;
  }

  .options {
    display: grid;
    grid-template-columns: 1fr;
    gap: 4px;
  }

  .primary {
    margin-top: 18px;
    width: 100%;
    padding: 14px 18px;
    border: none;
    border-radius: 12px;
    background: linear-gradient(120deg, var(--accent), #f4b266);
    color: #fff;
    font-weight: 700;
    cursor: pointer;
    transition: transform 140ms ease, box-shadow 140ms ease;
  }

  .primary:hover:not(:disabled) {
    transform: translateY(-1px);
    box-shadow: 0 10px 20px rgba(176, 87, 32, 0.28);
  }

  .primary:disabled,
  .ghost:disabled {
    opacity: 0.65;
    cursor: not-allowed;
  }

  .error {
    margin-top: 14px;
    padding: 10px 12px;
    border-radius: 10px;
    color: var(--danger);
    background: var(--danger-bg);
    border: 1px solid #f2b9a4;
    font-weight: 500;
  }

  .result-head {
    margin-bottom: 6px;
  }

  .action-row {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
    margin-top: 14px;
  }

  .ghost {
    border-radius: 10px;
    border: 1px solid #b3c5d7;
    background: #f4f9ff;
    color: #153754;
    font-weight: 600;
    padding: 10px 12px;
    cursor: pointer;
  }

  .qr-card {
    margin-top: 18px;
    border-radius: 16px;
    border: 1px solid #d6e2ec;
    background: #fff;
    padding: 14px;
    display: grid;
    place-items: center;
  }

  .qr-card img {
    width: min(100%, 320px);
    height: auto;
    aspect-ratio: 1 / 1;
    image-rendering: pixelated;
  }

  .meta {
    color: var(--muted);
    margin-bottom: 0;
  }

  .empty-state {
    min-height: 100%;
    border: 2px dashed #bed0e2;
    border-radius: 14px;
    padding: 22px;
    display: grid;
    place-content: center;
    text-align: center;
    color: var(--muted);
  }

  @media (max-width: 640px) {
    .page {
      padding: 18px;
    }

    .action-row {
      flex-direction: column;
    }

    .ghost {
      width: 100%;
    }
  }

  @keyframes rise {
    from {
      opacity: 0;
      transform: translateY(16px);
    }

    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
</style>
