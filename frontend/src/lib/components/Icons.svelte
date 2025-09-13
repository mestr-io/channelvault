<script lang="ts" module>
export { chevron, circle, import_icon }
</script>

// CHEVRON
{#snippet chevron(rotated = false)}
    <svg
        class:chevron-rotated={rotated}
        viewBox="0 0 20 20"
        fill="currentColor"
    >
        <path
            fill-rule="evenodd"
            d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z"
            clip-rule="evenodd"
        />
    </svg>
{/snippet}

// CIRCLE
{#snippet circle(radius = "4px", solid = false)}
    <svg
        width={radius}
        height={radius}
        viewBox="0 0 24 24"
        fill-rule="evenodd"
        xmlns="http://www.w3.org/2000/svg"
    >
        <path
            fill-rule={solid ? "nonzero" : "evenodd"}
            clip-rule="evenodd"
            d="M12 17C14.7614 17 17 14.7614 17 12C17 9.23858 14.7614 7 12 7C9.23858 7 7 9.23858 7 12C7 14.7614 9.23858 17 12 17ZM12 20C16.4183 20 20 16.4183 20 12C20 7.58172 16.4183 4 12 4C7.58172 4 4 7.58172 4 12C4 16.4183 7.58172 20 12 20Z"
            fill="currentColor"
        />
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
            d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
        ></path>
    </svg>
{/snippet}

<style>
    .chevron-rotated {
        transform: rotate(90deg);
    }
</style>
