<script lang="ts">
import { GetChannels, GetMessages, GetThreadByTs } from '../../../wailsjs/go/app/App'
import { Pane, Splitpanes } from 'svelte-splitpanes'
import { onDestroy } from 'svelte'
import { VList, type VListHandle } from 'virtua/svelte'
import DayDivider from './DayDivider.svelte'
import Message from './Message.svelte'
import { models } from '../../../wailsjs/go/models'
import { onMount, tick } from 'svelte'
import Welcome from './Welcome.svelte'
import { setContext } from 'svelte'
import { Carta } from 'carta-md'
import { emoji } from '@cartamd/plugin-emoji'

let channels = $state<models.Channel[]>([] as models.Channel[])
let messages = $state<models.Message[]>([] as models.Message[])
let threadMessages = $state<models.Message[]>([] as models.Message[])
let selectedChannel = $state<models.Channel | null>(null)
let isLoadingMessages = $state(false)
let hasMoreMessages = $state(true)

let showThread = $state(false)
let selectedThread = $state<string | null>(null)
let isInitialLoad = $state(true)

// svelte-ignore non_reactive_update
let vlist_messages: VListHandle
// svelte-ignore non_reactive_update
let vlist_thread: VListHandle

const startResize = () => {
  document.body.classList.add('is-resizing-panes')
}

const endResize = () => {
  document.body.classList.remove('is-resizing-panes')
}

onDestroy(() => {
  document.body.classList.remove('is-resizing-panes')
})

const carta = new Carta({
  extensions: [emoji()],
})

setContext('carta', carta)

// init the cache
setContext('cache', new Map<string, string>())

$effect(() => {
  if (!showThread) {
    selectedThread = null
  }
})

$effect(() => {
  if (selectedThread && selectedChannel) {
    const loadThread = async () => {
      threadMessages = await GetThreadByTs(selectedThread!, selectedChannel!.id!)
      if (vlist_thread && threadMessages.length > 0) {
        await tick()
        vlist_thread.scrollToIndex(threadMessages.length - 1, {
          align: 'end',
        })
      }
    }
    loadThread()
    showThread = true
  }
})

// Sidebar resizing logic
onMount(async () => {
  channels = await GetChannels()
})

async function selectChannel(channel: models.Channel) {
  if (selectedChannel?.id === channel.id) return

  selectedChannel = channel
  messages = []
  hasMoreMessages = true
  isInitialLoad = true
  showThread = false
  await loadMoreMessages()
  if (vlist_messages && messages.length > 0) {
    await tick()
    vlist_messages.scrollToIndex(messages.length - 1, { align: 'end' })
  }
}

async function loadMoreMessages() {
  if (isLoadingMessages || !selectedChannel || !hasMoreMessages) return

  isLoadingMessages = true
  try {
    const cursor = messages[0]?.ts || ''

    // The backend returns messages newest-to-oldest
    const newMessages: models.Message[] = await GetMessages(selectedChannel.id, cursor, 200)

    if (newMessages?.length > 0) {
      // Prepend older messages. Reverse the fetched chunk so it's in chronological order.
      messages = [...newMessages.reverse(), ...messages]
    } else {
      hasMoreMessages = false
    }
  } catch (error) {
    console.error(`Failed to get messages for channel ${selectedChannel.id}:`, error)
  } finally {
    isLoadingMessages = false
  }
}

function handleScroll(offset: number) {
  if (offset === 0 && hasMoreMessages && !isLoadingMessages) {
    loadMoreMessages()
  }
}
</script>

