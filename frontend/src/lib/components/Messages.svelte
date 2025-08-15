<script lang="ts">
    import { GetBase64Image } from '../../../wailsjs/go/app/App';
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

    <div class="message" data-ts={msg.ts}>
      <img class="avatar" src={msg.avatar_url} alt={msg.user_name} />
      <div class="message-content">
        <div class="message-header">
          <strong>{msg.user_name}</strong>
          <span class="timestamp">
            {formatMessageTime(msg.ts)}
          </span>
        </div>
        <p>{msg.text}</p>
        
        {#if msg.files}
          <div class="files-container">
            {#each msg.files as f}
              {#if f.mimetype.startsWith('image/')}
                <a href={f.url_private} target="_blank" rel="noopener noreferrer" class="file-preview image-preview">
                  {#await GetBase64Image(f.id) then image}
                    <img src={image} alt={f.title || f.name} />
                  {/await}
                </a>
              {:else}
                <a href={f.url_private} target="_blank" rel="noopener noreferrer" download={f.name} class="file-preview generic-file">
                  <div class="file-icon">
                    </div>
                  <div class="file-info">
                    <div class="file-name">{f.name || "file_not_found"}</div>
                    <div class="file-title">{f.title}</div>
                  </div>
                </a>
              {/if}
            {/each}
          </div>
        {/if}
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

  /* START: New Styles for File Attachments */
  .files-container {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-top: 8px;
  }

  .file-preview {
    text-decoration: none;
    color: inherit;
    display: inline-flex;
    border: 1px solid #ddd;
    border-radius: 8px;
    overflow: hidden;
    max-width: 360px;
  }
  
  .file-preview:hover {
    border-color: #999;
  }

  .image-preview img {
    max-width: 100%;
    max-height: 250px;
    display: block;
    object-fit: cover;
  }

  .generic-file {
    align-items: center;
    padding: 10px;
    background-color: #f8f8f8;
  }

  .file-icon {
    width: 32px;
    height: 32px;
    background-color: #e8e8e8; /* Placeholder for a real icon */
    border-radius: 6px;
    margin-right: 12px;
    flex-shrink: 0;
    /* In a real app, you'd set a background-image with an SVG icon */
  }

  .file-info {
    display: flex;
    flex-direction: column;
    overflow: hidden; /* Prevents long file names from breaking layout */
  }

  .file-name {
    font-weight: 700;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .file-title {
    font-size: 0.85rem;
    color: #555;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  /* END: New Styles for File Attachments */
</style>