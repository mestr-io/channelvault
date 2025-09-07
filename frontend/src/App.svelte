<script lang="ts">
import * as runtime from '../wailsjs/runtime'
import { onMount, setContext, onDestroy, tick } from 'svelte'
import { GetViewState, GetChannels, GetMessages, GetThreadByTs } from '../wailsjs/go/app/App'
import { Pane, Splitpanes } from 'svelte-splitpanes'
import { VList, type VListHandle } from 'virtua/svelte'
import DayDivider from './lib/components/DayDivider.svelte'
import Message from './lib/components/Message.svelte'
import Welcome from './lib/components/Welcome.svelte'
import ImportManager from './lib/components/ImportManager.svelte'
import ImageOverlay from './lib/components/ImageOverlay.svelte'
import { models } from '../wailsjs/go/models'

type ViewState = 'setup' | 'dashboard' | 'loading'

let imageOverlayComponent: ImageOverlay
let importsOverlayComponent: ImportManager

let viewState = $state<ViewState>('loading')

// Dashboard state
let channels = $state<models.Channel[]>([] as models.Channel[])
let selectedChannel = $state<models.Channel | null>(null)
let threadMessages = $state<models.Message[]>([] as models.Message[])
let messages = $state<models.Message[]>([] as models.Message[])
let isLoadingMessages = $state(false)
let hasMoreMessages = $state(true)
let showThread = $state(false)
let selectedThread = $state<string | null>(null)
let currentView = $state<'welcome' | 'channel' | 'imports'>('welcome')
let vlist_messages: VListHandle
let vlist_thread: VListHandle

$effect(() => {
  if (!showThread) selectedThread = null
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

onMount(async () => {
  try {
    viewState = (await GetViewState()) as ViewState
    runtime.EventsOn('backend:viewState', (newState: string) => {
      viewState = newState as ViewState
    })

    channels = await GetChannels()
    runtime.EventsOn('importComplete', async () => {
      const fresh = await GetChannels()
      channels = [...fresh]
    })
    runtime.EventsOn('importDeleted', async () => {
      const fresh = await GetChannels()
      channels = [...fresh]
    })
  } catch (error) {
    console.error('Failed to initialize view state or set up listener:', error)
  }
})

onDestroy(() => {
  document.body.classList.remove('is-resizing-panes')
})

setContext('cache', new Map<string, string>())
setContext('imageOverlay', { openImageOverlay })

function openImageOverlay(src: string) {
  if (imageOverlayComponent) {
    imageOverlayComponent.showOverlay(src)
  }
}

function openImportsOverlay() {
  if (importsOverlayComponent) {
    importsOverlayComponent.showOverlay()
  }
}

async function selectChannel(channel: models.Channel) {
  if (selectedChannel?.id === channel.id && currentView === 'channel') return
  selectedChannel = channel
  messages = []
  hasMoreMessages = true
  showThread = false
  currentView = 'channel'
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
    const newMessages: models.Message[] = await GetMessages(selectedChannel.id, cursor, 100)
    if (newMessages?.length > 0) {
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

function startResize() {
  document.body.classList.add('is-resizing-panes')
}
function endResize() {
  document.body.classList.remove('is-resizing-panes')
}
</script>

<ImageOverlay bind:this={imageOverlayComponent} />
<ImportManager bind:this={importsOverlayComponent} />

<main class="slack-clone">
  <Splitpanes on:resize={startResize} on:resized={endResize}>
    <Pane minSize={10} maxSize={30} size={25}>
      <div class="sidebar">
        <h4>Channels</h4>
        {#each channels as channel (channel.id)}
          <div class="channel-item">
            <button
              class:active={channel.id === selectedChannel?.id &&
                currentView === "channel"}
              onclick={() => selectChannel(channel)}># {channel.name}</button
            >
          </div>
        {/each}
        <div class="imports" onclick={() => openImportsOverlay()}>
          <svg
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
            ><path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
            ></path></svg
          >
          Imports
        </div>
      </div>
    </Pane>
    <Pane>
      <div class="main-content" id="main-content">
        {#if currentView === "channel" && selectedChannel}
          <header class="channel-header">
            <h2># {selectedChannel.name}</h2>
          </header>
          <VList
            bind:this={vlist_messages}
            data={messages}
            getKey={(msg: models.Message, _) => `${msg.ts}${msg.ImportID}`}
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
            getKey={(msg: models.Message, _) => `${msg.ts}${msg.ImportID}`}
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
  :global(body) {
    font-family:
      "Lato",
      -apple-system,
      BlinkMacSystemFont,
      "Segoe UI",
      Roboto,
      Helvetica,
      Arial,
      sans-serif;
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }
  .app-main {
    background-color: transparent;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    overflow: hidden;
  }
  :global(body.is-resizing-panes) {
    cursor: col-resize !important;
    user-select: none !important;
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
    display: flex;
    flex-direction: column;
  }

  .sidebar .imports {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: #dbc9e0;
    padding: 0.5rem 1.25rem;
    margin: auto 0.5rem 0.5rem;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.1s ease-in-out;
  }
  .sidebar .imports:hover {
    background-color: #663867;
  }
  .sidebar .imports svg {
    height: 1rem;
  }

  .sidebar .imports:active {
    background-color: #f9edff;
    color: #39063a;
  }

  .sidebar h4 {
    padding: 1rem 1.25rem 0.5rem;
    font-size: 1.2rem;
    color: #dbc9e0;
  }

  .channel-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 1.25rem;
    border-bottom: 1px solid #3b3d40;
    flex-shrink: 0;
  }

  .close-thread-btn {
    background: none;
    border: none;
    cursor: pointer;
    padding: 4px;
    border-radius: 6px;
    color: #dbc9e0;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.1s ease-in-out;
  }

  .close-thread-btn:hover {
    background-color: #663867;
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

  .channel-header h2 {
    margin: 0;
  }
</style>
