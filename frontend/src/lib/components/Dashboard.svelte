<script lang="ts">
  import { GetChannels, GetMessages } from "../../../wailsjs/go/app/App";
  import {models} from "../../../wailsjs/go/models";
  // ===================================================================================
  // 1. BACKEND INTERFACES
  // These TypeScript interfaces must match the JSON tags of your Go structs.
  // ===================================================================================


  // ===================================================================================
  // 2. COMPONENT STATE
  // State variables to manage the UI and data.
  // ===================================================================================

  let channels = $state<models.Channel[]>([]);
  let messages = $state<models.Message[]>([]);
  let selectedChannel = $state<models.Channel | null>(null);
  let isLoadingMessages = $state(false);
  let hasMoreMessages = $state(true); // Tracks if there are older messages to fetch
  let newMessageText = $state("");

  // ===================================================================================
  // 3. DATA FETCHING & ACTIONS
  // Functions to interact with the Go backend.
  // ===================================================================================

  // Fetch initial channel list on component mount
  $effect(() => {
    (async () => {
      try {
        channels = await GetChannels();
        // Automatically select the first channel if available
        if (channels.length > 0 && !selectedChannel) {
          await selectChannel(channels[0]);
        }
      } catch (error) {
        console.error("Failed to get channels:", error);
      }
    })();
  });

  // Function to switch channels and load its messages
  async function selectChannel(channel: Channel) {
    if (selectedChannel?.id === channel.id) return; // Avoid reloading same channel

    selectedChannel = channel;
    messages = []; // Clear previous messages
    hasMoreMessages = true; // Reset pagination state for the new channel
    await loadMoreMessages();
  }

  // Function to load a page of messages using the oldest message's timestamp as a cursor
  async function loadMoreMessages() {
    if (isLoadingMessages || !selectedChannel || !hasMoreMessages) return;

    isLoadingMessages = true;
    try {
      // The cursor is now the timestamp of the oldest message we have.
      // The list is sorted with newest at the bottom, so the oldest is at index 0.
      const lastTimestamp = messages[0]?.timestamp || null;
      
      // Backend now returns Message[] directly, using the timestamp as the cursor.
      const newMessages: Message[] = await GetMessages(selectedChannel.id, lastTimestamp);
      
      if (newMessages && newMessages.length > 0) {
          // Prepend older messages to the top of the list
          messages = [...newMessages, ...messages];
      } else {
          // No more messages were returned, so we've reached the end.
          hasMoreMessages = false;
      }

    } catch (error) {
        console.error(`Failed to get messages for channel ${selectedChannel.id}:`, error);
    } finally {
        isLoadingMessages = false;
    }
  }
  
  // Function to handle sending a new message
  async function handleSendMessage() {
    if (!newMessageText.trim() || !selectedChannel) return;
    
    try {
      await SendMessage(selectedChannel.id, newMessageText);
      newMessageText = ""; // Clear the input box
      
      // Refresh messages to show the new one. A more advanced implementation
      // might use WebSockets or optimistic UI updates.
      messages = []; 
      hasMoreMessages = true; // Reset pagination
      await loadMoreMessages();

    } catch (error) {
      console.error("Failed to send message:", error);
    }
  }
</script>

<!-- =================================================================================== -->
<!-- 4. HTML STRUCTURE & STYLING                                                         -->
<!-- =================================================================================== -->

<main class="slack-clone">
  <!-- Sidebar -->
  <div class="sidebar">
    
    <div class="channels-list-container">
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
  </div>

  <!-- Main Content Area -->
  <div class="main-content">
    {#if selectedChannel}
      <!-- Channel Header -->
      <header class="channel-header">
        <h2># {selectedChannel.name}</h2>
        <p>This is where the magic happens. Send a message to get started!</p>
      </header>

      <!-- Messages Display -->
      <div class="messages-container">  
        <!-- Messages List -->
        <ul>
          {#each messages as msg (msg.ts)}
            <li class="message">
              <img class="avatar" src={msg.authorAvatarUrl} alt={msg.authorName} />
              <div class="message-content">
                <div class="message-header">
                  <strong>{msg.authorName}</strong>
                  <span class="timestamp">
                    {new Date(msg.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                  </span>
                </div>
                <p>{msg.text}</p>
              </div>
            </li>
          {/each}
        </ul>
      </div>
    {:else}
      <!-- Welcome Screen -->
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
    --slack-sidebar-bg: #1A1D21;
    --slack-sidebar-text: #D1D2D3;
    --slack-sidebar-hover-bg: #272A2E;
    --slack-sidebar-active-bg: #1164A3;
    --slack-sidebar-active-text: #FFFFFF;
    --slack-border-color: #3B3D40;
    --slack-main-bg: #F8F8F8;
    --slack-main-text: #1D1C1D;
    --slack-input-bg: #FFFFFF;
    --font-family: 'Lato', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  }

  .slack-clone {
    display: flex;
    height: 100vh;
    width: 100vw;
    font-family: var(--font-family);
    background-color: var(--slack-main-bg);
    color: var(--slack-main-text);
  }
  
  /* --- Sidebar --- */
  .sidebar {
    display: flex;
    flex-direction: column;
    background-color: var(--slack-sidebar-bg);
    color: var(--slack-sidebar-text);
    min-width: 240px;
    flex-shrink: 0;
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

  .channels-list-container {
    padding: 1rem 0;
  }

  .channels-list-container h4 {
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

  .loader, .load-more-btn {
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
    border: 1px solid #BDBDBD;
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
