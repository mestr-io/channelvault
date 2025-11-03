<script lang="ts">
import { onMount } from 'svelte'
import * as runtime from '../../../wailsjs/runtime'
import { models } from '../../../wailsjs/go/models'
import { RenameImport, DownloadBinaries } from '../../../wailsjs/go/app/App'
import Prompt from './Prompt.svelte'

import {
  chevron,
  downloadIcon,
  calendarIcon,
  doneIcon,
  envelopeIcon,
  burgerMenuIcon,
  renameIcon,
  deleteIcon,
  errorIcon,
  retryIcon,
} from './icons'

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

$inspect(isCollapsed)

$effect(() => {
  isCollapsed = legacyIsCollapsed(slackImportCard)
})

function legacyIsCollapsed(importState: models.SlackImport): boolean {
  // Legacy mode assings isCollapsed to =>  slackImportCard.Finished !== "" && slackImportCard.Error === ""

  return importState.Finished !== '' && importState.Error === ''
}

function detectPrimarySeparator(path: string): '/' | '\\' | null {
  const lastFwd = path.lastIndexOf('/')
  const lastBack = path.lastIndexOf('\\')
  if (lastFwd === -1 && lastBack === -1) return null
  return lastFwd > lastBack ? '/' : '\\'
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

        <button
            class="card-title"
            class:truncated={isCollapsed}
            onclick={toggleCollapse}
        >
            {slackImportCard.Name || "Untitled Import"}
        </button>

        <div class="card-actions">
            <div
                class="status-badge"
                style="background-color: {getStatusColor(
                    progress?.Progress ?? 0,
                    slackImportCard.Error,
                )}"
            >
                {getStatusText(
                    progress?.Progress ?? 0,
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
                {@render renameIcon()}
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
                        <div class="spinner"></div>
                    {:else}
                        {@render downloadIcon()}
                    {/if}
                </button>
            {/if}
            <button
                aria-label="delete import"
                class="action-button delete-button"
                onclick={() => handleDelete(slackImportCard.ID)}
                title="Delete import"
            >
                {@render deleteIcon()}
            </button>
        </div>
    </div>

    <div class="collapsible {isCollapsed ? 'collapsed' : ''}">
        {#if slackImportCard.Error}
            <div class="error-section">
                <div class="error-header">
                    {@render errorIcon()}
                    <span class="error-title">Import Error</span>
                </div>
                <p class="error-message">{slackImportCard.Error}</p>
                <div class="error-actions">
                    <button
                        class="retry-button"
                        onclick={() => handleRetry(slackImportCard.ID)}
                        title="Retry import"
                    >
                        {@render retryIcon()}
                        Retry Import
                    </button>
                </div>
            </div>
        {/if}

        <div class="card-content">
            <div class="info-row zebra">
                <div class="info-item">
                    {@render downloadIcon()}
                    <div class="info-text">
                        <span class="info-label">Original Location</span>
                        <span title={slackImportCard.Path} class="info-value"
                            >{slackImportCard.Path}</span
                        >
                    </div>
                </div>
            </div>

            <div class="info-row">
                <div class="info-item">
                    {@render calendarIcon()}
                    <div class="info-text">
                        <span class="info-label">Started</span>
                        <span class="info-value"
                            >{formatDate(slackImportCard.Created)}</span
                        >
                    </div>
                </div>

                {#if slackImportCard.Finished}
                    <div class="info-item">
                        {@render doneIcon()}
                        <div class="info-text">
                            <span class="info-label">Completed</span>
                            <span class="info-value"
                                >{formatDate(slackImportCard.Finished)}</span
                            >
                        </div>
                    </div>
                {/if}
            </div>

            <div class="info-row zebra">
                <div class="info-item">
                    {@render envelopeIcon()}
                    <div class="info-text">
                        <span class="info-label">JSON files</span>
                        <span class="info-value"
                            >{slackImportCard.NumFiles?.toLocaleString() ||
                                "0"}</span
                        >
                    </div>
                </div>

                <div class="info-item">
                    {@render burgerMenuIcon()}
                    <div class="info-text">
                        <span class="info-label">Import ID</span>
                        <span class="info-value">{slackImportCard.ID}</span>
                    </div>
                </div>
            </div>
            <div class="info-row">
                <div class="info-item">
                    {@render doneIcon()}
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
                            0}%; background-color: {getStatusColor(
                            progress?.Progress ?? 0,
                            slackImportCard.Error,
                        )}"
                    ></div>
                </div>
            </div>
        {/if}
    </div>
</main>

<style>
    /* * =============================================
     * Design System: Variables (Design Tokens)
     * =============================================
     * Using variables makes the design consistent and easy to change.
    */
    main {
        --c-surface: #ffffff;
        --c-surface-soft: #f9fafb; /* Softer background for zebra stripes */
        --c-border: #e5e7eb;

        --c-text-primary: #1f2937;
        --c-text-secondary: #6b7280;
        --c-text-label: #9ca3af;

        --c-brand-blue: #3b82f6;

        --c-status-error-bg: #fef2f2;
        --c-status-error-border: #fecaca;
        --c-status-error-text: #b91c1c;
        --c-status-error-icon: #ef4444;

        --shadow-sm: 0 1px 2px 0 rgb(0 0 0 / 0.05);
        --shadow-md:
            0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1);
        --radius: 0.75rem; /* A more generous border radius */
        --transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    }

    /* * =============================================
     * Spinner Animation
     * =============================================
    */
    .spinner {
        width: 1.125rem; /* 18px */
        height: 1.125rem;
        border: 2px solid var(--c-brand-blue);
        border-bottom-color: transparent;
        border-radius: 50%;
        display: inline-block;
        box-sizing: border-box;
        animation: spin 1s linear infinite;
    }

    @keyframes spin {
        0% {
            transform: rotate(0deg);
        }
        100% {
            transform: rotate(360deg);
        }
    }

    /* * =============================================
     * Main Card Structure
     * =============================================
    */
    .card {
        background: var(--c-surface);
        border-radius: var(--radius);
        border: 1px solid var(--c-border);
        box-shadow: var(--shadow-sm);
        padding: 1rem; /* More breathing room */
        width: 100%;
        overflow: hidden; /* Important for child border-radius */
    }

    .card-header {
        display: flex;
        align-items: center; /* Better vertical alignment */
        gap: 0.75rem; /* Consistent spacing */
        min-width: 0;
    }

    .card-title {
        flex-grow: 1;
        color: var(--c-text-primary);
        font-size: 1rem;
        font-weight: 600;
        line-height: 1.5;
        margin: 0; /* Reset margin */
        min-width: 0;
        white-space: normal;
        border: 0;
        background-color: transparent;
        text-align: left;
        cursor: pointer;
        padding: 0;
    }

    .card-title.truncated {
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .card-actions {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        flex-shrink: 0;
    }

    /* * =============================================
     * Status Badge & Action Buttons
     * =============================================
    */
    .status-badge {
        color: #fff;
        font-size: 0.75rem;
        font-weight: 500;
        padding: 0.25rem 0.75rem;
        border-radius: 9999px;
        line-height: 1.25;
    }

    .action-button,
    .chevron-button {
        width: 2.25rem; /* 36px */
        height: 2.25rem;
        display: grid;
        place-items: center;
        border: 1px solid var(--c-border);
        border-radius: 0.5rem;
        background: var(--c-surface);
        color: var(--c-text-secondary);
        cursor: pointer;
        transition: var(--transition);
    }

    .action-button :global(svg),
    .chevron-button :global(svg) {
        width: 1.25rem;
        height: 1.25rem;
    }

    .action-button:hover,
    .chevron-button:hover {
        background-color: var(--c-surface-soft);
        color: var(--c-text-primary);
        border-color: #d1d5db;
        transform: translateY(-1px);
        box-shadow: var(--shadow-sm);
    }

    .action-button:disabled {
        cursor: not-allowed;
        background-color: var(--c-surface-soft);
        color: #d1d5db;
        transform: none;
        box-shadow: none;
    }

    /* * =============================================
     * Collapsible Section & Animation
     * =============================================
    */
    .collapsible {
        overflow: hidden;
        display: grid;
        grid-template-rows: 1fr;
        transition:
            grid-template-rows 0.3s ease-out,
            margin-top 0.3s ease-out;
        margin-top: 1rem;
    }

    .collapsible.collapsed {
        grid-template-rows: 0fr;
        margin-top: 0;
    }

    /* This direct child is necessary for the grid animation to work correctly */
    .collapsible > * {
        min-height: 0;
    }

    /* * =============================================
     * Error Section
     * =============================================
    */
    .error-section {
        background-color: var(--c-status-error-bg);
        border: 1px solid var(--c-status-error-border);
        border-radius: var(--radius);
        padding: 1rem;
        margin-bottom: 1rem;
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
        color: var(--c-status-error-icon);
        flex-shrink: 0;
    }
    .error-title {
        font-size: 0.875rem;
        font-weight: 600;
        color: var(--c-status-error-text);
    }
    .error-message {
        font-size: 0.875rem;
        color: #991b1b;
        margin: 0 0 1rem 0;
        line-height: 1.5;
        padding-left: 1.75rem;
    }
    .error-actions {
        display: flex;
        justify-content: flex-end; /* Align button to the right */
    }
    .retry-button {
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.5rem 1rem;
        background-color: var(--c-status-error-icon);
        color: white;
        border: none;
        border-radius: 0.5rem;
        font-size: 0.875rem;
        font-weight: 500;
        cursor: pointer;
        transition: var(--transition);
    }
    .retry-button:hover {
        background-color: var(--c-status-error-text);
        box-shadow: var(--shadow-md);
        transform: translateY(-2px);
    }
    .retry-button svg {
        width: 1rem;
        height: 1rem;
    }

    /* * =============================================
     * Card Content Details
     * =============================================
    */
    .card-content {
        margin: 0;
        display: flex;
        flex-direction: column;
        gap: 0.5rem; /* Space out the rows */
    }

    .info-row {
        display: flex;
        gap: 1rem; /* More space between items */
        padding: 0.75rem 1rem;
        border-radius: 0.5rem;
    }

    .zebra {
        background-color: var(--c-surface-soft);
    }

    .info-item {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        flex: 1;
        min-width: 0;
    }

    .info-item :global(svg) {
        width: 1.25rem;
        height: 1.25rem;
        color: var(--c-text-secondary);
        flex-shrink: 0;
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
        color: var(--c-text-label);
        text-transform: uppercase;
        letter-spacing: 0.05em;
        margin-bottom: 0.125rem;
    }

    .info-value {
        font-size: 0.875rem;
        color: var(--c-text-secondary);
        font-weight: 500;
        word-break: break-all;
        line-height: 1.4;
    }

    .info-value[title] {
        /* Add subtle indicator for truncated text */
        cursor: help;
    }

    /* * =============================================
     * Progress Bar
     * =============================================
    */
    .progress-section {
        padding-top: 1rem;
        margin-top: 1rem;
        border-top: 1px solid var(--c-border);
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
        color: var(--c-text-label);
        text-transform: uppercase;
        letter-spacing: 0.05em;
    }
    .progress-percentage {
        font-size: 0.875rem;
        font-weight: 600;
        color: var(--c-text-primary);
    }
    .progress-bar {
        width: 100%;
        height: 0.5rem;
        background-color: var(--c-surface-soft);
        border: 1px solid var(--c-border);
        border-radius: 9999px;
        overflow: hidden;
    }
    .progress-fill {
        height: 100%;
        border-radius: 9999px;
        transition: width 0.3s ease;
    }
</style>
