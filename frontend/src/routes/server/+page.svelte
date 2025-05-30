<script>
    import Server from './servers.svelte'

    const port = 2727
    const servers = [
        'yoga',
        'nura',
        'hackmack',
        'bereisen'
    ]

</script>

<div class="main">
    {#each servers as server}
    <Server {server} {port} />
    {/each}
</div>

<style>
.main {
    grid-template-rows: repeat(4, 1fr);
}
</style>