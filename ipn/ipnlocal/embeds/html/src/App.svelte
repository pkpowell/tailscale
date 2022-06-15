<script lang="ts">
    import { onMount } from "svelte"
    // import type { 
    //     SSEMessage 
    // } from "./types/types"
    import NodeInfo from "./lib/Nodeinfo.svelte"
    import Peers from './lib/Peers.svelte'
    // import { local, peers } from "./store/sse"

    // onMount(async()=> {
    //     let sse = new EventSource(`http://100.100.100.100/events/`)

    //     sse.onmessage = event => {
    //         let response: SSEMessage = JSON.parse(event.data)
    //         if(!response.length) {
    //             switch (response.type) {
    //                 case "ping":

    //                     break;

    //                 case "local":
    //                     console.log("local", response.payload)
    //                     local.set(response.payload)
    //                     break;

    //                 case "peer":
    //                     console.log("peer", response.payload)
    //                     // $peers = [...$peers, response.payload]
    //                     // peers.set(response.payload)
    //                     break;
                
    //                 default:
    //                     console.warn("unknown payload", response.payload)
    //                     break;
    //             }
    //         } else {
    //             console.warn("empty message", event)
    //         }
    //     }

    //     sse.onerror = event => {
    //         console.error("SSE error", event)
    //     }

    //     sse.onopen = event => {
    //         console.log("on open", event)
    //     }
    // })
</script>

<header class="flex bg-gray-100 items-center py-2 mb-8">
    <NodeInfo />
    <!-- <Peers /> -->
</header>

<style>
  @import "../../local.css";

  :root {
    font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
    font-size: 100%;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  }

</style>
