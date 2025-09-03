<script lang="ts">
import { GetBase64Image } from '../../../../wailsjs/go/app/App'

import { onMount, getContext } from 'svelte'
let { fileId, alt, url } = $props()

let shouldLoad = $state(false)
let container: HTMLDivElement

const delay = 100
const threshold = 0.05 // Load when 5% of the image is visible
const { openOverlay } = getContext<{ openOverlay: (src: string) => void }>('overlay')

onMount(() => {
  // Ensure the container element is mounted before proceeding
  if (!container) return

  // A variable to hold the timer ID
  let timerId: number | null = null

  const observer = new IntersectionObserver(
    entries => {
      const entry = entries[0]
      if (entry.isIntersecting) {
        // When the element enters the viewport, start a timer.
        // Check for timerId to avoid setting multiple timers.
        if (!timerId) {
          timerId = setTimeout(() => {
            shouldLoad = true
            // Once we've decided to load, we no longer need the observer.
            observer.disconnect()
          }, delay)
        }
      } else {
        // If the element leaves the viewport before the timer finishes, cancel it.
        if (timerId) {
          clearTimeout(timerId)
          timerId = null
        }
      }
    },
    {
      // A small threshold makes intersection detection more reliable.
      threshold: threshold,
    },
  )

  observer.observe(container)

  // Cleanup: disconnect the observer and clear any pending timers
  return () => {
    observer.disconnect()
    if (timerId) {
      clearTimeout(timerId)
    }
  }
})
</script>

<div>
  <div bind:this={container} class="wrapper">
    <div class="placeholder">
      {#if shouldLoad}
        {#await GetBase64Image(fileId) then src}
          <img
            {src}
            {alt}
            onclick={() => {
              openOverlay(src);
            }}
          />
        {/await}
      {/if}
    </div>
  </div>
</div>

<style>
  /* The wrapper must be a layout element for the observer to see it. */
  .wrapper {
    display: flex;
  }

  .placeholder {
    display: block;
    height: 300px;
    width: 100%;
  }

  .placeholder > img {
    display: block;
    max-width: 90%;
    height: 300px;
    object-fit: contain;
  }
</style>
