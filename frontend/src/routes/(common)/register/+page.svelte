<script lang="ts">
	import { goto } from "$app/navigation";
	import { registerMutation } from "$lib/client/@tanstack/svelte-query.gen";
	import TextInput from "$lib/components/ui/form/TextInput.svelte";
	import { registerSchema } from "$lib/zod-schemas";
	import { createMutation } from "@tanstack/svelte-query";
	import { superForm } from "sveltekit-superforms";
	import { zod } from "sveltekit-superforms/adapters";

	import { authState } from "$lib/states.svelte";

	import PasswordInput from "$lib/components/ui/form/PasswordInput.svelte";
	import * as toast from "$lib/toast";

	if (authState.valid) {
		goto("/app");
	}
	const mutation = createMutation({
		...registerMutation(),
		onSuccess() {
			goto("/verify");
		},
		onError(error) {
			toast.error(error.detail);
		}
	});

	// TODO form defaults
	const superform = superForm(
		{
			email: "chernyakoff@gmail.com",
			password: "12345aA-",
			confirmPassword: "12345aA-"
		},
		{
			SPA: true,
			validators: zod(registerSchema),
			onUpdate({ form }) {
				if (form.valid) {
					// eslint-disable-next-line @typescript-eslint/no-unused-vars
					const { confirmPassword, ...data } = form.data;
					localStorage.setItem("verify", data.email);
					$mutation.mutate({ body: data });
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
				<PasswordInput {superform} confirm={true} />
				<div class="form-control mt-3">
					<button type="submit" class="btn btn-primary text-white">Register</button>
				</div>
			</form>
		</div>
		<div class="mt-2 text-center">
			<span class="label-text-alt"
				>Already have an account? <a href="/login" class="link-hover link link-success">Login</a
				></span
			>
		</div>
	</div>
</div>
