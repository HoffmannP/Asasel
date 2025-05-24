<script>
    import slugify from 'slugify'

    let points = [
        "nura",
        "yoga"
    ]
</script>

<div class="main" style="--rows: {points.length}">
    {#each points as point}
    <a href="/{slugify(point)}">{point}</a>
    {/each}
</div>
