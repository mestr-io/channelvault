<script lang="ts">
import type { models } from '../../../../wailsjs/go/models'
import LazyImage from './LazyImage.svelte'
import GenericFile from './File.svelte'
import { onMount, onDestroy } from 'svelte'

const { attachment } = $props()

const src = $derived(attachment.image_url)
let imgEl: HTMLImageElement
let observer: IntersectionObserver
let controller: AbortController | null = null
let objectUrl: string | null = null
let visible = false

async function loadImage() {
  if (src || objectUrl) return
  controller = new AbortController()

  try {
    const res = await fetch(src, { signal: controller.signal })
    const blob = await res.blob()
    objectUrl = URL.createObjectURL(blob)
    if (imgEl) imgEl.src = objectUrl
    console.log(objectUrl)
  } catch (e) {
    if (e instanceof Error && (e as any).name !== 'AbortError') console.error('Image load failed:', e)
  }
}

function unloadImage() {
  controller?.abort()
  controller = null
  if (objectUrl) {
    URL.revokeObjectURL(objectUrl)
    objectUrl = null
    if (imgEl) imgEl.src = '' // clear
  }
}

onMount(() => {
  observer = new IntersectionObserver(
    ([entry]) => {
      if (entry.isIntersecting) {
        // delay to avoid triggering mid-scroll
        const t = setTimeout(() => {
          visible = true
          loadImage()
        }, 150)

        return () => clearTimeout(t)
      } else {
        visible = false
        unloadImage()
      }
    },
    { rootMargin: '200px' },
  )

  observer.observe(imgEl)
  return () => {
    observer.disconnect()
    unloadImage()
  }
})

onDestroy(() => {
  observer?.disconnect()
  unloadImage()
})
</script>

<div class="attachment">
    <div class="attachment-header">
        <img
            class="service-icon"
            src={attachment.service_icon}
            alt={attachment.service_name}
            loading="lazy"
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
            bind:this={imgEl}
            class="external-image"
            src={attachment.image_url}
            alt={attachment.title || attachment.fallback}
            width={Math.min(attachment.image_width, 400)}
            height="auto"
            loading="lazy"
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
