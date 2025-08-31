<script lang="ts">
const { attachment } = $props()
function getYouTubeId(url: string): string | null {
  const match = url.match(/(?:youtube\.com\/.*v=|youtu\.be\/)([a-zA-Z0-9_-]{11})/)
  return match ? match[1] : null
}
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
        <iframe
            class="youtube-embed"
            src={`https://www.youtube.com/embed/${getYouTubeId(
                attachment.original_url,
            )}`}
            title={attachment.title || attachment.fallback}
            frameborder="0"
            allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
            allowfullscreen
        ></iframe>
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
    .youtube-embed {
        width: 400px;
        height: 300px;
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
