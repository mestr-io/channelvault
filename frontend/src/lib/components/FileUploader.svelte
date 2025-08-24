<script lang="ts">
  import { onMount } from 'svelte';
  import { ListImports } from '../../../wailsjs/go/app/App'
  import { models } from '../../../wailsjs/go/models'
  
  let isDragging = $state(false);
  let dragCounter = $state(0);
  let imports: models.SlackImport[] = $state([]);

  onMount(async () => {
    imports = await ListImports();
    // Add a dummy error card for demonstration
    imports.unshift({
      ID: 'demo-error',
      Name: 'Failed Import Example',
      Path: '/path/to/failed/slack-export.zip',
      Created: new Date().toISOString(),
      Finished: '',
      NumMessages: 0,
      Progress: 45,
      Error: 'Failed to parse channel data: Invalid JSON structure in channels.json. The file appears to be corrupted or incomplete.'
    });
  });

  // --- Drag & Drop Event Handlers ---
  function handleDragEnter(event: DragEvent): void {
    event.preventDefault();
    dragCounter++;
    if (dragCounter === 1) {
      isDragging = true;
    }
  }

  function handleDragOver(event: DragEvent): void {
    event.preventDefault();
    if (!isDragging) {
      isDragging = true;
    }
  }

  function handleDragLeave(event: DragEvent): void {
    event.preventDefault();
    dragCounter--;
    if (dragCounter === 0) {
      isDragging = false;
    }
  }

  function handleDrop(event: DragEvent): void {
    event.preventDefault();
    dragCounter = 0;
    isDragging = false;
  }

  function formatDate(dateString: string): string {
    if (!dateString) return 'Not started';
    const date = new Date(dateString);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'});
  }

  function getStatusColor(progress: number, error: string): string {
    if (error) return '#ef4444'; // red for errors
    if (progress === 0) return '#6b7280'; // gray
    if (progress < 50) return '#f59e0b'; // amber
    if (progress < 100) return '#3b82f6'; // blue
    return '#10b981'; // green
  }

  function getStatusText(progress: number, finished: string, error: string): string {
    if (error) return 'Error';
    if (finished) return 'Completed';
    if (progress === 0) return 'Pending';
    if (progress < 100) return 'In Progress';
    return 'Processing';
  }

  function handleRename(importId: string, currentName: string) {
    const newName = prompt('Enter new name for import:', currentName);
    if (newName && newName !== currentName) {
      // Here you would call your rename function
      // RenameImport(importId, newName);
      console.log(`Renaming import ${importId} to "${newName}"`);
      
      // Update the local state for demonstration
      imports = imports.map(imp => 
        imp.ID === importId ? { ...imp, Name: newName } : imp
      );
    }
  }

  function handleDelete(importId: string, importName: string) {
    if (confirm(`Are you sure you want to delete "${importName}"? This action cannot be undone.`)) {
      // Here you would call your delete function
      // DeleteImport(importId);
      console.log(`Deleting import ${importId}`);
      
      // Update the local state for demonstration
      imports = imports.filter(imp => imp.ID !== importId);
    }
  }
</script>

