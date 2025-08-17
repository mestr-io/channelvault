<script lang="ts">
  import { models } from "../../../wailsjs/go/models";
  import ThreadInfo from "./ThreadInfo.svelte";
  import LazyImage from "./LazyImage.svelte";
  import DayDivider from "./DayDivider.svelte";
  import Message from "./Message.svelte";

  import { VList } from "virtua/svelte";

  type Props = {
    messages: models.Message[];
    isLoadingMessages: boolean;
    hasMoreMessages: boolean;
    selectedThread: string | null;
  };
  let {
    messages,
    isLoadingMessages,
    hasMoreMessages,
    selectedThread = $bindable(null),
  }: Props = $props();

  /**
   * Formats the message timestamp to show only the time (e.g., "2:03 PM").
   */
  const formatMessageTime = (ts: string): string => {
    const date = new Date(parseFloat(ts) * 1000);
    // Using 'en-US' locale to get AM/PM format
    return date.toLocaleTimeString("en-US", {
      hour: "numeric",
      minute: "2-digit",
    });
  };
</script>

{#if hasMoreMessages}
  <div class="loader">
    {#if isLoadingMessages}<span>Loading history...</span>{/if}
  </div>
{/if}

<VList data={messages} getKey={(msg, _) => msg.ts}>
  {#snippet children(msg, i)}
    <DayDivider current={msg} previous={messages[i - 1]} />
    <Message {msg} {selectedThread} />  
  {/snippet}
</VList>


<style>
  .loader {
    text-align: center;
    padding: 1rem;
    min-height: 40px;
  }
  .messages-container {
    margin: 1rem;
    margin-top: auto;
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
  .message-content {
    flex: 1;
    max-width: calc(100% - 50px);
  }
  .message-content p {
    margin: 0;
    line-height: 1.45;
    display: block;
    max-width: 98%;
    word-wrap: break-word;
    overflow-wrap: break-word;
    white-space: normal;
  }
  .message-header strong {
    font-weight: 900;
  }
  .timestamp {
    color: #9e9ea6;
    font-size: 0.8rem;
    margin-left: 0.5rem;
  }

  /* START: New Styles for File Attachments */
  .files-container {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-top: 8px;
  }

  .file-preview {
    text-decoration: none;
    color: inherit;
    display: inline-flex;
    border: 1px solid #ddd;
    border-radius: 8px;
    overflow: hidden;
    max-width: 360px;
  }

  .file-preview:hover {
    border-color: #999;
  }

  .image-preview img {
    max-width: 100%;
    max-height: 250px;
    display: block;
    object-fit: cover;
  }

  .generic-file {
    align-items: center;
    padding: 10px;
    background-color: #f8f8f8;
  }

  .file-icon {
    width: 32px;
    height: 32px;
    background-color: #e8e8e8; /* Placeholder for a real icon */
    border-radius: 6px;
    margin-right: 12px;
    flex-shrink: 0;
    /* In a real app, you'd set a background-image with an SVG icon */
  }

  .file-info {
    display: flex;
    flex-direction: column;
    overflow: hidden; /* Prevents long file names from breaking layout */
  }

  .file-name {
    font-weight: 700;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .file-title {
    font-size: 0.85rem;
    color: #555;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  /* END: New Styles for File Attachments */
</style>
