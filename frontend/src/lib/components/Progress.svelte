<!-- 
  File: src/lib/Progress.svelte
  This version includes a console.log to diagnose the data from Go.
-->
<script>
import { onMount } from 'svelte'
import * as runtime from '../../../wailsjs/runtime'

let progress = {
  FileName: 'Waiting for file drop...',
  Progress: 0,
  TotalFiles: 0,
  CurrentFile: 0,
}
let errorMessage = ''
let importComplete = false

onMount(() => {
  runtime.EventsOn('importProgress', p => {
    // --- DIAGNOSTIC LOG ---
    // This is the most important line. It will show us the raw data.
    console.log('Received progress object from Go:', p)

    if (!importComplete) {
      progress = p
      errorMessage = ''
    }
  })

  runtime.EventsOn('importError', err => {
    errorMessage = err
    importComplete = false
  })

  runtime.EventsOn('importComplete', () => {
    importComplete = true
    progress = {
      FileName: 'Finished',
      TotalFiles: progress.TotalFiles,
      CurrentFile: progress.TotalFiles,
      Progress: 1,
    }
  })
})

// To be extra safe, we clamp the progress value between 0 and 1 before multiplying.
$: clampedProgress = Math.max(0, Math.min(1, progress.Progress || 0))

// Use the clamped value for all calculations.
$: percentageValue = clampedProgress * 100
$: formattedPercentage = percentageValue.toFixed(percentageValue >= 99.9 ? 0 : 1)
</script>

<!-- HTML Markup -->
<div class="progress-wrapper">
  <div class="header">
    <span class="filename">{progress.FileName}</span>
    <span class="file-count">{progress.CurrentFile} / {progress.TotalFiles}</span>
  </div>

  <div class="progress-bar-container">
    <div class="progress-bar-inner" style="width: {percentageValue}%"></div>
  </div>

  <div class="footer">
    {#if errorMessage}
      <div class="status-text error">
        <svg xmlns="http://www.w3.org/2000/svg" class="icon" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" /></svg>
        Error: {errorMessage}
      </div>
    {:else}
      <div class="status-text success">
        {#if importComplete}
          ✅ Import Complete!
        {:else}
          Importing...
        {/if}
      </div>
    {/if}
    <span class="percentage">{formattedPercentage}%</span>
  </div>
</div>

<!-- Component-specific styles (Plain CSS) -->
<style>
  .progress-wrapper {
    max-height: 200px;
    background-color: #4a0e4e;
    margin-left: 10%;
    margin-right: 10%;
    margin-top: 20%;
    box-sizing: border-box;
    background-color: #334155;
    color: #cbd5e1;
    border-radius: 0.5rem;
    box-shadow: 0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1);
    padding: 1.5rem;
    flex-grow: 1;
    transition: all 0.3s ease;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
  }
  .header, .footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .header {
    margin-bottom: 0.5rem;
  }
  .filename {
    font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
    font-size: 0.875rem;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    padding-right: 1rem;
  }
  .file-count {
    font-size: 0.875rem;
    font-weight: 600;
  }
  .progress-bar-container {
    width: 100%;
    background-color: #1e293b;
    border-radius: 9999px;
    height: 1rem;
    margin-bottom: 0.75rem;
    box-shadow: inset 0 2px 4px 0 rgb(0 0 0 / 0.05);
    overflow: hidden;
  }
  .progress-bar-inner {
    background-color: #0ea5e9;
    height: 100%;
    border-radius: 9999px;
    transition: width 0.3s ease-in-out;
  }
  .footer {
    font-size: 0.75rem;
    color: #94a3b8;
  }
  .status-text {
    display: flex;
    align-items: center;
    font-weight: 600;
  }
  .status-text.success {
    color: #38bdf8;
  }
  .status-text.error {
    color: #f87171;
  }
  .icon {
    height: 1rem;
    width: 1rem;
    display: inline-block;
    margin-right: 0.25rem;
  }
  .percentage {
    font-weight: 700;
    font-size: 1.125rem;
  }
</style>
