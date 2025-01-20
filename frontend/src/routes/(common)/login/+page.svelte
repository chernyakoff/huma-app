<script lang="ts">
	import Meta from '$lib/components/layout/Meta.svelte';
	import TextInput from '$lib/components/ui/form/TextInput.svelte';
	import { loginSchema } from '$lib/zod-schemas';
	import SuperDebug, { superForm, setMessage, setError, defaults } from 'sveltekit-superforms';
	import { zod } from 'sveltekit-superforms/adapters';
	import { Field, Control, Label, FieldErrors } from 'formsnap';
	import Icon from '@iconify/svelte';
	import { loginMutation } from '$lib/client/@tanstack/svelte-query.gen';

	import { createMutation } from '@tanstack/svelte-query';

	const mutation = createMutation({ ...loginMutation() });
	let showPassword = $state(false);

	const superform = superForm(
		{
			email: 'chernyakoff@gmail.com',
			password: '12345aA-'
		},
		{
			SPA: true,
			validators: zod(loginSchema),
			onUpdate({ form }) {
				$mutation.mutate({ body: form.data });
			}
		}
	);
	const { form: formData, enhance } = superform;
</script>

<Meta title="Логин" />
{#if $mutation.isSuccess}
	{$mutation.data}
{/if}
<div class="hero flex flex-grow items-center justify-center">
	<div class="card w-full max-w-sm shrink-0 bg-base-100 shadow-2xl">
		<form class="card-body" method="POST" use:enhance>
			<TextInput {superform} name="email" type="email" />
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
									type={showPassword ? 'text' : 'password'}
									placeholder="Password"
									bind:value={$formData.password}
								/>
								<label class="swap swap-rotate">
									<input type="checkbox" class="hidden" bind:checked={showPassword} />
									<Icon icon="ph:eye" class="icon24 swap-off" />
									<Icon icon="ph:eye-slash" class="icon24  swap-on" />
								</label>
							</label>
						{/snippet}
					</Control>
					<div class="pl-1 pt-1 text-xs text-red-400"><FieldErrors /></div>
				</Field>
			</div>
			<div class="form-control">
				<label class="label" for="forgot">
					<a href="/" class="link-hover link label-text-alt">Forgot password?</a>
				</label>
			</div>
			<div class="form-control mt-6">
				<button type="submit" class="btn btn-primary">Login</button>
			</div>
		</form>
	</div>
</div>
<SuperDebug data={$formData} />
