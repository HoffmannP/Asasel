<script>
    import { enhance } from '$app/forms';

    export let data
    export let form

    $: commands = data?.commands
</script>

<div>
{#if form}
<div class="banner">{form}</div>
{/if}
{#each Object.entries(commands) as [name, parameters]}
    <form method="POST" use:enhance>
        {#each parameters as parameter, i}
        <label>{parameter} <input name={`param${i}`}></label>
        {/each}
        <button name="command" value="{name}">{name}</button>
    </form>
{/each}
</div>

