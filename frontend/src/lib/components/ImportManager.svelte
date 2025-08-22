<script lang="ts">
  import { onMount } from 'svelte'

  interface Import {
    id: number
    name: string
    progress: number
    status: 'uploading' | 'complete' | 'error'
  }

  let imports = $state<Import[]>([])
  let dragCounter = $state(0)
  let renamingId = $state<number | null>(null)
  let renameValue = $state('')

  // --- Mock API Functions ---
  const getImports = async (): Promise<Import[]> => {
    await new Promise((res) => setTimeout(res, 500))
    return [
      { id: 1, name: 'archive-2024-01.zip', progress: 100, status: 'complete' },
      { id: 2, name: 'project-files.zip', progress: 75, status: 'uploading' },
      { id: 3, name: 'customer-data.csv', progress: 0, status: 'error' },
    ]
  }

  const deleteImport = async (id: number) => {
    console.log(`Deleting import ${id}...`)
    await new Promise((res) => setTimeout(res, 300))
    imports = imports.filter((imp) => imp.id !== id)
  }

  const uploadFile = (file: File) => {
    const newId = Math.max(0, ...imports.map((i) => i.id)) + 1
    const newImport: Import = { id: newId, name: file.name, progress: 0, status: 'uploading' }
    imports = [...imports, newImport]

    // Simulate upload progress
    const interval = setInterval(() => {
      const imp = imports.find((i) => i.id === newId)
      if (imp) {
        imp.progress += 10
        if (imp.progress >= 100) {
          imp.progress = 100
          imp.status = 'complete'
          clearInterval(interval)
        }
        imports = [...imports]
      } else {
        clearInterval(interval)
      }
    }, 250)
  }

  // --- Component Logic ---
  onMount(async () => {
    imports = await getImports()
  })

  const handleDragEnter = (e: DragEvent) => {
    e.preventDefault()
    dragCounter++
  }

  const handleDragLeave = (e: DragEvent) => {
    e.preventDefault()
    dragCounter--
  }

  const handleDragOver = (e: DragEvent) => {
    e.preventDefault()
  }

  const handleDrop = (e: DragEvent) => {
    e.preventDefault()
    dragCounter = 0
    if (e.dataTransfer?.files) {
      ;[...e.dataTransfer.files].forEach(uploadFile)
    }
  }

  const startRename = (imp: Import) => {
    renamingId = imp.id
    renameValue = imp.name
  }

  const commitRename = (e: Event) => {
    const target = e.target as HTMLInputElement
    const imp = imports.find((i) => i.id === renamingId)
    if (imp) {
      imp.name = target.value
    }
    renamingId = null
  }

  const triggerFileInput = () => {
    const input = document.createElement('input')
    input.type = 'file'
    input.multiple = true
    input.onchange = (e) => {
      const files = (e.target as HTMLInputElement).files
      if (files) {
        ;[...files].forEach(uploadFile)
      }
    }
    input.click()
  }

  const getFileIcon = (filename: string) => {
    const ext = filename.split('.').pop()?.toLowerCase()
    switch (ext) {
      case 'zip': case 'rar': case '7z':
        return '📦'
      case 'csv': case 'xlsx': case 'xls':
        return '📊'
      case 'pdf':
        return '📄'
      case 'jpg': case 'jpeg': case 'png': case 'gif':
        return '🖼️'
      default:
        return '📁'
    }
  }
</script>

