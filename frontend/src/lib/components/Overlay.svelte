<script lang="ts">
    let isOverlayActive = $state(false);
    let overlayImageSrc = $state("");

    let overlayElement: HTMLDivElement | null = null;

    // Note the export!!!
    export function showOverlay(src: string) {
        isOverlayActive = true;
        overlayImageSrc = src;
    }

    function closeOverlay() {
        isOverlayActive = false;
    }

    $effect(() => {
        if (isOverlayActive && overlayElement) {
            // to handle ESC key
            overlayElement.focus();
        }
    });
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
    role="button"
>
    <img alt="Fullscreen Image" src={overlayImageSrc} />
</div>
