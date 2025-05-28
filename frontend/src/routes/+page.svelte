<script>
    import slugify from 'slugify'

    let points = [
        "nura",
        "yoga"
    ]
</script>

<div class="main">
    {#each points as point}
    <a href="/{slugify(point)}">{point}</a>
    {/each}
</div>


<style>
.main {
    padding: 1rem;
    width: calc(100vw - 2.67rem);
    height: calc(100vh - 2.67rem);
    background-color: lightgrey;
    display: grid;
    grid-template-rows: repeat(2, 1fr);
    gap: 1.67rem;

}
</style>