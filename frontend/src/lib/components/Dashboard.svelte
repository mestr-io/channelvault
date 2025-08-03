<script lang="ts">

</script>

<main class="main">
  <div class="left-column">
    <h3>Channels</h3>
    <ul id="channels-list">
        <li>Channel 1</li>
        <li>Channel 2</li>
        <li>Channel 3</li>
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
    gap: 10px;
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
