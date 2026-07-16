<script>
    import { onMount } from 'svelte'
    import { STATUS_POLL_INTERVAL_MS } from '$lib/constants'
    import Loginstate from './loginstate.svelte'
    import Lockstate from './lockstate.svelte'
    import Timeout from './timeout.svelte'

    let user = $state('')
    const endpoint = '/api'
    let mode = $state('local')
    let controllerUrl = $state('')
    let status = $state({ lockstate: false, duration: -1, remaining: -1 })
    let statusLoading = $state(true)
    let statusError = $state(false)
    let inFlight = false

    async function loadClientConfig () {
        const request = await fetch('/api/config/client')
        const data = await request.json()
        if (data.mode) {
            mode = data.mode
        }
        if (data.controllerUrl) {
            controllerUrl = data.controllerUrl.replace(/\/$/, '')
        }
        if (data.account) {
            user = data.account
        }
    }

    async function loadStatus () {
        if (!user) {
            return
        }
        const request = await fetch(`${endpoint}/accounts/state/${user}`)
        const data = await request.json()
        status = {
            lockstate: !!data.lockstate,
            duration: Number.isInteger(data.duration) ? data.duration : -1,
            remaining: Number.isInteger(data.remaining) ? data.remaining : -1,
        }
    }

    async function updateStatus () {
        if (inFlight || !user) {
            return
        }
        inFlight = true
        try {
            await loadStatus()
            statusError = false
        } catch (err) {
            console.error(err)
            statusError = true
        } finally {
            statusLoading = false
            inFlight = false
        }
    }

    onMount(() => {
        let timer
        loadClientConfig().then(async () => {
            await updateStatus()
            timer = setInterval(updateStatus, STATUS_POLL_INTERVAL_MS)
        })
        return () => {
            if (timer) {
                clearInterval(timer)
            }
        }
    })
</script>

<div class="main">
    {#if user === ''}
    <div class="component silent">loading account...</div>
    {:else if statusLoading}
    <div class="component loading">... loading ...</div>
    {:else if statusError}
    <div class="component error">Error loading remote ressource</div>
    {:else}
    <Lockstate {endpoint} {user} {status} onRefresh={updateStatus} />
    <Loginstate {endpoint} {user} {status} onRefresh={updateStatus} />
    <Timeout {endpoint} {user} {status} onRefresh={updateStatus} />
    {/if}
    {#if mode === 'control'}
    <div class="component footer-nav"><a class="footer-link" href="/server">‹ Server</a></div>
    {:else if controllerUrl !== ''}
    <div class="component footer-nav"><a class="footer-link" href={controllerUrl}>‹ Controller GUI</a></div>
    {:else}
    <div class="component footer-nav"><a class="footer-link" href="/">‹ Controller GUI</a></div>
    {/if}
</div>

<style>
.main {
    grid-template-rows: repeat(3, 1fr) 5rem;
}

.footer-nav a {
    font-size: clamp(1.2rem, 2.6vmin, 2rem);
    line-height: 1.1;
}
</style>