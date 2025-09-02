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

  runtime.EventsOn('downloadBinariesDone', async _ => {
    const fresh = await ListImports()
    imports = [...fresh]
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
      aria-label={isDropZoneCollapsed
        ? "Expand drop zone"
        : "Collapse drop zone"}
    >
      <svg
        class="chevron"
        class:collapsed={isDropZoneCollapsed}
        xmlns="http://www.w3.org/2000/svg"
        width="20"
        height="20"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path d="m6 9 6 6 6-6" />
      </svg>
      <span>Add import</span>
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
      <div class="icons-container">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          x="0px"
          y="0px"
          width="100"
          height="100"
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
        <div>Drop here your Slack export zip file</div>
        <div>
          Check this <a
            href="#"
            onclick={() =>
              runtime.BrowserOpenURL(
                "https://slack.com/intl/en-gb/help/articles/201658943-Export-your-workspace-data",
              )}>link</a
          > for more information about the export process.
        </div>
        <div>
          <small>*Only public channels are supported for the moment</small>
        </div>
      </div>
    </div>
  </div>

  <div class="imports">
    {#each imports as slackImportcard (slackImportcard.ID)}
      <SlackImport
        {RetryImportHandle}
        {slackImportcard}
        isCollapsed={slackImportcard.Finished !== "" &&
          slackImportcard.Error === ""}
        {DeleteImportHandle}
      />
    {/each}
  </div>
</main>

<style>
  main {
    margin: 0;
    color: #e2e8f0;
    display: flex;
    flex-direction: column;
    place-items: center;
    min-height: 100vh;
    background-color: #552958;
    overflow: auto;
    padding: 0 1rem;
  }

  .imports {
    margin-top: 3rem;
    margin-bottom: 3rem;
    display: flex;
    flex-direction: column;
    gap: 2rem;
    width: 100%;
    max-width: 1200px;
  }

  .drop-zone-container {
    margin-top: 1rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
  }

  .collapse-toggle {
    background: rgba(45, 55, 72, 0.8);
    border: 1px solid #4a5568;
    border-radius: 0.5rem;
    padding: 0.5rem 1rem;
    color: #e2e8f0;
    font-size: 0.875rem;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    transition: all 200ms ease-in-out;
    backdrop-filter: blur(4px);
  }

  .collapse-toggle:hover {
    background: rgba(45, 55, 72, 1);
    border-color: #718096;
    transform: translateY(-1px);
  }

  .collapse-toggle:active {
    transform: translateY(0);
  }

  .chevron {
    transition: transform 200ms ease-in-out;
  }

  .chevron.collapsed {
    transform: rotate(-90deg);
  }

  .drop-zone {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 2rem;
    border-radius: 0.75rem;
    box-shadow:
      0 10px 15px -3px rgba(0, 0, 0, 0.1),
      0 4px 6px -2px rgba(0, 0, 0, 0.05);
    transition: all 300ms ease-in-out;
    border: 2px dashed #6b7280;
    background-color: #2d3748;
    min-width: 350px;
    min-height: 220px;
    position: relative;
    overflow: hidden;
    opacity: 1;
    max-height: 400px;
  }

  .drop-zone.is-collapsed {
    min-height: 0;
    max-height: 0;
    padding: 0 2rem;
    opacity: 0;
    pointer-events: none;
    margin: 0;
  }

  .drop-zone.is-disabled {
    opacity: 0.5;
    pointer-events: none;
    border-color: #4b5563;
    background-color: #1f2937;
  }

  .drop-zone.is-collapsed.is-disabled {
    opacity: 0;
  }

  .drop-zone.is-dragging {
    border-color: #3b82f6;
    box-shadow:
      0 0 30px rgba(59, 130, 246, 0.7),
      0 0 15px rgba(59, 130, 246, 0.5);
    transform: scale(1.05);
    opacity: 1;
    max-height: 400px;
    min-height: 220px;
    padding: 2rem;
    pointer-events: all;
  }

  .icons-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    margin-bottom: 1rem;
  }

  /* Responsive adjustments */
  @media (max-width: 768px) {
    .imports {
      grid-template-columns: 1fr;
      gap: 1.5rem;
    }

    .drop-zone {
      min-width: 300px;
    }
  }

  a {
    color: #768998;
    text-decoration: none;
    cursor: pointer;
  }
</style>
