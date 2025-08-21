<script lang="ts">
import { onMount } from 'svelte'
import { ListImports } from '../../../wailsjs/go/app/App'
import { models } from '../../../wailsjs/go/models'

// svelte-ignore non_reactive_update
enum Status {
  Ready = 'ready',
  Processing = 'processing',
  Error = 'error',
}

let imports: models.SlackImport[] = $state([])
let isOpen = $state(false)

onMount(async () => {
  const fn = async () => {
    imports = await ListImports()
  }
  fn()
})

function removeImport(id: number) {
  imports = imports.filter(imp => imp.ID !== id)
}

function viewImport(id: number) {
  console.log('View import', id)
}

function toggleOpen() {
  isOpen = !isOpen
}

function getStatus(imp: models.SlackImport): Status {
  if (imp.Error) {
    return Status.Error
  } else if (imp.Finished) {
    return Status.Ready
  } else {
    return Status.Processing
  }
}
</script>

<div class="imports">
  <button class="header" onclick={toggleOpen}>
    <svg
      class="w-5 h-5 mr-3"
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
      xmlns="http://www.w3.org/2000/svg"
      ><path
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="2"
        d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
      ></path></svg
    >
    <span>Imports ({imports.length})</span>
    <svg
      class:rotated={!isOpen}
      xmlns="http://www.w3.org/2000/svg"
      width="14"
      height="14"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
    >
      <polyline points="6 9 12 15 18 9"></polyline>
    </svg>
  </button>

  {#if isOpen}
    <div class="import-list">
      <div class="import-header">
        <span>Add new import</span>
        <button class="add-btn" title="Add import">＋</button>
      </div>

      {#each imports as imp (imp.ID)}
        <div class="import-row">
          <div class="status-col">
            <span class="dot {getStatus(imp)}"></span>
          </div>

          <div class="name-col">{imp.Path}</div>

          <div class="meta-col">
            {#if getStatus(imp) === Status.Ready}
              <span>{imp.NumMessages} messages</span>
            {:else if getStatus(imp) === Status.Processing}
              <span class="processing">Processing…</span>
            {:else}
              <span class="error">Error</span>
            {/if}
          </div>

          <div class="actions">
            <!-- svelte-ignore a11y_consider_explicit_label -->
            <button
              class="icon-btn"
              title="View"
              onclick={() => viewImport(imp.id)}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
                <circle cx="12" cy="12" r="3"></circle>
              </svg>
            </button>
            <!-- svelte-ignore a11y_consider_explicit_label -->
            <button
              class="icon-btn"
              title="Remove"
              onclick={() => removeImport(imp.id)}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <polyline points="3 6 5 6 21 6"></polyline>
                <path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"></path>
                <path d="M10 11v6"></path>
                <path d="M14 11v6"></path>
              </svg>
            </button>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .imports {
    color: #dbc9e0;
    font-size: 0.9rem;
  }

  .header {
    width: 100%;
    background: none;
    border: none;
    color: #dbc9e0;
    font-weight: 600;
    cursor: pointer;
    padding: 0.75rem 1rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  .header:hover {
    background-color: #663867;
  }
  .header svg {
    transition: transform 0.2s ease;
  }
  .header svg.rotated {
    transform: rotate(180deg);
  }

  .import-list {
    padding: 0 0.25rem 0.5rem;
  }

  .import-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.3rem 0.5rem;
    font-size: 0.85rem;
    color: #cbbcd0;
  }

  .add-btn {
    background: #663867;
    color: #fff;
    border: none;
    border-radius: 4px;
    padding: 0.2rem 0.5rem;
    cursor: pointer;
    font-size: 0.8rem;
  }
  .add-btn:hover {
    background: #814e83;
  }

  .import-row {
    display: grid;
    grid-template-columns: 20px 1fr auto auto;
    align-items: center;
    padding: 0.35rem 0.5rem;
    border-radius: 4px;
    margin: 0.15rem 0;
    gap: 0.5rem;
  }
  .import-row:hover {
    background: #663867;
  }

  .status-col {
    display: flex;
    justify-content: center;
    align-items: center;
  }
  .dot {
    display: inline-block;
    width: 8px;
    height: 8px;
    border-radius: 50%;
  }
  .dot.ready {
    background: #98c379;
  }
  .dot.processing {
    background: #f0c674;
  }
  .dot.error {
    background: #e06c75;
  }

  .name-col {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    color: #f9edff;
    font-size: 0.9rem;
  }

  .meta-col {
    font-size: 0.8rem;
    color: #cbbcd0;
    padding-right: 0.5rem;
  }
  .meta-col .processing {
    color: #f0c674;
  }
  .meta-col .error {
    color: #e06c75;
  }

  .actions {
    display: flex;
    gap: 0.3rem;
  }
  .icon-btn {
    background: none;
    border: none;
    cursor: pointer;
    color: #dbc9e0;
    padding: 0.2rem;
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .icon-btn svg {
    stroke: currentColor;
  }
  .icon-btn:hover {
    background-color: #814e83;
    color: #fff;
  }
</style>
