<script lang="ts">
import FileUploader from './lib/components/FileUploader.svelte'
import Dashboard from './lib/components/Dashboard.svelte'
import Progress from './lib/components/Progress.svelte'
import { onMount } from 'svelte'
import { GetViewState } from '../wailsjs/go/app/App'
import * as runtime from '../wailsjs/runtime'

type ViewState = 'setup' | 'dashboard' | 'loading'

let viewState = $state<ViewState>('loading')

// onMount is the perfect lifecycle hook to set up listeners.
// It runs once after the component is first added to the DOM.
onMount(async () => {
  try {
    viewState = (await GetViewState()) as ViewState
    console.log(`Initial state set to: ${viewState}`)

    // We pass a function that directly reassigns the `viewState`.
    runtime.EventsOn('backend:viewState', (newState: string) => {
      console.log(`Received event from backend. New state: ${newState}`)
      viewState = newState as ViewState
    })
  } catch (error) {
    console.error('Failed to initialize view state or set up listener:', error)
  }
})
</script>

<main class="app-main">
  <header></header>
  <!-- Routing -->
  {#if viewState === "loading"}
    <Progress />
  {:else if viewState === "setup"}
    <FileUploader />
  {:else if viewState === "dashboard"}
    <Dashboard />
  {/if}
  <footer></footer>
</main>

<style>
  /* Inter font from Google Fonts (or similar system font for local development) */
  :global(body) {
    font-family: "Inter", sans-serif;
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  .app-main {
    padding: 1rem;
    background-color: #1a202c;
    color: #f7fafc;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    overflow: hidden;
  }
  
</style>
