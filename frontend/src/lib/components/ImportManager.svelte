<script lang="ts">
import { onMount } from 'svelte'
import { ListImports, ResumeImport } from '../../../wailsjs/go/app/App'
import { models } from '../../../wailsjs/go/models'
import SlackImport from './SlackImport.svelte'
import * as runtime from '../../../wailsjs/runtime/runtime'
import { DeleteImport } from '../../../wailsjs/go/app/App'

let isDragging = $state(false)
let dragCounter = $state(0)
let imports: models.SlackImport[] = $state([])
let isImporting = $state(false)
let isDropZoneCollapsed = $state(true)

onMount(async () => {
  // subscribe to new import notification
  runtime.EventsOn('startImport', async (imp: models.SlackImport) => {
    const fresh = await ListImports()
    imports = [...fresh]
    isImporting = true
  })

  // reload the
  runtime.EventsOn('importComplete', async () => {
    const fresh = await ListImports()
    imports = [...fresh]
    isImporting = false
  })

  runtime.EventsOn('importError', async () => {
    const fresh = await ListImports()
    imports = [...fresh]
    isImporting = false
  })

  // load imports
  const fresh = await ListImports()
  imports = [...fresh]

  if (imports.length === 0) {
    isDropZoneCollapsed = false
  }
})

let DeleteImportHandle = async (importID: number) => {
  try {
    await DeleteImport(importID)
    const fresh = await ListImports()
    imports = [...fresh]
  } catch (error) {}
}

let RetryImportHandle = async (importID: number) => {
  await ResumeImport(importID)
}

function toggleDropZone() {
  isDropZoneCollapsed = !isDropZoneCollapsed
}

// --- Drag & Drop Event Handlers ---
function handleDragEnter(event: DragEvent): void {
  event.preventDefault()
  dragCounter++
  if (dragCounter === 1) {
    isDragging = true
    // Auto-expand when dragging over collapsed zone
    if (isDropZoneCollapsed) {
      isDropZoneCollapsed = false
    }
  }
}

function handleDragOver(event: DragEvent): void {
  event.preventDefault()
  if (!isDragging) {
    isDragging = true
    // Auto-expand when dragging over collapsed zone
    if (isDropZoneCollapsed) {
      isDropZoneCollapsed = false
    }
  }
}

function handleDragLeave(event: DragEvent): void {
  event.preventDefault()
  dragCounter--
  if (dragCounter === 0) {
    isDragging = false
  }
}

function handleDrop(event: DragEvent): void {
  event.preventDefault()
  if (isImporting) return

  dragCounter = 0
  isDragging = false
}
</script>

