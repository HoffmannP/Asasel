<script>
const { server, port } = $props()

function check() {
    return fetch(`http://${server}:${port}/`, {
        mode: 'no-cors'
    })
}

function local () {
    return (location.hostname == server) && (location.port == port)
}
</script>

<div class="component">
    {#await check()}
    <a href="http://{server}:{port}/" class=loading>{server}</a>
    {:then}
    <a href="http://{server}:{port}/" class:local={local()}>{server}</a>
    {:catch}
    <a href="/" disabled>{server}</a>
    {/await}
</div>

<style>
a.loading {
    color: #999;
}
a.local {
    color: blue;
}
</style>