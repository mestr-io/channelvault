<script lang="ts">
import { models } from '../../../wailsjs/go/models'
import ThreadInfo from './ThreadInfo.svelte'
import Attachments from './Attachments.svelte'
import Reactions from './Reactions.svelte'
import { Carta, Markdown } from 'carta-md'
import { getContext } from 'svelte'

type Props = {
  msg: models.Message
  selectedThread: string | null
}

let { msg, selectedThread = $bindable() }: Props = $props()

const carta: Carta = getContext('carta')

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

<div class="message">
    <img class="avatar" src={msg.avatar_url} alt={msg.user_name} />
    <div class="message-content">
        <div class="message-header">
            <strong>{msg.user_name}</strong>
            <span class="timestamp">
                {formatMessageTime(msg.ts)}
            </span>
        </div>

        <Markdown {carta} value={msg.text || ""} />
        <Attachments files={msg.files} />

        {#if msg.thread_ts}
            <ThreadInfo
                replyUsers={msg.reply_users}
                replyCount={msg.reply_count}
                threadTs={msg.thread_ts}
                bind:selectedThread
            />
        {/if}
        <Reactions reactions={msg.reactions} />
    </div>
</div>

<style>
    :global(.carta-viewer > :first-child) {
        margin-top: 2px;
    }
    .message {
        display: flex;
        padding: 0.75rem 0;
        transform: translateZ(0);
        -webkit-transform: translateZ(0);
        /* Optimize compositing */
        contain: layout style paint;
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

    .message-header strong {
        font-weight: 900;
    }
    .timestamp {
        color: #9e9ea6;
        font-size: 0.8rem;
        margin-left: 0.5rem;
    }
</style>
