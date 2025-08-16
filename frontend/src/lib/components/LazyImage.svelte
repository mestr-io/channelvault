<script lang="ts">
import { GetBase64Image } from '../../../wailsjs/go/app/App'

let {
  file_id,
  alt,
  delay = 250,
  class: className = '',
  threshold = 0.05, // Load when 5% of the image is visible
} = $props()

let shouldLoad = $state(false)
let container: HTMLDivElement

$effect(() => {
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

<div bind:this={container} class="wrapper">
    {#if shouldLoad}
        {#await GetBase64Image(file_id)}
            <slot name="placeholder">
                <div class="placeholder {className}"></div>
            </slot>
        {:then src}
            <img {src} {alt} class={className} />
        {/await}
    {:else}
        <slot name="placeholder">
            <div class="placeholder {className}"></div>
        </slot>
    {/if}
</div>

<style>
    /* The wrapper must be a layout element for the observer to see it. */
    .wrapper {
        display: block;
        line-height: 0; /* Prevents extra space below the image */
        min-height: 1px; /* Ensures it has a dimension even when empty */
    }

    /* Apply styles to both the placeholder and the final image */
    :global(.wrapper > .placeholder),
    :global(.wrapper > img) {
        display: block;
        height: 300px;
        width: 100%;
        object-fit: scale-down;
    }

    .placeholder {
        /* The initial background color */
        background-color: #f0f0f0;
        height: 300px;
        width: 100%;
        /* Apply the new pulse animation */
        animation: pulse 1.5s cubic-bezier(0.4, 0, 0.6, 1) infinite;
    }

    /* --- Add this keyframes rule --- */
    @keyframes pulse {
        0%,
        100% {
            /* Starting and ending color */
            background-color: #f0f0f0;
        }
        50% {
            /* Mid-point color */
            background-color: #e0e0e0;
        }
    }
</style>
