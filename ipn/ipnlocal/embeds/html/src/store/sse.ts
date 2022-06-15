import { writable } from 'svelte/store'
import type { 
      Peer, 
    AppData, 
    SSEMessage
} from "../types/types"

const local = writable({})
// const local = writable<AppData>()
// local.update(value => Object.assign(value, {HostName: ""}))
const peers = writable([])
// const peers = writable<Peer[]>()

let sse = new EventSource(`http://100.100.100.100/events/`)

sse.onmessage = event => {
    let response: SSEMessage = JSON.parse(event.data)
    if(!response.length) {
        switch (response.type) {
            case "ping":

                break;

            case "local":
                console.log("local", response.payload)
                local.set(response.payload)
                // local.update(value => Object.assign(value, response.payload))
                break;

            case "peer":
                console.log("peer", response.payload)
                // $peers = [...$peers, response.payload]
                // peers.set(response.payload)
                break;
        
            default:
                console.warn("unknown payload", response.payload)
                break;
        }
    } else {
        console.warn("empty message", event)
    }
}

sse.onerror = event => {
    console.error("SSE error", event)
}

sse.onopen = event => {
    console.log("on open", event)
}


export {
    local, peers
}