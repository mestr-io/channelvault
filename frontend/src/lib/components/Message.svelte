<script lang="ts">
    import { models } from "../../../wailsjs/go/models";
    import ThreadInfo from "./ThreadInfo.svelte";
    import Attachments from "./Attachments.svelte";
    import Reactions from "./Reactions.svelte";
    import { Carta, Markdown } from "carta-md";
    import { getContext } from "svelte";

    type Props = {
        msg: models.Message;
        selectedThread: string | null;
    };

    let { msg, selectedThread = $bindable() }: Props = $props();

    const carta: Carta = getContext("carta");

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

<div class="message" data-ts={msg.ts}>
    <img class="avatar" src={msg.avatar_url} alt={msg.user_name} />
    <div class="message-content">
        <div class="message-header">
            <strong>{msg.user_name}</strong>
            <span class="timestamp">
                {formatMessageTime(msg.ts)}
            </span>
        </div>
        <Markdown {carta} value={msg.text!} />
        <Attachments files={msg.files} />

        {#if msg.thread_ts}
            <ThreadInfo
                replyUsers={msg.reply_users}
                replyCount={msg.reply_count}
                threadTs={msg.thread_ts}
                bind:selectedThread
            />
        {/if}
        <Reactions reactions={msg.reactions} />
    </div>
</div>

<style>
    :global(.carta-viewer > :first-child) {
        margin-top: 2px;
    }
    .message {
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
