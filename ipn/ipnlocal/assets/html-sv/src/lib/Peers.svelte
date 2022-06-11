
<div class="peers">
    <div class="py-8 text-3xl font-semibold tracking-tight leading-tight">Peers</div>
    <table class="tb">
        <thead class="stick opaque py-2">
            <tr class="w-full md:text-base">
                <th on:click={() =>sort("ID")} class="pointer w-8 pr-3 flex-auto md:flex-initial md:shrink-0 w-0 ">ID</th>
                <th on:click={() =>sort("HostName")} class="pointer md:w-1/8 flex-auto md:flex-initial md:shrink-0 w-0 text-ellipsis">machine</th>
                <th on:click={() =>sort("IPv4")} class="pointer hidden md:block md:w-1/8">IP</th>
                <th on:click={() =>sort("OS")} class="pointer hidden md:block md:w-1/12">OS</th>
                <th on:click={() =>sort("LastSeen")} class="pointer hidden md:block md:w-1/12">Last Seen</th>
                <th class="hidden md:block md:w-1/12">Relay</th>
                <th on:click={() =>sort("DNSName")} class="pointer hidden md:block md:w-1/8">DNS</th>
                <th on:click={() =>sort("RX")} class="pointer hidden md:block md:w-1/12 text-right">rx</th>
                <th on:click={() =>sort("TX")} class="pointer hidden md:block md:w-1/12 text-right">tx</th>
                <th on:click={() =>sort("Created")} class="pointer hidden md:block md:w-1/12 text-right">Created</th>
            </tr>
        </thead>

        <tbody class="table-body">
            {#each data as p}
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
                <td class="md:w-1/8 flex-auto md:flex-initial md:shrink-0 w-0 text-ellipsis">
                    <div class="relative">
                        <div class="items-center text-gray-900">
                            <h3 class="font-semibold hover:text-blue-500">
                                {p.HostName}
                            </h3>
                        </div>
                    </div>
                    
                </td>
                <td class="hidden md:block md:w-1/8">
                    <ul>
                    <li class="pr-6">
                        <div class="flex relative min-w-0">
                            <div class="truncate">
                                <span>{p.IPv4}</span>
                            </div>
                        </div>
                    </li>
                    <li class="pr-6">
                        <div class="flex relative min-w-0">
                            <div class="truncate">
                                <span>{p.IPv6}</span>
                            </div>
                        </div>
                    </li>
                    </ul>
                </td>
                <td class="hidden md:block md:w-1/12 ">{p.OS}</td>
                <td class="hidden md:block md:w-1/12" title="{new Date(p.LastSeen).toLocaleDateString("en-US", options)}">{ago(p.LastSeen, p.Unseen)}</td>
                <td class="hidden md:block md:w-1/12">{p.Connection}</td>
                <td class="hidden md:block md:w-1/8 truncate">{p.DNSName}</td>
                <td class="hidden md:block md:w-1/12 text-right">{bytes(p.RXb)}</td>
                <td class="hidden md:block md:w-1/12 text-right">{bytes(p.TXb)}</td>
                <td class="hidden md:block md:w-1/12 text-right"><span>
                    <div>{p.CreatedDate}</div>
                    <div>{p.CreatedTime}</div>
                </span></td>
            </tr>
            {/each}
        </tbody>
    </table>
</div>

<style>
@import "../../../local.css";

.pointer {
    cursor: pointer;
}
</style>

<script lang="ts">
    import type { Peer, Base } from "../types/types"
    import { onMount } from "svelte"
    import dayjs from 'dayjs'
    import relativeTime from 'dayjs/plugin/relativeTime'
    dayjs.extend(relativeTime)
    import { FormatBytes } from "../js/lib"


    export let data: Peer[] = []

    const options: Intl.DateTimeFormatOptions = { 
        weekday: 'short', 
        year: 'numeric', 
        month: 'short', 
        day: 'numeric' 
    }

    let sortBy = {
        col: "HostName", 
        asc: true
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
		
		if (sortBy.col == column) {
			sortBy.asc = !sortBy.asc
		} else {
			sortBy.col = column
			sortBy.asc = true
		}
		
		let sortModifier = (sortBy.asc) ? 1 : -1;
		
		let sorter = (a: Peer, b: Peer) => {
            let x = a[column].toLowerCase()
            let y = b[column].toLowerCase()
			return (x < y) 
			? -1 * sortModifier 
			: (x > y) 
			? 1 * sortModifier 
			: 0;
        }
		
		data = data.sort(sorter)
	}

    onMount( () => {
        sort("ID")

        const sse = new EventSource(`http://100.100.100.100/events/`)
        console.log("EventSource", sse)
        sse.onmessage = event => {
            let response = JSON.parse(event.data)
            if(!response.length) {
                console.log("sse response", response)
            }
        }
        
        // return () => {
        //     if(sse.readyState === 1) {
        //         console.log("sse closing")
        //         sse.close()
        //     }
        // })
    })
</script>