<main>
  <div class="drop-zone-container">
    <button 
      class="collapse-toggle" 
      onclick={toggleDropZone}
      aria-label={isDropZoneCollapsed ? "Expand drop zone" : "Collapse drop zone"}
    >
      <span>Add New Import</span>
       <svg 
        class="chevron" 
        class:collapsed={isDropZoneCollapsed}
        xmlns="http://www.w3.org/2000/svg" 
        width="20" 
        height="20" 
        viewBox="0 0 24 24" 
        fill="none" 
        stroke="currentColor" 
        stroke-width="2.5" 
        stroke-linecap="round" 
        stroke-linejoin="round"
      >
        <path d="m6 9 6 6 6-6"/>
      </svg>
    </button>
    
    <div
      class="drop-zone"
      class:is-dragging={isDragging}
      class:is-disabled={isImporting}
      class:is-collapsed={isDropZoneCollapsed}
      role="button"
      tabindex={isDropZoneCollapsed ? -1 : 0}
      ondragenter={handleDragEnter}
      ondragleave={handleDragLeave}
      ondragover={handleDragOver}
      ondrop={handleDrop}
    >
      <div class="drop-zone-content">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          x="0px"
          y="0px"
          width="80"
          height="80"
          viewBox="0 0 48 48"
        >
           <path
            fill="#33d375"
            d="M33,8c0-2.209-1.791-4-4-4s-4,1.791-4,4c0,1.254,0,9.741,0,11c0,2.209,1.791,4,4,4s4-1.791,4-4 C33,17.741,33,9.254,33,8z"
          ></path><path
            fill="#33d375"
            d="M43,19c0,2.209-1.791,4-4,4c-1.195,0-4,0-4,0s0-2.986,0-4c0-2.209,1.791-4,4-4S43,16.791,43,19z"
          ></path><path
            fill="#40c4ff"
            d="M8,14c-2.209,0-4,1.791-4,4s1.791,4,4,4c1.254,0,9.741,0,11,0c2.209,0,4-1.791,4-4s-1.791-4-4-4 C17.741,14,9.254,14,8,14z"
          ></path><path
            fill="#40c4ff"
            d="M19,4c2.209,0,4,1.791,4,4c0,1.195,0,4,0,4s-2.986,0-4,0c-2.209,0-4-1.791-4-4S16.791,4,19,4z"
          ></path><path
            fill="#e91e63"
            d="M14,39.006C14,41.212,15.791,43,18,43s4-1.788,4-3.994c0-1.252,0-9.727,0-10.984  c0-2.206-1.791-3.994-4-3.994s-4,1.788-4,3.994C14,29.279,14,37.754,14,39.006z"
          ></path><path
            fill="#e91e63"
            d="M4,28.022c0-2.206,1.791-3.994,4-3.994c1.195,0,4,0,4,0s0,2.981,0,3.994c0,2.206-1.791,3.994-4,3.994  S4,30.228,4,28.022z"
          ></path><path
            fill="#ffc107"
            d="M39,33c2.209,0,4-1.791,4-4s-1.791-4-4-4c-1.254,0-9.741,0-11,0c-2.209,0-4,1.791-4,4s1.791,4,4,4  C29.258,33,37.746,33,39,33z"
          ></path><path
            fill="#ffc107"
            d="M28,43c-2.209,0-4-1.791-4-4c0-1.195,0-4,0-4s2.986,0,4,0c2.209,0,4,1.791,4,4S30.209,43,28,43z"
          ></path>
        </svg>
        <div class="drop-zone-text">Drop your Slack export `.zip` file here</div>
        <div class="drop-zone-subtext">
          Learn more about the 
          <a
            href="#"
            onclick={() =>
              runtime.BrowserOpenURL(
                "https://slack.com/intl/en-gb/help/articles/201658943-Export-your-workspace-data",
              )}>export process</a
          >.
        </div>
        <div class="drop-zone-disclaimer">*Only public channels are supported for now</div>
      </div>
    </div>
  </div>

  <div class="imports-list">
    {#if imports.length > 0}
      <h2 class="imports-title">Your Imports</h2>
    {/if}
    {#each imports as slackImportcard (slackImportcard.ID)}
      <SlackImport {RetryImportHandle} {slackImportcard} isCollapsed={slackImportcard.Finished !== "" && slackImportcard.Error === ""} {DeleteImportHandle}/>
    {/each}
  </div>
</main>

<style>
  main {
    margin: 0;
    color: #d1d2d3; /* Slack primary text */
    display: flex;
    flex-direction: column;
    align-items: center;
    min-height: 100vh;
    background-color: #1a1d21; /* Slack content background */
    overflow-y: auto;
    padding: 2rem 1rem;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  }

  .imports-list {
    margin-top: 2rem;
    margin-bottom: 3rem;
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    width: 100%;
    max-width: 800px;
  }

  .imports-title {
    color: #86898e; /* Slack muted text */
    font-size: 0.875rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    padding-bottom: 0.5rem;
    border-bottom: 1px solid #2c2f33;
    margin-bottom: 1rem;
  }

  .drop-zone-container {
    margin-top: 2rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    width: 100%;
    max-width: 800px;
  }

  .collapse-toggle {
    background: transparent;
    border: none;
    border-radius: 0.5rem;
    padding: 0.5rem 1rem;
    color: #86898e;
    font-size: 1rem;
    font-weight: 500;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    transition: color 200ms ease-in-out;
  }

  .collapse-toggle:hover {
    color: #d1d2d3;
  }

  .chevron {
    transition: transform 300ms cubic-bezier(0.4, 0, 0.2, 1);
    color: #86898e;
  }

  .chevron.collapsed {
    transform: rotate(-90deg);
  }

  .drop-zone {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 2rem;
    border-radius: 0.75rem;
    transition: all 300ms cubic-bezier(0.4, 0, 0.2, 1);
    border: 1px solid #35373B;
    background-color: #22252A; /* Slightly lighter than main bg */
    width: 100%;
    min-height: 250px;
    position: relative;
    overflow: hidden;
    opacity: 1;
    max-height: 400px;
  }

  .drop-zone.is-collapsed {
    max-height: 0; min-height: 0; padding-top: 0; padding-bottom: 0;
    opacity: 0; pointer-events: none; border-color: transparent; margin: 0;
  }
  
  .drop-zone.is-disabled {
    opacity: 0.4; pointer-events: none; background-color: #2a2e34;
  }

  .drop-zone.is-dragging {
    border-color: #1264a3; /* Slack accent blue */
    background-color: #2a2f35;
    transform: scale(1.02);
    box-shadow: 0 0 25px rgba(18, 100, 163, 0.2);
  }

  .drop-zone-content {
    display: flex; flex-direction: column; align-items: center; gap: 0.75rem; text-align: center;
  }
  
  .drop-zone-text {
    font-size: 1.125rem; font-weight: 500; color: #f8f8f8; margin-top: 0.5rem;
  }
  
  .drop-zone-subtext,
  .drop-zone-disclaimer {
    font-size: 0.875rem; color: #86898e;
  }

  a {
    color: #1264a3; /* Slack accent blue */
    text-decoration: none; cursor: pointer; font-weight: 500;
  }
  
  a:hover {
    text-decoration: underline;
  }
</style>
