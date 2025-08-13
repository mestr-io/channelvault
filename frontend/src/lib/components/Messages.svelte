<script lang="ts">
// This part remains the same.
// NOTE: I'm assuming 'msg.reply_user_avatars' is an array of avatar URLs for the users in the thread.
// You will need to make sure your data structure provides this.
import { GetUserAvatar } from '../../../wailsjs/go/app/App'
import { models } from '../../../wailsjs/go/models'

type Props = {
  messages: models.Message[]
  isLoadingMessages: boolean
  hasMoreMessages: boolean
  selectedThread: string | undefined
}
let { messages, isLoadingMessages, hasMoreMessages, selectedThread = $bindable() }: Props = $props()
const cache = $state(new Map<string, string>())

async function getUser(userId: string): Promise<string> {
  if (cache.has(userId)) {
    return cache.get(userId)!
  }

  try {
    const url = await GetUserAvatar(userId)
    cache.set(userId, url)

    return url
  } catch (error) {
    console.error(`Failed to get user info for ${userId}:`, error)
    return ''
  }
}
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
          <div class="thread-info">
            <div class="reply-avatars">
              {#each (new Set(msg.reply_users) || [])
                .values()
                .toArray()
                .slice(0, 5) as user_id}
                {#await getUser(user_id)}
                  <div class="reply-avatar-placeholder"></div>
                {:then avatar_url}
                  <img class="reply-avatar" src={avatar_url} alt="Reply user" />
                {:catch error}
                  <div class="reply-avatar-placeholder error">!</div>
                {/await}
              {/each}
            </div>

            <a href="#" class="reply-link" onclick={() => selectedThread = msg.thread_ts}>
              {msg.reply_count}
              {msg.reply_count === 1 ? "reply" : "replies"}
            </a>

            <span class="view-thread-prompt">View thread</span>
          </div>
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

  /* --- NEW STYLES FOR THE THREAD --- */

  .thread-info {
    display: flex;
    align-items: center;
    margin-top: 0.5rem;
    cursor: pointer; /* Make the whole area feel clickable */
  }

  .reply-avatars {
    display: flex;
  }

  .reply-avatar {
    width: 24px;
    height: 24px;
    border-radius: 50%;
    border: 2px solid white; /* Creates separation between avatars */
    object-fit: cover;
  }

  /* This selector makes the avatars overlap */
  .reply-avatar:not(:first-child) {
    margin-left: -8px;
  }

  .reply-link {
    color: #1264a3; /* A nice blue, similar to Slack's */
    font-weight: 600;
    font-size: 0.9rem;
    text-decoration: none;
    margin-left: 0.5rem;
  }

  .reply-link:hover {
    text-decoration: underline;
  }

  .view-thread-prompt {
    color: #9e9ea6;
    font-size: 0.9rem;
    margin-left: 0.5rem;
  }
</style>
