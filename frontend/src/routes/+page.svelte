<script>
    import Loginstate from './loginstate.svelte'
    import Lockstate from './lockstate.svelte'
    import Timeout from './timeout.svelte'

    const USER = 'linus'
    // const PORT = 8888
    // let { data } = $props()
    // let endpoint = `http://{data.point}:${PORT}/api`
    const endpoint = '/api'

    function swipeLeft(event) {
        direction = alert("Swipe")
    }
</script>

<div class="main" use:swipe={{ timeframe: 300, minSwipeDistance: 100, touchAction: 'pan-y' }} on:swipe={swipeLeft}>
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