<script>
    import { onMount } from 'svelte'

    let { endpoint, user } = $props()
    let state = $state({ lockstate: false })

    async function lock (lock) {
        await fetch(`${endpoint}/accounts/lock/${user}`, {
            method: "POST",
            body: JSON.stringify({ lockstate: lock }),
        })
        state = lockstate()
    }

    async function lockstate () {
        const request = await fetch(`${endpoint}/accounts/lock/${user}`)
        return request.json()
    }

    onMount(() => {
        // state = lockstate()
    })
</script>

<div class="component">
    {#await state}
    <div class="loading">... loading ...</div>
    {:then userlock}
    {#if userlock.lockstate}
    <button class="locked" onclick={_ => lock(false)}>Unlock</button>
    {:else}
    <button class="" onclick={_ => lock(true)}>Lock</button>
    {/if}
    {:catch err}
    {console.error(err)}
    <div class="error">Error loading remote ressource</div>
    {/await}
</div>

<style>
button {
    font-size: 0;
}

button::after {
    margin-top: -12vmin;
    font-size: calc(30vmin);
    content: "🔓";
}

button.locked::after {
    content: "🔒";
}
</style>