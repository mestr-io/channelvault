<script lang="ts">
    import { BrowserOpenURL } from "../../../wailsjs/runtime";
    const { href, title, rel, target } = $props();
    const truncatedHref = $derived(
        href.length > 60 ? href.substring(0, 60) + "…" : href,
    );
</script>

<a href={"#"} {title} {rel} {target} onclick={() => BrowserOpenURL(href)}>
    {truncatedHref}
</a>

<style>
    a {
        color: #1264a3;
        text-decoration: none;
    }
</style>
