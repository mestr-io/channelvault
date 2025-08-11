<script lang="ts">
import { GetChannels, GetMessages } from '../../../wailsjs/go/app/App'
import { models } from '../../../wailsjs/go/models'
import { onMount } from 'svelte'

// ===================================================================================
// 1. BACKEND INTERFACES
// These TypeScript interfaces must match the JSON tags of your Go structs.
// ===================================================================================

// ===================================================================================
// 2. COMPONENT STATE
// State variables to manage the UI and data.
// ===================================================================================

let channels = $state<models.Channel[]>([])
let messages = $state<models.Message[]>([])
let selectedChannel = $state<models.Channel | null>(null)
let isLoadingMessages = $state(false)
let hasMoreMessages = $state(true) // Tracks if there are older messages to fetch
let newMessageText = $state('')

let startX: number
let startWidth: number
const initialDefaultWidth = 200
// Add these two lines to your component state (section 2)
let scrollTriggerEl: HTMLLIElement | undefined = $state()
let messagesContainerEl: HTMLDivElement | undefined = $state()

// PRESERVE SCROLL STATE
// Add this new state variable
let scrollHeightBeforeLoad = $state(0)

$effect(() => {
  // When messages are prepended, adjust scroll to keep the view stable
  if (scrollHeightBeforeLoad && messagesContainerEl) {
    messagesContainerEl.scrollTop +=
      messagesContainerEl.scrollHeight - scrollHeightBeforeLoad
    scrollHeightBeforeLoad = 0 // Reset after adjusting
  }
})

// Add this $effect to create and manage the Intersection Observer
$effect(() => {
  if (!scrollTriggerEl || !messagesContainerEl) return

  const observer = new IntersectionObserver(
    (entries) => {
      // If the trigger element is intersecting (visible) and we're not already loading...
      if (entries[0].isIntersecting && !isLoadingMessages) {
        // ...save the current scroll height before loading new messages
        scrollHeightBeforeLoad = messagesContainerEl.scrollHeight
        loadMoreMessages()
      }
    },
    {
      root: messagesContainerEl, // Observe within the messages container
      threshold: 1.0, // Trigger when 100% of the element is visible
    },
  )

  // Start observing the trigger element
  observer.observe(scrollTriggerEl)

  // Return a cleanup function to disconnect the observer when the component unmounts
  // or the trigger element is removed from the DOM.
  return () => {
    observer.disconnect()
  }
})

// Fetch initial channel list on component mount
$effect(() => {
  ;(async () => {
    try {
      channels = await GetChannels()
      // Automatically select the first channel if available
      if (channels.length > 0 && !selectedChannel) {
        await selectChannel(channels[0])
      }
    } catch (error) {
      console.error('Failed to get channels:', error)
    }
  })()
})
onMount(() => {
  const separator = document.getElementById('separator')!
  const sidebar = document.getElementById('sidebar')!
  const mainContent = document.getElementById('main-content')!

  // Retrieve the saved width from localStorage or set to default
  const savedWidth = localStorage.getItem('sidebarWidth')
  const initialWidth = savedWidth ? parseInt(savedWidth, 10) : initialDefaultWidth

  sidebar.style.width = `${initialWidth}px`
  mainContent.style.width = `calc(100% - ${initialWidth}px)`

  separator.addEventListener('mousedown', startDrag)

  function startDrag(e) {
    startX = e.pageX
    startWidth = sidebar.offsetWidth

    document.addEventListener('mousemove', drag)
    document.addEventListener('mouseup', stopDrag)
  }

  function drag(e) {
    e.preventDefault()
    const diffX = e.pageX - startX
    const newWidth = Math.max(200, startWidth + diffX) // Minimum width of 200px

    sidebar.style.width = newWidth + 'px'
    mainContent.style.width = `calc(100% - ${newWidth}px)`

    // Save the new width to localStorage, so user will experience always last width
    localStorage.setItem('sidebarWidth', newWidth.toString())
  }

  function stopDrag() {
    document.removeEventListener('mousemove', drag)
    document.removeEventListener('mouseup', stopDrag)
  }
})

