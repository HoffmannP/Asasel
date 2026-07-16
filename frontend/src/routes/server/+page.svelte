<script>
    import { onMount } from 'svelte'
    import Server from './servers.svelte'

    let servers = $state([])
    let mode = $state('local')
    let controllerUrl = $state('')

    async function loadServers () {
        const request = await fetch('/api/config/servers')
        const data = await request.json()
        return data.server || []
    }

    async function loadClientConfig () {
        const request = await fetch('/api/config/client')
        const data = await request.json()
        if (data.mode) {
            mode = data.mode
        }
        if (data.controllerUrl) {
            controllerUrl = data.controllerUrl.replace(/\/$/, '')
        }
    }

    onMount(() => {
        loadClientConfig().then(() => {
            if (mode === 'control') {
                loadServers().then(result => { servers = result })
            }
        })
    })

</script>

<div class="main">
    {#if mode === 'control'}
        {#each servers as server}
        <Server {server} />
        {/each}
    {:else if controllerUrl !== ''}
        <div class="component footer-nav"><a class="footer-link" href={controllerUrl}>‹ Controller GUI</a></div>
    {:else}
        <div class="component footer-nav"><a class="footer-link" href="/">‹ Controller GUI</a></div>
    {/if}
</div>

<style>
.main {
    grid-template-rows: repeat(4, 1fr);
}

.footer-nav a {
    font-size: clamp(1.2rem, 2.6vmin, 2rem);
    line-height: 1.1;
}
</style>
