<script>
  // Example data (you can replace this with your actual data)
  let files = [
    {
      name: "Report_Q4.pdf",
      start: "2023-11-15",
      finish: "2023-11-16",
      progress: "99/290",
      status: "in-progress"
    },
    {
      name: "Budget_2024.xlsx",
      start: "2023-11-17",
      finish: "-",
      progress: "100/100",
      status: "complete"
    },
    {
      name: "Meeting_Notes.docx",
      start: "2023-11-18",
      finish: "-",
      progress: "55/200",
      status: "in-progress"
    }
  ];
</script>

<style>
  :global(body) {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
    background-color: #f8f9fa;
    color: #24292e;
  }

  .container {
    padding: 2rem;
    max-width: 1200px;
    margin: 0 auto;
  }

  .drop-zone {
    background: #ffffff;
    border: 2px dashed #d0d7dd;
    border-radius: 12px;
    padding: 2rem;
    text-align: center;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
    transition: box-shadow 0.2s ease;
  }

  .drop-zone:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }

  .drop-zone p {
    margin: 0 0 1.5rem 0;
    color: #6e7680;
  }

  .import-list {
    margin-top: 2rem;
  }

  .list-header {
    margin-bottom: 1.5rem;
  }

  .list-header h4 {
    font-size: 1.125rem;
    font-weight: 600;
    color: #24292e;
  }

  .items-container {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .import-item {
    background: #ffffff;
    border: 1px solid #eaecef;
    border-radius: 12px;
    padding: 1.5rem;
    display: flex;
    align-items: center;
    gap: 1.5rem;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
    transition: box-shadow 0.2s ease;
  }

  .import-item:hover {
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
  }

  .item-icon {
    font-size: 1.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 48px;
    height: 48px;
    background: #f3f4f6;
    border-radius: 10px;
    flex-shrink: 0;
    color: #3b82f6;
  }

  .item-details {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .item-name {
    font-weight: 600;
    color: #1f2937;
  }

  .dates {
    display: flex;
    justify-content: space-between;
    color: #6b7280;
    font-size: 0.875rem;
  }

  .progress-container {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: #1f2937;
    font-size: 0.875rem;
  }

  .progress-bar {
    flex: 1;
    height: 6px;
    background-color: #e5e7eb;
    border-radius: 3px;
    overflow: hidden;
    position: relative;
  }

  .progress-fill {
    height: 100%;
    width: 60%;
    background-color: #3b82f6;
    transition: width 0.3s ease;
  }

  .progress-complete {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    color: #10b981;
  }

  .progress-complete svg {
    width: 16px;
    height: 16px;
  }

  .action-buttons {
    display: flex;
    gap: 0.5rem;
  }

  .action-btn {
    padding: 0.5rem 0.75rem;
    border: none;
    border-radius: 6px;
    font-weight: 500;
    cursor: pointer;
    transition: 0.2s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.25rem;
  }

  .action-btn-primary {
    background-color: #3b82f6;
    color: #ffffff;
  }

  .action-btn-primary:hover {
    background-color: #2563eb;
  }

  .action-btn-secondary {
    background-color: #f3f4f6;
    color: #3b82f6;
    border: 1px solid #d0d7dd;
  }

  .action-btn-secondary:hover {
    background-color: #e5e7eb;
  }
</style>

<div class="container">
  <div class="drop-zone">
    <p>Drag and drop your files here</p>
  </div>

  <div class="import-list">
    <div class="list-header">
      <h4>Files</h4>
    </div>

    <div class="items-container">
      {#each files as file}
        <div class="import-item">
          <div class="item-icon">📄</div>
          <div class="item-details">
            <div class="item-name">{file.name}</div>
            <div class="dates">
              <span>{file.start}</span>
              <span>{file.finish}</span>
            </div>
            <div class="progress-container">
              {#if file.status === "complete"}
                <div class="progress-complete">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <polyline points="20 6 9 17 4 12"></polyline>
                  </svg>
                  <span>Completed</span>
                </div>
              {:else}
                <span>{file.progress}</span>
              {/if}
            </div>
          </div>
          <div class="action-buttons">
            <button class="action-btn action-btn-primary">
              <span>_rename</span>
            </button>
            <button class="action-btn action-btn-secondary">
              <span>delete</span>
            </button>
          </div>
        </div>
      {/each}
    </div>
  </div>
</div>
