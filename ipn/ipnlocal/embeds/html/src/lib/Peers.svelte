<div class="peers">
    <div class="py-8 text-3xl font-semibold tracking-tight leading-tight">Peers</div>
    <!-- svelte-ignore a11y-autofocus -->
    <input
        placeholder="Search..."
        bind:value ={searchTerm}
        autofocus
    >
    <div class="table">

        <div on:click={() => sort("ID")} class="sticky header">ID</div>
        <div on:click={() => sort("HostName")} class="sticky header">name</div>
        <div on:click={() => sort(ip.k)} class="sticky header text-right">IP<button class="ip-toggle" on:click|stopPropagation={toggleIP}>{ip.f}</button></div>
        <div on:click={() => sort("LastSeen")} class="sticky header">Last Seen</div>
        <div on:click={() => sort("RXb")} class="sticky header text-right">rx</div>
        <div on:click={() => sort("TXb")} class="sticky header text-right">tx</div>


            {#if $peersReady}
            {#each peers() as p}

            <!-- <div on:click={() => toggleDetails(parseInt(p.ID))} class="row table-row w-full px-0.5 hover:bg-gray-0 py-1"> -->
                <div class="">
                    <span>
                        {p.ID}
                    </span>
                </div>

                <div class="text-ellipsis">
                    <h3 class="font-semibold">
                        {p.HostName}
                    </h3>
                </div>

                <div class="text-right mono">
                        {#if ip.v === IP.v4}
                            <div class="">
                                <div class="">{p.IPv4}</div>
                            </div>
                        {:else}
                            <div class="">
                                <span class="">{p.IPv6}</span>
                            </div>
                        {/if}
                </div>

                <div class="" title="{new Date(p.LastSeen).toLocaleDateString("en-US", dateShort)}">{ago(p.LastSeen, p.Unseen)}</div>

                <div class="text-right">{bytes(p.RXb)}</div>

                <div class="text-right">{bytes(p.TXb)}</div>
            <!-- </div> -->

            <!-- <div style="height:{visible[p.ID] ? '100%' : '0'}" class="details">
                <div class="keyval text-sm">
                    <span class="key md:w-1/12">DNS</span>
                    <span class="val md:w-1/3">{p.DNSName}</span>
                </div>

                <div class="keyval text-sm">
                    <span class="key md:w-1/12">OS</span>
                    <span class="val md:w-1/3">{p.OS}</span>
                </div>

                <div class="keyval text-sm">
                    <span class="key md:w-1/12">Relay</span>
                    <span class="val md:w-1/3 {p.RelayActive ? 'red' : 'gray'}">{p.Relay}</span>
                </div>

                <div class="keyval text-sm">
                    <span class="key md:w-1/12">Created</span>
                    <span class="val md:w-1/3">{new Date(p.Created).toLocaleDateString("en-US", dateShort)}</span>
                </div>

                <div class="keyval text-sm">
                    <span class="key md:w-1/12">Node Key</span>
                    <span class="val md:w-1/3">{p.NodeKey}</span>
                </div>

                <div class="keyval text-sm">
                    <span class="key md:w-1/12">Last seen</span>
                    <span class="val md:w-1/3">{new Date(p.LastSeen).toLocaleDateString("en-US", dateTime)}</span>
                </div>

                {#if p.PeerAPIPort > 0}
                <div class="keyval text-sm">
                    <span class="key md:w-1/12">API Port</span>
                    <span class="val md:w-1/3">{p.PeerAPIPort}</span>
                </div>
                {/if}
            </div> -->
            {/each}
            {/if}
        <!-- </tbody> -->
    </div>
</div>

<style lang="scss">
@import "../../../local.css";
.peers {
    font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
}
// .mono {
//     font-family: monospace;
// }
.table {
    display: grid;
    grid-template-columns: 20px 100px 200px 200px  200px 10%;
    column-gap: 50px;
    row-gap: 2%;
    max-width: 80%;
}

input, input:focus {
    outline: none;
}
// .details {
//     height: 0;
//     overflow: hidden;
//     padding-left: 32px;
// }
// .row {
//     &:hover {
//         cursor: pointer;
//     }
// }
.header {
    cursor: pointer;
    font-weight: 500;
    font-size: large;
    position: sticky;
    top: 0;
    text-transform: uppercase;
}
// .red {
//     color: red;
// }
// .gray {
//     color: gray;
// }
.ip-toggle {
    z-index: 100;
    position: relative;
    transition: color .3s;
    &:hover {
        color: rgb(112 110 109);
    }
}
</style>

<script lang="ts">
    import { 
        peerMap, 
        peersReady, 
        // peers 
    } from "../store/sse"

    import type { 
        Peer,
        Base,
        //  SSEMessage 
    } from "../types/types"

    import { 
        onMount,
        //  onDestroy 
    } from "svelte"

    import { FormatBytes } from "../js/lib"
    import dayjs from 'dayjs'
    import relativeTime from 'dayjs/plugin/relativeTime'
    dayjs.extend(relativeTime)

    let searchTerm:string = ""

    let visible = {}

    const toggleDetails = (id:number) => {
        visible[id] = !visible[id]
    }

    const IP = Object.freeze({
        v4: Symbol("v4"),
        v6: Symbol("v6")
    })

    let ip = {
        v: IP.v4,
        k: "IPv4Num",
        f: "v4",
    }

    const toggleIP = (e:any) => {
        console.log("toggling ip",e)
        if (ip.v === IP.v4) {
            ip = {
                v: IP.v6,
                k: "IPv6Num",
                f: "v6",
            }
        } else {
            ip = {
                v: IP.v4,
                k: "IPv4Num",
                f: "v4",
            }
        }
    }

    const filterKeys = ["HostName", "IPv4", "IPv6", "DNSName", "OS", "ID", "NodeKey"]

    const dateShort: Intl.DateTimeFormatOptions = { 
        // weekday: 'short', 
        year: 'numeric', 
        month: 'short', 
        day: 'numeric' 
    }
    const dateTime: Intl.DateTimeFormatOptions = { 
        year: 'numeric', 
        month: 'short', 
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    }

    let sorter: (a: Peer, b: Peer) => number

    $: peers = () => {
        // hash to array
        let p = [...$peerMap.entries()].map(x=>x[1]).filter(o=>{
            let hit = false
            for (let k of filterKeys) {
                if (o[k].toLowerCase().includes(searchTerm.toLowerCase())) hit = true
            }
            return hit
        })
        return p.sort(sorter)
    }

    let copy = (st: string) => {
        if (typeof window !== 'undefined') {
            navigator.clipboard.writeText(st).then(() => {
                console.log("copied ", st)
            }, err => {
                console.log("copy failed ", err)
            })
        }
    }

    let sortBy = {
        col: "HostName", 
        asc: false
    }

    const ago = (t: Date, u: boolean) => {
        if (!u) return dayjs(t).fromNow()
        return "–"
    }

    const bytes = (b: number) => {
        var base: Base = {
            suffix: "b",
            factor: 1024

        }
        return FormatBytes(b, base)
    }

    $: sort = (column: string) => {
       
        //     case "IPv4":
        //         break;
        
        //     default:
            console.log("Sorting by %s", column)
            if (sortBy.col == column) {
                sortBy.asc = !sortBy.asc
            } else {
                sortBy.col = column
                sortBy.asc = true
            }
            
            let sortModifier = (sortBy.asc) ? 1 : -1;
            
            sorter = (a: Peer, b: Peer) => {
                // console.log ("column",column) 
                let x:number | string, y:number | string
                if (typeof a[column] === "number") {
                    x = a[column]
                    y = b[column]
                } else {
                    x = a[column].toLowerCase()
                    y = b[column].toLowerCase()
                }

                return (x < y) 
                ? -1 * sortModifier 
                : (x > y) 
                ? 1 * sortModifier 
                : 0;
            }
            // break;
        // }

	}
    onMount(async() =>{
        sort("HostName")
    })

</script>