<script lang="ts">
  import FileUploader from "./lib/components/FileUploader.svelte";
  import Dashboard from "./lib/components/Dashboard.svelte";
  import Progress from "./lib/components/Progress.svelte";
  import { onMount, setContext } from "svelte";
  import { GetViewState } from "../wailsjs/go/app/App";
  import * as runtime from "../wailsjs/runtime";

  type ViewState = "setup" | "dashboard" | "loading";

  let viewState = $state<ViewState>("loading");

  let overlayElement: HTMLDivElement | null = null;

  const openOverlay = (src: string) => {
    console.log("Open");
    console.log(overlayElement.classList);
    if (overlayElement) {
      overlayElement.classList.add("active");
      overlayElement.querySelector("img")!.src = src; // Set the image source
    }
  };

  const closeOverlay = () => {
    console.log("CLose");
    if (overlayElement) {
      overlayElement.classList.remove("active");
      overlayElement.querySelector("img")!.src = ""; // Clear the image source
    }
  };

  // onMount is the perfect lifecycle hook to set up listeners.
  // It runs once after the component is first added to the DOM.
  onMount(async () => {
    try {
      viewState = (await GetViewState()) as ViewState;
      console.log(`Initial state set to: ${viewState}`);

      // We pass a function that directly reassigns the `viewState`.
      runtime.EventsOn("backend:viewState", (newState: string) => {
        console.log(`Received event from backend. New state: ${newState}`);
        viewState = newState as ViewState;
      });
    } catch (error) {
      console.error(
        "Failed to initialize view state or set up listener:",
        error,
      );
    }
  });
  setContext("overlay", { openOverlay, closeOverlay });
</script>

<!-- Overlay -->
<div
  class="fullscreen-overlay"
  bind:this={overlayElement}
  tabindex="0"
  onkeyup={(event) => {
    if (event.key === "Escape") closeOverlay();
  }}
  onclick={closeOverlay}
>
  <img alt="Fullscreen Image" />
</div>
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
