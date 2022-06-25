
    <div class="width-80">
        {#if $localReady}
        <div class="flex py-8">

            <img src={logo} alt="TS Logo" class="w-8 pr-3"/>

            <div class="flex items-center space-x-2 ">
                <div class="text-right truncate leading-4">
                    {#if $local.Profile}
                    <div class="font-semibold truncate leading-normal">{$local.Profile.LoginName}</div>
                    {/if}
                </div>
            </div>
        </div>
        <div class="flex ">
            <div class="border border-gray-200 bg-gray-0 rounded-lg flex items-center justify-between">
                <div class="flex items-center min-width-0">

                    <img src={device} alt="Device Logo" class="w-8 pr-3"/>

                    <div class="truncate mr-2">
                        <div class="">
                            <span class="font-semibold">{$local.HostName}</span> 
                            <span class="text-sm">{$local.StableID}</span><span> – </span>
                            <span class="text-sm">{$local.NodeKey}</span>
                        </div>
                        <div class="text-sm">
                            <span>Created </span><span>{new Date($local.Created).toLocaleString("en-US", options)}</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="info py-8 text-sm">
            <div class="status keyval">
                <span class="key md:w-1/8">STATUS</span>
                <span class="val md:w-1/3">
                    {#if $local.Health == null}
                    <div class="green">OK</div>
                    {:else}
                        {#each $local.Health as h}
                        <div>{h}</div>
                        {/each}
                    {/if}
                </span>
            </div>

            <div class="ipv4 keyval">
                <span class="key md:w-1/8">IPv4</span>
                <span class="val md:w-1/3 font-semibold">{$local.IPv4}</span>
            </div>
            <div class="ipv6 keyval">
                <span class="key md:w-1/8">IPv6</span>
                <span class="val md:w-1/3 font-semibold">{$local.IPv6}</span>
            </div>
            <div class="ipv6 keyval">
                <span class="key md:w-1/8">API Port</span>
                <span class="val md:w-1/3 font-semibold">{$local.PeerAPIPort}</span>
            </div>

            <div class="backend keyval">
                <span class="key md:w-1/8">Server URL</span>
                <span class="val md:w-1/3">{$local.ServerURL}</span>
            </div>


            <div class="os keyval">
                <span class="key md:w-1/8">OS</span>
                <span class="val md:w-1/3">{$local.OS}</span>
            </div>

            <div class="version keyval">
                <span class="key md:w-1/8">Version</span> 
                <span class="val md:w-1/3">{$local.Version}</span>
            </div>

            <div class="dns keyval">
                <span class="key md:w-1/8">DNS</span> 
                <span class="val md:w-1/3">{$local.Name}</span>
            </div>

            <div class="services keyval">
                <span class="key md:w-1/8">Services</span> 
                <span class="val md:w-1/3">
                    {#if $local.Services}
                    {#each $local.Services as s}
                    <div class=" flex">
                        <span class="font-semibold md:w-1/3">{s.Description}</span> 
                        <span class="md:w-1/4">{s.Proto} <span>{s.Port}</span></span> 
                    </div>
                    {/each} 
                    {/if}
                </span>
            </div>
        </div>
        {/if}
    </div>

<script lang="ts">
    import { local, localReady } from "../store/sse"
    // import { onMount } from "svelte"
    import logo from '../assets/logo.svg'
    import device from '../assets/device.svg'


    const options: Intl.DateTimeFormatOptions = { 
        weekday: 'short', 
        year: 'numeric', 
        month: 'short', 
        day: 'numeric' 
    }
</script>

<style>
@import "../../../local.css";
</style>