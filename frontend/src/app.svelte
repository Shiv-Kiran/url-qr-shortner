<script>
  import Page from './routes/+page.svelte';
</script>

<Page />

<style global>
  :global(:root) {
    color-scheme: light;
    --bg: 230 26% 96%;
    --bg-soft: 230 24% 94%;
    --surface: 0 0% 100%;
    --surface-2: 220 18% 92%;
    --foreground: 220 20% 12%;
    --muted: 220 12% 45%;
    --border: 220 12% 80%;
    --accent-1: 200 100% 60%;
    --accent-2: 265 86% 68%;
    --accent-3: 160 72% 55%;
    --glow: 200 95% 62%;
    --shadow: 0 20px 50px rgba(20, 24, 34, 0.18);
  }

  :global(.dark) {
    color-scheme: dark;
    --bg: 220 40% 8%;
    --bg-soft: 220 36% 10%;
    --surface: 220 28% 14%;
    --surface-2: 220 26% 18%;
    --foreground: 220 20% 96%;
    --muted: 220 14% 70%;
    --border: 220 14% 26%;
    --accent-1: 196 100% 70%;
    --accent-2: 265 100% 75%;
    --accent-3: 165 86% 62%;
    --glow: 196 100% 70%;
    --shadow: 0 24px 70px rgba(8, 12, 22, 0.6);
  }

  :global(body) {
    margin: 0;
    min-height: 100vh;
    font-family: 'Space Grotesk', 'Segoe UI', system-ui, sans-serif;
    background: hsl(var(--bg));
    color: hsl(var(--foreground));
  }

  :global(*) {
    box-sizing: border-box;
  }

  :global(a) {
    color: inherit;
    text-decoration: none;
  }
</style>
