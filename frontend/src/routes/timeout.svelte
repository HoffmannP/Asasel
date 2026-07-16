<script>
    let { endpoint, user, status, onRefresh } = $props()
    let varTimeout = $state(15)

    async function setTimeout (duration) {
        await fetch(`${endpoint}/timeouts/${user}`, {
            method: "POST",
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ duration }),
        })
        await onRefresh()
    }

    async function delTimeout () {
        await fetch(`${endpoint}/timeouts/${user}`, {
            method: 'DELETE'
        })
        await onRefresh()
    }
</script>

<div class="component timeout">
    {#if status.remaining == -1}
    <button class="min" style="grid-area: preset1" onclick={_ => setTimeout(30)}>30</button>
    <button class="min" style="grid-area: preset2" onclick={_ => setTimeout(60)}>60</button>
    <button disabled={ varTimeout <= 5 } style="grid-area: minus" onclick={() => varTimeout -= 5}>−</button>
    <button class="min" style="grid-area: valset" onclick={_ => setTimeout(varTimeout)}>{varTimeout}</button>
    <button style="grid-area: plus" onclick={() => varTimeout += 5}>+</button>
    {:else}
    <button class="allgrid removeTimeout" onclick={delTimeout}>
        <div class="remaining">{status.remaining}</div>
        <div>remove</div>
    </button>
    {/if}
</div>

<style>
.timeout {
    display: grid;
    grid-template: repeat(2, 1fr) / repeat(4, 1fr);
    grid-template-areas:
        "preset1 preset1 preset2 preset2"
        "minus   valset  valset   plus  ";
    gap: 0;
}

.min::after {
    content: "min";
    font-size: 50%;
}

.allgrid {
    grid-area: 1 / 5 / 3 / 1;
    font-size: 200%;
}

.remaining {
    font-size: 35%;
}

.remaining::before {
    content: "Timeout in "
}

.remaining::after {
    content: "min"
}
</style>