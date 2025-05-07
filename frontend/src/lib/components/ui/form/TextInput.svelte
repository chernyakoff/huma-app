<script lang="ts">
	import { ucfirst } from "$lib/string";
	import { Control, Field, FieldErrors, Label } from "formsnap";
	import type { HTMLInputTypeAttribute } from "svelte/elements";
	import { formFieldProxy, type SuperForm } from "sveltekit-superforms";
	let {
		superform,
		type,
		name
	}: {
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		superform: SuperForm<any>;
		type: HTMLInputTypeAttribute;
		name: string;
	} = $props();

	const label = ucfirst(name);
	const placeholder = ucfirst(name);
	const { value } = formFieldProxy(superform, name);
</script>

<div class="form-control">
	<Field {name} form={superform}>
		<Control>
			{#snippet children({ props })}
				<Label class="label">
					<span class="label-text">{label}</span>
				</Label>
				<input class="input input-bordered" {...props} {type} {placeholder} bind:value={$value} />
			{/snippet}
		</Control>
		<div class="pl-1 pt-1 text-xs text-red-400"><FieldErrors /></div>
	</Field>
</div>
