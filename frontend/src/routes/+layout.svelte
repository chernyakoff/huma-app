<script lang="ts">
	import { browser } from "$app/environment";
	import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
	import { SvelteQueryDevtools } from "@tanstack/svelte-query-devtools";

	import "../app.css";

	import Title from "$lib/components/layout/Title.svelte";
	import { SvelteToast } from "@zerodevx/svelte-toast";
	import type { Snippet } from "svelte";

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
</script>

<Title />
<QueryClientProvider client={queryClient}>
	<div class="bg-base-200">
		{@render children?.()}
	</div>
	<SvelteQueryDevtools />
</QueryClientProvider>
<SvelteToast />