<main>
  <div
    class="drop-zone"
    class:is-dragging={isDragging}
    role="button"
    tabindex="0"
    ondragenter={handleDragEnter}
    ondragleave={handleDragLeave}
    ondragover={handleDragOver}
    ondrop={handleDrop}
  >
    <div class="add-file">
      <svg
        width="24"
        height="24"
        viewBox="0 0 24 24"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          d="M12 6C12.5523 6 13 6.44772 13 7V11H17C17.5523 11 18 11.4477 18 12C18 12.5523 17.5523 13 17 13H13V17C13 17.5523 12.5523 18 12 18C11.4477 18 11 17.5523 11 17V13H7C6.44772 13 6 12.5523 6 12C6 11.4477 6.44772 11 7 11H11V7C11 6.44772 11.4477 6 12 6Z"
          fill="currentColor"
        />
        <path
          fill-rule="evenodd"
          clip-rule="evenodd"
          d="M5 22C3.34315 22 2 20.6569 2 19V5C2 3.34315 3.34315 2 5 2H19C20.6569 2 22 3.34315 22 5V19C22 20.6569 20.6569 22 19 22H5ZM4 19C4 19.5523 4.44772 20 5 20H19C19.5523 20 20 19.5523 20 19V5C20 4.44772 19.5523 4 19 4H5C4.44772 4 4 4.44772 4 5V19Z"
          fill="currentColor"
        />
      </svg>
    </div>
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
    </div>
  </div>
  
  <div class="imports">
    {#each imports as card}
    <div class="card">
      <div class="card-header">
        <h3 class="card-title">{card.Name || 'Untitled Import'}</h3>
        <div class="card-actions">
          <div class="status-badge" style="background-color: {getStatusColor(card.Progress || 0, card.Error)}">
            {getStatusText(card.Progress || 0, card.Finished, card.Error)}
          </div>
          <button 
            class="action-button rename-button" 
            onclick={() => handleRename(card.ID, card.Name)}
            title="Rename import"
          >
            <svg viewBox="0 0 20 20" fill="currentColor">
              <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
            </svg>
          </button>
          <button 
            class="action-button delete-button" 
            onclick={() => handleDelete(card.ID, card.Name)}
            title="Delete import"
          >
            <svg viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
            </svg>
          </button>
        </div>
      </div>
      
      {#if card.Error}
      <div class="error-section">
        <div class="error-header">
          <svg class="error-icon" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
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
              <path fill-rule="evenodd" d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z" clip-rule="evenodd" />
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
              <path fill-rule="evenodd" d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z" clip-rule="evenodd" />
            </svg>
            <div class="info-text">
              <span class="info-label">Started</span>
              <span class="info-value">{formatDate(card.Created)}</span>
            </div>
          </div>
          
          {#if card.Finished}
          <div class="info-item">
            <svg class="info-icon" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
            <div class="info-text">
              <span class="info-label">Completed</span>
              <span class="info-value">{formatDate(card.Finished)}</span>
            </div>
          </div>
          {/if}
        </div>

        <div class="info-row">
          <div class="info-item">
            <svg class="info-icon" viewBox="0 0 20 20" fill="currentColor">
              <path d="M2.003 5.884L10 9.882l7.997-3.998A2 2 0 0016 4H4a2 2 0 00-1.997 1.884z" />
              <path d="M18 8.118l-8 4-8-4V14a2 2 0 002 2h12a2 2 0 002-2V8.118z" />
            </svg>
            <div class="info-text">
              <span class="info-label">Messages</span>
              <span class="info-value">{card.NumMessages?.toLocaleString() || '0'}</span>
            </div>
          </div>
          
          <div class="info-item">
            <svg class="info-icon" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M3 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd" />
            </svg>
            <div class="info-text">
              <span class="info-label">Import ID</span>
              <span class="info-value">{card.ID}</span>
            </div>
          </div>
        </div>
      </div>
      
      <div class="progress-section">
        <div class="progress-header">
          <span class="progress-label">Progress</span>
          <span class="progress-percentage">{card.Progress || 0}%</span>
        </div>
        <div class="progress-bar">
          <div 
            class="progress-fill" 
            style="width: {card.Progress || 0}%; background-color: {getStatusColor(card.Progress || 0)}"
          ></div>
        </div>
      </div>
    </div>
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
    display: grid;
    gap: 2rem;
    grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
    width: 100%;
    max-width: 1200px;
  }

  .card {
    background: linear-gradient(135deg, #ffffff 0%, #f8fafc 100%);
    border-radius: 1rem;
    padding: 1.5rem;
    box-shadow: 
      0 4px 6px -1px rgba(0, 0, 0, 0.1),
      0 2px 4px -1px rgba(0, 0, 0, 0.06),
      0 0 0 1px rgba(0, 0, 0, 0.05);
    transition: all 0.3s ease;
    border: 1px solid #e2e8f0;
  }

  .card:hover {
    transform: translateY(-2px);
    box-shadow: 
      0 10px 25px -5px rgba(0, 0, 0, 0.1),
      0 10px 10px -5px rgba(0, 0, 0, 0.04),
      0 0 0 1px rgba(0, 0, 0, 0.05);
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
    align-items: flex-start;
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
    background: linear-gradient(90deg, currentColor 0%, currentColor 100%);
  }

  .drop-zone {
    margin-top: 5rem;
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
  }

  .add-file {
    position: absolute;
    top: 1rem;
    right: 1rem;
    cursor: pointer;
  }

  .drop-zone.is-dragging {
    border-color: #3b82f6;
    box-shadow:
      0 0 30px rgba(59, 130, 246, 0.7),
      0 0 15px rgba(59, 130, 246, 0.5);
    transform: scale(1.05);
  }

  .icons-container {
    display: flex;
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
    
    .card {
      padding: 1.25rem;
    }
    
    .info-row {
      flex-direction: column;
      gap: 0.75rem;
    }
    
    .card-actions {
      gap: 0.25rem;
    }
    
    .card-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 0.75rem;
    }
  }
</style>
