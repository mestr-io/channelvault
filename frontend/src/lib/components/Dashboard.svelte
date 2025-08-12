<script lang="ts">
import { GetChannels, GetMessages } from '../../../wailsjs/go/app/App'
import { models } from '../../../wailsjs/go/models'
import { onMount, tick } from 'svelte'
import Messages from './Messages.svelte'
import Welcome from './Welcome.svelte'

let channels = $state<models.Channel[]>([] as models.Channel[])
let messages = $state<models.Message[]>([] as models.Message[])
let selectedChannel = $state<models.Channel | null>(null)
let isLoadingMessages = $state(false)
let hasMoreMessages = $state(true)

// State for sidebar resizing
let startX: number
let startWidth: number
const initialDefaultWidth = 250

// State and refs for scroll management
// svelte-ignore non_reactive_update
let messagesContainer: HTMLDivElement
let isInitialLoad = $state(true)

$effect(() => {
  // This runs after messages are rendered. If it was the first load, scroll down.
  if (messages.length > 0 && isInitialLoad && messagesContainer) {
    tick().then(() => {
      messagesContainer.scrollTop = messagesContainer.scrollHeight
      isInitialLoad = false
    })
  }
})

// Sidebar resizing logic
onMount(async () => {
  const separator = document.getElementById('separator')!
  const sidebar = document.getElementById('sidebar')!
  const mainContent = document.getElementById('main-content')!
  const savedWidth = localStorage.getItem('sidebarWidth')
  const initialWidth = savedWidth ? parseInt(savedWidth, 10) : initialDefaultWidth

  sidebar.style.width = `${initialWidth}px`
  mainContent.style.width = `calc(100% - ${initialWidth}px)`

  separator.addEventListener('mousedown', startDrag)

  function startDrag(e: MouseEvent) {
    startX = e.pageX
    startWidth = sidebar.offsetWidth
    document.addEventListener('mousemove', drag)
    document.addEventListener('mouseup', stopDrag)
  }

  function drag(e: MouseEvent) {
    e.preventDefault()
    separator.style.backgroundColor = '#1d9bd1' // Change color on drag
    const diffX = e.pageX - startX
    const newWidth = Math.max(200, startWidth + diffX)
    sidebar.style.width = newWidth + 'px'
    mainContent.style.width = `calc(100% - ${newWidth}px)`
    // TODO: store that in BE
    localStorage.setItem('sidebarWidth', newWidth.toString())
  }

  function stopDrag() {
    // restore color
    separator.style.backgroundColor = '#5b6065'
    document.removeEventListener('mousemove', drag)
    document.removeEventListener('mouseup', stopDrag)
  }
  channels = await GetChannels()
  // commented, just to test if we can add a splash screen
  // if (channels.length > 0 && !selectedChannel) {
  //   await selectChannel(channels[0]);
  // }
})

async function selectChannel(channel: models.Channel) {
  if (selectedChannel?.id === channel.id) return

  selectedChannel = channel
  messages = []
  hasMoreMessages = true
  isInitialLoad = true
  await loadMoreMessages()
}

async function loadMoreMessages() {
  if (isLoadingMessages || !selectedChannel || !hasMoreMessages) return

  isLoadingMessages = true
  try {
    const cursor = messages[0]?.ts || ''
    const oldScrollHeight = messagesContainer?.scrollHeight || 0

    // The backend returns messages newest-to-oldest
    const newMessages: models.Message[] = await GetMessages(selectedChannel.id, cursor)

    if (newMessages?.length > 0) {
      // Prepend older messages. Reverse the fetched chunk so it's in chronological order.
      messages = [...newMessages.reverse(), ...messages]

      // If this was NOT the initial load, preserve scroll position
      if (!isInitialLoad && messagesContainer) {
        await tick() // Wait for DOM to update
        messagesContainer.scrollTop = messagesContainer.scrollHeight - oldScrollHeight
      }
    } else {
      hasMoreMessages = false
    }
  } catch (error) {
    console.error(`Failed to get messages for channel ${selectedChannel.id}:`, error)
  } finally {
    isLoadingMessages = false
  }
}

function handleScroll() {
  if (messagesContainer.scrollTop === 0 && hasMoreMessages && !isLoadingMessages) {
    loadMoreMessages()
  }
}
</script>

<main class="slack-clone">
  <div class="sidebar" id="sidebar">
    <h4>Channels</h4>
    {#each channels as channel (channel.id)}
      <div class="channel-item">
        <button
          class:active={channel.id === selectedChannel?.id}
          onclick={() => selectChannel(channel)}># {channel.name}</button
        >
      </div>
    {/each}
  </div>

  <div class="separator" id="separator"></div>

  <div class="main-content" id="main-content">
    {#if selectedChannel}
      <header class="channel-header">
        <h2># {selectedChannel.name}</h2>
      </header>
      <div
        class="messages-container"
        bind:this={messagesContainer}
        onscroll={handleScroll}
      >
        <Messages {messages} {isLoadingMessages} {hasMoreMessages} />
      </div>
    {:else}
      <Welcome />
    {/if}
  </div>
</main>

<style>
  .slack-clone {
    display: flex;
    height: 100vh;
  }
  .sidebar {
    flex-shrink: 0;
    overflow-y: auto;
    overflow-x: hidden;
    background-color: #552958;
  }

  .sidebar h4 {
    padding: 1rem 1.25rem 0.5rem;
    font-size: 1.2rem;
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
  .separator {
    width: 5px;
    cursor: ew-resize;
    background-color: #5b6065;
    border: none;
    flex-shrink: 0;
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
</style>
