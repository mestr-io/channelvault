<script lang="ts">
import type { models } from '../../../../wailsjs/go/models'
import LazyImage from './LazyImage.svelte'
import GenericFile from './File.svelte'

const { attachment } = $props()
</script>

<div class="attachment">
    <div class="attachment-header">
        <img
            class="service-icon"
            src={attachment.service_icon}
            alt={attachment.service_name}
        />
        <span class="service-name">{attachment.service_name}</span>
    </div>
    <div class="attachment-title">
        <a
            href={attachment.title_link}
            target="_blank"
            rel="noopener noreferrer"
        >
            {attachment.title || attachment.fallback}
        </a>
    </div>
    <div class="attachment-image-row">
        <img
            class="external-image"
            src={attachment.image_url}
            alt={attachment.title || attachment.fallback}
            width={Math.min(attachment.image_width, 400)}
            height="auto"
        />
    </div>
</div>

<style>
    .attachment {
        display: flex;
        flex-direction: column;
        gap: 0.5em;
        padding: 0.2em 0 0 0.7em;
        border-left: 6px solid #cfcfcf;
    }

    .attachment-header {
        display: flex;
        align-items: center;
        gap: 0.5em;
    }

    .service-icon {
        width: 20px;
        height: 20px;
        object-fit: contain;
    }

    .service-name {
        font-weight: 700;
        font-size: 0.8em;
    }
    .external-image {
        max-width: 400px;
        max-height: 300px;
        border-radius: 6px;
    }
    .attachment-title {
        font-weight: 600;
        font-size: 0.8em;
    }
    .attachment-title a {
        color: #1264a3;
        text-decoration: none;
    }
</style>
