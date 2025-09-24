<script lang="ts" module>
    export { closeIcon };
</script>

{#snippet closeIcon(size = "1.25rem")}
    <svg
        viewBox="0 0 20 20"
        fill="currentColor"
        width={size}
        height={size}
        aria-hidden="true"
        focusable="false"
    >
        <path
            fill-rule="evenodd"
            clip-rule="evenodd"
            d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
        />
    </svg>
{/snippet}
