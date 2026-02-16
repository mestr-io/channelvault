<script lang="ts">
import { SearchMessages } from '../../../wailsjs/go/app/App'
import type { models } from '../../../wailsjs/go/models'

let show = $state(false)
let query = $state('')
let results = $state<models.SearchResult[]>([])
let isLoading = $state(false)
let error = $state('')

export function showOverlay() {
  show = true
  query = ''
  results = []
  error = ''
}

export function closeOverlay() {
  show = false
}

async function handleSearch() {
  if (!query.trim()) return
  isLoading = true
  error = ''
  results = []
  try {
    const res = await SearchMessages(query)
    results = res || []
  } catch (err: any) {
    error = err.message || 'Failed to search messages'
  } finally {
    isLoading = false
  }
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter') {
    handleSearch()
  } else if (e.key === 'Escape') {
    closeOverlay()
  }
}

// Exposed callback for navigation
let { onNavigate } = $props()

function selectResult(result: models.SearchResult) {
  if (onNavigate) {
    onNavigate(result)
    closeOverlay()
  }
}
</script>

{#if show}
  <div class="overlay" role="dialog" aria-modal="true">
    <div class="modal">
      <div class="modal-header">
        <h2>Search Messages</h2>
        <button class="close-btn" onclick={closeOverlay} aria-label="Close search">
          <svg viewBox="0 0 24 24" width="24" height="24" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>

      <div class="search-box">
        <input
          type="text"
          placeholder="Search..."
          bind:value={query}
          onkeydown={handleKeydown}
          autofocus
        />
        <button onclick={handleSearch} disabled={isLoading}>
          {#if isLoading}
            Searching...
          {:else}
            Search
          {/if}
        </button>
      </div>

      {#if error}
        <div class="error-message">{error}</div>
      {/if}

      <div class="results-list">
        {#if results.length === 0 && !isLoading && query}
          <div class="no-results">No results found.</div>
        {:else}
          {#each results as result}
            <button class="result-item" onclick={() => selectResult(result)}>
              <div class="result-header">
                <span class="user-name">{result.user_name || result.user_id}</span>
                <span class="channel-name">#{result.channel_name}</span>
                <span class="timestamp">{new Date(parseFloat(result.ts) * 1000).toLocaleString()}</span>
              </div>
              <div class="result-text">{result.text}</div>
            </button>
          {/each}
        {/if}
      </div>
    </div>
  </div>
{/if}

<style>
  .overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }

  .modal {
    background-color: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    width: 600px;
    max-width: 90%;
    max-height: 80vh;
    display: flex;
    flex-direction: column;
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
  }

  .modal-header h2 {
    margin: 0;
    font-size: 1.5rem;
    color: #333;
  }

  .close-btn {
    background: none;
    border: none;
    cursor: pointer;
    padding: 0.5rem;
    color: #666;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    transition: background-color 0.2s;
  }

  .close-btn:hover {
    background-color: #f0f0f0;
  }

  .search-box {
    display: flex;
    gap: 10px;
    margin-bottom: 1rem;
  }

  input {
    flex-grow: 1;
    padding: 0.8rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
  }

  button {
    padding: 0.8rem 1.5rem;
    background-color: #552958;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-weight: 600;
    transition: background-color 0.2s;
  }

  button:hover {
    background-color: #3e0e42;
  }

  button:disabled {
    background-color: #ccc;
    cursor: not-allowed;
  }

  .error-message {
    color: #d32f2f;
    margin-bottom: 1rem;
    padding: 0.5rem;
    background-color: #ffebee;
    border-radius: 4px;
  }

  .results-list {
    flex-grow: 1;
    overflow-y: auto;
    border-top: 1px solid #eee;
    margin-top: 1rem;
  }

  .no-results {
    padding: 2rem;
    text-align: center;
    color: #666;
  }

  .result-item {
    display: block;
    width: 100%;
    text-align: left;
    padding: 1rem;
    border: none;
    border-bottom: 1px solid #eee;
    background: none;
    cursor: pointer;
    transition: background-color 0.1s;
  }

  .result-item:hover {
    background-color: #f9f9f9;
  }

  .result-header {
    display: flex;
    gap: 10px;
    font-size: 0.85rem;
    color: #666;
    margin-bottom: 0.3rem;
  }

  .user-name {
    font-weight: bold;
    color: #333;
  }

  .channel-name {
    color: #552958;
  }

  .result-text {
    font-size: 0.95rem;
    color: #333;
    white-space: pre-wrap; /* Preserve line breaks */
    display: -webkit-box;
    -webkit-line-clamp: 2; /* Limit to 2 lines */
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
</style>
