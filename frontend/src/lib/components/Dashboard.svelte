<script lang="ts">
  import {
    GetChannels,
    GetMessages,
    GetThreadByTs,
  } from "../../../wailsjs/go/app/App";
  import { Pane, Splitpanes } from "svelte-splitpanes";

  import { models } from "../../../wailsjs/go/models";
  import { onMount, tick } from "svelte";
  import Messages from "./Messages.svelte";
  import Welcome from "./Welcome.svelte";

  let channels = $state<models.Channel[]>([] as models.Channel[]);
  let messages = $state<models.Message[]>([] as models.Message[]);
  let threadMessages = $state<models.Message[]>([] as models.Message[]);
  let selectedChannel = $state<models.Channel | null>(null);
  let isLoadingMessages = $state(false);
  let hasMoreMessages = $state(true);
  let selectedThread = $state<string | undefined>(undefined);

  // State and refs for scroll management
  // svelte-ignore non_reactive_update
  let messagesContainer: HTMLDivElement;
  let threadContainer: HTMLDivElement;

  let isInitialLoad = $state(true);

  $effect(() => {
    // This runs after messages are rendered. If it was the first load, scroll down.
    if (messages.length > 0 && isInitialLoad && messagesContainer) {
      tick().then(() => {
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
        isInitialLoad = false;
      });
    }
  });

  $effect(() => {
    console.log("Effect is running!", {
      thread: selectedThread,
      channel: selectedChannel,
    });
    if (selectedThread && selectedChannel) {
      const loadThread = async () => {
        threadMessages = await GetThreadByTs(
          selectedThread,
          selectedChannel.id,
        );
        await tick();
        threadContainer.scrollTop = threadContainer.scrollHeight;
      };
      loadThread();
    }
  });

  // Sidebar resizing logic
  onMount(async () => {
    channels = await GetChannels();
  });

  async function selectChannel(channel: models.Channel) {
    if (selectedChannel?.id === channel.id) return;

    selectedChannel = channel;
    messages = [];
    hasMoreMessages = true;
    isInitialLoad = true;
    selectedThread = undefined;
    await loadMoreMessages();
  }

  async function loadMoreMessages() {
    if (isLoadingMessages || !selectedChannel || !hasMoreMessages) return;

    isLoadingMessages = true;
    try {
      const cursor = messages[0]?.ts || "";
      const oldScrollHeight = messagesContainer?.scrollHeight || 0;

      // The backend returns messages newest-to-oldest
      const newMessages: models.Message[] = await GetMessages(
        selectedChannel.id,
        cursor,
      );

      if (newMessages?.length > 0) {
        // Prepend older messages. Reverse the fetched chunk so it's in chronological order.
        messages = [...newMessages.reverse(), ...messages];

        // If this was NOT the initial load, preserve scroll position
        if (!isInitialLoad && messagesContainer) {
          await tick(); // Wait for DOM to update
          messagesContainer.scrollTop =
            messagesContainer.scrollHeight - oldScrollHeight;
        }
      } else {
        hasMoreMessages = false;
      }
    } catch (error) {
      console.error(
        `Failed to get messages for channel ${selectedChannel.id}:`,
        error,
      );
    } finally {
      isLoadingMessages = false;
    }
  }

  function handleScroll(event: Event) {
    const element = event.currentTarget as HTMLDivElement;

    if (element.scrollTop === 0) {
      // Figure out WHICH element is scrolling and call the correct loader
      if (
        element.id === "messages-container" &&
        hasMoreMessages &&
        !isLoadingMessages
      ) {
        loadMoreMessages();
      }
    }
  }
</script>

<main class="slack-clone">
  <Splitpanes>
    <Pane minSize={10} maxSize={30}>
      <div class="sidebar">
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
    </Pane>
    <Pane>
      <div class="main-content" id="main-content">
        {#if selectedChannel}
          <header class="channel-header">
            <h2># {selectedChannel.name}</h2>
          </header>
          <div
            class="messages-container"
            id="messages-container"
            bind:this={messagesContainer}
            onscroll={handleScroll}
          >
            <Messages
              {messages}
              {isLoadingMessages}
              {hasMoreMessages}
              bind:selectedThread
            />
          </div>
        {:else}
          <Welcome />
        {/if}
      </div>
    </Pane>

    {#if selectedThread != undefined}
      <Pane>
        <div class="main-content">
          <header class="channel-header">
            <h2>Thread</h2>
            <button
              class="close-thread-btn"
              aria-label="Close thread"
              onclick={() => (selectedThread = undefined)}
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
          <div
            class="thread-container"
            id="thread-container"
            bind:this={threadContainer}
          >
            <Messages
              messages={threadMessages}
              isLoadingMessages={false}
              hasMoreMessages={false}
              bind:selectedThread
            />
          </div>
        </div>
      </Pane>
    {/if}
  </Splitpanes>
</main>

<style>
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
