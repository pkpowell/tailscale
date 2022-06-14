import { writable } from 'svelte/store'
import type { 
      Peer, 
    AppData, 
} from "../types/types"

const local = writable<AppData>()
const peers = writable<Peer[]>()

export {
    local, peers
}