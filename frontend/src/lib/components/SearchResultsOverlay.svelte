<script lang="ts">
    import { models } from '../../../wailsjs/go/models'


    type Props = {
        selectedThread: string
        selectChannel: (name: string) => Promise<void>
        currentTs: string
        searchText: string
        onClose: () => void
        searchResults: models.SearchResult[]
    }

    let { onClose, searchText, searchResults, selectedThread = $bindable(), selectChannel , currentTs = $bindable()}: Props = $props()

    async function selectHandler(channel: string, ts: string, thread_ts: string) {
        await selectChannel(channel)
        currentTs = ts
        selectedThread = thread_ts
        onClose()
    }

</script>

<div class="overlay-backdrop">
    <div class="overlay-content">
        <div class="overlay-header">
            <h3>Search Results for "{searchText}"</h3>
            <button class="close-button" onclick={onClose}>&times;</button>
        </div>
        <div class="overlay-body">
            {#if searchResults.length > 0}
                {#each searchResults as result}
                    <div class="search-result-item">
                        <button onclick={() => selectHandler(result.ChannelID, result.Ts, result.ThreadTs)}>Go</button>
                        <p class="matched-text">
                            {@html result.MatchedText}
                        </p>
                        <p class="metadata">
                            Channel: {result.ChannelID} | User: {result.UserId} |
                            Timestamp: {result.Ts}
                        </p>
                    </div>
                {/each}
            {:else}
                <p>No results found.</p>
            {/if}
        </div>
    </div>
</div>

<style>
    .overlay-backdrop {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background-color: rgba(0, 0, 0, 0.7);
        display: flex;
        justify-content: center;
        align-items: center;
        z-index: 1000;
    }

    .overlay-content {
        background-color: #3e0e42;
        color: #dbc9e0;
        border-radius: 8px;
        width: 80%;
        max-width: 800px;
        max-height: 90%;
        display: flex;
        flex-direction: column;
        overflow: hidden;
    }

    .overlay-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1rem 1.5rem;
        border-bottom: 1px solid #663867;
    }

    .overlay-header h3 {
        margin: 0;
        font-size: 1.2rem;
    }

    .close-button {
        background: none;
        border: none;
        color: #dbc9e0;
        font-size: 1.8rem;
        cursor: pointer;
    }

    .close-button:hover {
        color: #f9edff;
    }

    .overlay-body {
        padding: 1.5rem;
        overflow-y: auto;
        flex-grow: 1;
    }

    .search-result-item {
        background-color: #552958;
        border-radius: 4px;
        padding: 1rem;
        margin-bottom: 1rem;
    }

    .search-result-item:last-child {
        margin-bottom: 0;
    }

    .matched-text {
        font-size: 1rem;
        margin-bottom: 0.5rem;
    }

    .metadata {
        font-size: 0.8rem;
        color: #a08fa6;
    }

    mark {
        background-color: #a08fa6;
        color: #3e0e42;
        padding: 0 2px;
        border-radius: 2px;
    }
</style>
