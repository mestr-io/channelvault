<script lang="ts">
let {
  open = false, // controls visibility
  title = 'Enter value',
  message = '',
  defaultValue = '',
  onConfirm = (v: string) => {},
  onCancel = () => {},
} = $props()

let value: string = $state(defaultValue)

function handleConfirm() {
  onConfirm(value)
  open = false
}

function handleCancel() {
  onCancel()
  open = false
}

// Handle Escape key to close the modal
function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Escape') {
        handleCancel();
    }
}
</script>

{#if open}
  <div class="backdrop" on:click={handleCancel}>
    <div class="modal" on:click|stopPropagation on:keydown={handleKeydown}>
      <h2 class="modal-title">{title}</h2>
      {#if message}<p class="modal-message">{message}</p>{/if}
      <input
        class="modal-input"
        type="text"
        bind:value
        autofocus
        on:keydown={e => { if (e.key === 'Enter') handleConfirm() }}
      />
      <div class="modal-actions">
        <button class="btn cancel" on:click={handleCancel}>Cancel</button>
        <button class="btn confirm" on:click={handleConfirm}>OK</button>
      </div>
    </div>
  </div>
{/if}

<style>
  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  @keyframes scaleUp {
    from { transform: scale(0.95); opacity: 0; }
    to { transform: scale(1); opacity: 1; }
  }

  .backdrop {
    position: fixed;
    inset: 0;
    background: rgba(18, 20, 23, 0.7);
    backdrop-filter: blur(4px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 50;
    animation: fadeIn 150ms ease-out;
  }

  .modal {
    background: #2d3748; /* Dark background from card component */
    color: #e2e8f0;
    border-radius: 0.75rem;
    padding: 1.5rem;
    width: 90%;
    max-width: 400px;
    box-shadow: 0 10px 25px rgba(0,0,0,0.5);
    display: flex;
    flex-direction: column;
    gap: 1rem;
    border: 1px solid #4a5568;
    animation: scaleUp 150ms ease-out;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  }

  .modal-title {
    font-size: 1.25rem;
    font-weight: 600;
    margin: 0;
    color: #f7fafc;
  }

  .modal-message {
    font-size: 0.9rem;
    color: #a0aec0; /* Lighter, secondary text color */
    margin: -0.5rem 0 0;
    line-height: 1.5;
  }

  .modal-input {
    background: #25292f;
    border: 1px solid #4a5568;
    color: #e2e8f0;
    border-radius: 0.5rem;
    padding: 0.6rem 0.75rem;
    font-size: 1rem;
    outline: none;
    transition: all 150ms ease-in-out;
  }

  .modal-input:focus {
    border-color: #4a90e2; /* Slack blue accent */
    box-shadow: 0 0 0 3px rgba(74, 144, 226, 0.3);
  }

  .modal-actions {
    margin-top: 0.5rem;
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
  }

  .btn {
    padding: 0.6rem 1rem;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 600;
    cursor: pointer;
    border: none;
    transition: background-color 200ms ease;
  }

  .cancel {
    background: #4a5568;
    color: #f7fafc;
  }
  .cancel:hover {
    background: #5a667a;
  }

  .confirm {
    background: #4a90e2; /* Slack blue */
    color: white;
  }
  .confirm:hover {
    background: #357ABD;
  }
</style>
