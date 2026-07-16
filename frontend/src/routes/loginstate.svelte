<script>
    let { endpoint, user, status, onRefresh } = $props()

    async function kill () {
        await fetch(`${endpoint}/accounts/killall/${user}`, {
            method: 'POST'
        })
        await (new Promise(_ => setTimeout(_, 1000)))
        await onRefresh()
    }
</script>

<div class="component">
    {#if status.duration == -1}
    <div class="silent">no login</div>
    {:else}
    <button onclick={kill}>
        <div class="usertime">{status.duration}</div>
        <div>kill</div>
    </button>
    {/if}
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
    content: "Logged in for "
}

.usertime::after {
    content: "min"
}
</style>