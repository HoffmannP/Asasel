<script>
    import { onMount } from 'svelte'

    let { endpoint, user } = $props()
    let state = $state({ duration: -1 })

    async function kill () {
        await fetch(`${endpoint}/accounts/kill/${user}`)
        state = logintime()
    }

    async function logintime () {
        const request = await fetch(`${endpoint}/accounts/time/${user}`)
        return request.json()
    }

    onMount(() => {
        state = logintime()
    })
</script>

<div class="component">
    {#await state}
    <div class="loading">... loading ...</div>
    {:then usertime}
    {#if usertime.duration == 1}
    <div class="silent">no login</div>
    {:else}
    <button onclick={kill}>
        <div class="usertime">{usertime.duration}</div>
        <div>Kill</div>
    </button>
    {/if}
    {:catch err}
    {console.error(err)}
    <div class="error">Error loading remote ressource</div>
    {/await}
</div>

<style>
button {
    display: flex;
    font-size: 200%;
}

.usertime {
    font-size: 35%;
}

.usertime::before {
    content: "Login since "
}

.usertime::after {
    content: "min"
}
</style>