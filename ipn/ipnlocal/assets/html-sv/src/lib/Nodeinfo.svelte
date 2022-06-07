
    <div class="width-80">
        <div class="flex py-8">
            <img src={logo} alt="Svelte Logo" />

            <div class="flex items-center space-x-2 ">
                <div class="text-right truncate leading-4">
                    {#if data.Profile}
                    <div class="font-semibold truncate leading-normal">{data.Profile.LoginName}</div>
                    {/if}
                </div>
            </div>
        </div>
        <div class="flex ">
            <div class="border border-gray-200 bg-gray-0 rounded-lg flex items-center justify-between">
                <div class="flex items-center min-width-0">
                    <img src={device} alt="Svelte Logo" />
                    <div class="truncate mr-2">
                        <div class="">
                            <span class="font-semibold">{data.HostName}</span> 
                            <span class="text-sm">{data.StableID}</span><span> – </span>
                            <span class="text-sm">{data.NodeKey}</span>
                        </div>
                        <div class="text-sm">
                            <span>Created </span><span>{new Date(data.Created).toLocaleString("en-US", options)}</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="info py-8 text-sm">
            <div class="status keyval">
                <span class="key md:w-1/8">STATUS</span>
                <span class="val md:w-1/3">
                    {#if data.Health == null}
                    <div class="green">OK</div>
                    {:else}
                        {#each data.Health as h}
                        <div>{h}</div>
                        {/each}
                    {/if}
                </span>
            </div>

            <div class="ipv4 keyval">
                <span class="key md:w-1/8">IPv4</span>
                <span class="val md:w-1/3 font-semibold">{data.IPv4}</span>
            </div>
            <div class="ipv6 keyval">
                <span class="key md:w-1/8">IPv6</span>
                <span class="val md:w-1/3 font-semibold">{data.IPv6}</span>
            </div>

            <div class="backend keyval">
                <span class="key md:w-1/8">Server URL</span>
                <span class="val md:w-1/3">{data.ServerURL}</span>
            </div>


            <div class="os keyval">
                <span class="key md:w-1/8">OS</span>
                <span class="val md:w-1/3">{data.OS}</span>
            </div>

            <div class="version keyval">
                <span class="key md:w-1/8">Version</span> 
                <span class="val md:w-1/3">{data.Version}</span>
            </div>

            <div class="dns keyval">
                <span class="key md:w-1/8">DNS</span> 
                <span class="val md:w-1/3">{data.Name}</span>
            </div>

            <div class="services keyval">
                <span class="key md:w-1/8">Services</span> 
                <span class="val md:w-1/3">
                    {#each data.Services as s}
                    <div class=" flex">
                        <span class="font-semibold md:w-1/3">{s.Description}</span> 
                        <span class="md:w-1/4">{s.Proto} <span>{s.Port}</span></span> 
                    </div>
                    {/each} 
                </span>
            </div>

            <Peers data={data.Peers} />

        </div>

    </div>


<script>
    import { onMount } from "svelte"
    import logo from '../assets/logo.svg'
    import device from '../assets/device.svg'
    import Peers from './Peers.svelte'

    const endpoint = "http://100.100.100.100/json/"
    const options = { weekday: 'short', year: 'numeric', month: 'short', day: 'numeric' }

    let data = {
        Peers: [],
        Services: [],
    }

    onMount(async () => {
        const response = await fetch(endpoint)
        data = await response.json()
        console.log(data)
    })

    // const compReg = (a, b) => {
    //     a = a.toLowerCase()
    //     b = b.toLowerCase()
    //     if (a > b) {
    //         return 1;
    //     } else if (b > a) {
    //         return -1;
    //     } else {
    //         return 0;
    //     }
    // }

    // const newGrid = d => {
    //     const options = { weekday: 'short', year: 'numeric', month: 'short', day: 'numeric' }
    //     new Grid({
    //         columns: [
    //             {
    //                 id: 'HostName',
    //                 name: 'Hostname',
    //                 sort: {
    //                     compare: compReg
    //                 }
    //             }, 
    //             {
    //                 id: 'IPs',
    //                 name: 'IPs',
    //                 // formatter: cell => `${cell.join(", ")}`
    //                 formatter: html(cell => cell.forEach(c => `<div>${c}</div>`))
    //             }, 
    //             {
    //                 id: 'OS',
    //                 name: 'OS'
    //             },
    //             {
    //                 id: 'LastSeen',
    //                 name: 'Last Seen',
    //                 formatter: cell => new Date(cell).toLocaleDateString("en-US", options)
    //             }, 
    //             {
    //                 id: 'DNSName',
    //                 name: 'DNS'
    //             }, 
    //             {
    //                 id: 'RX',
    //                 name: 'RX'
    //             },
    //             {
    //                 id: 'TX',
    //                 name: 'TX'
    //             },
    //         ],
    //         sort: true,
    //         // search: true,
    //         data: d
    //     }).render(document.getElementById("wrapper"));
    // }

</script>

<style>
@import "../../../local.css";
</style>