<script lang="ts">
  import FileUploader from "./lib/components/FileUploader.svelte";
  import Dashboard from "./lib/components/Dashboard.svelte";
  import Progress from "./lib/components/Progress.svelte";
  import {viewState} from "./lib/store"
	import { onMount } from 'svelte';
  import { GetViewState } from '../wailsjs/go/app/App';
  import * as runtime from "../wailsjs/runtime"

  // onMount is the perfect lifecycle hook to set up listeners.
  // It runs once after the component is first added to the DOM.
  onMount(async () => {
      const initialState = await GetViewState();
      viewState.set(initialState as 'setup' | 'dashboard' | 'loading');
      console.log(`Initial state set to: ${initialState}`);
      runtime.EventsOn('backend:viewState', (newState) => {
        console.log(`Received event from backend. New state: ${newState}`);
        viewState.set(newState as 'setup' | 'dashboard' | 'loading');
      });
  });

</script>

<header></header>
<main class="app-main">
  <!-- Routing -->
  {#if $viewState === 'loading'}
    <Progress />
  {:else if $viewState === 'setup'}
    <FileUploader />
  {:else if $viewState === 'dashboard'}
    <Dashboard />
  {/if}
</main>
<footer></footer>

<style>
   /* Inter font from Google Fonts (or similar system font for local development) */
  :global(body) {
    font-family: 'Inter', sans-serif;
    margin: 0;
    padding: 0;
    box-sizing: border-box;
    overflow: hidden; /* Added to remove scrollbars */
  }

  .app-main {
    width: 100vw;
    height: 100vh;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    background-color: #1a202c;
    color: #f7fafc;
    padding: 1rem;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  }

</style>
