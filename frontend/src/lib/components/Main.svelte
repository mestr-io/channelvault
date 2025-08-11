<script lang="ts">
    import { GetChannels } from "../../../wailsjs/go/app/App";
    import { models } from "../../../wailsjs/go/models";
    import History from './ChatHistory.svelte'
    import Sidebar from "./Sidebar.svelte";

</script>

<main>
    <Sidebar />
    <History channelId="quien-apizza-no-es-traidor"/>
</main>

<style>
</style>
