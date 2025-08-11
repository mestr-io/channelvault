<script lang="ts">
import { GetChannels, GetMessages } from '../../../wailsjs/go/app/App'
import { models } from '../../../wailsjs/go/models'
import { onMount, tick } from 'svelte'

let channels = $state<models.Channel[]>([])
let messages = $state<models.Message[]>([])
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

// Fetch initial channel list
$effect(() => {
  ;(async () => {
    try {
      channels = await GetChannels()
      if (channels.length > 0 && !selectedChannel) {
        await selectChannel(channels[0])
      }
    } catch (error) {
      console.error('Failed to get channels:', error)
    }
  })()
})

$effect(() => {
  // This runs after messages are rendered. If it was the first load, scroll down.
  if (messages.length > 0 && isInitialLoad && messagesContainer) {
    tick().then(() => {
      messagesContainer.scrollTop = messagesContainer.scrollHeight
      isInitialLoad = false
    })
  }
})

// TODO: move this to a component!!! <SideBar> ... </>
// Sidebar resizing logic
onMount(() => {
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
    const diffX = e.pageX - startX
    const newWidth = Math.max(200, startWidth + diffX)
    sidebar.style.width = newWidth + 'px'
    mainContent.style.width = `calc(100% - ${newWidth}px)`
    localStorage.setItem('sidebarWidth', newWidth.toString())
  }

  function stopDrag() {
    document.removeEventListener('mousemove', drag)
    document.removeEventListener('mouseup', stopDrag)
  }
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
    <ul class="channels-list">
      {#each channels as channel (channel.id)}
        <li
          onclick={() => selectChannel(channel)}
          class:active={channel.id === selectedChannel?.id}
        >
          # {channel.name}
        </li>
      {/each}
    </ul>
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
        {#if hasMoreMessages}
          <div class="loader">
            {#if isLoadingMessages}<span>Loading history...</span>{/if}
          </div>
        {/if}

        <ul>
          {#each messages as msg (msg.ts)}
            <li class="message">
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
                <p>{msg.reply_count} replies <button>View thread</button></p>
                {/if}
              </div>
            </li>
          {/each}
        </ul>
      </div>
    {:else}
      <div class="welcome-screen">
        <h2>Welcome!</h2>
        <p>Select a channel to view messages.</p>
      </div>
    {/if}
  </div>
</main>

<style>
  /* --- MODIFIED: Updated color variables for a Slack-like dark theme --- */
  :root {
    --slack-sidebar-bg: #19171d; /* Slack's dark purple/eggplant */
    --slack-main-bg: #1d1c1d; /* A slightly lighter dark gray for content */
    --slack-main-text: #d1d2d3; /* Light gray text for readability on dark BG */
    --slack-sidebar-text: #d1d2d3;
    --slack-sidebar-hover-bg: #272a2e;
    --slack-sidebar-active-bg: #1164a3; /* Keep the highlight blue */
    --slack-border-color: #3b3d40;
    --font-family: "Lato", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
      Helvetica, Arial, sans-serif;
  }
  .slack-clone {
    display: flex;
    height: 100vh;
    font-family: var(--font-family);
    background-color: var(--slack-main-bg);
    color: var(--slack-main-text); /* MODIFIED: Set global text color */
  }
  .sidebar {
    flex-shrink: 0;
    overflow-y: auto;
    background-color: var(--slack-sidebar-bg);
    color: var(--slack-sidebar-text);
  }
  .sidebar h4 {
    padding: 1rem 1.25rem 0.5rem;
    font-size: 0.85rem;
  }
  .channels-list {
    list-style: none;
    padding: 0;
    margin: 0;
  }
  .channels-list li {
    padding: 0.4rem 1.25rem;
    cursor: pointer;
    border-radius: 4px;
    margin: 0 0.5rem;
    white-space: nowrap;
  }
  .channels-list li:hover {
    background-color: var(--slack-sidebar-hover-bg);
  }
  .channels-list li.active {
    background-color: var(--slack-sidebar-active-bg);
    color: white;
  }
  .separator {
    width: 5px;
    cursor: ew-resize;
    background-color: var(
      --slack-border-color
    ); /* MODIFIED: Darker separator */
    border: none;
    flex-shrink: 0;
  }
  .main-content {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    height: 100vh;
  }
  .welcome-screen {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    height: 100%;
  }
  .channel-header {
    padding: 1rem 1.25rem;
    border-bottom: 1px solid var(--slack-border-color); /* MODIFIED: Darker border */
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
  .messages-container ul {
    list-style: none;
    padding: 1rem;
    margin: 0;
    margin-top: auto;
  }
  .loader {
    text-align: center;
    padding: 1rem;
    min-height: 40px;
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
  .message-content p {
    margin: 0;
    line-height: 1.45;
    white-space: pre-wrap; /* ADDED: This is the key to fixing text formatting */
  }
  .message-header strong {
    font-weight: 900;
  }
  .timestamp {
    color: #9e9ea6; /* MODIFIED: Lighter gray for timestamp */
    font-size: 0.8rem;
    margin-left: 0.5rem;
  }
</style>
