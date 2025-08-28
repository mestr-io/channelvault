<script lang="ts">
import * as runtime from '../wailsjs/runtime'
import { onMount, setContext } from 'svelte'
import { GetViewState } from '../wailsjs/go/app/App'

import Dashboard from './lib/components/Dashboard.svelte'
import Overlay from './lib/components/Overlay.svelte'

type ViewState = 'setup' | 'dashboard' | 'loading'

let overlayComponent: Overlay
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

function openOverlay(src: string) {
  if (overlayComponent) {
    overlayComponent.showOverlay(src)
  }
}

setContext('overlay', { openOverlay })
</script>

<Overlay bind:this={overlayComponent} />
<main class="app-main">
  <header></header>
    <Dashboard />
  <footer></footer>
</main>

<style>
  /* Inter font from Google Fonts (or similar system font for local development) */
  :global(body) {
    font-family:
      "Lato",
      -apple-system,
      BlinkMacSystemFont,
      "Segoe UI",
      Roboto,
      Helvetica,
      Arial,
      sans-serif;
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  .app-main {
    background-color: transparent;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    overflow: hidden;
  }
</style>
