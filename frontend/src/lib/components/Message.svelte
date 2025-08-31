<script lang="ts">
import { component } from '@cartamd/plugin-component'
import { initializeComponents, svelte } from '@cartamd/plugin-component/svelte'
import { emoji } from '@cartamd/plugin-emoji'
import { Carta, Markdown, type Options } from 'carta-md'
import DOMPurify from 'isomorphic-dompurify'
import type { models } from '../../../wailsjs/go/models'
import Attachments from './Attachments.svelte'
import Link from './Link.svelte'
import Reactions from './Reactions.svelte'
import ThreadInfo from './ThreadInfo.svelte'

type Props = {
  msg: models.Message
  selectedThread: string | null
}

let { msg, selectedThread = $bindable() }: Props = $props()

// Markdown
const mapped = [svelte('a', Link)]
const options: Options = {
  extensions: [emoji(), component(mapped, initializeComponents)],
  sanitizer: DOMPurify.sanitize,
}
export const carta = new Carta(options)
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
        <Attachments message={msg} />
        {#if msg.thread_ts}
            <ThreadInfo
                replyUsers={msg.reply_users}
                replyCount={msg.reply_count}
                threadTs={msg.thread_ts}
                bind:selectedThread
            />
        {/if}
        {#if msg.reactions}
            <Reactions reactions={msg.reactions} />
        {/if}
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
