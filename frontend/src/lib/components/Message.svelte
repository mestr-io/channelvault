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
import { circle } from './icons'
import type { ThreadID } from '../types'
import { GetUserHandle } from '../../../wailsjs/go/app/App'
import { getContext } from 'svelte'

type Props = {
  group: models.Message[]
  selectedThread: ThreadID | null
  highlightTs?: string
}

const cache: Map<string, string> = getContext('userhandle_cache')

let { group, selectedThread = $bindable(), highlightTs = $bindable() }: Props = $props()
let isHighlighted = $derived(group.some(msg => msg.ts === highlightTs))
let timeoutId: ReturnType<typeof setTimeout> | undefined
let current_version = $state(0)
let msg = $derived(group[current_version])
let isTransitioning = $state(false)
let message_text = $derived(replaceUserIDWithHandle(msg.text || ''))

// Markdown
const mapped = [svelte('a', Link)]
const options: Options = {
  extensions: [emoji(), component(mapped, initializeComponents)],
  sanitizer: DOMPurify.sanitize,
}

export const carta = new Carta(options)

async function getUserHandle(userId: string): Promise<string> {
  if (cache.has(userId)) {
    return cache.get(userId)!
  }

  try {
    const handle = await GetUserHandle(userId)
    cache.set(userId, handle)
    return handle
  } catch (error) {
    return userId
  }
}

async function replaceUserIDWithHandle(text: string): Promise<string> {
  let result = text

  for (const match of text.matchAll(/<@(U[A-Z0-9]+)>/g)) {
    const userId = match[1]
    const handle = await GetUserHandle(userId)
    result = result.replace(match[0], handle ? `__@${handle}__` : `__@${userId}__`)
  }

  return result
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

/**
 * Handle version selection with animation
 */
const selectVersion = (index: number) => {
  if (index === current_version || isTransitioning) return

  isTransitioning = true

  // Start fade out
  setTimeout(() => {
    current_version = index
    // End transition after content updates
    setTimeout(() => {
      isTransitioning = false
    }, 150)
  }, 150)
}

$effect(() => {
  if (selectedThread && selectedThread?.thread_ts == group[current_version].thread_ts) {
    selectedThread.import_id = group[current_version].ImportID
  }
})

$effect(() => {
  // When hl becomes true, start/reset the timer
  if (isHighlighted) {
    if (timeoutId) clearTimeout(timeoutId)
    timeoutId = setTimeout(() => {
      isHighlighted = false
      highlightTs = undefined
      timeoutId = undefined
    }, 2000)
  } else {
    // If hl becomes false externally, cancel timer and turn off highlight
    if (timeoutId) {
      clearTimeout(timeoutId)
      timeoutId = undefined
    }
    isHighlighted = false
  }

  // Cleanup function: clear timeout if component unmounts or effect re-runs
  return () => {
    if (timeoutId) {
      clearTimeout(timeoutId)
      highlightTs = undefined
      timeoutId = undefined
    }
  }
})
</script>

<div class="message" class:highlight={isHighlighted}>
    <img class="avatar" src={msg.avatar_url} alt={msg.user_name} />
    <div class="message-content">
        <div class="message-header">
            <div>
                <strong>{msg.user_name}</strong>
                <span class="timestamp">
                    {formatMessageTime(msg.ts)}
                </span>
            </div>
            {#if group.length > 1}
                <div class="version-container">
                    {#each group as version, i}
                        <button
                            class="version-button"
                            class:active={i === current_version}
                            onclick={() => selectVersion(i)}
                            title="Import ID: {version.ImportID}"
                        >
                            {@render circle("1rem", i === current_version)}
                        </button>
                    {/each}
                </div>
            {/if}
        </div>

        <div class="content-wrapper" class:transitioning={isTransitioning}>
            {#await message_text then value}
                <Markdown {carta} {value} />
            {/await}
            <Attachments message={msg} />
            {#if msg.thread_ts}
                <ThreadInfo
                    replyUsers={msg.reply_users}
                    replyCount={msg.reply_count}
                    thread={{
                        thread_ts: msg.thread_ts,
                        import_id: msg.ImportID,
                    }}
                    bind:selectedThread
                />
            {/if}
            {#if msg.reactions}
                <Reactions reactions={msg.reactions} />
            {/if}
        </div>
    </div>
</div>

<style>
    :global(.carta-viewer > :first-child) {
        margin-top: 2px;
    }

    .version-container {
        padding-right: 20px;
        display: flex;
        align-items: center;
    }

    .version-button {
        background: none;
        border: none;
        padding: 0;
        cursor: pointer;
        transition: all 0.2s ease;
        color: rgba(134, 134, 134, 0.8);
    }

    .version-button:hover {
        transform: scale(1.1);
    }

    .version-button.active {
        background-color: rgba(255, 255, 255, 0.15);
    }

    .content-wrapper {
        transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
        transform: translateX(0);
        opacity: 1;
    }

    .content-wrapper.transitioning {
        opacity: 0;
    }

    @keyframes highlight-blink {
        0% { background-color: palegoldenrod; }
        20% { background-color: #fff3b0; } /* lighter pulse */
        40% { background-color: palegoldenrod; }
        60% { background-color: #fff3b0; }
        80% { background-color: palegoldenrod; }
        100% { background-color: palegoldenrod; } /* end at highlight color */
    }
    .message.highlight {
        animation: highlight-blink 5s ease;
    }

    .message {
        margin-left: 8px;
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

    .message-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .message-header strong {
        font-weight: 700;
    }

    .timestamp {
        color: #9e9ea6;
        font-size: 0.8rem;
        margin-left: 0.5rem;
    }
</style>
