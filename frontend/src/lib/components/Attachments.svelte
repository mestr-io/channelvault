<script lang="ts">
const { files } = $props()

import LazyImage from './AttachmentLazyImage.svelte'
import AttachmentFile from './AttachmentFile.svelte'
import type { models } from '../../../wailsjs/go/models'

function detectFileType(file: models.File) {
  if (file.mimetype.startsWith('image/')) {
    return 'image'
  } else if (file.mimetype.startsWith('video/')) {
    return 'video'
  } else if (file.mimetype.startsWith('audio/')) {
    return 'audio'
  } else {
    return 'generic'
  }
}
</script>

{#if files}
    <div>
        {#each files as file}
            {@const filetype = detectFileType(file)}
            {#if filetype === "image"}
                <LazyImage
                    url={file.url_private}
                    fileId={file.id}
                    alt={file.title || file.name}
                />
            {:else if filetype === "generic"}
                <AttachmentFile {file} />
            {/if}
        {/each}
    </div>
{/if}

<style>
</style>