// Function to switch channels and load its messages
async function selectChannel(channel: Channel) {
  if (selectedChannel?.id === channel.id) return // Avoid reloading same channel

  selectedChannel = channel
  messages = [] // Clear previous messages
  hasMoreMessages = true // Reset pagination state for the new channel
  await loadMoreMessages()
}

// Function to load a page of messages using the oldest message's timestamp as a cursor
async function loadMoreMessages() {
  if (isLoadingMessages || !selectedChannel || !hasMoreMessages) return

  isLoadingMessages = true
  try {
    const lastTimestamp = messages.at(-1)?.ts || ""

    const newMessages: models.Message[] = await GetMessages(selectedChannel.id, lastTimestamp)

    if (newMessages && newMessages.length > 0) {
      messages = [...newMessages, ...messages]
    } else {
      hasMoreMessages = false
    }
  } catch (error) {
    console.error(`Failed to get messages for channel ${selectedChannel.id}:`, error)
  } finally {
    isLoadingMessages = false
  }
}
</script>

<!-- =================================================================================== -->
<!-- 4. HTML STRUCTURE & STYLING                                                         -->
<!-- =================================================================================== -->

<main class="slack-clone">
  <!-- Sidebar -->
  <div class="sidebar" id="sidebar">
    <h4>Channels</h4>
    <ul class="channels-list">
      {#each channels as channel (channel.id)}
        <li
          on:click={() => selectChannel(channel)}
          class:active={channel.id === selectedChannel?.id}
        >
          # {channel.name}
        </li>
      {/each}
    </ul>
  </div>
  <div class="separator" id="separator"></div>
  <!-- Main Content Area -->
 <div class="main-content" id="main-content">
  {#if selectedChannel}
    <header class="channel-header">
      <h2># {selectedChannel.name}</h2>
    </header>

    <div class="messages-container" bind:this={messagesContainerEl}>
      <ul>
        {#if hasMoreMessages && !isLoadingMessages}
          <li class="scroll-trigger" bind:this={scrollTriggerEl}></li>
        {/if}

        {#if isLoadingMessages}
          <li class="loader-li">
            <div class="loader">Loading...</div>
          </li>
        {/if}

        {#each messages as msg (msg.ts)}
          <li class="message">
            <img class="avatar" src={msg.avatar_url} alt={msg.user_name} />
            <div class="message-content">
              <div class="message-header">
                <strong>{msg.user_name}</strong>
                <span class="timestamp">
                  {new Date(parseFloat(msg.ts) * 1000).toLocaleTimeString(
                    [],
                    {
                      hour: '2-digit',
                      minute: '2-digit',
                      day: '2-digit',
                      month: '2-digit',
                      year: 'numeric',
                    },
                  )}
                </span>
              </div>
              <p>{msg.text}</p>
            </div>
          </li>
        {/each}
      </ul>
    </div>
  {:else}
    <div class="welcome-screen">
      <h2>Welcome to your Slack Clone!</h2>
      <p>Select a channel from the left sidebar to start chatting.</p>
    </div>
  {/if}
</div>
</main>

<style>
  /* CSS Variables for easy theming */
  :root {
    --slack-sidebar-bg: #1a1d21;
    --slack-sidebar-text: #d1d2d3;
    --slack-sidebar-hover-bg: #272a2e;
    --slack-sidebar-active-bg: #1164a3;
    --slack-sidebar-active-text: #ffffff;
    --slack-border-color: #3b3d40;
    --slack-main-bg: #f8f8f8;
    --slack-main-text: #1d1c1d;
    --slack-input-bg: #ffffff;
    --font-family: "Lato", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
      Helvetica, Arial, sans-serif;
  }

  .slack-clone {
    display: flex;
    height: 100vh;
    /* width: 100vw; */
    font-family: var(--font-family);
    background-color: var(--slack-main-bg);
    color: var(--slack-main-text);
  }

  .separator {
    width: 20px; /* Just for debugging process*/
    cursor: ew-resize;
    background-color: #eee;
    border-left: 1px solid #ccc;
    border-right: 1px solid #ccc;
  }


.scroll-trigger {
  /* This makes the trigger element take up space but be invisible */
  height: 1px;
  visibility: hidden;
}

.loader-li {
  /* This ensures the loader is centered without list-item styling */
  list-style: none;
  padding: 1rem 0;
  text-align: center;
}
  /* --- Sidebar --- */
  .sidebar {
    overflow: auto;
    flex-grow: 0;
    background-color: var(--slack-sidebar-bg);
    color: var(--slack-sidebar-text);
  }

  .workspace-header {
    padding: 1rem 1.25rem;
    border-bottom: 1px solid var(--slack-border-color);
  }

  .workspace-header h3 {
    margin: 0;
    font-size: 1.1rem;
    color: white;
  }

  /* .channels-list-container {
    padding: 1rem 0;
  } */

  .sidebar h4 {
    padding: 0 1.25rem 0.5rem;
    margin: 0;
    font-size: 0.85rem;
    text-transform: uppercase;
    font-weight: 700;
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
  }

  .channels-list li:hover {
    background-color: var(--slack-sidebar-hover-bg);
  }

  .channels-list li.active {
    background-color: var(--slack-sidebar-active-bg);
    color: var(--slack-sidebar-active-text);
    font-weight: bold;
  }

  /* --- Main Content --- */
  .main-content {
    overflow: auto;
    flex-grow: 1;
    height: 100vh;
  }

  .welcome-screen {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    height: 100%;
    color: #6c757d;
  }

  .channel-header {
    padding: 1rem 1.25rem;
    border-bottom: 1px solid #e0e0e0;
    flex-shrink: 0;
  }

  .channel-header h2 {
    margin: 0 0 0.25rem 0;
    font-size: 1.25rem;
  }
  .channel-header p {
    margin: 0;
    font-size: 0.9rem;
    color: #6c757d;
  }

  /* --- Messages --- */
  .messages-container {
    flex-grow: 1;
    overflow-y: auto;
    padding: 1rem;
    display: flex;
    flex-direction: column-reverse; /* Key for chat UI */
  }

  .messages-container ul {
    list-style: none;
    padding: 0;
  }

  .loader,
  .load-more-btn {
    text-align: center;
    padding: 0.75rem;
    margin: 0.5rem auto;
  }

  .load-more-btn {
    background: #007a5a;
    color: white;
    border: none;
    border-radius: 20px;
    cursor: pointer;
    font-size: 0.9rem;
  }

  .load-more-btn:hover {
    background: #148567;
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
    flex-shrink: 0;
    object-fit: cover;
  }

  .message-content p {
    margin: 0;
    line-height: 1.4;
  }

  .message-header strong {
    font-weight: 900;
  }

  .timestamp {
    color: #616061;
    font-size: 0.8rem;
    margin-left: 0.5rem;
  }

  /* --- Message Input --- */
  .message-input-area {
    padding: 1rem;
    flex-shrink: 0;
  }

  .message-input-area form {
    display: flex;
    align-items: center;
    border: 1px solid #bdbdbd;
    border-radius: 8px;
    background-color: var(--slack-input-bg);
    padding: 0.25rem;
  }

  .message-input-area input {
    flex-grow: 1;
    border: none;
    outline: none;
    padding: 0.75rem;
    font-size: 1rem;
    background: transparent;
  }

  .message-input-area button {
    background-color: #007a5a;
    border: none;
    border-radius: 6px;
    color: white;
    padding: 0.5rem;
    margin-right: 0.25rem;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .message-input-area button:disabled {
    background-color: #e0e0e0;
    cursor: not-allowed;
  }

  .message-input-area button:not(:disabled):hover {
    background-color: #148567;
  }
</style>
