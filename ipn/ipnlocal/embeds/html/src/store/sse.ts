import { writable } from 'svelte/store'

import type { 
    Peer, 
    AppData, 
    SSEMessage
} from "../types/types"

const local = writable<AppData>()
const peers = writable<Peer[]>()
const peerMap = writable<Map<string, Peer>>(new Map<string, Peer>())

const updatePeers = (p: Peer) => {
    peerMap.update(records => records.set(p.ID, p))
}

// const appendPeer =  (p: Peer) => {
//     peers.update(current => {
//         if (current === null ||  typeof current === "undefined") current = []
//         console.log("current", current)
//         current = [...current, p]
//         return current
//     })
// }

const localReady = writable(false)
const peersReady = writable(false)

let sseUrl: string

if (import.meta.env.DEV) {
    sseUrl = `//100.100.100.100/events/`
} else {
    sseUrl = `/events/`
}

let sse = new EventSource(sseUrl)

sse.onmessage = event => {
    let response: SSEMessage = JSON.parse(event.data)
    if(!response.length) {
        switch (response.type) {
            case "ping":

                break;

            case "local":
                // console.log("local", response.payload)
                local.set(response.payload)
                localReady.set(true)
                break;
                
                case "peer":
                    // console.log("peer", response.payload)
                    updatePeers(response.payload)
                    peersReady.set(true)
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
    local, peers, localReady, peersReady, peerMap
}