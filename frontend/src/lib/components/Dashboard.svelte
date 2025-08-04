<script lang="ts">
  import { GetChannels } from "../../../wailsjs/go/app/App";

  let channels = $state<string[]>([]);

  $effect(() => {
    (async () => {
      channels = await GetChannels();
    })();
  });
</script>

<main class="main">
  <div class="left-column">
    <h3>Channels</h3>
    <ul id="channels-list">
      {#each channels as channel}
        <li>{channel}</li>
      {/each}
    </ul>
  </div>
  <div class="right-column">
    <section class="section">
      <h2>Dashboard</h2>
      <p>Welcome to the dashboard! Here you can view and manage your data.</p>
      <p>Use the left column to navigate through different channels.</p>
    </section>
  </div>
</main>

<style>
  .main {
    display: flex;
    flex-flow: row nowrap;
    align-items: flex-start;
    height: 100vh;
    width: 100vw;
    /* gap: 10px; */
    background-color: white;
    color: #333333;
  }

  .left-column {
    order: 1;
    flex-grow: 0;
    text-align: left;
  }

  .right-column {
    order: 2;
    flex-grow: 4;
    align-content: flex-start;
  }

  .section {
    text-align: left;
  }
</style>
