<script lang="ts">
import { models } from '../../../wailsjs/go/models'
import ThreadInfo from './ThreadInfo.svelte'

type Props = {
  messages: models.Message[]
  isLoadingMessages: boolean
  hasMoreMessages: boolean
  selectedThread: string | null
}
let { messages, isLoadingMessages, hasMoreMessages, selectedThread = $bindable(null) }: Props = $props()

/**
 * Helper function to check if the day pill should be displayed.
 * It returns true for the first message or when the day is different from the previous message.
 */
const shouldShowDayPill = (currentMsg: models.Message, previousMsg: models.Message | undefined): boolean => {
  if (!previousMsg) {
    return true // Always show for the very first message
  }
  const currentDate = new Date(parseFloat(currentMsg.ts) * 1000)
  const previousDate = new Date(parseFloat(previousMsg.ts) * 1000)

  return currentDate.toDateString() !== previousDate.toDateString()
}

/**
 * Formats the timestamp for the day pill (e.g., "Friday, August 8").
 */
const formatDayPillDate = (ts: string): string => {
  const date = new Date(parseFloat(ts) * 1000)
  // Using 'en-US' locale for consistent formatting like "Friday, August 8"
  return date.toLocaleDateString('en-US', {
    weekday: 'long',
    month: 'long',
    day: 'numeric',
  })
}

/**
 * Formats the message timestamp to show only the time (e.g., "2:03 PM").
 */
const formatMessageTime = (ts: string): string => {
  const date = new Date(parseFloat(ts) * 1000)
  // Using 'en-US' locale to get AM/PM format
  return date.toLocaleTimeString('en-US', {
    hour: 'numeric',
    minute: '2-digit',
  })
}
</script>

{#if hasMoreMessages}
  <div class="loader">
    {#if isLoadingMessages}<span>Loading history...</span>{/if}
  </div>
{/if}

<div class="messages-container" id="messages-container">
  {#each messages as msg, i}
    {#if shouldShowDayPill(msg, messages[i - 1])}
      <div class="day-divider">
        <hr />
        <span>{formatDayPillDate(msg.ts)}</span>
        <hr />
      </div>
    {/if}

    <div class="message">
      <img class="avatar" src={msg.avatar_url} alt={msg.user_name} />
      <div class="message-content">
        <div class="message-header">
          <strong>{msg.user_name}</strong>
          <span class="timestamp">
            {formatMessageTime(msg.ts)}
          </span>
        </div>
        <p>{msg.text}</p>

        {#if msg.thread_ts}
          <ThreadInfo
            replyUsers={msg.reply_users}
            replyCount={msg.reply_count}
            threadTs={msg.thread_ts}
            bind:selectedThread
          />
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
  
  /* STYLES FOR THE DAY PILL */
  .day-divider {
    display: flex;
    align-items: center;
    text-align: center;
    color: #555;
    padding: 1rem 0;
    font-size: 0.8rem;
  }

  .day-divider hr {
    flex-grow: 1;
    border: none;
    border-top: 1px solid #ddd;
    margin: 0;
  }

  .day-divider span {
    padding: 0.25rem 0.75rem;
    border: 1px solid #ddd;
    border-radius: 24px;
    background-color: white;
    font-weight: bold;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  }
</style>