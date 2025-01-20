<script lang="ts">
	type Props = {
		title: string;
		decription?: string;
	};
	const { title, description = '' } = $props();
</script>

<svelte:head>
	<title>{title}</title>
	<meta name="description" content={description} />
</svelte:head>
