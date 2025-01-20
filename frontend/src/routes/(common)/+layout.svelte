<script lang="ts">
	import ThemeSwap from '$lib/components/ui/ThemeSwap.svelte';
	interface Props {
		children?: import('svelte').Snippet;
	}
	let { children }: Props = $props();
</script>

<div class="flex h-screen flex-col">
	<div class="flex w-full justify-end">
		<div class="flex p-2">
			<ThemeSwap />
		</div>
	</div>
	{@render children?.()}
</div>
