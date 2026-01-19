<script>
  import { onMount } from 'svelte';
  import { shortenURL } from '../lib/api';
  
  let originalUrl = '';
  let shortUrl = '';
  let shortCode = '';
  let loading = false;
  let error = '';
  let success = false;

  async function handleShorten() {
    if (!originalUrl.trim()) {
      error = 'Please enter a URL';
      return;
    }

    loading = true;
    error = '';
    success = false;

    try {
      const response = await shortenURL(originalUrl);
      const data = response.data;
      
      shortCode = data.short_code;
      shortUrl = data.short_url;
      success = true;
      originalUrl = '';
    } catch (err) {
      error = err.response?.data?.error || 'Failed to shorten URL';
    } finally {
      loading = false;
    }
  }

  function copyToClipboard() {
    navigator.clipboard.writeText(shortUrl);
    alert('Copied to clipboard!');
  }

  function handleKeypress(e) {
    if (e.key === 'Enter') {
      handleShorten();
    }
  }
</script>

<main>
  <div class="container">
    <h1>🔗 URL Shortener</h1>
    
    <div class="form-container">
      <input
        type="text"
        placeholder="Enter your URL"
        bind:value={originalUrl}
        on:keypress={handleKeypress}
        disabled={loading}
      />
      <button on:click={handleShorten} disabled={loading}>
        {loading ? 'Shortening...' : 'Shorten'}
      </button>
    </div>

    {#if error}
      <div class="error">{error}</div>
    {/if}

    {#if success}
      <div class="result">
        <p><strong>Short URL:</strong></p>
        <div class="url-display">
          <input type="text" value={shortUrl} readonly />
          <button on:click={copyToClipboard}>Copy</button>
        </div>
        <p style="margin-top: 10px; font-size: 0.9em; color: #666;">
          Short code: <code>{shortCode}</code>
        </p>
      </div>
    {/if}
  </div>
</main>

<style>
  main {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    padding: 20px;
  }

  .container {
    background: white;
    border-radius: 10px;
    padding: 40px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
    max-width: 500px;
    width: 100%;
  }

  h1 {
    text-align: center;
    color: #333;
    margin-bottom: 30px;
  }

  .form-container {
    display: flex;
    gap: 10px;
    margin-bottom: 20px;
  }

  input[type="text"] {
    flex: 1;
    padding: 12px 15px;
    border: 1px solid #ddd;
    border-radius: 5px;
    font-size: 14px;
    transition: border-color 0.3s;
  }

  input[type="text"]:focus {
    outline: none;
    border-color: #667eea;
  }

  button {
    padding: 12px 24px;
    background: #667eea;
    color: white;
    border: none;
    border-radius: 5px;
    font-weight: bold;
    cursor: pointer;
    transition: background 0.3s;
  }

  button:hover:not(:disabled) {
    background: #5568d3;
  }

  button:disabled {
    background: #ccc;
    cursor: not-allowed;
  }

  .error {
    background: #fee;
    color: #c33;
    padding: 12px;
    border-radius: 5px;
    margin-bottom: 15px;
    border-left: 4px solid #c33;
  }

  .result {
    background: #f0f9ff;
    padding: 20px;
    border-radius: 5px;
    border-left: 4px solid #667eea;
  }

  .result p {
    margin: 10px 0;
    color: #333;
  }

  .url-display {
    display: flex;
    gap: 10px;
  }

  .url-display input {
    flex: 1;
  }

  .url-display button {
    padding: 12px 16px;
  }

  code {
    background: #f5f5f5;
    padding: 2px 6px;
    border-radius: 3px;
    font-family: monospace;
  }
</style>
