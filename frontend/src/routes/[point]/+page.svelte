<script>
    const PORT = 8888
    const USER = 'linus'

    let { data } = $props()
    let endpoint = `http://${data.point}:${PORT}`
    let varTimeout = $state(15)

    function kill () {
        fetch(`${endpoint}/accounts/kill/${USER}`)
    }

    function lock (lock) {
        fetch(`${endpoint}/accounts/lock/${USER}`, {
            method: "POST",
            body: JSON.stringify({ lockstate: lock }),
        })
    }

    async function logintime () {
        const request = await fetch(`${endpoint}/accounts/time/${USER}`)
        return request.json()
    }

    async function lockstate () {
        const request = await fetch(`${endpoint}/accounts/lock/${USER}`)
        return request.json()
    }

    async function loadtimeout () {
        const request = await fetch(`${endpoint}/timeouts/${USER}`)
        return request.json()
    }

    function decTimeout () {
        varTimeout -= 5
    }
    function incTimeout () {
        varTimeout += 5
    }

    function setTimeout (duration) {
        fetch(`${endpoint}/timeouts/${USER}`, { method: 'DELETE' })
    }

    function delTimeout (duration) {
        fetch(`${endpoint}/timeouts/${USER}`, {
            method: "POST",
            body: JSON.stringify({ duration }),
        })
    }
</script>

<div class="main">
    {#await logintime()}
    <div class="loading">... loading ...</div>
    {:then usertime}
    {#if usertime.Duration == -1}
    <div class="usertime">Not logged in</div>
    {:else}
    <div class="usertime">{usertime.Duration}</div>
    <button onclick={kill}>Kill</button>
    {/if}
    {/await}

    {#await lockstate()}
    <div class="loading">... loading ...</div>
    {:then userlock}
    {#if userlock.lockstate}
    <button onclick={_ => lock(false)}>Unlock</button>
    {:else}
    <button onclick={_ => lock(true)}>Lock</button>
    {/if}
    {/await}

    {#await loadtimeout()}
    <div class="loading">... loading ...</div>
    {:then timeout}
    {#if timeout.remaining == -1}
    <button onclick={_ => setTimeout(30)}>Timeout 30min</button>
    <button onclick={_ => setTimeout(60)}>Timeout 60min</button>
    <div>
        <button onclick={decTimeout}>-</button>
        <button onclick={_ => setTimeout({varTimeout})}>Timeout {varTimeout}min</button>
        <button onclick={incTimeout}>+</button>
    </div>
    {:else}
    <button onclick={delTimeout}>Del Timeout</button>
    {/if}
    {/await}


    <div><a href="/">back</a></div>
</div>
