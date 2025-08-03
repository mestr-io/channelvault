<script lang="ts">
  import FileUploader from "./lib/components/FileUploader.svelte";
  import Dashboard from "./lib/components/Dashboard.svelte";
  import Progress from "./lib/components/Progress.svelte";
  // import {viewState} from "./lib/store"
	import { onMount } from 'svelte';
  import { GetViewState } from '../wailsjs/go/main/App';
  import * as runtime from "../wailsjs/runtime"


  type ViewState = 'setup' | 'dashboard' | 'loading';

  let viewState = $state<ViewState>('loading');

  // onMount is the perfect lifecycle hook to set up listeners.
  // It runs once after the component is first added to the DOM.
  onMount(async () => {
    try {
      // Correctly update the top-level `viewState` variable by direct assignment.
      // Do not redeclare it with `const`.
      viewState = (await GetViewState()) as ViewState;
      console.log(`Initial state set to: ${viewState}`);

      // Set up the event listener.
      // We pass a function that directly reassigns the `viewState`.
      runtime.EventsOn('backend:viewState', (newState: string) => {
        console.log(`Received event from backend. New state: ${newState}`);
        viewState = newState as ViewState;
      });
    } catch (error) {
      console.error('Failed to initialize view state or set up listener:', error);
      // Optional: handle the error gracefully, e.g., set a specific state
      // viewState = 'error';
    }
  });

</script>

<header></header>
<main class="app-main">
  <!-- Routing -->
  {#if viewState === 'loading'}
    <Progress />
  {:else if viewState === 'setup'}
    <FileUploader />
  {:else if viewState === 'dashboard'}
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
