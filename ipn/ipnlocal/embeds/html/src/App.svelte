<script lang="ts">
    // import { onMount } from "svelte"

    import NodeInfo from "./lib/Nodeinfo.svelte"
    import Peers from './lib/Peers.svelte'
</script>

<!-- <header class="flex bg-gray-100 items-center py-2 mb-8"> -->
    <NodeInfo />
    <Peers />
<!-- </header> -->

<style>
  @import "../../local.css";

  /* :root {
    font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
    font-size: 100%;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  } */

</style>
