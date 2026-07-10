<script>
const { server } = $props()
</script>

<div class="component">
    <a href="/remote?server={encodeURIComponent(server)}">{server}</a>
</div>

<style>
a {
    text-transform: none;
}
</style>