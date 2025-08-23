<script lang="ts">
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
        cursor: zoom-out;
    }

    :global(.fullscreen-overlay.active) {
        opacity: 1 !important;
        pointer-events: auto !important;
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
</style>
