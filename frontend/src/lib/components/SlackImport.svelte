<script lang="ts">
import { onMount } from 'svelte'
import * as runtime from '../../../wailsjs/runtime'
import { models } from '../../../wailsjs/go/models'
import { RenameImport, DownloadBinaries } from '../../../wailsjs/go/app/App'
import Prompt from './Prompt.svelte'

// import { chevron } from './Icons.svelte'
import { chevron } from './icons/chevron.svelte'

let {
  slackImportCard,
  DeleteImportHandle,
  RetryImportHandle,
}: {
  slackImportCard: models.SlackImport
  DeleteImportHandle: (id: number) => Promise<void>
  RetryImportHandle: (id: number) => Promise<void>
} = $props()

// State variables
let progress: models.Progress | null = $state(null)
let pendingBinaries = $state(slackImportCard.PendingBinaries)
let isDownloading: boolean = $state(false)
let showPrompt = $state(false)
let isCollapsed = $state(true)
let promptDefault = ''
let importId: number

function legacyIsCollapsed(importState: models.SlackImport): boolean {
  // Legacy mode assings isCollapsed to =>  slackImportCard.Finished !== "" && slackImportCard.Error === ""

  return importState.Finished === '' && importState.Error === ''
}

onMount(() => {
  // --- Event Listeners ---
  const unsubImportProgress = runtime.EventsOn('importProgress', (p: models.Progress) => {
    if (slackImportCard.ID == p.ImportID) {
      progress = p
    }
  })

  const unsubDownloadDone = runtime.EventsOn('downloadBinariesDone', (importID: number) => {
    if (slackImportCard.ID == importID) {
      isDownloading = false
      pendingBinaries = 0
    }
  })

  const unsubDownloadStarted = runtime.EventsOn('downloadBinariesStarted', (importID: number) => {
    if (slackImportCard.ID == importID) {
      isDownloading = true
    }
  })

  const unsubDownloadedBinary = runtime.EventsOn('downloadedBinary', (importID: number, pending: number) => {
    if (slackImportCard.ID == importID) {
      pendingBinaries = pending
    }
  })

  // Check initial download state on mount
  if (slackImportCard.DownloadStarted.Valid && !slackImportCard.DownloadFinished.Valid) {
    isDownloading = true
  }

  // --- Cleanup function ---
  // This is crucial to prevent memory leaks
  return () => {
    unsubImportProgress()
    unsubDownloadDone()
    unsubDownloadStarted()
    unsubDownloadedBinary()
  }
})

