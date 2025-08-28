<script lang="ts">
import { onMount } from 'svelte'
import * as runtime from '../../../wailsjs/runtime'
import { models } from '../../../wailsjs/go/models'
import { RenameImport } from '../../../wailsjs/go/app/App'
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
  })
})

function formatDate(dateString: string): string {
  if (!dateString) return 'Not started'
  const date = new Date(dateString)
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

function getStatusColor(progress: number, error: string): string {
  if (error) return '#e53e3e' // Red
  if (progress === 0 && slackImportcard.Finished === '') return '#718096' // Gray
  if (progress < 100 && slackImportcard.Finished === '') return '#4a90e2' // Blue
  return '#38a169' // Green
}

function getStatusText(progress: number, finished: string, error: string): string {
  if (error) return 'Error'
  if (finished) return 'Completed'
  if (progress === 0) return 'Pending'
  return 'In Progress'
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

async function handleDelete(importId: number) {
  await DeleteImportHandle(importId)
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
<article class="card">
    <header class="card-header" onclick={toggleCollapse}>
        <button
            aria-label="collapse card"
            class="chevron-button"
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

        <div class="card-actions" onclick={e => e.stopPropagation()}>
            <div
                class="status-badge"
                style="--status-color: {getStatusColor(
                    progress?.Progress ?? (slackImportcard.Finished ? 100 : 0),
                    slackImportcard.Error,
                )}"
            >
                {getStatusText(
                    progress?.Progress ?? (slackImportcard.Finished ? 100 : 0),
                    slackImportcard.Finished,
                    slackImportcard.Error,
                )}
            </div>
            <button
                aria-label="rename import"
                class="action-button"
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
            
            {#if false }
            <button
                aria-label="delete import"
                class="action-button delete-button"
                onclick={() => handleDelete(slackImportcard.ID)}
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
            {/if}
        </div>
    </header>

    <div class="collapsible" class:collapsed={isCollapsed}>
        <div class="collapsible-content">
            {#if slackImportcard.Error}
                <section class="error-section">
                    <div class="error-header">
                        <svg class="error-icon" viewBox="0 0 20 20" fill="currentColor">
                            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>
                        </svg>
                        <h4 class="error-title">Import Error</h4>
                    </div>
                    <p class="error-message">{slackImportcard.Error}</p>
                    <div class="error-actions">
                        <button class="retry-button" onclick={() => handleRetry(slackImportcard.ID)}>
                             <svg viewBox="0 0 20 20" fill="currentColor">
                                <path fill-rule="evenodd" d="M4 2a1 1 0 011 1v2.101a7.002 7.002 0 0111.601 2.566 1 1 0 11-1.885.666A5.002 5.002 0 005.999 7H9a1 1 0 010 2H4a1 1 0 01-1-1V3a1 1 0 011-1zm.008 9.057a1 1 0 011.276.61A5.002 5.002 0 0014.001 13H11a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0v-2.101a7.002 7.002 0 01-11.601-2.566 1 1 0 01.61-1.276z" clip-rule="evenodd"/>
                            </svg>
                            Retry Import
                        </button>
                    </div>
                </section>
            {/if}

            <section class="card-content">
                <div class="info-row">
                    <span class="info-label">Source Path</span>
                    <span class="info-value path">{slackImportcard.Path}</span>
                </div>
                <div class="info-grid">
                     <div class="info-item">
                        <span class="info-label">Started</span>
                        <span class="info-value">{formatDate(slackImportcard.Created)}</span>
                    </div>
                     {#if slackImportcard.Finished}
                        <div class="info-item">
                            <span class="info-label">Completed</span>
                            <span class="info-value">{formatDate(slackImportcard.Finished)}</span>
                        </div>
                    {/if}
                    <div class="info-item">
                        <span class="info-label">Total Files</span>
                        <span class="info-value">{slackImportcard.NumFiles?.toLocaleString() || "N/A"}</span>
                    </div>
                     <div class="info-item">
                        <span class="info-label">Import ID</span>
                        <span class="info-value">{slackImportcard.ID}</span>
                    </div>
                </div>
            </section>

            {#if !slackImportcard.Finished && !slackImportcard.Error}
                <section class="progress-section">
                    <div class="progress-header">
                        <span class="progress-label">Progress</span>
                        <span class="progress-percentage">{progress?.Progress ?? 0}%</span>
                    </div>
                    <div class="progress-bar">
                        <div
                            class="progress-fill"
                            style="width: {progress?.Progress ?? 0}%; background-color: {getStatusColor(
                                progress?.Progress ?? 0,
                                slackImportcard.Error
                            )}"
                        ></div>
                    </div>
                </section>
            {/if}
        </div>
    </div>
</article>

<style>
    .card {
        background: #2d3748; /* Dark card background */
        border-radius: 0.75rem;
        border: 1px solid #4a5568;
        color: #e2e8f0;
        overflow: hidden;
        transition: border-color 0.2s ease;
    }
    .card:hover {
        border-color: #5a667a;
    }

    .card-header {
        display: flex;
        align-items: center;
        padding: 1rem 1.25rem;
        cursor: pointer;
        gap: 1rem;
    }

    .card-title {
        font-size: 1.125rem;
        font-weight: 600;
        margin: 0;
        flex-grow: 1;
        color: #f7fafc;
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

    .status-badge {
        font-size: 0.75rem;
        font-weight: 600;
        padding: 0.25rem 0.75rem;
        border-radius: 9999px;
        background-color: color-mix(in srgb, var(--status-color) 25%, transparent);
        color: var(--status-color);
        border: 1px solid var(--status-color);
    }
    
    .action-button, .chevron-button {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 2.25rem;
        height: 2.25rem;
        border: none;
        border-radius: 0.5rem;
        cursor: pointer;
        transition: background-color 0.2s ease, color 0.2s ease;
        background-color: transparent;
        color: #a0aec0; /* Icon color */
    }

    .action-button:hover { background-color: #4a5568; color: #ffffff; }
    .chevron-button { margin-left: -0.5rem; }
    .delete-button:hover { background-color: #e53e3e; color: #ffffff; }

    .action-button svg, .chevron-button svg { width: 1.125rem; height: 1.125rem; }

    .chevron-button svg { transition: transform 0.3s ease; }
    .chevron-button svg.rotated { transform: rotate(90deg); }

    .collapsible {
        display: grid;
        grid-template-rows: 0fr;
        transition: grid-template-rows 0.3s ease-out;
    }
    .collapsible:not(.collapsed) {
        grid-template-rows: 1fr;
    }
    .collapsible-content {
        overflow: hidden;
        padding: 0 1.25rem 1.25rem 1.25rem;
        border-top: 1px solid #4a5568;
        margin-top: -1px; /* Overlap border */
    }

    .error-section {
        background-color: #442121;
        border: 1px solid #e53e3e;
        border-radius: 0.5rem;
        padding: 1rem;
        margin-top: 1rem;
    }
    .error-header { display: flex; align-items: center; gap: 0.5rem; margin-bottom: 0.5rem; }
    .error-icon { width: 1.25rem; height: 1.25rem; color: #e53e3e; }
    .error-title { font-size: 0.875rem; font-weight: 600; color: #fed7d7; margin: 0; }
    .error-message { font-size: 0.875rem; color: #e2e8f0; margin: 0 0 1rem; line-height: 1.5; }
    
    .retry-button {
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.5rem 1rem;
        background-color: #e53e3e;
        color: white;
        border: none;
        border-radius: 0.375rem;
        font-size: 0.875rem;
        font-weight: 500;
        cursor: pointer;
        transition: background-color 0.2s ease;
    }
    .retry-button:hover { background-color: #c53030; }
    .retry-button svg { width: 1rem; height: 1rem; }

    .card-content { display: flex; flex-direction: column; gap: 1.25rem; margin-top: 1rem; }
    .info-row { display: flex; flex-direction: column; gap: 0.25rem; }
    .info-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap: 1.25rem; }
    .info-item { display: flex; flex-direction: column; gap: 0.25rem; }
    .info-label { font-size: 0.75rem; font-weight: 500; color: #a0aec0; text-transform: uppercase; letter-spacing: 0.025em; }
    .info-value { font-size: 0.9rem; color: #f7fafc; font-weight: 500; }
    .info-value.path { font-family: monospace; word-break: break-all; }

    .progress-section { margin-top: 1rem; }
    .progress-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem; }
    .progress-label { font-size: 0.75rem; font-weight: 500; color: #a0aec0; text-transform: uppercase; }
    .progress-percentage { font-size: 0.875rem; font-weight: 600; color: #f7fafc; }
    .progress-bar { width: 100%; height: 0.5rem; background-color: #4a5568; border-radius: 9999px; overflow: hidden; }
    .progress-fill { height: 100%; border-radius: 9999px; transition: width 0.3s ease; }
</style>
