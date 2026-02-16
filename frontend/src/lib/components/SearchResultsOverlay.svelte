<script lang="ts">
    import { models } from "../../../wailsjs/go/models";
    import type { ThreadID } from "../types";

    type Props = {
        selectedThread: string;
        selectChannel: (name: string) => Promise<void>;
        currentTs: string;
        searchText: string;
        onClose: () => void;
        searchResults: models.SearchResult[];
    };

    let {
        onClose,
        searchText,
        searchResults,
        selectedThread = $bindable(),
        selectChannel,
        currentTs = $bindable(),
    }: Props = $props();
    
    async function selectHandler(result: models.SearchResult) {
        let thread: ThreadID | null = null;
        if (result.ThreadTs) {
            thread = {
                thread_ts: result.ThreadTs,
                import_id: result.ImportID,
                current_ts: result.Ts,
            };
        }
        await selectChannel(result.ChannelID, result.Ts, thread);
        onClose();
    }

    // Adjusted to match the "Dec 7th, 2023" style format seen in the screenshot
    const formatMessageDate = (ts: string): string => {
        const date = new Date(parseFloat(ts) * 1000);
        return date.toLocaleDateString("en-US", {
            month: "short",
            day: "numeric",
            year: "numeric",
        });
    };
</script>

<div
    class="overlay-backdrop"
    role="button"
    tabindex="0"
    onclick={onClose}
    onkeydown={(e) => {
        if (e.key === "Escape" || e.key === "Enter" || e.key === " ") {
            onClose();
        }
    }}
>
    <div
        class="overlay-content"
        role="dialog"
        aria-modal="true"
        tabindex="-1"
        onclick={(e) => e.stopPropagation()}
        onkeydown={(e) => e.stopPropagation()}
    >
        <div class="overlay-header">
            <h3>Search: <strong>"{searchText}"</strong></h3>
            <button class="close-button" onclick={onClose}>&times;</button>
        </div>

        <div class="overlay-body">
            {#if searchResults.length > 0}
                {#each searchResults as result}
                    <div
                        class="search-result-item"
                        role="button"
                        tabindex="0"
                        onclick={() => selectHandler(result)}
                        onkeydown={(e) => {
                            if (e.key === "Enter" || e.key === " ") {
                                e.preventDefault();
                                selectHandler(result);
                            }
                        }}
                    >
                        <div class="avatar-container">
                            <img
                                class="avatar"
                                src={result.ProfileUrl}
                                alt="Avatar"
                            />
                        </div>

                        <div class="content-container">
                            <div class="message-header">
                                <span class="user-name">{result.Username}</span>

                                <span class="context">
                                    {#if result.ThreadTs}
                                        Thread in <strong
                                            ># {result.ChannelID}</strong
                                        >
                                    {:else}
                                        <strong># {result.ChannelID}</strong>
                                    {/if}
                                </span>

                                <span class="timestamp">
                                    {formatMessageDate(result.Ts)}
                                </span>
                            </div>

                            <div class="matched-text">
                                {@html result.MatchedText}
                            </div>
                        </div>
                    </div>
                {/each}
            {:else}
                <div class="no-results">
                    <p>No results found.</p>
                </div>
            {/if}
        </div>
    </div>
</div>

<style>
    /* Base typography applied to the overlay */
    .overlay-content {
        font-family:
            -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica,
            Arial, sans-serif;
    }

    .overlay-backdrop {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background-color: rgba(0, 0, 0, 0.5); /* Lighter backdrop */
        display: flex;
        justify-content: center;
        align-items: center;
        z-index: 1000;
    }

    .overlay-content {
        background-color: #ffffff; /* Changed to white */
        color: #1d1c1d; /* Standard dark text */
        border-radius: 8px;
        width: 85%;
        max-width: 900px;
        height: 85vh;
        max-height: 800px;
        display: flex;
        flex-direction: column;
        overflow: hidden;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    }

    .overlay-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 20px 24px;
        border-bottom: 1px solid #e2e2e2;
        background-color: #f8f8f8;
    }

    .overlay-header h3 {
        margin: 0;
        font-size: 1.25rem;
        font-weight: 400;
        color: #1d1c1d;
    }

    .close-button {
        background: none;
        border: none;
        color: #616061;
        font-size: 1.5rem;
        cursor: pointer;
        padding: 4px;
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 4px;
    }

    .close-button:hover {
        background-color: #e8e8e8;
        color: #1d1c1d;
    }

    .overlay-body {
        overflow-y: auto;
        flex-grow: 1;
        background-color: #ffffff;
    }

    .search-result-item {
        display: flex;
        padding: 16px 24px;
        border-bottom: 1px solid #f2f2f2;
        cursor: pointer;
        transition: background-color 0.1s ease;
    }

    .search-result-item:hover {
        background-color: #f8f8f8;
    }

    .avatar-container {
        flex-shrink: 0;
        margin-right: 12px;
    }

    .avatar {
        width: 36px;
        height: 36px;
        border-radius: 4px; /* Slack style slightly rounded square */
        object-fit: cover;
    }

    .content-container {
        flex-grow: 1;
        min-width: 0; /* Prevents flex children from overflowing */
    }

    .message-header {
        display: flex;
        align-items: baseline;
        margin-bottom: 4px;
        gap: 8px;
    }

    .user-name {
        font-weight: 900;
        color: #1d1c1d;
        font-size: 0.95rem;
    }

    .context {
        font-size: 0.85rem;
        color: #616061;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .context strong {
        font-weight: 600;
    }

    .timestamp {
        margin-left: auto;
        font-size: 0.8rem;
        color: #616061;
        flex-shrink: 0;
    }

    .matched-text {
        font-size: 0.95rem;
        color: #1d1c1d;
        line-height: 1.46;
        word-break: break-word;
    }

    /* Style for highlighted text matching the yellow in the screenshot */
    :global(.matched-text mark) {
        background-color: #fbf3d0;
        color: inherit;
        padding: 0 3px;
        border-radius: 3px;
        font-weight: inherit;
    }

    /* Additional styling for inline code blocks (like the pink ones in the screenshot) */
    :global(.matched-text code) {
        color: #e01e5a;
        background-color: rgba(29, 28, 29, 0.04);
        padding: 2px 4px;
        border-radius: 3px;
        border: 1px solid rgba(29, 28, 29, 0.13);
        font-family: Monaco, Menlo, Consolas, "Courier New", monospace;
        font-size: 0.85rem;
    }

    .no-results {
        padding: 40px;
        text-align: center;
        color: #616061;
        font-size: 1rem;
    }
</style>
