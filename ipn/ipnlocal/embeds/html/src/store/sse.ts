import { writable } from 'svelte/store'
// import type { Writable } from 'svelte/store'

import type { 
    Peer, 
    AppData, 
    SSEMessage
} from "../types/types"

const local = writable<AppData>()
// const local = writable<AppData>()
// local.update(value => Object.assign(value, {HostName: ""}))
const peers = writable<Peer[]>()
const peerMap = writable<Map<string, Peer>>(new Map<string, Peer>())

// const setPeers = () => {
//     peerMap.set((x)=>x=m)
// }
const updatePeers = (p: Peer) => {
    peerMap.update(records =>{
        // if (typeof records !== "undefined") 
        console.log("records",records)
        return records.set(p.ID, p)
        // records.set(p.ID,  p)
    })
}

const appendPeer =  (p: Peer) => {
    peers.update(current => {
        if (current === null ||  typeof current === "undefined") current = []
        console.log("current", current)
        current = [...current, p]
        return current
    })
}

const localReady = writable(false)
const peersReady = writable(false)

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
                localReady.set(true)
                // local.update(value => Object.assign(value, response.payload))
                break;
                
                case "peer":
                    console.log("peer", response.payload)
                    updatePeers(response.payload)
                    peersReady.set(true)
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
    local, peers, localReady, peersReady, peerMap
}