<script lang="ts">
    import NodeInfo from "./lib/Nodeinfo.svelte"
    import Peers from './lib/Peers.svelte'
</script>

<header class=" bg-gray-100 items-center py-2 mb-8">
    <NodeInfo />
    <Peers />
</header>

<style>
  @import "../../local.css";
</style>
