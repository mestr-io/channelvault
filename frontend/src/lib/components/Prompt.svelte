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
</script>

{#if open}
  <div class="backdrop">
    <div class="modal">
      <h2 class="modal-title">{title}</h2>
      {#if message}<p class="modal-message">{message}</p>{/if}
      <input
        class="modal-input"
        type="text"
        bind:value
        autofocus
      />
      <div class="modal-actions">
        <button class="cancel" onclick={handleCancel}>Cancel</button>
        <button class="confirm" onclick={handleConfirm}>OK</button>
      </div>
    </div>
  </div>
{/if}

<style>
  .backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0,0,0,0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 50;
  }

  .modal {
    background: #fff;
    border-radius: 0.75rem;
    padding: 1.25rem;
    width: 90%;
    max-width: 400px;
    box-shadow: 0 10px 25px rgba(0,0,0,0.2);
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .modal-title {
    font-size: 1.125rem;
    font-weight: 600;
    margin: 0;
    color: #374151;
  }

  .modal-message {
    font-size: 0.9rem;
    color: #374151;
  }

  .modal-input {
    border: 1px solid #d1d5db;
    border-radius: 0.5rem;
    padding: 0.5rem;
    font-size: 0.9rem;
  }

  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
  }

  .cancel, .confirm {
    padding: 0.5rem 0.75rem;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    border: none;
  }

  .cancel {
    background: #f3f4f6;
    color: #374151;
  }
  .cancel:hover {
    background: #e5e7eb;
  }

  .confirm {
    background: #2563eb;
    color: white;
  }
  .confirm:hover {
    background: #1d4ed8;
  }
</style>
