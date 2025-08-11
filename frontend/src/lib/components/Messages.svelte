<script lang="ts">
  const { messages, isLoadingMessages, hasMoreMessages } = $props();
</script>

{#if hasMoreMessages}
  <div class="loader">
    {#if isLoadingMessages}<span>Loading history...</span>{/if}
  </div>
{/if}
<div class="messages-container" id="messages-container">
  {#each messages as msg}
    <div class="message">
      <img class="avatar" src={msg.avatar_url} alt={msg.user_name} />
      <div class="message-content">
        <div class="message-header">
          <strong>{msg.user_name}</strong>
          <span class="timestamp">
            {new Date(parseFloat(msg.ts) * 1000).toLocaleString()}
          </span>
        </div>
        <p>{msg.text}</p>
        {#if msg.thread_ts}
          <p>{msg.reply_count} replies <button>View thread</button></p>
        {/if}
      </div>
    </div>
  {/each}
</div>

<style>
  .messages-container {
    margin: 1rem;
    margin-top: auto;
  }

  .message {
    display: flex;
    padding: 0.75rem 0;
  }
  .avatar {
    width: 40px;
    height: 40px;
    border-radius: 6px;
    margin-right: 1rem;
    object-fit: cover;
  }
  .message-content p {
    margin: 0;
    line-height: 1.45;
    white-space: pre-wrap; /* ADDED: This is the key to fixing text formatting */
  }
  .message-header strong {
    font-weight: 900;
  }
  .timestamp {
    color: #9e9ea6; /* MODIFIED: Lighter gray for timestamp */
    font-size: 0.8rem;
    margin-left: 0.5rem;
  }
</style>