<main class="slack-clone">
  <Splitpanes on:resize={startResize} on:resized={endResize}>
    <Pane minSize={10} maxSize={30} size={25}>
      <div class="sidebar">
        <h4 class="title">Channels</h4>
        {#each channels as channel (channel.id)}
          <div class="channel-item">
            <button
              class:active={channel.id === selectedChannel?.id}
              onclick={() => selectChannel(channel)}># {channel.name}</button
            >
          </div>
        {/each}
      </div>
    </Pane>
    <Pane>
      <div class="main-content" id="main-content">
        {#if selectedChannel}
          <header class="channel-header">
            <h2># {selectedChannel.name}</h2>
          </header>
          <VList
            bind:this={vlist_messages}
            data={messages}
            getKey={(msg: models.Message, _) => msg.ts}
            overscan={10}
            shift={true}
            onscroll={handleScroll}
          >
            {#snippet children(msg, i)}
              <DayDivider current={msg} previous={messages[i - 1]} />
              <Message {msg} bind:selectedThread />
            {/snippet}
          </VList>
        {:else}
          <Welcome />
        {/if}
      </div>
    </Pane>

    {#if showThread}
      <Pane size={30}>
        <div class="main-content">
          <header class="channel-header">
            <h2>Thread</h2>
            <button
              class="close-thread-btn"
              aria-label="Close thread"
              onclick={() => (showThread = false)}
            >
              <svg
                viewBox="0 0 20 20"
                width="20"
                height="20"
                stroke="currentColor"
                stroke-width="1.5"
                stroke-linecap="round"
                stroke-linejoin="round"
                fill="none"
                shape-rendering="geometricPrecision"
              >
                <path d="M15 5L5 15"></path>
                <path d="M5 5L15 15"></path>
              </svg>
            </button>
          </header>
          <VList
            bind:this={vlist_thread}
            data={threadMessages}
            getKey={(msg: models.Message, _) => msg.ts}
            overscan={10}
            shift={true}
            onscroll={handleScroll}
          >
            {#snippet children(msg, i)}
              <DayDivider current={msg} previous={threadMessages[i - 1]} />
              <Message {msg} bind:selectedThread />
            {/snippet}
          </VList>
        </div>
      </Pane>
    {/if}
  </Splitpanes>
</main>

<style>
  :global(body.is-resizing-panes) {
    cursor: col-resize !important; /* Force the cursor */
    user-select: none !important; /* Force disable selection */
    -webkit-user-select: none !important;
    -moz-user-select: none !important;
    -ms-user-select: none !important;
  }
  .slack-clone {
    display: flex;
    height: 100vh;
  }
  .sidebar {
    flex-shrink: 0;
    overflow-y: auto;
    height: 100%;
    overflow-x: hidden;
    background-color: #552958;
  }

  .sidebar h4 {
    padding: 1rem 1.25rem 0.5rem;
    font-size: 1.2rem;
  }

  h4.title {
    color: #dbc9e0;
  }
  .channel-header {
    /* Add these properties */
    display: flex;
    justify-content: space-between;
    align-items: center;

    /* Keep existing properties */
    padding: 1rem 1.25rem;
    border-bottom: 1px solid #3b3d40;
    flex-shrink: 0;
  }
  /* Add this new style block */
  .close-thread-btn {
    background: none;
    border: none;
    cursor: pointer;
    padding: 4px;
    border-radius: 6px;
    color: #dbc9e0; /* Match sidebar icon/text color */
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.1s ease-in-out;
  }

  .close-thread-btn:hover {
    background-color: #663867; /* Use existing hover color */
  }
  .channel-item button {
    background: none;
    border-radius: 4px;
    border: none;
    color: #dbc9e0;
    cursor: pointer;
    margin: 0 0.5rem;
    padding: 0.5rem 1.25rem;
    text-align: left;
    white-space: nowrap;
    width: 100%;
    font-size: 1rem;
    font-weight: 500;
  }

  .channel-item button:hover {
    background-color: #663867;
  }

  .channel-item button.active {
    background-color: #f9edff;
    color: #39063a;
  }

  .main-content {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    height: 100vh;
  }

  .channel-header {
    padding: 1rem 1.25rem;
    border-bottom: 1px solid #3b3d40;
    flex-shrink: 0;
  }
  .channel-header h2 {
    margin: 0;
  }
  .messages-container {
    flex-grow: 1;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
  }

  .thread-container {
    flex-grow: 1; /* Key property */
    overflow-y: auto; /* Key property */
    display: flex;
    flex-direction: column;
  }
</style>
