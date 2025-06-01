<script>
    import { onMount } from 'svelte'

    let { endpoint, user } = $props()
    let state = $state({ remaining: -1 })
    let varTimeout = $state(15)

    async function loadtimeout () {
        const request = await fetch(`${endpoint}/timeouts/${user}`)
        return request.json()
    }

    async function setTimeout (duration) {
        await fetch(`${endpoint}/timeouts/${user}`, {
            method: "POST",
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ duration }),
        })
        state = loadtimeout()
    }

    async function delTimeout (duration) {
        await fetch(`${endpoint}/timeouts/${user}`, {
            method: 'DELETE'
        })
        state = loadtimeout()
    }

    async function update () {
        loadtimeout().then(result => { state = result })
    }

    onMount(() => {
        state = loadtimeout()
        setInterval(update, 5000)
    })
</script>

<div class="component timeout">
    {#await state}
    <div class="allgrid loading">... loading ...</div>
    {:then timeout}
    {#if timeout.remaining == -1}
    <button class="min" style="grid-area: preset1" onclick={_ => setTimeout(30)}>30</button>
    <button class="min" style="grid-area: preset2" onclick={_ => setTimeout(60)}>60</button>
    <button disabled={ varTimeout <= 5 } style="grid-area: minus" onclick={() => varTimeout -= 5}>−</button>
    <button class="min" style="grid-area: valset" onclick={_ => setTimeout(varTimeout)}>{varTimeout}</button>
    <button style="grid-area: plus" onclick={() => varTimeout += 5}>+</button>
    {:else}
    <button class="allgrid removeTimeout" onclick={delTimeout}>
        <div class="remaining">{timeout.remaining}</div>
        <div>remove</div>
    </button>
    {/if}
    {:catch err}
    {console.error(err)}
    <div class="allgrid error">Error loading remote ressource</div>
    {/await}
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