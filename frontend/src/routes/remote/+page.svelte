<script>
    import { onMount } from 'svelte'
    import Loginstate from '../loginstate.svelte'
    import Lockstate from '../lockstate.svelte'
    import Timeout from '../timeout.svelte'

    let user = $state('linus')
    let server = $state('')
    let endpoint = $state('/api')

    async function loadClientConfig () {
        const request = await fetch('/api/config/client')
        const data = await request.json()
        if (data.account) {
            user = data.account
        }
    }

    onMount(() => {
        const params = new URLSearchParams(location.search)
        server = params.get('server') || ''
        if (server !== '') {
            endpoint = `/api/remote/${encodeURIComponent(server)}`
        }
        loadClientConfig()
    })
</script>

<div class="main">
    {#if server === ''}
    <div class="component silent">No remote server selected</div>
    {:else}
    <Lockstate {endpoint} {user} />
    <Loginstate {endpoint} {user} />
    <Timeout {endpoint} {user} />
    {/if}
    <div class="component"><a href="/server">‹ Server</a></div>
</div>

<style>
.main {
    grid-template-rows: repeat(3, 1fr) 5rem;
}
</style>
