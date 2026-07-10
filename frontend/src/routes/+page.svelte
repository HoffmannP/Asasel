<script>
    import { onMount } from 'svelte'
    import Loginstate from './loginstate.svelte'
    import Lockstate from './lockstate.svelte'
    import Timeout from './timeout.svelte'

    let user = $state('')
    const endpoint = '/api'

    async function loadClientConfig () {
        const request = await fetch('/api/config/client')
        const data = await request.json()
        if (data.account) {
            user = data.account
        }
    }

    onMount(() => {
        loadClientConfig()
    })
</script>

<div class="main">
    <Lockstate {endpoint} {user} />
    <Loginstate {endpoint} {user} />
    <Timeout {endpoint} {user} />
    <div class="component"><a href="/server">‹ Server</a></div>
</div>

<style>
.main {
    grid-template-rows: repeat(3, 1fr) 5rem;
}
</style>