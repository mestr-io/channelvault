<script lang="ts" module>
export { doneIcon }
</script>

{#snippet doneIcon(size = "1.25rem", color = "#6b7280")}
    <svg
        viewBox="0 0 20 20"
        fill="currentColor"
        width={size}
        height={size}
        style={`color:${color}`}
    >
        <path
            fill-rule="evenodd"
            d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
            clip-rule="evenodd"
        />
    </svg>
{/snippet}
