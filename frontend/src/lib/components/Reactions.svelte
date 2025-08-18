<!-- Reactions.svelte -->
<script lang="ts">
import { models } from '../../../wailsjs/go/models'
import EmojiConvertor from 'emoji-js';
type Props = {
  reactions: models.Reaction[]
}

let { reactions }: Props = $props()
let converter: EmojiConvertor = new EmojiConvertor()

// local mutable state (Svelte 5 rune)
let state = $state({ reactions })
const labelFor = (r: models.Reaction) => `:${r.name}: ${r.count} ${r.count === 1 ? 'reaction' : 'reactions'}`

</script>

<div class="reactions">
    {#each state.reactions as r (r.name)}
        <button class="reaction" title={labelFor(r)}>
            {@html converter.replace_colons(r.name)}
            <span class="count">{r.count}</span>
        </button>
    {/each}
</div>

<style>
    .reactions {
        display: flex;
        flex-wrap: wrap;
        gap: 4px;
        margin-top: 6px;
    }
    .reaction {
        display: inline-flex;
        align-items: center;
        gap: 6px;
        padding: 2px 8px;
        border: 1px solid #ddd;
        border-radius: 12px;
        background: #f9f9f9;
        cursor: pointer;
        font-size: 0.85rem;
        line-height: 1.2;
        transition:
            background 0.15s,
            border-color 0.15s;
    }
    .reaction:hover {
        background: #f1f1f1;
        border-color: #aaa;
    }
    .emoji-text {
        font-size: 1rem;
        line-height: 1;
    }
    .count {
        font-size: 0.78rem;
    }
</style>
