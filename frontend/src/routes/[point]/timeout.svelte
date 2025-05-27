<script>
    import { onMount } from 'svelte'

    let { endpoint, user } = $props()
    let state = $state({ remaining: 1 })
    let varTimeout = $state(15)

    async function loadtimeout () {
        const request = await fetch(`${endpoint}/timeouts/${user}`)
        return request.json()
    }

    async function setTimeout (duration) {
        await fetch(`${endpoint}/timeouts/${user}`, {
            method: 'DELETE'
        })
        state = loadtimeout()
    }

    async function delTimeout (duration) {
        await fetch(`${endpoint}/timeouts/${user}`, {
            method: "POST",
            body: JSON.stringify({ duration }),
        })
        state = loadtimeout()
    }

    onMount(() => {
        // state = loadtimeout()
    })
</script>

<div class="component timeout">
    {#await state}
    <div class="allgrid loading">... loading ...</div>
    {:then timeout}
    {#if timeout.remaining == -1}
    <button class="min" style="grid-area: preset1" onclick={_ => setTimeout(30)}>30</button>
    <button class="min" style="grid-area: preset2" onclick={_ => setTimeout(60)}>60</button>
    <button style="grid-area: minus" onclick={() => varTimeout -= 5}>-</button>
    <button class="min" style="grid-area: valset" onclick={_ => setTimeout({varTimeout})}>{varTimeout}</button>
    <button style="grid-area: plus" onclick={() => varTimeout += 5}>+</button>
    {:else}
    <button class="allgrid removeTimeout" onclick={delTimeout}>remove Timeout</button>
    {/if}
    {:catch err}
    {console.error(err)}
    <div class="allgrid error">Error loading remote ressource</div>
    {/await}
</div>

<style>
.timeout {
    display: grid;
    grid-template: repeat(2, 1fr) / repeat(6, 1fr);
    grid-template-areas:
        "preset1 preset1 preset1 preset2 preset2 preset2"
        "minus   minus   valset  valset   plus  plus";
    gap: 0;
}

.min::after {
    content: "min";
    font-size: 50%;
}

.allgrid {
    grid-area: 1 / 7 / 3 / 1;
}
</style>