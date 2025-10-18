<script lang="ts">
import { onMount } from 'svelte'
import { FetchChangelog } from '../../../wailsjs/go/app/App'

let changelog = 'Loading changelog...'

onMount(async () => {
  try {
    changelog = await FetchChangelog()
  } catch (e) {
    changelog = 'Failed to load changelog.'
  }
})
</script>

<div class="welcome-screen">
    <div class="app-icon-row">
        <img
            src="./src/assets/images/appicon.png"
            alt="ChannelVault icon"
            class="app-icon"
        />
    </div>
    <div class="app-name-row">
        <h1>ChannelVault</h1>
        Select a channel or import a new workspace to get started.
    </div>
    <div class="changelog-row">
        <h2>Changelog</h2>
        <div class="changelog-content">
            {@html changelog}
        </div>
    </div>
</div>

<style>
    .welcome-screen {
        display: flex;
        flex-direction: column;
        align-items: center;
        height: 100%;
        padding: 2em 0;
        gap: 2em;
    }

    .app-icon-row {
        display: flex;
        justify-content: center;
        align-items: center;
    }
    .app-icon {
        width: 200px;
        border-radius: 12px;
    }
    .app-name-row h1 {
        font-size: 2em;
        font-weight: 700;
        margin: 0;
        color: #552958;
        text-align: center;
    }

    .changelog-row {
        width: 90vw;
        max-width: 600px;
        border-radius: 8px;
        padding: 1em;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
        overflow-x: auto;
    }
    .changelog-content {
        white-space: pre-wrap;
        font-family: "Fira Mono", "Consolas", monospace;
        font-size: 1em;
        color: #333;
        margin: 0;
    }
</style>
