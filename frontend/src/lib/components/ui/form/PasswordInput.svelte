<script lang="ts">
	import { Field, Control, Label, FieldErrors } from "formsnap";
	import { formFieldProxy, type SuperForm } from "sveltekit-superforms";

	import Icon from "@iconify/svelte";
	let {
		superform,
		confirm = false
	}: {
		superform: SuperForm<any>;
		confirm?: boolean;
	} = $props();

	let showPassword = $state(false);

	const { value: passwordValue } = formFieldProxy(superform, "password");

	const { value: confirmValue } = formFieldProxy(superform, "confirmPassword");
</script>

<div class="form-control">
	<Field name="password" form={superform}>
		<Control>
			{#snippet children({ props })}
				<Label class="label ">
					<span class="label-text">Password</span>
				</Label>
				<label class="input input-bordered flex items-center gap-2">
					<input
						class="grow"
						{...props}
						type={showPassword ? "text" : "password"}
						placeholder="Password"
						bind:value={$passwordValue}
					/>
					<label class="swap swap-rotate">
						<input type="checkbox" class="hidden" bind:checked={showPassword} />
						<Icon icon="ph:eye" class="icon24 swap-off" />
						<Icon icon="ph:eye-slash" class="icon24 swap-on" />
					</label>
				</label>
			{/snippet}
		</Control>
		<div class="pl-1 pt-1 text-xs text-red-400"><FieldErrors /></div>
	</Field>
</div>
{#if confirm}
	<div class="form-control">
		<Field name="confirmPassword" form={superform}>
			<Control>
				{#snippet children({ props })}
					<Label class="label ">
						<span class="label-text">Confirm password</span>
					</Label>
					<label class="input input-bordered flex items-center gap-2">
						<input
							class="grow"
							{...props}
							type={showPassword ? "text" : "password"}
							placeholder="Password"
							bind:value={$confirmValue}
						/>
						<label class="swap swap-rotate">
							<input type="checkbox" class="hidden" bind:checked={showPassword} />
							<Icon icon="ph:eye" class="icon24 swap-off" />
							<Icon icon="ph:eye-slash" class="icon24 swap-on" />
						</label>
					</label>
				{/snippet}
			</Control>
			<div class="pl-1 pt-1 text-xs text-red-400"><FieldErrors /></div>
		</Field>
	</div>
{/if}