// --- Helper Functions ---
function formatDate(dateString: string): string {
  if (!dateString) return 'Not started'
  const date = new Date(dateString)
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

function getStatusColor(progress: number, error: string): string {
  if (error) return '#DC2626' // red
  if (progress === 0) return '#6B7280' // gray
  if (progress < 50) return '#D97706' // amber
  if (progress < 100) return '#2563EB' // blue
  if (isDownloading) return '#D97706'
  return '#059669' // green
}

function getStatusText(progress: number, finished: string, error: string): string {
  if (error) return 'Error'
  if (finished && progress < 100) return 'Completed'
  if (progress === 0) return 'Pending'
  if (progress < 100) return 'In Progress'
  if (isDownloading) return 'Downloading binaries'
  return 'Completed'
}

// --- Action Handlers ---
async function handleRename(importId_: number, currentName: string) {
  importId = importId_
  promptDefault = currentName
  showPrompt = true
}

async function confirmRename(newName: string) {
  if (!newName) return
  await RenameImport(importId, newName)
  slackImportCard.Name = newName
  showPrompt = false
}

async function handleDelete(importId: number) {
  await DeleteImportHandle(importId)
}

async function handleDownload(importId: number) {
  await DownloadBinaries(importId, true)
}

async function handleRetry(importId: number) {
  await RetryImportHandle(importId)
}

function toggleCollapse() {
  isCollapsed = !isCollapsed
}
</script>

<Prompt
    open={showPrompt}
    title="Rename Import"
    message="Enter new name:"
    defaultValue={promptDefault}
    onConfirm={confirmRename}
    onCancel={() => (showPrompt = false)}
/>
<main class="card">
    <div class="card-header">
        <button
            aria-label="collapse card"
            class="chevron-button"
            onclick={toggleCollapse}
            title={isCollapsed ? "Expand" : "Collapse"}
        >
            {@render chevron(!isCollapsed)}
        </button>

        <div class="card-title">
            {slackImportCard.Name || "Untitled Import"}
        </div>

        <div class="card-actions">
            <div
                class="status-badge"
                style="background-color: {getStatusColor(
                    progress?.Progress ?? slackImportCard.Progress,
                    slackImportCard.Error,
                )}"
            >
                {getStatusText(
                    progress?.Progress ?? slackImportCard.Progress,
                    slackImportCard.Finished,
                    slackImportCard.Error,
                )}
            </div>
            <button
                aria-label="rename import"
                class="action-button rename-button"
                onclick={() =>
                    handleRename(slackImportCard.ID, slackImportCard.Name)}
                title="Rename import"
            >
                <svg viewBox="0 0 20 20" fill="currentColor">
                    <path
                        d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z"
                    />
                </svg>
            </button>

            {#if slackImportCard.Finished && !slackImportCard.DownloadFinished.Valid}
                <button
                    aria-label="download binaries"
                    class="action-button download-button"
                    onclick={() => handleDownload(slackImportCard.ID)}
                    title="Download binaries"
                    disabled={isDownloading}
                >
                    {#if isDownloading}
                        <svg
                            class="spinner"
                            width="24"
                            height="24"
                            viewBox="0 0 24 24"
                            fill="none"
                            xmlns="http://www.w3.org/2000/svg"
                        >
                            <path
                                d="M12 2.99988V5.99988M12 18V21M4.92896 4.92883L7.05028 7.05016M16.9497 16.9497L19.071 19.071M2.99994 12H5.99994M18 12H21M4.92896 19.071L7.05028 16.9497M16.9497 7.05016L19.071 4.92883"
                                stroke="currentColor"
                                stroke-width="2"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            />
                        </svg>
                    {:else}
                        <svg
                            width="24"
                            height="24"
                            viewBox="0 0 24 24"
                            fill="none"
                            xmlns="http://www.w3.org/2000/svg"
                        >
                            <path
                                d="M11 5C11 4.44772 11.4477 4 12 4C12.5523 4 13 4.44772 13 5V12.1578L16.2428 8.91501L17.657 10.3292L12.0001 15.9861L6.34326 10.3292L7.75748 8.91501L11 12.1575V5Z"
                                fill="currentColor"
                            />
                            <path
                                d="M4 14H6V18H18V14H20V18C20 19.1046 19.1046 20 18 20H6C4.89543 20 4 19.1046 4 18V14Z"
                                fill="currentColor"
                            />
                        </svg>
                    {/if}
                </button>
            {/if}
            <button
                aria-label="delete import"
                class="action-button delete-button"
                onclick={() => handleDelete(slackImportCard.ID)}
                title="Delete import"
            >
                <svg viewBox="0 0 20 20" fill="currentColor">
                    <path
                        fill-rule="evenodd"
                        d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z"
                        clip-rule="evenodd"
                    />
                </svg>
            </button>
        </div>
    </div>

    <div class="collapsible {isCollapsed ? 'collapsed' : ''}">
        {#if slackImportCard.Error}
            <div class="error-section">
                <div class="error-header">
                    <svg
                        class="error-icon"
                        viewBox="0 0 20 20"
                        fill="currentColor"
                    >
                        <path
                            fill-rule="evenodd"
                            d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z"
                            clip-rule="evenodd"
                        />
                    </svg>
                    <span class="error-title">Import Error</span>
                </div>
                <p class="error-message">{slackImportCard.Error}</p>
                <div class="error-actions">
                    <button
                        class="retry-button"
                        onclick={() => handleRetry(slackImportCard.ID)}
                        title="Retry import"
                    >
                        <svg viewBox="0 0 20 20" fill="currentColor">
                            <path
                                fill-rule="evenodd"
                                d="M4 2a1 1 0 011 1v2.101a7.002 7.002 0 0111.601 2.566 1 1 0 11-1.885.666A5.002 5.002 0 005.999 7H9a1 1 0 010 2H4a1 1 0 01-1-1V3a1 1 0 011-1zm.008 9.057a1 1 0 011.276.61A5.002 5.002 0 0014.001 13H11a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0v-2.101a7.002 7.002 0 01-11.601-2.566 1 1 0 01.61-1.276z"
                                clip-rule="evenodd"
                            />
                        </svg>
                        Retry Import
                    </button>
                </div>
            </div>
        {/if}

        <div class="card-content">
            <div class="info-row">
                <div class="info-item">
                    <svg
                        class="info-icon"
                        viewBox="0 0 20 20"
                        fill="currentColor"
                    >
                        <path
                            fill-rule="evenodd"
                            d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z"
                            clip-rule="evenodd"
                        />
                    </svg>
                    <div class="info-text">
                        <span class="info-label">Source Path</span>
                        <span class="info-value">{slackImportCard.Path}</span>
                    </div>
                </div>
            </div>

            <div class="info-row">
                <div class="info-item">
                    <svg
                        class="info-icon"
                        viewBox="0 0 20 20"
                        fill="currentColor"
                    >
                        <path
                            fill-rule="evenodd"
                            d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z"
                            clip-rule="evenodd"
                        />
                    </svg>
                    <div class="info-text">
                        <span class="info-label">Started</span>
                        <span class="info-value"
                            >{formatDate(slackImportCard.Created)}</span
                        >
                    </div>
                </div>

                {#if slackImportCard.Finished}
                    <div class="info-item">
                        <svg
                            class="info-icon"
                            viewBox="0 0 20 20"
                            fill="currentColor"
                        >
                            <path
                                fill-rule="evenodd"
                                d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                                clip-rule="evenodd"
                            />
                        </svg>
                        <div class="info-text">
                            <span class="info-label">Completed</span>
                            <span class="info-value"
                                >{formatDate(slackImportCard.Finished)}</span
                            >
                        </div>
                    </div>
                {/if}
            </div>

            <div class="info-row">
                <div class="info-item">
                    <svg
                        class="info-icon"
                        viewBox="0 0 20 20"
                        fill="currentColor"
                    >
                        <path
                            d="M2.003 5.884L10 9.882l7.997-3.998A2 2 0 0016 4H4a2 2 0 00-1.997 1.884z"
                        />
                        <path
                            d="M18 8.118l-8 4-8-4V14a2 2 0 002 2h12a2 2 0 002-2V8.118z"
                        />
                    </svg>
                    <div class="info-text">
                        <span class="info-label">JSON files</span>
                        <span class="info-value"
                            >{slackImportCard.NumFiles?.toLocaleString() ||
                                "0"}</span
                        >
                    </div>
                </div>

                <div class="info-item">
                    <svg
                        class="info-icon"
                        viewBox="0 0 20 20"
                        fill="currentColor"
                    >
                        <path
                            fill-rule="evenodd"
                            d="M3 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z"
                            clip-rule="evenodd"
                        />
                    </svg>
                    <div class="info-text">
                        <span class="info-label">Import ID</span>
                        <span class="info-value">{slackImportCard.ID}</span>
                    </div>
                </div>
            </div>
            <div class="info-row">
                <div class="info-item">
                    <svg
                        class="info-icon"
                        viewBox="0 0 24 24"
                        fill="none"
                        xmlns="http://www.w3.org/2000/svg"
                    >
                        <path
                            d="M12 2C6.477 2 2 6.477 2 12s4.477 10 10 10 10-4.477 10-10S17.523 2 12 2zm4.293 8.293l-5 5a1 1 0 01-1.414 0l-3-3a1 1 0 011.414-1.414L11 12.586l4.293-4.293a1 1 0 011.414 1.414z"
                            fill="currentColor"
                        />
                    </svg>
                    <div class="info-text">
                        <span class="info-label">Pending downloads</span>
                        <span class="info-value">{pendingBinaries}</span>
                    </div>
                </div>
            </div>
        </div>

        {#if !slackImportCard.Finished}
            <div class="progress-section">
                <div class="progress-header">
                    <span class="progress-label">Progress</span>
                    <span class="progress-percentage"
                        >{progress?.Progress ?? slackImportCard.Progress}%</span
                    >
                </div>
                <div class="progress-bar">
                    <div
                        class="progress-fill"
                        style="width: {progress?.Progress ??
                            slackImportCard.Progress}%; background-color: {getStatusColor(
                            progress?.Progress ?? slackImportCard.Progress,
                            slackImportCard.Error,
                        )}"
                    ></div>
                </div>
            </div>
        {/if}
    </div>
</main>

<style>
    /* Add this keyframes rule for the spinner animation */
    @keyframes spin {
        from {
            transform: rotate(0deg);
        }
        to {
            transform: rotate(360deg);
        }
    }

    /* Add this class for the spinner SVG */
    .spinner {
        animation: spin 1s linear infinite;
        width: 1rem;
        height: 1rem;
    }

    .card {
        background: #ffffff;
        border-radius: 0.5rem;
        padding: 0.4rem;
        /*box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);*/
        border: 1px solid #e5e7eb;
        border: 1px solid red;
    }

    .card-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        /* BUG FIX: Removed fixed margin-bottom to let the content inside the collapsible define the spacing*/
        /*gap: 1rem;*/
        border: 1px solid green;
    }
    .collapsible:not(.collapsed) .card-header {
        /*margin-bottom: 1rem;*/
    }

    .card-title {
        color: #111827;
        font-size: 1.2rem;
        font-weight: 600;
        margin: 0;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        text-align: left;
        /* (Optional) define a width or let a flex/grid container constrain it */
        /* max-width: 260px; */
    }

    .card-actions {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        flex-shrink: 0;
    }

    .status-badge {
        color: #fff;
        font-size: 0.75rem;
        font-weight: 600;
        padding: 0.25rem 0.75rem;
        border-radius: 9999px;
        text-transform: uppercase;
        letter-spacing: 0.025em;
        flex-shrink: 0;
    }

    .action-button,
    .chevron-button {
        min-width: 30px;
        width: 30px;
        aspect-ratio: 1;
        display: grid;
        place-items: center;
        padding: 0;
        border: 1px solid #ccc;
        border-radius: 6px;
        background: #fff;
        cursor: pointer;
    }

    /* Style for disabled download button */
    .action-button:disabled {
        cursor: not-allowed;
        background-color: #f3f4f6;
        color: #9ca3af;
    }
    .action-button:disabled:hover {
        transform: none;
        box-shadow: none;
    }

    .action-button:hover,
    .chevron-button:hover {
        background-color: #f3f4f6;
        color: #111827;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    }

    /*.action-button svg,
    .chevron-button svg {
        width: 1rem;
        height: 1rem;
    }*/
    /*
    .chevron-button svg.rotated {
        transform: rotate(90deg);
    }*/

    .collapsible {
        overflow: hidden;
    }

    .collapsible.collapsed {
        max-height: 0;
        /* remove all margins from children when collapsed */
        margin-top: -1.5rem; /* compensate for card padding */
    }

    .error-section {
        background-color: #fef2f2;
        border: 1px solid #fecaca;
        border-radius: 0.5rem;
        padding: 1rem;
        margin-bottom: 1.25rem;
    }
    .error-header {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        margin-bottom: 0.5rem;
    }
    .error-icon {
        width: 1.25rem;
        height: 1.25rem;
        color: #dc2626;
        flex-shrink: 0;
    }
    .error-title {
        font-size: 0.875rem;
        font-weight: 600;
        color: #b91c1c;
    }
    .error-message {
        font-size: 0.875rem;
        color: #7f1d1d;
        margin: 0 0 0.75rem 0;
        line-height: 1.5;
    }

    .error-actions {
        display: flex;
        gap: 0.5rem;
        margin-top: 0.75rem;
    }

    .retry-button {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.5rem 0.75rem;
        background-color: #dc2626;
        color: white;
        border: none;
        border-radius: 0.375rem;
        font-size: 0.875rem;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.2s ease;
    }

    .retry-button:hover {
        background-color: #b91c1c;
        transform: translateY(-1px);
        box-shadow: 0 2px 4px rgba(220, 38, 38, 0.2);
    }

    .retry-button svg {
        width: 1rem;
        height: 1rem;
    }

    .card-content {
        display: flex;
        flex-direction: column;
        gap: 1.25rem;
    }
    .info-row {
        display: flex;
        gap: 1.25rem;
    }
    .info-item {
        display: flex;
        align-items: flex-start; /* Changed to flex-start for better alignment */
        gap: 0.75rem;
        flex: 1;
        min-width: 0;
    }
    .info-icon {
        width: 1.25rem;
        height: 1.25rem;
        color: #6b7280;
        margin-top: 0.125rem;
        flex-shrink: 0; /* Prevent icon from shrinking */
    }
    .info-text {
        display: flex;
        flex-direction: column;
        flex: 1;
        min-width: 0;
    }
    .info-label {
        font-size: 0.75rem;
        font-weight: 500;
        color: #6b7280;
        text-transform: uppercase;
        margin-bottom: 0.25rem;
    }
    .info-value {
        font-size: 0.9rem;
        color: #111827;
        font-weight: 500;
        word-break: break-all; /* Changed to break-all for long paths */
    }

    .progress-section {
        border-top: 1px solid #e5e7eb;
        padding-top: 1.25rem;
        margin-top: 1.25rem;
    }
    .progress-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 0.5rem;
    }
    .progress-label {
        font-size: 0.75rem;
        font-weight: 500;
        color: #6b7280;
        text-transform: uppercase;
    }
    .progress-percentage {
        font-size: 0.875rem;
        font-weight: 600;
        color: #111827;
    }
    .progress-bar {
        width: 100%;
        height: 0.625rem;
        background-color: #f3f4f6;
        border-radius: 9999px;
        overflow: hidden;
    }
    .progress-fill {
        height: 100%;
        border-radius: 9999px;
        transition: width 0.3s ease;
    }
</style>
