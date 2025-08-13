<script lang="ts">
// This part remains the same.
// NOTE: I'm assuming 'msg.reply_user_avatars' is an array of avatar URLs for the users in the thread.
// You will need to make sure your data structure provides this.
import { models } from '../../../wailsjs/go/models'
import Thread from './Thread.svelte'

type Props = {
  messages: models.Message[]
  isLoadingMessages: boolean
  hasMoreMessages: boolean
}
const { messages, isLoadingMessages, hasMoreMessages }: Props = $props()
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
          <Thread replyUsers={msg.reply_users} replyCount={msg.reply_count} />
        {/if}
      </div>
    </div>
  {/each}
</div>

<style>
  .loader {
    text-align: center;
    padding: 1rem;
    min-height: 40px;
  }
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
  .message-content {
    flex: 1;
    max-width: calc(100% - 50px);
  }
  .message-content p {
    margin: 0;
    line-height: 1.45;
    display: block;
    max-width: 98%;
    word-wrap: break-word;
    overflow-wrap: break-word;
    white-space: normal;
  }
  .message-header strong {
    font-weight: 900;
  }
  .timestamp {
    color: #9e9ea6;
    font-size: 0.8rem;
    margin-left: 0.5rem;
  }
</style>
