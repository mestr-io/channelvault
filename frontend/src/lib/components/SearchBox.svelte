<script lang="ts">
    import { models } from '../../../wailsjs/go/models'
    import SearchResultsOverlay from './SearchResultsOverlay.svelte';
    import { SearchMessages } from '../../../wailsjs/go/app/App'

    let { selectedThread = $bindable(), selectedChannel = $bindable() , currentTs = $bindable()} = $props()

    let searchText = $state("");
    let searchResults = $state<models.SearchResult[]>([]);
    let showOverlay = $state(false);

    function handleSearch() {
        console.log("Searching for:", searchText);
        SearchMessages(searchText).then((result) => {
            searchResults = result;
            showOverlay = result.length > 0;
        });
    }

    function closeOverlay() {
        showOverlay = false;
        searchResults = [];
    }

    function handleKeyPress(event: KeyboardEvent) {
        if (event.key === "Enter") {
            handleSearch();
        }
    }
</script>

<header class="search-box-container">
    <input
        type="text"
        placeholder="Search messages..."
        bind:value={searchText}
        onkeydown={handleKeyPress}
        class="search-input"
    />
    <button onclick={handleSearch} class="search-button">Search</button>
</header>

{#if showOverlay}
    <SearchResultsOverlay
        bind:selectedChannel bind:selectedThread bind:currentTs
        {searchResults}
        {searchText}
        onClose={closeOverlay}
    />
{/if}

<style>
    .search-box-container {
        padding: 1rem 1.25rem 0.5rem;
        display: flex;
        gap: 10px;
        align-items: center;
        background-color: #552958; /* Match sidebar background */
        border-bottom: 1px solid #7a4b87; /* Match sidebar header border */
    }

    .search-input {
        flex-grow: 1;
        background-color: #3e0e42; /* Similar to channel-filter-input */
        border: 1px solid #663867;
        border-radius: 4px;
        color: #dbc9e0;
        padding: 0.5rem 0.75rem;
        font-size: 0.9rem;
    }

    .search-input::placeholder {
        color: #a08fa6;
    }

    .search-input:focus {
        outline: none;
        border-color: #dbc9e0;
        box-shadow: 0 0 0 1px #dbc9e0;
    }

    .search-button {
        background-color: #7a4b87; /* A shade of purple from the theme */
        color: #f9edff;
        border: none;
        border-radius: 4px;
        padding: 0.5rem 1rem;
        font-size: 0.9rem;
        cursor: pointer;
        transition: background-color 0.2s ease-in-out;
    }

    .search-button:hover {
        background-color: #8c5a96; /* Slightly lighter purple on hover */
    }

    .search-button:active {
        background-color: #663867; /* Darker purple on active */
    }
</style>
