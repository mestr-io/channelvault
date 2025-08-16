<script lang="ts">
import { models } from '../../../wailsjs/go/models'
const { current, previous } = $props()

const shouldShowDayPill = (currentMsg: models.Message, previousMsg: models.Message | undefined): boolean => {
  if (!previousMsg) {
    return true // Always show for the very first message
  }
  const currentDate = new Date(parseFloat(currentMsg.ts) * 1000)
  const previousDate = new Date(parseFloat(previousMsg.ts) * 1000)

  return currentDate.toDateString() !== previousDate.toDateString()
}

const formatDayPillDate = (ts: string): string => {
  const date = new Date(parseFloat(ts) * 1000)
  // Using 'en-US' locale for consistent formatting like "Friday, August 8"
  return date.toLocaleDateString('en-US', {
    weekday: 'long',
    month: 'long',
    day: 'numeric',
  })
}
</script>

{#if shouldShowDayPill(current, previous)}
    <div class="day-divider">
        <div></div>
        <div class="day-pill">{formatDayPillDate(current.ts)}</div>
        <div></div>
    </div>
{/if}

<style>
    .day-divider {
        display: flex;
        align-items: center;
        flex-direction: row;
        text-align: center;
        color: #555;
        padding: 1rem 0;
        font-size: 0.8rem;
    }

    .day-divider div {
        flex-grow: 1;
        border: none;
        border-top: 1px solid #ddd;
    }

    .day-divider .day-pill {
        flex-grow: 0;
        padding: 0.25rem 0.75rem;
        border: 1px solid #ddd;
        border-radius: 24px;
        background-color: white;
        font-weight: bold;
        box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
    }
</style>