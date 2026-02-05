<script>
  import Page from './routes/+page.svelte';
</script>

<Page />

<style global>
  :global(body) {
    margin: 0;
    min-height: 100vh;
  }

  :global(*) {
    box-sizing: border-box;
  }
</style>
