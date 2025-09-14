<script lang="ts" module>
export { chevron }
</script>

{#snippet chevron(rotated = false)}
    {@const rotation = `rotate(${rotated ? 90 : 0} 10 10)`}
    <svg class:rotated viewBox="0 0 20 20" fill="currentColor">
        <path
            transform={rotation}
            fill-rule="evenodd"
            d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z"
            clip-rule="evenodd"
        />
    </svg>
{/snippet}
