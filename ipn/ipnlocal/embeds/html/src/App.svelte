<script lang="ts">
    import NodeInfo from "./lib/Nodeinfo.svelte"
    import Peers from './lib/Peers.svelte'
</script>

<header class="bg-gray-100 items-center py-2 mb-8">
  <div class="width-80">
    <NodeInfo />
    <Peers />
  </div>
</header>

<style>
@import "./assets/css/local.css";
</style>
