<script lang="ts" module>
export { import_icon, search_icon }
</script>

// SEARCH_ICON
{#snippet search_icon(size = "1rem")}
    <svg
        style="height: {size}; width: {size};"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
    >
        <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            vector-effect="non-scaling-stroke"
            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
        ></path>
    </svg>
{/snippet}

// IMPORT_ICON
{#snippet import_icon(size = "1rem")}
    <svg
        style="height: {size}; width: {size};"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
    >
        <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            vector-effect="non-scaling-stroke"
            d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
        ></path>
    </svg>
{/snippet}