<div class="import-manager">
  <div class="header-section">
    <div class="header-content">
      <h2>Import Manager</h2>
      <p class="subtitle">Manage your file imports and uploads</p>
    </div>
    <button class="browse-btn" onclick={triggerFileInput}>
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M12 5v14m-7-7h14"/>
      </svg>
      Browse Files
    </button>
  </div>

  <div
    class="drop-zone"
    class:dragging-over={dragCounter > 0}
    ondragenter={handleDragEnter}
    ondragover={handleDragOver}
    ondragleave={handleDragLeave}
    ondrop={handleDrop}
  >
    <div class="drop-zone-content">
      <div class="upload-icon">
        <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
          <polyline points="14,2 14,8 20,8"/>
          <line x1="16" y1="13" x2="8" y2="13"/>
          <line x1="16" y1="17" x2="8" y2="17"/>
          <polyline points="10,9 9,9 8,9"/>
        </svg>
      </div>
      <h3>Drop files here to upload</h3>
      <p>or click Browse Files to select manually</p>
      <div class="supported-formats">
        <span>Supports ZIP, CSV, PDF, Images and more</span>
      </div>
    </div>
  </div>

  <div class="import-list">
    <div class="list-header">
      <h4>Recent Imports ({imports.length})</h4>
    </div>
    
    {#if imports.length === 0}
      <div class="empty-state">
        <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
          <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/>
          <polyline points="13,2 13,9 20,9"/>
        </svg>
        <p>No imports yet</p>
        <small>Upload your first file to get started</small>
      </div>
    {:else}
      <div class="items-container">
        {#each imports as item (item.id)}
          <div class="import-item" class:uploading={item.status === 'uploading'} class:error={item.status === 'error'}>
            <div class="item-icon">
              <span class="file-icon">{getFileIcon(item.name)}</span>
            </div>
            
            <div class="item-details">
              {#if renamingId === item.id}
                <input
                  type="text"
                  class="rename-input"
                  bind:value={renameValue}
                  onblur={commitRename}
                  onkeydown={(e) => e.key === 'Enter' && commitRename(e)}
                  autofocus
                />
              {:else}
                <div class="item-name-container">
                  <span class="item-name">{item.name}</span>
                  <div class="item-status">
                    {#if item.status === 'complete'}
                      <span class="status-complete">
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <polyline points="20,6 9,17 4,12"/>
                        </svg>
                        Complete
                      </span>
                    {:else if item.status === 'uploading'}
                      <span class="status-uploading">
                        <div class="spinner"></div>
                        Uploading {item.progress}%
                      </span>
                    {:else if item.status === 'error'}
                      <span class="status-error">
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <circle cx="12" cy="12" r="10"/>
                          <line x1="15" y1="9" x2="9" y2="15"/>
                          <line x1="9" y1="9" x2="15" y2="15"/>
                        </svg>
                        Failed
                      </span>
                    {/if}
                  </div>
                </div>
              {/if}
              
              {#if item.status === 'uploading'}
                <div class="progress-container">
                  <div class="progress-bar">
                    <div class="progress-fill" style="width: {item.progress}%"></div>
                  </div>
                  <span class="progress-text">{item.progress}%</span>
                </div>
              {/if}
            </div>
            
            <div class="item-actions">
              <button class="action-btn secondary" onclick={() => startRename(item)} aria-label="Rename">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                  <path d="m18.5 2.5 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                </svg>
              </button>
              <button class="action-btn danger" onclick={() => deleteImport(item.id)} aria-label="Delete">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3,6 5,6 21,6"/>
                  <path d="m19,6v14a2,2 0 0,1-2,2H7a2,2 0 0,1-2-2V6m3,0V4a2,2 0 0,1,2-2h4a2,2 0 0,1,2,2v2"/>
                </svg>
              </button>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .import-manager {
    background: #552958;
    color: #dbc9e0;
    padding: 2rem;
    display: flex;
    height: 100%;
    flex-direction: column;
    gap: 2rem;
    min-height: 600px;
  }

  .header-section {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    padding-bottom: 1.5rem;
    border-bottom: 1px solid rgba(148, 163, 184, 0.1);
  }

  .header-content h2 {
    margin: 0;
    font-size: 1.875rem;
    font-weight: 700;
    background: linear-gradient(135deg, #60a5fa, #a78bfa);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .subtitle {
    margin: 0.5rem 0 0 0;
    color: #94a3b8;
    font-size: 0.875rem;
  }

  .browse-btn {
    background: linear-gradient(135deg, #3b82f6, #8b5cf6);
    border: none;
    color: white;
    padding: 0.75rem 1.5rem;
    border-radius: 12px;
    cursor: pointer;
    font-weight: 600;
    font-size: 0.875rem;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    transition: all 0.2s ease;
    box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
  }

  .browse-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 25px rgba(59, 130, 246, 0.4);
  }

  .drop-zone {
    border: 2px dashed #374151;
    border-radius: 16px;
    padding: 3rem 2rem;
    text-align: center;
    background: rgba(15, 23, 42, 0.5);
    backdrop-filter: blur(10px);
    transition: all 0.3s ease;
    position: relative;
    overflow: hidden;
  }

  .drop-zone::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(135deg, rgba(59, 130, 246, 0.1), rgba(139, 92, 246, 0.1));
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  .drop-zone.dragging-over {
    border-color: #60a5fa;
    background: rgba(15, 23, 42, 0.8);
    transform: scale(1.02);
    box-shadow: 0 20px 40px rgba(59, 130, 246, 0.2);
  }

  .drop-zone.dragging-over::before {
    opacity: 1;
  }

  .drop-zone-content {
    position: relative;
    z-index: 1;
  }

  .upload-icon {
    color: #60a5fa;
    margin-bottom: 1rem;
  }

  .drop-zone h3 {
    margin: 0 0 0.5rem 0;
    font-size: 1.25rem;
    font-weight: 600;
  }

  .drop-zone p {
    margin: 0 0 1.5rem 0;
    color: #94a3b8;
  }

  .supported-formats {
    padding: 0.75rem 1.5rem;
    background: rgba(30, 41, 59, 0.8);
    border-radius: 8px;
    display: inline-block;
  }

  .supported-formats span {
    font-size: 0.75rem;
    color: #cbd5e1;
  }

  .import-list {
    flex: 1;
  }

  .list-header {
    margin-bottom: 1.5rem;
  }

  .list-header h4 {
    margin: 0;
    font-size: 1.125rem;
    font-weight: 600;
    color: #f1f5f9;
  }

  .empty-state {
    text-align: center;
    padding: 3rem 2rem;
    color: #64748b;
  }

  .empty-state svg {
    margin-bottom: 1rem;
    opacity: 0.5;
  }

  .empty-state p {
    margin: 0 0 0.5rem 0;
    font-size: 1.125rem;
    font-weight: 500;
  }

  .empty-state small {
    font-size: 0.875rem;
  }

  .items-container {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .import-item {
    background: rgba(30, 41, 59, 0.6);
    border: 1px solid rgba(148, 163, 184, 0.1);
    border-radius: 12px;
    padding: 1.25rem;
    display: flex;
    align-items: center;
    gap: 1rem;
    transition: all 0.2s ease;
    backdrop-filter: blur(10px);
  }

  .import-item:hover {
    background: rgba(30, 41, 59, 0.8);
    border-color: rgba(148, 163, 184, 0.2);
    transform: translateY(-1px);
  }

  .import-item.uploading {
    border-left: 4px solid #3b82f6;
  }

  .import-item.error {
    border-left: 4px solid #ef4444;
  }

  .item-icon {
    font-size: 1.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 48px;
    height: 48px;
    background: rgba(15, 23, 42, 0.8);
    border-radius: 10px;
    border: 1px solid rgba(148, 163, 184, 0.1);
  }

  .item-details {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .item-name-container {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .item-name {
    font-weight: 500;
    font-size: 0.95rem;
    color: #f1f5f9;
  }

  .rename-input {
    background: rgba(15, 23, 42, 0.8);
    border: 2px solid #3b82f6;
    border-radius: 8px;
    color: #f1f5f9;
    padding: 0.5rem;
    font-size: 0.95rem;
    width: 100%;
  }

  .rename-input:focus {
    outline: none;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
  }

  .item-status {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.75rem;
  }

  .status-complete {
    color: #10b981;
    display: flex;
    align-items: center;
    gap: 0.25rem;
    background: rgba(16, 185, 129, 0.1);
    padding: 0.25rem 0.5rem;
    border-radius: 6px;
  }

  .status-uploading {
    color: #3b82f6;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: rgba(59, 130, 246, 0.1);
    padding: 0.25rem 0.5rem;
    border-radius: 6px;
  }

  .status-error {
    color: #ef4444;
    display: flex;
    align-items: center;
    gap: 0.25rem;
    background: rgba(239, 68, 68, 0.1);
    padding: 0.25rem 0.5rem;
    border-radius: 6px;
  }

  .spinner {
    width: 12px;
    height: 12px;
    border: 2px solid rgba(59, 130, 246, 0.3);
    border-top: 2px solid #3b82f6;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }

  .progress-container {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .progress-bar {
    flex: 1;
    height: 6px;
    background: rgba(148, 163, 184, 0.2);
    border-radius: 3px;
    overflow: hidden;
  }

  .progress-fill {
    height: 100%;
    background: linear-gradient(90deg, #3b82f6, #8b5cf6);
    border-radius: 3px;
    transition: width 0.3s ease;
    position: relative;
  }

  .progress-fill::after {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
    animation: shimmer 2s infinite;
  }

  @keyframes shimmer {
    0% { transform: translateX(-100%); }
    100% { transform: translateX(100%); }
  }

  .progress-text {
    font-size: 0.75rem;
    color: #94a3b8;
    min-width: 35px;
    text-align: right;
  }

  .item-actions {
    display: flex;
    gap: 0.5rem;
  }

  .action-btn {
    background: rgba(30, 41, 59, 0.8);
    border: 1px solid rgba(148, 163, 184, 0.2);
    color: #cbd5e1;
    padding: 0.5rem;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .action-btn:hover {
    background: rgba(30, 41, 59, 1);
    border-color: rgba(148, 163, 184, 0.3);
  }

  .action-btn.secondary:hover {
    background: rgba(59, 130, 246, 0.2);
    border-color: #3b82f6;
    color: #60a5fa;
  }

  .action-btn.danger:hover {
    background: rgba(239, 68, 68, 0.2);
    border-color: #ef4444;
    color: #f87171;
  }
</style>
