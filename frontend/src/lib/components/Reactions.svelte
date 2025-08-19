<script lang="ts">
import { models } from '../../../wailsjs/go/models'
import EmojiConvertor from 'emoji-js'

type Props = {
  reactions: models.Reaction[]
}

let { reactions }: Props = $props()

const converter = new EmojiConvertor()
const labelFor = (r: models.Reaction) => `:${r.name}: ${r.count} ${r.count === 1 ? 'reaction' : 'reactions'}`

function emojify(txt: string) {
  return converter.replace_colons(txt)
}
</script>

<div class="reactions">
    {#each reactions as r (r.name)}
        <button class="reaction" title={labelFor(r)}>
            {@html emojify(':' + r.name + ':')}
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
    .count {
        font-size: 0.78rem;
    }
</style>