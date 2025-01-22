<script lang="ts">
	import '../app.css';
	import { themeChange } from 'theme-change';
	import { browser } from '$app/environment';
	import { SvelteQueryDevtools } from '@tanstack/svelte-query-devtools';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';

	import { SvelteToast } from '@zerodevx/svelte-toast';
	import type { Snippet } from 'svelte';
	import Title from '$lib/components/layout/Title.svelte';

	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				enabled: browser
			}
		}
	});
	interface Props {
		children?: Snippet;
	}

	let { children }: Props = $props();
	$effect(() => {
		themeChange(false);
	});
</script>

<Title />
<QueryClientProvider client={queryClient}>
	<div class="bg-base-200">
		{@render children?.()}
	</div>
	<SvelteQueryDevtools />
</QueryClientProvider>
<SvelteToast />
