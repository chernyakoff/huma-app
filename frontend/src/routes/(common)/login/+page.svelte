<script lang="ts">
	import { goto } from "$app/navigation";
	import { loginMutation } from "$lib/client/@tanstack/svelte-query.gen";
	import TextInput from "$lib/components/ui/form/TextInput.svelte";
	import { authState } from "$lib/states.svelte";
	import { loginSchema } from "$lib/zod-schemas";
	import { createMutation } from "@tanstack/svelte-query";
	import { superForm } from "sveltekit-superforms";
	import { zod } from "sveltekit-superforms/adapters";

	import PasswordInput from "$lib/components/ui/form/PasswordInput.svelte";
	import * as toast from "$lib/toast";

	if (authState.valid) {
		goto("/app");
	}
	const mutation = createMutation({
		...loginMutation(),
		onSuccess() {
			goto("/app");
		},
		onError(error) {
			toast.error(error.detail);
		}
	});

	// TODO form defaults
	const superform = superForm(
		{
			email: "chernyakoff@gmail.com",
			password: "12345aA-"
		},
		{
			SPA: true,
			validators: zod(loginSchema),
			onUpdate({ form }) {
				if (form.valid) {
					$mutation.mutate({ body: form.data });
				}
			}
		}
	);
	const { enhance } = superform;
</script>

<div class="hero flex flex-grow items-center justify-center">
	<div class="w-full max-w-sm shrink-0">
		<div class="card bg-base-100 shadow-sm">
			<form class="card-body" method="POST" use:enhance>
				<TextInput {superform} name="email" type="email" />
				<PasswordInput {superform} />
				<div class="form-control">
					<label class="label" for="forgot">
						<a href="/" class="link-hover link label-text-alt">Forgot password?</a>
					</label>
				</div>
				<div class="form-control mt-3">
					<button type="submit" class="btn btn-primary text-white">Login</button>
				</div>
			</form>
		</div>
		<div class="mt-2 text-center">
			<span class="label-text-alt"
				>Don't have an account? <a href="/register" class="link-hover link link-success">Register</a
				></span
			>
		</div>
	</div>
</div>
