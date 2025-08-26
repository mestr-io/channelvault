<script lang="ts">
import { onMount } from 'svelte'
import * as runtime from '../../../wailsjs/runtime'
import { models } from '../../../wailsjs/go/models'
let { card }: { card: models.SlackImport } = $props()
let progress: models.Progress | null = $state(null)

onMount(() => {
  runtime.EventsOn('importProgress', (p: models.Progress) => {
    // dispatching
    if (card.ID == p.ImportID) {
      progress = p
      console.log(card, p)
    }
  })

  runtime.EventsOn('importError', err => {})

  runtime.EventsOn('importComplete', () => {})
})

function formatDate(dateString: string): string {
  if (!dateString) return 'Not started'
  const date = new Date(dateString)
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

function getStatusColor(progress: number, error: string): string {
  if (error) return '#ef4444' // red for errors
  if (progress === 0) return '#6b7280' // gray
  if (progress < 50) return '#f59e0b' // amber
  if (progress < 100) return '#3b82f6' // blue
  return '#10b981' // green
}

function getStatusText(progress: number, finished: string, error: string): string {
  if (error) return 'Error'
  if (finished) return 'Completed'
  if (progress === 0) return 'Pending'
  if (progress < 100) return 'In Progress'
  return 'Processing'
}

function handleRename(importId: number, currentName: string) {
  const newName = prompt('Enter new name for import:', currentName)
}

function handleDelete(importId: number, importName: string) {
  confirm(`Are you sure you want to delete "${importName}"? This action cannot be undone.`)
}
</script>

<main class="card">
    <div class="card-header">
        <h3 class="card-title">{card.Name || "Untitled Import"}</h3>
        <div class="card-actions">
            <div
                class="status-badge"
                style="background-color: {getStatusColor(
                    progress?.Progress ?? 0,
                    card.Error,
                )}"
            >
                {getStatusText(progress?.Progress ?? 0, card.Finished, card.Error)}
            </div>
            <button
                aria-label="rename import"
                class="action-button rename-button"
                onclick={() => handleRename(card.ID, card.Name)}
                title="Rename import"
            >
                <svg viewBox="0 0 20 20" fill="currentColor">
                    <path
                        d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z"
                    />
                </svg>
            </button>
            <button
                aria-label="delete import"
                class="action-button delete-button"
                onclick={() => handleDelete(card.ID, card.Name)}
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

    {#if card.Error}
        <div class="error-section">
            <div class="error-header">
                <svg class="error-icon" viewBox="0 0 20 20" fill="currentColor">
                    <path
                        fill-rule="evenodd"
                        d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z"
                        clip-rule="evenodd"
                    />
                </svg>
                <span class="error-title">Import Error</span>
            </div>
            <p class="error-message">{card.Error}</p>
        </div>
    {/if}

    <div class="card-content">
        <div class="info-row">
            <div class="info-item">
                <svg class="info-icon" viewBox="0 0 20 20" fill="currentColor">
                    <path
                        fill-rule="evenodd"
                        d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z"
                        clip-rule="evenodd"
                    />
                </svg>
                <div class="info-text">
                    <span class="info-label">Source Path</span>
                    <span class="info-value">{card.Path}</span>
                </div>
            </div>
        </div>

        <div class="info-row">
            <div class="info-item">
                <svg class="info-icon" viewBox="0 0 20 20" fill="currentColor">
                    <path
                        fill-rule="evenodd"
                        d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z"
                        clip-rule="evenodd"
                    />
                </svg>
                <div class="info-text">
                    <span class="info-label">Started</span>
                    <span class="info-value">{formatDate(card.Created)}</span>
                </div>
            </div>

            {#if card.Finished}
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
                            >{formatDate(card.Finished)}</span
                        >
                    </div>
                </div>
            {/if}
        </div>

        <div class="info-row">
            <div class="info-item">
                <svg class="info-icon" viewBox="0 0 20 20" fill="currentColor">
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
                        >{card.NumFiles?.toLocaleString() || "0"}</span
                    >
                </div>
            </div>

            <div class="info-item">
                <svg class="info-icon" viewBox="0 0 20 20" fill="currentColor">
                    <path
                        fill-rule="evenodd"
                        d="M3 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z"
                        clip-rule="evenodd"
                    />
                </svg>
                <div class="info-text">
                    <span class="info-label">Import ID</span>
                    <span class="info-value">{card.ID}</span>
                </div>
            </div>
        </div>
    </div>

    {#if !card.Finished}
    <div class="progress-section">
        <div class="progress-header">
            <span class="progress-label">Progress</span>
            <span class="progress-percentage">{progress?.Progress ?? 0}%</span>
        </div>
        <div class="progress-bar">
            <div
                class="progress-fill"
                style="width: {progress?.Progress ??
                    0}%; background-color: {getStatusColor(progress?.Progress ?? 0, card.Error)}"
            ></div>
        </div>
    </div>
    {/if}
</main>

<style>
    .card {
        background: linear-gradient(135deg, #ffffff 0%, #dfdfdf 100%);
        border-radius: 1rem;
        padding: 1.5rem;
        box-shadow:
            0 4px 6px -1px rgba(0, 0, 0, 0.1),
            0 2px 4px -1px rgba(0, 0, 0, 0.06),
            0 0 0 1px rgba(0, 0, 0, 0.05);
        border: 1px solid #e2e8f0;
    }

    .card-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 1.25rem;
        gap: 1rem;
    }

    .card-title {
        color: #1a202c;
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
        color: white;
        font-size: 0.75rem;
        font-weight: 600;
        padding: 0.25rem 0.75rem;
        border-radius: 9999px;
        text-transform: uppercase;
        letter-spacing: 0.025em;
        flex-shrink: 0;
    }

    .action-button {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 2rem;
        height: 2rem;
        border: none;
        border-radius: 0.375rem;
        cursor: pointer;
        transition: all 0.2s ease;
        background-color: #f3f4f6;
        color: #6b7280;
    }

    .action-button:hover {
        transform: translateY(-1px);
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }

    .action-button svg {
        width: 1rem;
        height: 1rem;
    }

    .rename-button:hover {
        background-color: #dbeafe;
        color: #3b82f6;
    }

    .delete-button:hover {
        background-color: #fef2f2;
        color: #ef4444;
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
        color: #ef4444;
        flex-shrink: 0;
    }

    .error-title {
        font-size: 0.875rem;
        font-weight: 600;
        color: #dc2626;
    }

    .error-message {
        font-size: 0.875rem;
        color: #7f1d1d;
        margin: 0;
        line-height: 1.5;
    }

    .card-content {
        display: flex;
        flex-direction: column;
        gap: 1rem;
        margin-bottom: 1.5rem;
    }

    .info-row {
        display: flex;
        gap: 1rem;
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
        flex-shrink: 0;
        margin-top: 0.125rem;
    }

    .info-text {
        display: flex;
        flex-direction: column;
        min-width: 0;
        flex: 1;
    }

    .info-label {
        font-size: 0.75rem;
        font-weight: 500;
        color: #6b7280;
        text-transform: uppercase;
        letter-spacing: 0.025em;
        margin-bottom: 0.25rem;
    }

    .info-value {
        font-size: 0.875rem;
        color: #374151;
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
        letter-spacing: 0.025em;
    }

    .progress-percentage {
        font-size: 0.875rem;
        font-weight: 600;
        color: #374151;
    }

    .progress-bar {
        width: 100%;
        height: 0.5rem;
        background-color: #e5e7eb;
        border-radius: 9999px;
        overflow: hidden;
    }

    .progress-fill {
        height: 100%;
        border-radius: 9999px;
        transition: all 0.3s ease;
    }
</style>
