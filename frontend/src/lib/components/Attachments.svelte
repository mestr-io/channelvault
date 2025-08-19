<script lang="ts">
    const { files } = $props();

    import LazyImage from "./LazyImage.svelte";
</script>

{#if files}
    <div class="files-container">
        {#each files as file}
            {#if file.mimetype.startsWith("image/")}
                <a
                    href={file.url_private}
                    target="_blank"
                    rel="noopener noreferrer"
                >
                    <LazyImage
                        file_id={file.id}
                        alt={file.title || file.name}
                    />
                </a>
            {:else}
                <a
                    href={file.url_private}
                    target="_blank"
                    rel="noopener noreferrer"
                    download={file.name}
                    class="file-preview generic-file"
                >
                    <div class="file-icon"></div>
                    <div class="file-info">
                        <div class="file-name">
                            {file.name || "file_not_found"}
                        </div>
                        <div class="file-title">{file.title}</div>
                    </div>
                </a>
            {/if}
        {/each}
    </div>
{/if}

<style>
    .attachments {
        border-top: 1px solid #eee;
        padding: 1rem 0;
    }
</style>
