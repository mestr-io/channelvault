<script lang="ts">
import { closeIcon } from './icons'

let isOverlayActive = $state(false)
let overlayImageSrc = $state('')

let overlayElement: HTMLDivElement | null = null

// Note the export!!!
export function showOverlay(src: string) {
  isOverlayActive = true
  overlayImageSrc = src
}

function closeOverlay() {
  isOverlayActive = false
}

$effect(() => {
  if (isOverlayActive && overlayElement) {
    // to handle ESC key
    overlayElement.focus()
  }
})
</script>

<!-- svelte-ignore a11y_no_noninteractive_tabindex, a11y_no_static_element_interactions, a11y_click_events_have_key_events, a11y_no_noninteractive_element_interactions -->
<div
    class="fullscreen-overlay"
    class:active={isOverlayActive}
    bind:this={overlayElement}
    tabindex="0"
    onclick={closeOverlay}
    onkeyup={(event) => {
        if (event.key === "Escape") closeOverlay();
    }}
>
    <button
        class="close-button"
        aria-label="Close thread"
        onclick={() => closeOverlay()}
    >
        {@render closeIcon("2rem")}
    </button>
    <img alt="Fullscreen Image" src={overlayImageSrc} />
</div>

<style>
    :global(.fullscreen-overlay) {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background-color: rgba(24, 22, 22, 0.9);
        z-index: 1000;
        opacity: 0;
        pointer-events: none; /* Ignore clicks when hidden */
        transition: opacity 0.3s ease;
        display: flex;
        justify-content: center;
        align-items: center;
        cursor: pointer;
    }

    :global(.fullscreen-overlay.active) {
        opacity: 1;
        pointer-events: auto;
    }

    :global(.fullscreen-overlay img) {
        max-height: 90vh;
        max-width: 90vw;
        object-fit: contain;
        transform: scale(0.95);
        transition: transform 0.3s ease;
    }

    :global(.fullscreen-overlay.active img) {
        transform: scale(1);
    }

    .close-button {
        position: absolute;
        top: 1rem;
        right: 1rem;
        width: 2rem;
        height: 2rem;
        padding: 0;
        border: none;
        color: #e6e6e6;
        background: transparent;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        border-radius: 1rem;
        z-index: 10; /* keep it above internal content if needed */
        transition:
            background-color 0.15s ease,
            color 0.15s ease;
    }
    .close-button:hover {
        background-color: rgba(255, 255, 255, 0.2);
        color: #dea1a1;
    }
</style>
