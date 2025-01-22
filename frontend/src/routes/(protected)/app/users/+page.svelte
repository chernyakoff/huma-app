<script lang="ts">
	import Title from "$lib/components/layout/Title.svelte";
	import { authState } from "$lib/states.svelte";
</script>

{#if authState.email}
	<h1>{authState.email}</h1>
{/if}
