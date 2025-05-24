<script>
    let { data } = $props()

    let commands = [
        "lock" ,
        "unlock",
        "time"
    ]

    function call(command) {
        return () => fetch(`/api/${data.point}/${command}`, { method: "POST" })
    }
</script>

<div class="main" style="--rows: {commands.length}">
    {#each commands as command}
    <button onclick={call(command)}>{command}</button>
    {/each}
    <a href="/">back</a>
</div>
