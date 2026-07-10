<script>
    import { onMount } from 'svelte'
    import Server from './servers.svelte'

    let servers = $state([])

    async function loadServers () {
        const request = await fetch('/api/config/servers')
        const data = await request.json()
        return data.server || []
    }

    onMount(() => {
        loadServers().then(result => { servers = result })
    })

</script>

<div class="main">
    {#each servers as server}
    <Server {server} />
    {/each}
</div>

<style>
.main {
    grid-template-rows: repeat(4, 1fr);
}
</style>
