<script lang="ts">
import { onMount } from 'svelte'

let isVisible = $state(false)

onMount(() => {
  setTimeout(() => {
    isVisible = true
  }, 300)
})
</script>

<div class="no-messages-container" class:is-visible={isVisible}>
    <h2>This channel has no messages.</h2>
    <p>Please select another channel or import more data.</p>
</div>

<style>
    .no-messages-container {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 100%;
        color: #9e9ea6;
        text-align: center;
        padding: 20px;
        opacity: 0;
        transition: opacity 0.3s ease-in-out;
    }

    .no-messages-container.is-visible {
        opacity: 1;
    }

    .no-messages-container h2 {
        font-size: 1.5rem;
        margin-bottom: 10px;
    }

    .no-messages-container p {
        font-size: 1rem;
    }
</style>
