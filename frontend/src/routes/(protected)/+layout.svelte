<script lang="ts">
	import { goto } from "$app/navigation";
	import type { GetUserByIdRow } from "$lib/client";
	import ThemeSwap from "$lib/components/ui/ThemeSwap.svelte";
	import { authState } from "$lib/states.svelte";
	import Icon from "@iconify/svelte";
	import type { Snippet } from "svelte";

	interface Props {
		children?: Snippet;
		data?: {
			user: GetUserByIdRow;
		};
	}

	let { children, data }: Props = $props();

	if (!data?.user) {
		goto("/login");
	} else {
		authState.set(data.user);
	}
</script>

<div class="navbar bg-base-300">
	<div class="navbar-start">
		<div class="dropdown">
			<div tabindex="0" role="button" class="btn btn-ghost lg:hidden">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-5 w-5"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M4 6h16M4 12h8m-8 6h16"
					/>
				</svg>
			</div>
			<ul class="menu dropdown-content menu-sm z-[1] mt-3 w-52 rounded-box bg-base-100 p-2 shadow">
				<li><a href="/app/users" tabindex="-1">Юзеры</a></li>
				<li>
					<a href="/" tabindex="-2">Parent</a>
					<ul class="p-2">
						<li><a href="/" tabindex="-1">Submenu 1</a></li>
						<li><a href="/" tabindex="-1">Submenu 2</a></li>
					</ul>
				</li>
				<li><a href="/" tabindex="-1">Item 3</a></li>
			</ul>
		</div>
		<a href="/" tabindex="-1" class="btn btn-ghost text-xl">HUMA</a>
	</div>
	<div class="navbar-center hidden lg:flex">
		<ul class="menu menu-horizontal px-1">
			<li><a href="/" tabindex="-1">Item 1</a></li>
			<li>
				<details>
					<summary>Parent</summary>
					<ul class="p-2">
						<li><a href="/" tabindex="-1">Submenu 1</a></li>
						<li><a href="/" tabindex="-1">Submenu 2</a></li>
					</ul>
				</details>
			</li>
			<li><a href="/about" tabindex="-1">About</a></li>
		</ul>
	</div>
	<div class="navbar-end">
		<div class="flex gap-4 pr-4">
			<a href="/logout" title="Logout">
				<Icon icon="ph:sign-out-light" class="icon24" />
			</a>
			<ThemeSwap />
		</div>
	</div>
</div>
{@render children?.()}
