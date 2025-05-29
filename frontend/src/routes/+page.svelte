<script>
    import Loginstate from './loginstate.svelte'
    import Lockstate from './lockstate.svelte'
    import Timeout from './timeout.svelte'

    const USER = 'linus'
    const endpoint = '/api'
</script>

<div class="main">
    <Lockstate {endpoint} user={USER} />
    <Loginstate {endpoint} user={USER} />
    <Timeout {endpoint} user={USER} />
    <div class="component"><a href="/server">‹ Server</a></div>
</div>

<style>
.main {
    grid-template-rows: repeat(3, 1fr) 5rem;
}
</style>