<script>
    let { endpoint, user, status, onRefresh } = $props()

    async function lock (lockState) {
        await fetch(`${endpoint}/accounts/lock/${user}`, {
            method: "POST",
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ lockstate: lockState }),
        })
        await onRefresh()
    }
</script>

<div class="component">
    {#if status.lockstate}
    <button class="locked" onclick={_ => lock(false)}>Unlock</button>
    {:else}
    <button class="" onclick={_ => lock(true)}>Lock</button>
    {/if}
</div>

<style>
.component {
    min-height: clamp(9rem, 18vmin, 16rem);
}

button {
    position: relative;
    display: grid;
    place-items: center;
    overflow: hidden;
    line-height: 1;
    font-size: 0;
    min-width: 0;
    min-height: 0;
}

button::after {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: clamp(4rem, 18vmin, 10rem);
    line-height: 1;
    width: 100%;
    height: 100%;
    content: "🔓";
}

button.locked::after {
    content: "🔒";
}
</style>