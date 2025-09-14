<script lang="ts" module>
export { burgerMenuIcon }
</script>

{#snippet burgerMenuIcon(size = "1.25rem", color = "#6b7280")}
    <svg
        viewBox="0 0 20 20"
        fill="currentColor"
        width={size}
        height={size}
        style={`color:${color}`}
    >
        <path
            fill-rule="evenodd"
            d="M3 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z"
            clip-rule="evenodd"
        />
    </svg>
{/snippet}
