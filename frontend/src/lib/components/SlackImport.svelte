<script lang="ts">
import { onMount } from 'svelte'
import * as runtime from '../../../wailsjs/runtime'
import { models } from '../../../wailsjs/go/models'
import { RenameImport, DownloadBinaries } from '../../../wailsjs/go/app/App'
import Prompt from './Prompt.svelte'

let {
  slackImportcard: slackImportcard,
  isCollapsed = true,
  DeleteImportHandle,
  RetryImportHandle,
}: {
  slackImportcard: models.SlackImport
  isCollapsed: boolean
  DeleteImportHandle: (id: number) => Promise<void>
  RetryImportHandle: (id: number) => Promise<void>
} = $props()
let progress: models.Progress | null = $state(null)

onMount(() => {
  runtime.EventsOn('importProgress', (p: models.Progress) => {
    if (slackImportcard.ID == p.ImportID) {
      progress = p
    }
    console.log(p)
  })

  runtime.EventsOn('downloadBinariesDone', () => {})
  runtime.EventsOn('downloadBinariesStarted', () => {})
  runtime.EventsOn('downloadedBinary', file => {
    console.log(file)
  })
})

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
  return '#059669' // green
}

function getStatusText(progress: number, finished: string, error: string): string {
  if (error) return 'Error'
  if (finished) return 'Completed'
  if (progress === 0) return 'Pending'
  if (progress < 100) return 'In Progress'
  return 'Processing'
}

let showPrompt = $state(false)
let promptDefault = ''
let importId: number

async function handleRename(importId_: number, currentName: string) {
  importId = importId_
  promptDefault = currentName
  showPrompt = true
}

async function confirmRename(newName: string) {
  if (!newName) return
  await RenameImport(importId, newName)
  slackImportcard.Name = newName
  showPrompt = false
}

async function handleDelete(importId: number, importName: string) {
  await DeleteImportHandle(importId)
}

async function handleDownload(importId: number, importName: string) {
  await DownloadBinaries(importId)
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
            <svg
                class:rotated={!isCollapsed}
                viewBox="0 0 20 20"
                fill="currentColor"
            >
                <path
                    fill-rule="evenodd"
                    d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z"
                    clip-rule="evenodd"
                />
            </svg>
        </button>

        <h3 class="card-title">{slackImportcard.Name || "Untitled Import"}</h3>

        <div class="card-actions">
            <div
                class="status-badge"
                style="background-color: {getStatusColor(
                    progress?.Progress ?? 0,
                    slackImportcard.Error,
                )}"
            >
                {getStatusText(
                    progress?.Progress ?? 0,
                    slackImportcard.Finished,
                    slackImportcard.Error,
                )}
            </div>
            <button
                aria-label="rename import"
                class="action-button rename-button"
                onclick={() =>
                    handleRename(slackImportcard.ID, slackImportcard.Name)}
                title="Rename import"
            >
                <svg viewBox="0 0 20 20" fill="currentColor">
                    <path
                        d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z"
                    />
                </svg>
            </button>
            {#if slackImportcard.Finished}
            <button
                aria-label="download binaries"
                class="action-button download-button"
                onclick={() =>
                    handleDownload(slackImportcard.ID, slackImportcard.Name)}
                title="Download binaries"
            >
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
            </button>
            {/if}
            <button
                aria-label="delete import"
                class="action-button delete-button"
                onclick={() =>
                    handleDelete(slackImportcard.ID, slackImportcard.Name)}
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
        {#if slackImportcard.Error}
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
                <p class="error-message">{slackImportcard.Error}</p>
                <div class="error-actions">
                    <button
                        class="retry-button"
                        onclick={() => handleRetry(slackImportcard.ID)}
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
            <!-- info rows -->
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
                        <span class="info-value">{slackImportcard.Path}</span>
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
                            d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a1 1 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z"
                            clip-rule="evenodd"
                        />
                    </svg>
                    <div class="info-text">
                        <span class="info-label">Started</span>
                        <span class="info-value"
                            >{formatDate(slackImportcard.Created)}</span
                        >
                    </div>
                </div>

                {#if slackImportcard.Finished}
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
                                >{formatDate(slackImportcard.Finished)}</span
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
                        <span class="info-label">Files</span>
                        <span class="info-value"
                            >{slackImportcard.NumFiles?.toLocaleString() ||
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
                        <span class="info-value">{slackImportcard.ID}</span>
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
                        <span class="info-label">Pending downloads</span>
                        <span class="info-value">{slackImportcard.PendingBinaries}</span>
                    </div>
                </div>
            </div>
        </div>

        {#if !slackImportcard.Finished}
            <div class="progress-section">
                <div class="progress-header">
                    <span class="progress-label">Progress</span>
                    <span class="progress-percentage"
                        >{progress?.Progress ?? 0}%</span
                    >
                </div>
                <div class="progress-bar">
                    <div
                        class="progress-fill"
                        style="width: {progress?.Progress ??
                            0}%; background-color: {getStatusColor(
                            progress?.Progress ?? 0,
                            slackImportcard.Error,
                        )}"
                    ></div>
                </div>
            </div>
        {/if}
    </div>
</main>

<style>
    .card {
        background: #ffffff;
        border-radius: 1rem;
        padding: 1.5rem;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
        border: 1px solid #e5e7eb;
    }

    .card-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 1rem;
        gap: 1rem;
    }

    .card-title {
        color: #111827;
        font-size: 1.25rem;
        font-weight: 600;
        margin: 0;
        line-height: 1.4;
        flex: 1;
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
        display: flex;
        align-items: center;
        justify-content: center;
        width: 2rem;
        height: 2rem;
        border: none;
        border-radius: 0.5rem;
        cursor: pointer;
        transition: all 0.2s ease;
        background-color: #f9fafb;
        color: #6b7280;
    }

    .action-button:hover,
    .chevron-button:hover {
        background-color: #f3f4f6;
        color: #111827;
        transform: translateY(-1px);
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    }

    .action-button svg,
    .chevron-button svg {
        width: 1rem;
        height: 1rem;
    }

    .chevron-button svg.rotated {
        transform: rotate(90deg);
    }

    .collapsible {
        max-height: 2000px;
        overflow: hidden;
        transition: max-height 0.3s ease;
    }
    .collapsible.collapsed {
        max-height: 0;
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
        margin-bottom: 1.5rem;
    }
    .info-row {
        display: flex;
        gap: 1.25rem;
    }
    .info-item {
        display: flex;
        align-items: flex-end;
        gap: 0.75rem;
        flex: 1;
        min-width: 0;
    }
    .info-icon {
        width: 1.25rem;
        height: 1.25rem;
        color: #6b7280;
        margin-top: 0.125rem;
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
        word-break: break-word;
    }

    .progress-section {
        border-top: 1px solid #e5e7eb;
        padding-top: 1rem;
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
