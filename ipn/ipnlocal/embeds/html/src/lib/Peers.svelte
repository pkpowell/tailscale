
<div class="peers">
    <div class="py-8 text-3xl font-semibold tracking-tight leading-tight">Peers</div>
    <input
        placeholder="Search..."
        bind:value ={searchTerm}
    >
    <table class="tb">
        <thead class="stick opaque py-2">
            <tr class="w-full md:text-base">
                <th on:click={() =>sort("ID")} class="pointer w-8 pr-3 flex-auto md:flex-initial md:shrink-0 w-0 ">ID</th>
                <th on:click={() =>sort("HostName")} class="pointer md:w-1/12 flex-auto md:flex-initial md:shrink-0 w-0 text-ellipsis">machine</th>
                <th on:click={() =>sort(ip.k)} class="pointer hidden md:block md:w-1/12">IP<button class="ip-toggle" on:click|stopPropagation={toggleIP}>{ip.f}</button></th>
                <!-- <th on:click={() =>sort("OS")} class="pointer hidden md:block md:w-1/12">OS</th> -->
                <th on:click={() =>sort("LastSeen")} class="pointer hidden md:block md:w-1/12">Last Seen</th>
                <!-- <th class="hidden md:block md:w-1/12">Relay</th> -->
                <!-- <th on:click={() =>sort("DNSName")} class="pointer hidden md:block md:w-1/8">DNS</th> -->
                <th on:click={() =>sort("RXb")} class="pointer hidden md:block md:w-1/12 text-right">rx</th>
                <th on:click={() =>sort("TXb")} class="pointer hidden md:block md:w-1/12 text-right">tx</th>
                <!-- <th on:click={() =>sort("Created")} class="pointer hidden md:block md:w-1/12 text-right">Created</th> -->
            </tr>
        </thead>

        <tbody class="table-body">

            {#if $peersReady}
            {#each peers() as p}
            <!-- {#each [...$peerMap.entries()] as [key, p], index (key)} -->

            <tr class="table-row w-full px-0.5 hover:bg-gray-0">
                <td class="w-8 pr-3 flex-auto md:flex-initial md:shrink-0 w-0 ">
                    <div class="relative">
                        <div class="flex items-center text-gray-600 text-sm">
                            <span>
                                {p.ID}
                            </span>
                        </div>
                    </div>
                    
                </td>
                <td class="md:w-1/12 flex-auto md:flex-initial md:shrink-0 w-0 text-ellipsis">
                    <div class="relative">
                        <div class="items-center text-gray-900">
                            <h3 class="font-semibold hover:text-blue-500">
                                {p.HostName}
                            </h3>
                        </div>
                    </div>
                    
                </td>
                <td class="hidden md:block md:w-1/12">
                    <ul>
                        {#if ip.v === IP.v4}
                        <li class="pr-6">
                            <div class="truncate pr-6">
                                <div class="">{p.IPv4}</div>
                            </div>
                        </li>
                        {:else}
                        <li class="pr-6">
                            <div class="truncate pr-6">
                                <span class="">{p.IPv6}</span>
                            </div>
                        </li>
                        {/if}
                    </ul>
                </td>
                <!-- <td class="hidden md:block md:w-1/12">{p.OS}</td> -->
                <td class="hidden md:block md:w-1/12" title="{new Date(p.LastSeen).toLocaleDateString("en-US", options)}">{ago(p.LastSeen, p.Unseen)}</td>
                <!-- <td class="hidden md:block md:w-1/12">{p.Connection}</td> -->
                <!-- <td on:click={()=>copy(p.DNSName)} class="hidden md:block md:w-1/8 truncate">{p.DNSName}</td> -->
                <td class="hidden md:block md:w-1/12 text-right">{bytes(p.RXb)}</td>
                <td class="hidden md:block md:w-1/12 text-right">{bytes(p.TXb)}</td>
            </tr>
            <div>more stuff...</div>
            {/each}
            {/if}
        </tbody>
    </table>
</div>

<style lang="scss">
@import "../../../local.css";

.pointer {
    cursor: pointer;
}

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

    let searchTerm: string = ""

    const IP = Object.freeze({
        v4: Symbol("v4"),
        v6: Symbol("v6")
    })

    let ip = {
        v: IP.v4,
        k: "IPv4Num",
        f: "v4",
    }

    const toggleIP=e=> {
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

    const filterKeys = ["HostName", "IPv4", "IPv6", "DNSName", "OS", "ID"]

    const options: Intl.DateTimeFormatOptions = { 
        // weekday: 'short', 
        year: 'numeric', 
        month: 'short', 
        day: 'numeric' 
    }

    let sorter: (a: Peer, b: Peer) => number

    $: peers = () => {
        // hash to array
        let p = [...$peerMap.entries()].map(x=>x[1]).filter(o=>{
            let hit = false
            for (let k of filterKeys) {
                if (o[k].includes(searchTerm)) hit = true
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