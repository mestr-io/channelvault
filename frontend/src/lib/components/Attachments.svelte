<script lang="ts">
import type { models } from '../../../wailsjs/go/models'
import LazyImage from './attachments/LazyImage.svelte'
import GenericFile from './attachments/File.svelte'
import ExternalImage from './attachments/ExternalImage.svelte'
import ExternalVideo from './attachments/ExternalVideo.svelte'

const { message } = $props()

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
function detectAttachment(attachment: models.Attachment) {
  if (attachment.image_url) {
    return 'image'
  }
  if (attachment.service_name.toLowerCase() === 'youtube') {
    return 'youtube'
  }

  return 'generic'
}

function showAttachments(msg) {
  $inspect(msg.attachments)
  return msg.attachments
}
</script>

{#if message.files}
  <div>
    {#each message.files as file}
      {@const filetype = detectFileType(file)}
      {#if filetype === "image"}
        <LazyImage
          url={file.url_private}
          fileId={file.id}
          alt={file.title || file.name}
        />
      {:else if filetype === "generic"}
        <GenericFile {file} />
      {/if}
    {/each}
  </div>
{/if}

{#if message.attachments}
  <div>
    {#each message.attachments as attachment}
      {@const attachmentType = detectAttachment(attachment)}
      {#if attachmentType === "image"}
        <ExternalImage {attachment} />
      {:else if attachmentType === "youtube"}
        <ExternalVideo {attachment} />
      {/if}
    {/each}
  </div>
{/if}

<style>
</style>
