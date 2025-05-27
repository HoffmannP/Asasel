<script>
    import Loginstate from './loginstate.svelte'
    import Lockstate from './lockstate.svelte'
    import Timeout from './timeout.svelte'

    const PORT = 8888
    const USER = 'linus'

    let { data } = $props()
    let endpoint = `http://${data.point}:${PORT}`
</script>

<div class="main">
    <Loginstate {endpoint} user={USER} />
    <Lockstate {endpoint} user={USER} />
    <Timeout {endpoint} user={USER} />
    <div><a href="/">back</a></div>
</div>

<style>
.main {
    padding: 1rem;
    width: calc(100vw - 2rem);
    height: calc(100vh - 2rem);
    background-color: lightgrey;
    display: grid;
    grid-template-rows: repeat(3, 1fr) 5rem;
    gap: 1rem;
}
</style>