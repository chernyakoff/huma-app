<script lang="ts">
	import { page } from "$app/state";
	import NotFound from "$lib/components/features/errors/NotFound.svelte";
	import ServerError from "$lib/components/features/errors/ServerError.svelte";

	const pages = {
		404: NotFound,
		500: ServerError
	} as const;

	type ErrorCode = keyof typeof pages;

	const status = +page.status;
	const index = Object.keys(pages)
		.map((x) => +x)
		.reduce((p, c) => (p < status ? c : p)) as ErrorCode;
	const component = pages[index];
</script>

<div class="flex h-screen flex-col">
	<div class="hero flex flex-grow items-center justify-center">
		<svelte:component this={component} />
	</div>
</div>
