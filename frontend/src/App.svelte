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
import NoMessages from './lib/components/NoMessages.svelte'
import { models } from '../wailsjs/go/models'
import { import_icon } from './lib/components/Icons.svelte'
import { chevron } from './lib/components/icons'

import type { ThreadID } from './lib/types'

type ViewState = 'setup' | 'dashboard' | 'loading'
type AggregatedMessage = models.Message[][]

let imageOverlayComponent: ImageOverlay
let importsOverlayComponent: ImportManager

let viewState = $state<ViewState>('loading')

// Dashboard state
let channels = $state<models.Channel[]>([] as models.Channel[])
let selectedChannel = $state<models.Channel | null>(null)
let threadMessages = $state<models.Message[]>([] as models.Message[])
let messages = $state<models.Message[]>([] as models.Message[])
let groups = $state<AggregatedMessage>([] as AggregatedMessage)
let isLoadingMessages = $state(false)
let hasMoreMessages = $state(true)
let showThread = $state(false)
let selectedThread = $state<ThreadID | null>(null)
let currentView = $state<'welcome' | 'channel' | 'imports'>('welcome')
let channelFilter = $state('') // New state variable for channel filter
let showChannelFilter = $state(false) // New state variable to control filter visibility
// svelte-ignore non_reactive_update
let vlist_messages: VListHandle
// svelte-ignore non_reactive_update
let vlist_thread: VListHandle

$effect(() => {
  if (!showThread) selectedThread = null
})

$effect(() => {
  if (!selectedThread) showThread = false
})

const filteredChannels = $derived(
  channels.filter(channel => channel.name.toLowerCase().includes(channelFilter.toLowerCase())),
)

$effect(() => {
  if (selectedThread && selectedChannel) {
    const loadThread = async () => {
      threadMessages = await GetThreadByTs(selectedThread!.thread_ts, selectedChannel!.id!, selectedThread!.import_id)
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

runtime.EventsOn('backend:viewState', (newState: string) => {
  viewState = newState as ViewState
})
runtime.EventsOn('importComplete', async () => {
  channels = await GetChannels()
})
runtime.EventsOn('importDeleted', async () => {
  channels = await GetChannels()
})

onMount(async () => {
  viewState = (await GetViewState()) as ViewState
  channels = await GetChannels()
})

onDestroy(() => {
  document.body.classList.remove('is-resizing-panes')
})

setContext('avatar_cache', new Map<string, string>())
setContext('userhandle_cache', new Map<string, string>())
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

function groupMessagesByDay(messages: models.Message[]): AggregatedMessage {
  if (!messages || messages.length === 0) {
    return []
  }

  const groups = messages.reduce((acc: AggregatedMessage, msg) => {
    const lastGroup = acc.length > 0 ? acc[acc.length - 1] : null

    if (lastGroup) {
      const lastGroupDate = lastGroup[0].ts
      if (lastGroupDate === msg.ts) {
        lastGroup.push(msg)
        return acc
      }
    }

    acc.push([msg])
    return acc
  }, [])

  return groups
}

async function loadMoreMessages() {
  if (isLoadingMessages || !selectedChannel || !hasMoreMessages) return
  isLoadingMessages = true
  try {
    const cursor = messages[0]?.ts || ''
    const newMessages: models.Message[] = await GetMessages(selectedChannel.id, cursor, 100)

    if (newMessages?.length > 0) {
      messages = [...newMessages.reverse(), ...messages]
      groups = groupMessagesByDay(messages)
    } else {
      hasMoreMessages = false
    }
  } catch (error) {
    console.error(`Failed to get messages for channel ${selectedChannel.id}:`, error)
  } finally {
    isLoadingMessages = false
  }
  await tick()
}

function handleScroll(offset: number) {
  if (offset === 0 && hasMoreMessages && !isLoadingMessages) {
    loadMoreMessages().then(() => {})
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

<main class="dashboard">
  <Splitpanes on:resize={startResize} on:resized={endResize}>
    <Pane minSize={10} maxSize={30} size={25}>
      <div class="sidebar">
        <div class="sidebar-header">
          <button
            class="sidebar-header-title"
            onclick={() => (showChannelFilter = !showChannelFilter)}
          >
            <div class="header-options-chevron">
              {@render chevron(showChannelFilter)}
            </div>
            Channels
          </button>
          {#if showChannelFilter}
            <input
              type="text"
              placeholder="Filter..."
              bind:value={channelFilter}
              class="channel-filter-input"
            />
          {/if}
        </div>
        <div class="channel-list">
          {#each filteredChannels as channel (channel.id)}
            <div class="channel-item">
              <button
                class:active={channel.id === selectedChannel?.id &&
                  currentView === "channel"}
                onclick={() => selectChannel(channel)}># {channel.name}</button
              >
            </div>
          {/each}
        </div>
        <div class="channel-item settings-button">
          <button onclick={() => openImportsOverlay()}>
            {@render import_icon("1rem")}
            Imports
          </button>
        </div>
      </div></Pane
    >
    <Pane>
      <div class="main-content" id="main-content">
        {#if currentView === "channel" && selectedChannel && messages.length > 0}
          <header class="channel-header">
            <h2># {selectedChannel.name}</h2>
          </header>
          <VList
            bind:this={vlist_messages}
            data={groups}
            getKey={(g: models.Message[], _) => `${g[0].ts}`}
            overscan={10}
            shift={true}
            onscroll={handleScroll}
          >
            {#snippet children(group, i)}
              <DayDivider
                current={group[0]}
                previous={groups[i - 1] ? groups[i - 1][0] : undefined}
              />
              <Message {group} bind:selectedThread />
            {/snippet}
          </VList>
        {:else if currentView === "channel" && selectedChannel && messages.length === 0}
          <NoMessages />
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
          >
            {#snippet children(msg, i)}
              <DayDivider current={msg} previous={threadMessages[i - 1]} />
              <Message group={[msg]} bind:selectedThread />
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

  :global(body.is-resizing-panes) {
    cursor: col-resize !important;
    user-select: none !important;
    -webkit-user-select: none !important;
    -moz-user-select: none !important;
    -ms-user-select: none !important;
  }
  .dashboard {
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
    justify-content: flex-start;
  }
  .channel-list {
    flex-grow: 1;
    overflow-y: auto;
    overflow-x: hidden;
  }

  .sidebar-header {
    padding: 1rem 1.25rem 0.5rem;
    font-size: 1.2rem;
    font-weight: 600;
    color: #dbc9e0;
    border-bottom: 1px solid #7a4b87;
    display: flex;
    flex-direction: column;
    gap: 10px;
  }
  .header-options-chevron {
    width: 1.6rem;
    height: 1.6rem;
  }

  .sidebar-header-title {
    display: flex;
    align-items: center;
    gap: 5px;
    background: none;
    border: none;
    color: #dbc9e0;
    font-size: 1.2rem;
    font-weight: 600;
    cursor: pointer;
    padding: 0;
    margin: 0;
  }

  .channel-filter-input {
    background-color: #3e0e42;
    border: 1px solid #663867;
    border-radius: 4px;
    color: #dbc9e0;
    padding: 0.3rem 0.4rem;
    width: calc(100% - 1.5rem);
    font-size: 0.8rem;
  }

  .channel-filter-input::placeholder {
    color: #a08fa6;
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

  .channel-item.settings-button {
    border-top: 1px solid #7a4b87;
    padding: 0 8px;
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
