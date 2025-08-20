<script lang="ts">
  import { GetBase64Image } from "../../../wailsjs/go/app/App";
  import { onMount } from "svelte";
  let { fileId, alt, url } = $props();

  let shouldLoad = $state(false);
  let container: HTMLDivElement;

  const delay = 250;
  const threshold = 0.05; // Load when 5% of the image is visible

  let overlayActive = $state(true);
  let overlaySrc = $state("");

  const openOverlay = (src: string) => {
    overlaySrc = src;
    overlayActive = true;
  };

  function closeOverlay() {
    overlayActive = false;
    overlaySrc = "";
  }

  onMount(() => {
    // Ensure the container element is mounted before proceeding
    if (!container) return;

    // A variable to hold the timer ID
    let timerId: number | null = null;

    const observer = new IntersectionObserver(
      (entries) => {
        const entry = entries[0];
        if (entry.isIntersecting) {
          // When the element enters the viewport, start a timer.
          // Check for timerId to avoid setting multiple timers.
          if (!timerId) {
            timerId = setTimeout(() => {
              shouldLoad = true;
              // Once we've decided to load, we no longer need the observer.
              observer.disconnect();
            }, delay);
          }
        } else {
          // If the element leaves the viewport before the timer finishes, cancel it.
          if (timerId) {
            clearTimeout(timerId);
            timerId = null;
          }
        }
      },
      {
        // A small threshold makes intersection detection more reliable.
        threshold: threshold,
      },
    );

    observer.observe(container);

    // Cleanup: disconnect the observer and clear any pending timers
    return () => {
      observer.disconnect();
      if (timerId) {
        clearTimeout(timerId);
      }
    };
  });
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
              console.log(overlayActive);
              openOverlay(src);
            }}
          />
        {/await}
      {/if}
    </div>
  </div>
  {#if overlayActive}
    <div
      class="fullscreen-overlay"
      tabindex="0"
      onkeyup={(event) => {
        if (event.key === "Escape") closeOverlay();
      }}
      onclick={closeOverlay}
    >
      <img src={overlaySrc} {alt} />
    </div>
  {/if}
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
    max-width: 100%;
    height: 300px;
    object-fit: scale-down;
  }

  .fullscreen-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.9);
    z-index: 1000;
    opacity: 0;
    pointer-events: none; /* Ignore clicks when hidden */
    transition: opacity 0.3s ease;
    display: flex;
    justify-content: center;
    align-items: center;
    cursor: zoom-out;
  }

  .fullscreen-overlay.active {
    opacity: 1;
    pointer-events: auto; /* Accept clicks again */
  }

  .fullscreen-overlay img {
    max-height: 90vh;
    max-width: 90vw;
    object-fit: contain;
    transform: scale(0.95);
    transition: transform 0.3s ease;
  }

  .fullscreen-overlay.active img {
    transform: scale(1);
  }
</style>
